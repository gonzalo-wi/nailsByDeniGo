package shared

import "math"

// RoundMoney redondea un monto a 2 decimales.
func RoundMoney(amount float64) float64 {
	return math.Round(amount*100) / 100
}
