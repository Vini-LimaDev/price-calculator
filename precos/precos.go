package precos

import (
	"fmt"
	"strings"

	"price-calculator/conversao"
	leituraArq "price-calculator/leitura_arq"
)

type CalculadoraDePrecos struct {
	FileManager          leituraArq.FileManager `json:"-"` // não será serializado no JSON
	Taxas                float64                `json:"taxa"`
	Precos               []float64              `json:"precos"`
	PrecosIncluindoTaxas []string               `json:"precos_incluindo_taxas"`
}

// Lê os preços de "precos.txt". Se não existir, pergunta ao usuário, cria o arquivo e então reabre p/ leitura.
func (c *CalculadoraDePrecos) LerPrecosDoArquivo() error {

	linhas, err := c.FileManager.LerArquivo()

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
		// adiciona ao slice formatado (duas casas decimais) que será serializado no JSON
		c.PrecosIncluindoTaxas = append(c.PrecosIncluindoTaxas, fmt.Sprintf("%.2f", taxaComPreco))
	}
	fmt.Println()
	c.FileManager.CriaJSON(c)
}

// CalcularComPrecos calcula taxas para uma lista de preços já fornecida (para uso em APIs)
func (c *CalculadoraDePrecos) CalcularComPrecos(precos []float64) {
	c.Precos = precos
	c.PrecosIncluindoTaxas = []string{} // limpa resultados anteriores

	for _, preco := range c.Precos {
		taxaComPreco := preco * (1 + c.Taxas)
		// adiciona ao slice formatado (duas casas decimais) que será serializado no JSON
		c.PrecosIncluindoTaxas = append(c.PrecosIncluindoTaxas, fmt.Sprintf("%.2f", taxaComPreco))
	}
}

func NovoCalculadoraDePrecos(fm leituraArq.FileManager, taxas float64) *CalculadoraDePrecos {
	return &CalculadoraDePrecos{
		FileManager:          fm,
		Taxas:                taxas,
		Precos:               []float64{},
		PrecosIncluindoTaxas: []string{},
	}
}

// NovoCalculadoraDePrecosParaAPI cria uma instância para uso em APIs (sem FileManager)
func NovoCalculadoraDePrecosParaAPI(taxas float64) *CalculadoraDePrecos {
	return &CalculadoraDePrecos{
		Taxas:                taxas,
		Precos:               []float64{},
		PrecosIncluindoTaxas: []string{},
	}
}
