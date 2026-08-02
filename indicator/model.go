package indicator

import "quantum-signal/models"

type Indicator interface {
	Name() string
	Calculate(candles []models.Candle) []float64
}

type SignalGenerator interface {
	GenerateSignal(candles []models.Candle, indicators map[string][]float64) models.Signal
}
