package leituraarq

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LerArquivo(nomeArquivo string) ([]string, error) {
	// 1) Tenta abrir
	arquivo, err := os.Open(nomeArquivo)

	if err != nil {
		// Se não existe, coleta do usuário e cria
		if os.IsNotExist(err) {
			fmt.Println("Não foi encontrado o arquivo. Defina os valores que deseja.")
			precosArq := coletarPrecosDoUsuario()
			if err := criarArquivo(precosArq); err != nil {
				return nil, fmt.Errorf("erro ao criar o arquivo: %w", err)
			}
			// Reabre para leitura
			arquivo, err = os.Open(nomeArquivo)
			if err != nil {
				return nil, errors.New("erro ao abrir o arquivo após criar")
			}
		} else {
			return nil, errors.New("erro ao abrir o arquivo")
		}
	}

	// 2) Lê as linhas
	var linhas []string
	scanner := bufio.NewScanner(arquivo)

	for scanner.Scan() {
		linhas = append(linhas, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		arquivo.Close()
		return nil, fmt.Errorf("erro ao ler o arquivo: %w", err)
	}

	arquivo.Close()
	return linhas, nil
}

func CriaJSON(nomeArquivo string, data interface{}) error {
	arq, err := os.Create(nomeArquivo)
	if err != nil {
		return errors.New("erro ao criar o arquivo JSON")
	}

	encoder := json.NewEncoder(arq)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)

	if err != nil {
		arq.Close()
		return errors.New("erro ao escrever dados no arquivo JSON")
	}
	arq.Close()
	return nil
}

func coletarPrecosDoUsuario() []float64 {
	fmt.Println("Digite os preços desejados (um por linha). Digite 'fim' para encerrar:")
	scanner := bufio.NewScanner(os.Stdin)
	var out []float64
	for {
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if strings.EqualFold(input, "fim") {
			break
		}
		v, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Println("Entrada inválida. Digite um número ou 'fim' para encerrar.")
			continue
		}
		out = append(out, v)
	}
	return out
}

func criarArquivo(precos []float64) error {
	f, err := os.Create("precos.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, p := range precos {
		if _, err := fmt.Fprintf(w, "%.2f\n", p); err != nil {
			return err
		}
	}
	return w.Flush()
}
