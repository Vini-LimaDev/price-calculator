package main

import (
	"fmt"
	leituraArq "price-calculator/leitura_arq"
	"price-calculator/precos"
)

func main() {
	taxas := []float64{0.1, 0.07, 0.15}

	for _, taxa := range taxas {
		fm := leituraArq.New("precos.txt", fmt.Sprintf("resultado_%.0f.json", taxa*100))
		prec := precos.NovoCalculadoraDePrecos(*fm, taxa)
		prec.Calcular()
	}
}
