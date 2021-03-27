package beeper

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var playing = false
var streamer beep.StreamSeekCloser

// Init initializes the speaker and loads the alarm file
func Init() {
	// Open the file in this directory called "alarm.mp3"
	f, err := os.Open("beeper/alarm.mp3")
	if err != nil {
		log.Fatal(err)
	}

	// Decode the contents of the file as an mp3
	var format beep.Format
	streamer, format, err = mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the speaker class with the sample rate of the mp3 file
	sr := format.SampleRate
	speaker.Init(sr, sr.N(time.Second/10))
}

// Play plays the beeping sound on a loop
func Play() {
	if !playing {
		speaker.Play(beep.Loop(-1, streamer))
		playing = true
	}
}

// Stop stops the beeping sound
func Stop() {
	if playing {
		speaker.Clear()
		playing = false
	}
}
