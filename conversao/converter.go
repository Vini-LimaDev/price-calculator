package conversao

import (
	"errors"
	"strconv"
)

// StringParaFloat64 converte uma string para float64.
func StringParaFloat64(strings []string) ([]float64, error) {
	var floats []float64
	for _, stringVal := range strings {
		floatVal, err := strconv.ParseFloat(stringVal, 64)
		if err != nil {
			return nil, errors.New("falha ao converter strings")
		}

		floats = append(floats, floatVal)
	}
	return floats, nil
}
