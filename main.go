package main

import (
	"fmt"
	"time"

	"./beeper"
	"./serial"
	"./triggers"
)

// SERIALPORT should be set to the path to the serial port used for communication with the load sensors
var SERIALPORT = "/dev/ttyUSB0"

// THRESHOLD controls the weight threshold used to determine when a person is in bed
var THRESHOLD = 1150000

// DURATION controls the amount of time after the alarm goes off to prevent the user from going into bed
var DURATION = "20s"

// TIME controls the time of the alarm. Leave this unset if you want to control the time with SleepAsAndroid
var TIME string // = "9:00AM"

func main() {
	// Preload the alarm audio file and initialize the class that plays it
	beeper.Init()

	// Connect to the load sensors / arduino over serial and start listening for changes in weight
	inBedChan := serial.Init(SERIALPORT, THRESHOLD)

	// Alarm events are either provided by an HTTP server that connects to the
	// SleepAsAndroid app, or by a daily scheduler that fires events at a specific
	// time indicated by the TIME variable
	var alarmChan chan bool
	if TIME == "" {
		alarmChan = triggers.Android()
	} else {
		alarmChan = triggers.Static(TIME)
	}

	fmt.Println("> initialized")
	handler(inBedChan, alarmChan)
}

func handler(inBed chan bool, alarm chan bool) {
	// The time at which the alarm shouldn't go off if the user gets into bed
	alarmExpiry := time.Now()
	// Whether the user is currently bed
	isInBed := false

	// Parse duration string into a time.Duration object
	duration, err := time.ParseDuration(DURATION)
	if err != nil {
		panic("Invalid duration specified: " + DURATION)
	}

	// Main loop of the program
	for {
		select {
		case enabled := <-alarm:
			// This block is entered with enabled = true if the alarm just went off
			// It is also entered with enabled = false if the alarm was just snoozed
			if enabled {
				// When the alarm goes off, update the alarm expiry time
				alarmExpiry = time.Now().Add(duration)
				// If the user is also in bed when the alarm goes off, start beeping
				if isInBed {
					beeper.Play()
				}
			} else {
				// If the alarm has just been snoozed, stop beeping
				beeper.Stop()
			}
			fmt.Println("> alarm enabled status:", enabled)
		case inBed := <-inBed:
			// This block is entered when the user gets in or out of bed
			// update the isInBed variable to match whether the user is in bed
			isInBed = inBed
			if isInBed {
				// If the user is in bed and its before the alarm expiry time, start beeping
				if time.Now().Before(alarmExpiry) {
					beeper.Play()
				}
			} else {
				// If the user isn't in bed, stop beeping
				beeper.Stop()
			}
			fmt.Println("> in bed status:", inBed)
		}
	}
}
