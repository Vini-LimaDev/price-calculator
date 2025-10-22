package precos

import (
	"fmt"
	"strings"

	"price-calculator/conversao"
	leituraArq "price-calculator/leitura_arq"
)

type CalculadoraDePrecos struct {
	Taxas                float64
	Precos               []float64
	PrecosIncluindoTaxas []string
}

const nome = "precos.txt"

// Lê os preços de "precos.txt". Se não existir, pergunta ao usuário, cria o arquivo e então reabre p/ leitura.
func (c *CalculadoraDePrecos) LerPrecosDoArquivo() error {

	linhas, err := leituraArq.LerArquivo(nome)
	if err != nil {
		return fmt.Errorf("erro ao ler os preços do arquivo: %w", err)
	}

	// 3) Converte para float64
	precos, err := conversao.StringParaFloat64(linhas)
	if err != nil {
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
	leituraArq.CriaJSON(fmt.Sprintf("resultado_%.0f.json", c.Taxas*100), c)
}

func NovoCalculadoraDePrecos(taxas float64) *CalculadoraDePrecos {
	return &CalculadoraDePrecos{Taxas: taxas, Precos: []float64{}}
}
