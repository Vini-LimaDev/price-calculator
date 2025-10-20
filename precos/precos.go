package precos

type CalculadoraDePrecos struct {
	Taxas          float64
	Precos         []float64
	TaxasComPrecos map[string]float64
}
