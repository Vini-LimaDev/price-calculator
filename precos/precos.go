package precos

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"price-calculator/conversao"
)

type CalculadoraDePrecos struct {
	Taxas  float64
	Precos []float64
}

// Lê os preços de "precos.txt". Se não existir, pergunta ao usuário, cria o arquivo e então reabre p/ leitura.
func (c *CalculadoraDePrecos) LerPrecosDoArquivo() error {
	const nome = "precos.txt"

	// 1) Tenta abrir
	arquivo, err := os.Open(nome)
	if err != nil {
		// Se não existe, coleta do usuário e cria
		if os.IsNotExist(err) {
			fmt.Println("Não foi encontrado o arquivo. Defina os valores que deseja.")
			precosArq := coletarPrecosDoUsuario()
			if err := criarArquivo(precosArq); err != nil {
				return fmt.Errorf("erro ao criar o arquivo: %w", err)
			}
			// Reabre para leitura
			arquivo, err = os.Open(nome)
			if err != nil {
				return fmt.Errorf("erro ao reabrir o arquivo após criar: %w", err)
			}
		} else {
			return fmt.Errorf("erro ao abrir o arquivo: %w", err)
		}
	}
	defer arquivo.Close()

	// 2) Lê as linhas
	var linhas []string
	scanner := bufio.NewScanner(arquivo)
	for scanner.Scan() {
		linhas = append(linhas, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("erro ao ler o arquivo: %w", err)
	}

	// 3) Converte para float64
	precos, err := conversao.StringParaFloat64(linhas)
	if err != nil {
		arquivo.Close()
		return fmt.Errorf("erro ao converter os preços: %w", err)
	}

	c.Precos = precos
	return nil
}

func (c *CalculadoraDePrecos) Calcular() {
	if err := c.LerPrecosDoArquivo(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Taxa aplicada: %.2f%%\n\n", c.Taxas*100)
	fmt.Printf("%-12s | %-15s\n", "Preço", "Preço com taxa")
	fmt.Println(strings.Repeat("-", 30))

	for _, preco := range c.Precos {
		taxaComPreco := preco * (1 + c.Taxas)
		fmt.Printf("%-12.2f | %-15.2f\n", preco, taxaComPreco)
	}
	fmt.Println()
}

func NovoCalculadoraDePrecos(taxas float64) *CalculadoraDePrecos {
	return &CalculadoraDePrecos{Taxas: taxas, Precos: []float64{}}
}

// === Helpers ===

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
