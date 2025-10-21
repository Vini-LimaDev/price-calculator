package precos

import (
	"fmt"
)

type CalculadoraDePrecos struct {
	Taxas          float64
	Precos         []float64
	TaxasComPrecos map[string]float64
}

func (c *CalculadoraDePrecos) Calcular() {
	resultado := make(map[string]float64)

	for _, preco := range c.Precos {
		resultado[fmt.Sprintf("%.2f", preco)] = preco * (1 + c.Taxas)
	}

	fmt.Println("Taxa:", c.Taxas)
	fmt.Println("Preços com taxa:", resultado)
}

func NovoCalculadoraDePrecos(taxas float64) *CalculadoraDePrecos {
	// Pegar os precos apartir de um arquivo

	return &CalculadoraDePrecos{
		Precos: []float64{10, 20, 30},
		Taxas:  taxas,
	}
}
