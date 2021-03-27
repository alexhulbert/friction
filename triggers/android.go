package triggers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type message struct {
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
	Event  string `json:"event"`
}

func processMsg(msg message, channel chan bool) {
	if msg.Event == "alarm_alert_start" || msg.Event == "alarm_alert_dismiss" {
		// "alert_start" and "alert_dismiss" indicate that an alarm just started
		channel <- true
	} else if msg.Event == "alarm_snooze_clicked" {
		// "snooze_clicked" indicates that the alarm should be temporarily paused
		channel <- false
	} else {
		// Unhandled events are logged for debugging purposes
		fmt.Println("> received unhandled event: ", msg.Event)
	}
}

// Android connects to SleepAsAndroid and returns a channel that
// emits true when an alarm occurs or false when the alarm is snoozed
func Android() chan bool {
	channel := make(chan bool)
	// This function executes when SleepAsAndroid sends a json payload to this http server
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// Decode the json sent to this http server into the msg variable
		var msg message
		err := json.NewDecoder(req.Body).Decode(&msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		// Send true or false through the alarm channel if an alarm was started or snoozed
		processMsg(msg, channel)
	})
	// Start an http server which listens on port 8090
	go http.ListenAndServe(":8090", nil)
	fmt.Println("> webhook server listening on port 8090")
	return channel
}
