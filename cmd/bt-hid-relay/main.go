package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type InputEvent struct {
	Time  [2]uint64
	Type  uint16
	Code  uint16
	Value int32
}

var (
	debug          bool
	mouseInput     string
	keyboardInput  string
	mouseOutput    string
	keyboardOutput string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "enable debug mode")
	flag.StringVar(&mouseInput, "mouse-input", "/dev/input/event0", "mouse input device")
	flag.StringVar(&keyboardInput, "keyboard-input", "/dev/input/event1", "keyboard input device")
	flag.StringVar(&mouseOutput, "mouse-output", "/dev/hidg0", "mouse output device")
	flag.StringVar(&keyboardOutput, "keyboard-output", "/dev/hidg1", "keyboard output device")
	flag.Parse()
}

func main() {
	log.Println("Bluetooth HID Relay starting...")

	if err := relayDeviceInputs(); err != nil {
		log.Fatalf("Error in relay: %v", err)
	}
}

func relayDeviceInputs() error {
	mouseHID, err := os.OpenFile(mouseOutput, os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open mouse HID gadget: %v", err)
	}
	defer mouseHID.Close()

	keyboardHID, err := os.OpenFile(keyboardOutput, os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open keyboard HID gadget: %v", err)
	}
	defer keyboardHID.Close()

	mouseFile, err := os.Open(mouseInput)
	if err != nil {
		return fmt.Errorf("failed to open mouse device: %v", err)
	}
	defer mouseFile.Close()

	keyboardFile, err := os.Open(keyboardInput)
	if err != nil {
		return fmt.Errorf("failed to open keyboard device: %v", err)
	}
	defer keyboardFile.Close()

	errChan := make(chan error, 2)
	go relayInput(mouseFile, mouseHID, "Mouse", errChan)
	go relayInput(keyboardFile, keyboardHID, "Keyboard", errChan)

	log.Println("Relaying device inputs (press Ctrl+C to exit)")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return fmt.Errorf("relay error: %v", err)
	case <-sigChan:
		log.Println("Received interrupt, shutting down...")
		return nil
	}
}

func relayInput(inputFile, outputFile *os.File, deviceName string, errChan chan<- error) {
	event := InputEvent{}
	for {
		err := binary.Read(inputFile, binary.LittleEndian, &event)
		if err != nil {
			errChan <- fmt.Errorf("error reading from %s: %v", deviceName, err)
			return
		}

		hidReport, err := convertToHIDReport(event, deviceName)
		if err != nil {
			if debug {
				log.Printf("[DEBUG] Error converting %s event: %v", deviceName, err)
			}
			continue
		}

		_, err = outputFile.Write(hidReport)
		if err != nil {
			errChan <- fmt.Errorf("error writing to %s HID: %v", deviceName, err)
			return
		}

		if debug {
			log.Printf("[DEBUG] %s event relayed - Type: %d, Code: %d, Value: %d", deviceName, event.Type, event.Code, event.Value)
		}
	}
}

func convertToHIDReport(event InputEvent, deviceName string) ([]byte, error) {
	switch deviceName {
	case "Mouse":
		return convertMouseEvent(event)
	case "Keyboard":
		return convertKeyboardEvent(event)
	default:
		return nil, fmt.Errorf("unknown device: %s", deviceName)
	}
}

func convertMouseEvent(event InputEvent) ([]byte, error) {
	buttons := byte(0)
	var x, y int8
	switch event.Type {
	case 1: // EV_KEY
		if event.Code <= 0x110 { // BTN_MOUSE
			if event.Value == 1 {
				buttons |= 1 << (event.Code - 0x110)
			}
		}
	case 2: // EV_REL
		switch event.Code {
		case 0: // REL_X
			x = int8(event.Value)
		case 1: // REL_Y
			y = int8(event.Value)
		}
	}
	return []byte{buttons, byte(x), byte(y), 0}, nil
}

func convertKeyboardEvent(event InputEvent) ([]byte, error) {
	if event.Type == 1 && event.Code <= 0x77 { // EV_KEY and valid key code
		return []byte{0, 0, byte(event.Code), 0, 0, 0, 0, 0}, nil
	}
	return nil, fmt.Errorf("unsupported keyboard event")
}
