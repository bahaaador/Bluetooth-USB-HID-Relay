package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bahaaador/bluetooth-usb-peripheral-relay/internal/device"
	"github.com/bahaaador/bluetooth-usb-peripheral-relay/internal/relay"
)

var (
	osOpenFile = os.OpenFile
	osStat     = os.Stat
)

const (
	checkMark = "\u2713" // ✓
	crossMark = "\u2717" // ✗
)

func main() {
	fmt.Println("Bluetooth Device Verification Tool")
	fmt.Println("=================================")

	if err := verifyUSBHostSupport(); err != nil {
		log.Fatal(err)
	}

	if err := verifyDevices(); err != nil {
		log.Fatal(err)
	}

	if err := echoDeviceInputs(); err != nil {
		log.Fatal(err)
	}
}

func verifyUSBHostSupport() error {
	fmt.Println("\nChecking USB Host Support:")
	hasHostCapability, isHostEnabled, err := device.CheckUSBHostSupport()
	if err != nil {
		return fmt.Errorf("failed to check USB host support: %v", err)
	}

	if !hasHostCapability {
		return fmt.Errorf("USB Host mode is not supported")
	} else {
		fmt.Printf("%s USB Host mode: supported\n", checkMark)
	}

	if !isHostEnabled {
		return fmt.Errorf("USB Host mode is not enabled")
	} else {
		fmt.Printf("%s USB Host mode: enabled\n", checkMark)
	}

	return nil
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
	if _, err := osStat(path); os.IsNotExist(err) {
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
		mouseFile, err := osOpenFile(mouseDevice, os.O_RDONLY, 0666)
		if err != nil {
			return fmt.Errorf("failed to open mouse device: %v", err)
		}
		defer mouseFile.Close()
		go readInput(mouseFile, "Mouse")
	}

	if keyboardDevice != "" {
		keyboardFile, err := osOpenFile(keyboardDevice, os.O_RDONLY, 0666)
		if err != nil {
			return fmt.Errorf("failed to open keyboard device: %v", err)
		}
		defer keyboardFile.Close()
		go readInput(keyboardFile, "Keyboard")
	}

	// Add a channel to handle program termination
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	
	// Wait for interrupt signal
	<-done
	return nil
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
