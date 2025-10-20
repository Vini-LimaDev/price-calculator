package main

import "fmt"

func main() {
	precos := []float64{100, 20, 30}
	taxas := []float64{0, 0.1, 0.07, 0.15}

	resultado := make(map[float64][]float64)

	for _, taxa := range taxas {
		precosComTaxa := make([]float64, len(precos))

		for i, preco := range precos {
			precosComTaxa[i] = preco * (1 + taxa)
		}
		resultado[taxa] = precosComTaxa
	}

	for taxa, precos := range resultado {
		formattedPrecos := make([]string, len(precos))
		for i, preco := range precos {
			formattedPrecos[i] = fmt.Sprintf("R$%.2f", preco)
		}
		fmt.Printf("Taxa: %.2f%% - Preços: %v\n", taxa*100, formattedPrecos)
	}
}
