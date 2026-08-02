package strategy

import (
	"quantum-signal/indicator"
	"quantum-signal/models"
)

type UTBotStrategy struct {
	ind indicator.Indicator // تزریق وابستگی به اینترفیس
}

func NewUTBotStrategy(ind indicator.Indicator) *UTBotStrategy {
	return &UTBotStrategy{ind: ind}
}

// GenerateSignals پیاده‌سازی اینترفیس استراتژی
func (s *UTBotStrategy) GenerateSignals(candles []models.Candle) []models.Signal {
	signals := make([]models.Signal, 0)

	// 1. گرفتن خط استاپ از اندیکاتور
	stops := s.ind.Calculate(candles)

	// 2. تعیین Source برای مقایسه (باید با Source اندیکاتور هماهنگ باشد)
	// نکته: اگر اندیکاتور UTBot با HeikinAshi=true ساخته شده باشد،
	// برای دقت ۱۰۰٪ مطابق تریدینگ‌ویو، بهتر است اینجا هم از HA Close استفاده شود.
	// برای سادگی و تطابق با اینترفیس خالص، اینجا از Close معمولی استفاده می‌کنیم.
	src := make([]float64, len(candles))
	for i := range candles {
		src[i] = candles[i].Close
	}

	// 3. منطق Crossover (ترجمه دقیق Pine Script)
	// ema(src, 1) در پاین اسکریپت دقیقاً برابر با خود src است.
	// above = crossover(ema, xATRTrailingStop)
	// buy = src > xATRTrailingStop and above
	for i := 1; i < len(candles); i++ {
		if stops[i] == 0 || stops[i-1] == 0 {
			continue // دیتای کافی برای مقایسه وجود ندارد
		}

		currSrc := src[i]
		prevSrc := src[i-1]
		currStop := stops[i]
		prevStop := stops[i-1]

		// منطق Buy: قیمت فعلی بالای استاپ است و قیمت قبلی پایین یا مساوی استاپ بوده (Crossover به بالا)
		if currSrc > currStop && prevSrc <= prevStop {
			signals = append(signals, models.Signal{
				Type:      models.BUY,
				Timestamp: candles[i+1].Time,
				Price:     candles[i+1].Close,
				Reason:    "UT Bot: Price crossed above Trailing Stop",
				Candle:    candles[i+1],
			})
		}

		// منطق Sell: قیمت فعلی پایین استاپ است و قیمت قبلی بالا یا مساوی استاپ بوده (Crossover به پایین)
		if currSrc < currStop && prevSrc >= prevStop {
			signals = append(signals, models.Signal{
				Type:      models.SELL,
				Timestamp: candles[i+1].Time,
				Price:     candles[i+1].Close,
				Reason:    "UT Bot: Price crossed below Trailing Stop",
				Candle:    candles[i+1],
			})
		}
	}

	return signals
}
