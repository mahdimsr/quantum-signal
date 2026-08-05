package scheduler

import (
	"fmt"
	"log"
	"quantum-signal/candles"
	"quantum-signal/indicator"
	"quantum-signal/strategy"
	"time"
)

func Start(minutes int) {

	fmt.Println("⏲️ Schedule started...")

	//ctx := context.Background()

	for {

		now := time.Now().UTC()
		nextRun := getNextRunMark(now, minutes)
		waitDuration := time.Until(nextRun)

		/*rabbitService, err := rabbitmq.NewRabbitMqService()
		if err != nil {
			log.Fatalf("Failed to create rabbitmq service: %s", err)
		}*/

		log.Printf("⌚ Next execution at: %s (waiting %v)", nextRun.Format(time.RFC3339), waitDuration)

		time.Sleep(waitDuration)

		log.Printf("🔄 Executing trading task at: %s", time.Now().UTC().Format(time.RFC3339))

		candleService := candles.NewCandleService(candles.EXCHANGE, "BTCUSDT", "15m", 100)
		candlesList := candleService.GetCandles()

		utbotValues := indicator.NewUTBot(3, 5, false)
		utbotStrategy := strategy.NewUTBotStrategy(utbotValues)
		signals := utbotStrategy.GenerateSignals(candlesList)

		lastSignal := signals[len(signals)-1]

		signalTime := time.UnixMilli(lastSignal.Timestamp).UTC()

		if IsWithinDuration(now, signalTime, 15) {

			// TODO: send signal to rabbit
		}

	}

}

func getNextRunMark(now time.Time, cycleMinutes int) time.Time {

	nowMinutes := now.Minute()

	nextMinutes := ((nowMinutes / cycleMinutes) + 1) * cycleMinutes

	nextRun := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		nextMinutes,
		0,
		0,
		time.UTC,
	)

	/*if nextMinutes >= 60 {
		nextRun = nextRun.Add(time.Hour)
		nextRun = time.Date(
			nextRun.Year(),
			nextRun.Month(),
			nextRun.Day(),
			nextRun.Hour(),
			0,
			0,
			0,
			time.UTC,
		)
	}*/

	return nextRun
}

func IsWithinDuration(t1, t2 time.Time, threshold time.Duration) bool {
	return t1.Sub(t2).Abs() < threshold
}
