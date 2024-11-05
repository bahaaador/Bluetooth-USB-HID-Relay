package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/bahaaador/bluetooth-usb-peripheral-relay/internal/device"
	"github.com/bahaaador/bluetooth-usb-peripheral-relay/internal/relay"
)

const (
	checkMark = "\u2713" // ✓
	crossMark = "\u2717" // ✗
)

func main() {
	fmt.Println("Bluetooth Device Verification Tool")
	fmt.Println("=================================")

	if err := verifyDevices(); err != nil {
		log.Fatal(err)
	}

	if err := echoDeviceInputs(); err != nil {
		log.Fatal(err)
	}
}

func verifyDevices() error {
	// Check HID gadget devices
	fmt.Println("\nChecking HID gadget devices:")
	checkDevice("/dev/hidg0", "Mouse HID gadget")
	checkDevice("/dev/hidg1", "Keyboard HID gadget")

	// Check input devices
	fmt.Println("\nChecking input devices:")

	mouse, err := device.FindInputDevice("mouse")
	if err != nil {
		fmt.Printf("%s Mouse input device: not found\n", crossMark)
	} else {
		fmt.Printf("%s Mouse input device: %s\n", checkMark, mouse)
	}

	keyboard, err := device.FindInputDevice("keyboard")
	if err != nil {
		fmt.Printf("%s Keyboard input device: not found\n", crossMark)
	} else {
		fmt.Printf("%s Keyboard input device: %s\n", checkMark, keyboard)
	}

	return nil
}

func checkDevice(path, description string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("%s %s: not found (%s)\n", crossMark, description, path)
	} else {
		fmt.Printf("%s %s: present (%s)\n", checkMark, description, path)
	}
}

func echoDeviceInputs() error {
	fmt.Println("\nAttempting to read device inputs:")

	// Find mouse device
	mouseDevice, err := device.FindInputDevice("mouse")
	if err != nil {
		fmt.Printf("%s Mouse input not available: %v\n", crossMark, err)
		mouseDevice = ""
	}

	// Find keyboard device
	keyboardDevice, err := device.FindInputDevice("keyboard")
	if err != nil {
		fmt.Printf("%s Keyboard input not available: %v\n", crossMark, err)
		keyboardDevice = ""
	}

	if mouseDevice == "" && keyboardDevice == "" {
		fmt.Println("\nNo input devices found, skipping input reading...")
		return nil
	}

	fmt.Println("\nListening for device inputs (press Ctrl+C to exit):")
	fmt.Println("================================================")

	// Only start readers for devices that were found
	if mouseDevice != "" {
		mouseFile, err := os.Open(mouseDevice)
		if err != nil {
			return fmt.Errorf("failed to open mouse device: %v", err)
		}
		defer mouseFile.Close()
		go readInput(mouseFile, "Mouse")
	}

	if keyboardDevice != "" {
		keyboardFile, err := os.Open(keyboardDevice)
		if err != nil {
			return fmt.Errorf("failed to open keyboard device: %v", err)
		}
		defer keyboardFile.Close()
		go readInput(keyboardFile, "Keyboard")
	}

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
