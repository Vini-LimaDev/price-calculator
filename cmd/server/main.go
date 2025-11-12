package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	leituraArq "price-calculator/leitura_arq"
)

type CalcularReq struct {
	Precos []float64 `json:"precos"`
	Taxas  []float64 `json:"taxas"` // em fração: 0.10, 0.07...
}

type ResultadoGo struct {
	Taxa                 float64   `json:"taxa"`
	Precos               []float64 `json:"precos"`
	PrecosIncluindoTaxas []string  `json:"precos_incluindo_taxas"`
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func round2(n float64) float64 {
	return math.Round(n*100) / 100
}

func calcular(precos []float64, taxa float64) ResultadoGo {
	out := make([]string, len(precos))
	for i, p := range precos {
		out[i] = fmt.Sprintf("%.2f", round2(p*(1+taxa)))
	}
	return ResultadoGo{
		Taxa:                 taxa,
		Precos:               precos,
		PrecosIncluindoTaxas: out,
	}
}

func salvarJSONResultado(res ResultadoGo) error {
	// nome no padrão resultado_XX.json
	percent := int(math.Round(res.Taxa * 100))
	outName := fmt.Sprintf("resultado_%d.json", percent)

	fm := leituraArq.New("", outName)
	return fm.CriaJSON(res) // reaproveita teu FileManager para persistir JSON
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"erro": "use POST"})
		return
	}

	// aceita multipart/form-data (file) ou text/plain no corpo
	ct := r.Header.Get("Content-Type")
	var body []byte
	var err error

	if strings.HasPrefix(ct, "multipart/form-data") {
		f, _, ferr := r.FormFile("file")
		if ferr != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"erro": "enviar arquivo no campo 'file'"})
			return
		}
		defer f.Close()
		body, err = io.ReadAll(f)
	} else {
		body, err = io.ReadAll(r.Body)
	}
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"erro": "falha ao ler arquivo"})
		return
	}

	// garante diretório atual (apenas por segurança se rodar de outro CWD)
	dest := filepath.Join(".", "precos.txt")
	if err := os.WriteFile(dest, body, 0o644); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"erro": "falha ao salvar precos.txt"})
		return
	}

	linhas := 0
	for _, ln := range strings.Split(string(body), "\n") {
		if strings.TrimSpace(ln) != "" {
			linhas++
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":           true,
		"arquivo":      "precos.txt",
		"linhas_lidas": linhas,
	})
}

func handleCalcular(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"erro": "use POST"})
		return
	}
	var req CalcularReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"erro": "JSON inválido"})
		return
	}
	if len(req.Precos) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"erro": "lista de precos vazia"})
		return
	}
	if len(req.Taxas) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"erro": "lista de taxas vazia"})
		return
	}
	for _, t := range req.Taxas {
		if t < 0 || math.IsNaN(t) || math.IsInf(t, 0) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"erro": "taxa inválida"})
			return
		}
	}

	var resultados []ResultadoGo
	for _, taxa := range req.Taxas {
		res := calcular(req.Precos, taxa)
		_ = salvarJSONResultado(res) // ignora erro para não travar resposta; logaria em prod
		resultados = append(resultados, res)
	}
	writeJSON(w, http.StatusOK, resultados)
}

func handleRaiz(w http.ResponseWriter, r *http.Request) {
	msg := []string{
		"Price Calculator API",
		"POST /upload   -> multipart (file) ou text/plain; salva precos.txt",
		"POST /calcular -> { precos: number[], taxas: number[] } (taxas em fração, ex.: 0.10)",
	}
	_, _ = w.Write([]byte(strings.Join(msg, "\n")))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRaiz)
	mux.HandleFunc("/upload", handleUpload)
	mux.HandleFunc("/calcular", handleCalcular)

	addr := ":8080"
	fmt.Println("API ouvindo em", addr)
	if err := http.ListenAndServe(addr, cors(mux)); err != nil {
		panic(err)
	}
}
