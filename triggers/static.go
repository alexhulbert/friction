package triggers

import (
	"time"

	"github.com/go-co-op/gocron"
)

// Static takes a time and emits true at that time each day
func Static(timeStr string) chan bool {
	alarmChannel := make(chan bool)
	militaryTimeStr := convertToMilitary(timeStr)
	scheduler := gocron.NewScheduler(time.Local)
	scheduler.Every(1).Day().At(militaryTimeStr).Do(func() {
		alarmChannel <- true
	})
	scheduler.StartAsync()
	return alarmChannel
}

func convertToMilitary(timeStr string) string {
	twelveHourLayout := "03:04PM"
	twentyFourHourLayout := "15:04"
	t, err := time.Parse(twelveHourLayout, timeStr)
	if err != nil {
		panic("Invalid time \"" + timeStr + "\" supplied to triggers.Static")
	}
	return t.Format(twentyFourHourLayout)
}
