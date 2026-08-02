package indicator

import (
	"math"
	"quantum-signal/models"
)

type UTBot struct {
	keyValue      float64
	atrPeriod     int
	useHeikinAshi bool
}

func NewUTBot(keyValue float64, atrPeriod int, useHeikinAshi bool) *UTBot {
	return &UTBot{
		keyValue:      keyValue,
		atrPeriod:     atrPeriod,
		useHeikinAshi: useHeikinAshi,
	}
}

func (u *UTBot) Name() string {
	return "UTBot"
}

// Calculate دقیقاً مطابق ساختار توافق شده: ورودی کندل، خروجی آرایه float64 (خط استاپ)
func (u *UTBot) Calculate(candles []models.Candle) []float64 {
	stops := make([]float64, len(candles))
	if len(candles) < u.atrPeriod+1 {
		return stops
	}

	// 1. تعیین Source (Heikin Ashi یا Close معمولی)
	src := make([]float64, len(candles))
	if u.useHeikinAshi {
		ha := u.calcHeikinAshi(candles)
		for i := range ha {
			src[i] = ha[i].Close
		}
	} else {
		for i := range candles {
			src[i] = candles[i].Close
		}
	}

	// 2. محاسبه ATR (ساده‌شده با RMA)
	atr := make([]float64, len(candles))
	trueRange := make([]float64, len(candles))
	for i := 1; i < len(candles); i++ {
		hl := candles[i].High - candles[i].Low
		hc := math.Abs(candles[i].High - candles[i-1].Close)
		lc := math.Abs(candles[i].Low - candles[i-1].Close)
		trueRange[i] = math.Max(hl, math.Max(hc, lc))
	}

	// مقداردهی اولیه ATR
	sum := 0.0
	for i := 1; i <= u.atrPeriod; i++ {
		sum += trueRange[i]
	}
	atr[u.atrPeriod] = sum / float64(u.atrPeriod)

	// محاسبه RMA برای بقیه
	for i := u.atrPeriod + 1; i < len(candles); i++ {
		atr[i] = (atr[i-1]*float64(u.atrPeriod-1) + trueRange[i]) / float64(u.atrPeriod)
	}

	// 3. محاسبه nLoss
	nLoss := make([]float64, len(candles))
	for i := range candles {
		nLoss[i] = u.keyValue * atr[i]
	}

	// 4. محاسبه xATRTrailingStop (ترجمه دقیق منطق Pine Script)
	for i := 1; i < len(candles); i++ {
		prevStop := stops[i-1]
		if prevStop == 0 {
			prevStop = src[i] // معادل nz(xATRTrailingStop[1], 0) برای شروع
		}
		prevSrc := src[i-1]

		if src[i] > prevStop && prevSrc > prevStop {
			stops[i] = math.Max(prevStop, src[i]-nLoss[i])
		} else if src[i] < prevStop && prevSrc < prevStop {
			stops[i] = math.Min(prevStop, src[i]+nLoss[i])
		} else if src[i] > prevStop {
			stops[i] = src[i] - nLoss[i]
		} else {
			stops[i] = src[i] + nLoss[i]
		}
	}

	return stops
}

// تابع کمکی محاسبه Heikin Ashi
func (u *UTBot) calcHeikinAshi(candles []models.Candle) []models.Candle {
	ha := make([]models.Candle, len(candles))
	if len(candles) == 0 {
		return ha
	}
	ha[0].Open = (candles[0].Open + candles[0].Close) / 2
	ha[0].Close = (candles[0].Open + candles[0].High + candles[0].Low + candles[0].Close) / 4
	ha[0].High = candles[0].High
	ha[0].Low = candles[0].Low

	for i := 1; i < len(candles); i++ {
		ha[i].Close = (candles[i].Open + candles[i].High + candles[i].Low + candles[i].Close) / 4
		ha[i].Open = (ha[i-1].Open + ha[i-1].Close) / 2
		ha[i].High = math.Max(candles[i].High, math.Max(ha[i].Open, ha[i].Close))
		ha[i].Low = math.Min(candles[i].Low, math.Min(ha[i].Open, ha[i].Close))
	}
	return ha
}
