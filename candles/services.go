package candles

import (
	"fmt"
	"quantum-signal/models"
	"time"
)

type CandleService struct {
	Repository Repository
	Count      int
	Timeframe  string
	Symbol     string
}

type Repository string

const (
	EXCHANGE   Repository = "exchange"
	METATRADER Repository = "metatrader"
)

func NewCandleService(repo Repository, symbol, timeframe string, count int) *CandleService {
	return &CandleService{
		Repository: repo,
		Count:      count,
		Symbol:     symbol,
		Timeframe:  timeframe,
	}
}

func (service CandleService) GetCandles() []models.Candle {

	var candles []models.Candle

	if service.Repository == EXCHANGE {

		candles, err := BinanceGetCandles(service.Symbol, service.Timeframe, service.Count)
		if err != nil {
			fmt.Printf("Error get candles from exchange: %s", err)
		}

		return candles
	}

	return candles
}

func calculateCandlesStartTime(candleCount int, timeframe string) (time.Time, error) {
	duration, err := parseTimeframe(timeframe)
	if err != nil {
		return time.Time{}, err
	}

	totalDuration := duration * time.Duration(candleCount)

	startTime := time.Now().UTC().Add(-totalDuration)

	return startTime, nil
}

func parseTimeframe(tf string) (time.Duration, error) {
	switch tf {

	case "1m":
		return 1 * time.Minute, nil
	case "3m":
		return 3 * time.Minute, nil
	case "5m":
		return 5 * time.Minute, nil
	case "15m":
		return 15 * time.Minute, nil
	case "30m":
		return 30 * time.Minute, nil

	case "1h", "60m":
		return 1 * time.Hour, nil
	case "2h", "120m":
		return 2 * time.Hour, nil
	case "4h", "240m":
		return 4 * time.Hour, nil
	case "6h", "360m":
		return 6 * time.Hour, nil
	case "12h", "720m":
		return 12 * time.Hour, nil

	case "1d", "1D":
		return 24 * time.Hour, nil
	case "1w", "1W":
		return 7 * 24 * time.Hour, nil

	default:
		return 0, fmt.Errorf("unsupported timeframe: %s", tf)
	}
}
