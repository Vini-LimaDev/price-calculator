package main

import (
	"price-calculator/precos"
)

func main() {
	taxas := []float64{0.1, 0.07, 0.15}

	for _, taxa := range taxas {
		prec := precos.NovoCalculadoraDePrecos(taxa)
		prec.Calcular()

	}
}
