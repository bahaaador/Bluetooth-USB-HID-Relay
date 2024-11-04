package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/bahaaador/bluetooth-usb-peripheral-relay/internal/relay"
)

func main() {
	fmt.Println("Bluetooth Device Verification Tool")

	if err := verifyDevices(); err != nil {
		log.Fatal(err)
	}

	if err := echoDeviceInputs(); err != nil {
		log.Fatal(err)
	}
}

func verifyDevices() error {
	// Check if the HID gadget devices exist
	if _, err := os.Stat("/dev/hidg0"); os.IsNotExist(err) {
		return fmt.Errorf("HID gadget device /dev/hidg0 not found")
	}
	if _, err := os.Stat("/dev/hidg1"); os.IsNotExist(err) {
		return fmt.Errorf("HID gadget device /dev/hidg1 not found")
	}

	fmt.Println("HID gadget devices are present")
	return nil
}

func echoDeviceInputs() error {
	// Open mouse and keyboard devices for reading
	mouseFile, err := os.Open("/dev/input/event0")
	if err != nil {
		return fmt.Errorf("failed to open mouse device: %v", err)
	}
	defer mouseFile.Close()

	keyboardFile, err := os.Open("/dev/input/event1")
	if err != nil {
		return fmt.Errorf("failed to open keyboard device: %v", err)
	}
	defer keyboardFile.Close()

	// Start goroutines to read from mouse and keyboard
	go readInput(mouseFile, "Mouse")
	go readInput(keyboardFile, "Keyboard")

	fmt.Println("Listening for device inputs (press Ctrl+C to exit):")

	// Keep the program running
	select {}
}

func readInput(file *os.File, deviceName string) {
	event := relay.InputEvent{}
	for {
		err := binary.Read(file, binary.LittleEndian, &event)
		if err != nil {
			fmt.Printf("Error reading from %s: %v\n", deviceName, err)
			return
		}

		fmt.Printf("%s event - Type: %d, Code: %d, Value: %d\n",
			deviceName, event.Type, event.Code, event.Value)
	}
}
