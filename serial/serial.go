package serial

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"

	serial "github.com/tarm/goserial"
)

func serialLoop(s io.ReadWriteCloser, threshold int64, channel chan bool) {
	wasInBed := false
	for {
		// Read a new line from the serial connection
		reader := bufio.NewReader(s)
		rawLine, err := reader.ReadBytes('\n')
		if err != nil {
			// TODO: Recover from unplugged cable here
			panic(err)
		}
		// Remove the newline and parse the line as an integer
		line := strings.ReplaceAll(string(rawLine), "\r\n", "")
		reading, _ := strconv.ParseInt(line, 10, 64)

		// Sometimes the serial connection skips a digit
		// so ignore sensor readings that are suspiciously low
		if reading > threshold/10 {
			// Emit a bool to the channel if the in-bed status has changed
			inBedNow := reading > threshold
			if inBedNow != wasInBed {
				channel <- inBedNow
				wasInBed = inBedNow
			}
		}
	}
}

// Init initializes the serial server
// Returns a channel that emits a bool indicating
// whether the user is in bed when the InBed status changes
func Init(port string, threshold int) chan bool {
	// Connect to the serial port
	cfg := &serial.Config{Name: port, Baud: 57600}
	s, err := serial.OpenPort(cfg)
	if err != nil {
		log.Fatal("Serial port is not plugged into " + port)
	}

	// Make the channel and start the loop that reads lines from the serial port
	channel := make(chan bool)
	go serialLoop(s, int64(threshold), channel)
	return channel
}
