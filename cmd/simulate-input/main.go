package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	mouseDevice    = "/dev/hidg0"
	keyboardDevice = "/dev/hidg1"
)

var (
	osOpenFile = os.OpenFile
	delay      = 300 * time.Millisecond
)

var openDevice = func(path string) (FileWriter, error) {
	return osOpenFile(path, os.O_WRONLY, 0666)
}

type FileWriter interface {
	Write(p []byte) (n int, err error)
	Close() error
}

func main() {
	for {
		fmt.Println("\nBluetooth HID Simulator")
		fmt.Println("======================")
		fmt.Println("1. Move mouse in circle")
		fmt.Println("2. Type test message")
		fmt.Println("Q. Quit")

		// Read a single character without waiting for Enter
		fmt.Print("\nSelect an option: ")
		var b = make([]byte, 1)
		os.Stdin.Read(b)

		input := string(b)

		switch input {
		case "1":
			moveMouseInCircle()
		case "2":
			typeTestMessage()
		case "q", "Q":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}

func moveMouseInCircle() {
	f, err := openDevice(mouseDevice)
	if err != nil {
		log.Printf("Error opening mouse device: %v", err)
		return
	}
	defer f.Close()

	fmt.Println("Moving mouse in a circle...")

	report := []byte{0, 40, 0, 0} // Move right
	f.Write(report)
	time.Sleep(delay)

	report = []byte{0, 0, 40, 0} // Move down
	f.Write(report)
	time.Sleep(delay)

	report = []byte{0, 216, 0, 0} // Move left  (-40 as byte = 256-40 = 216)
	f.Write(report)
	time.Sleep(delay)

	// Move up
	report = []byte{0, 0, 216, 0} // Move up    (-40 as byte = 256-40 = 216)
	f.Write(report)
}

func typeTestMessage() {
	f, err := openDevice(keyboardDevice)
	if err != nil {
		log.Printf("Error opening keyboard device: %v", err)
		return
	}
	defer f.Close()

	message := "this is a test"

	fmt.Printf("Typing: %s\n", message)

	// Direct ASCII to HID code mapping
	asciiToHID := map[rune]byte{
		't': 0x17,
		'h': 0x0b,
		'i': 0x0c,
		's': 0x16,
		' ': 0x2c,
		'a': 0x04,
		'e': 0x08,
	}

	for _, char := range message {
		if hidCode, exists := asciiToHID[char]; exists {
			// Send keypress
			report := make([]byte, 8)
			report[2] = hidCode
			f.Write(report)

			// Release key
			f.Write(make([]byte, 8))

			time.Sleep(delay)
		}
	}
}
