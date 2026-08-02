package scheduler

import (
	"fmt"
	"log"
	"time"
)

func Start(minutes int) {

	fmt.Println("⏲️ Schedule started...")

	//ctx := context.Background()

	for {

		now := time.Now().UTC()
		nextRun := getNextRunMark(now, minutes)
		waitDuration := time.Until(nextRun)

		log.Printf("⌚ Next execution at: %s (waiting %v)", nextRun.Format(time.RFC3339), waitDuration)

		time.Sleep(waitDuration)

		log.Printf("🔄 Executing trading task at: %s", time.Now().UTC().Format(time.RFC3339))

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

	if nextMinutes >= 60 {
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
	}

	return nextRun
}
