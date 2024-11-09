package main

import (
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/bahaaador/bluetooth-usb-peripheral-relay/internal/device"
)

func TestVerifyDevices(t *testing.T) {
	// Disable logging for tests
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	// Save original functions
	originalStat := osStat
	originalFindInputDevice := device.FindInputDeviceFunc
	defer func() {
		osStat = originalStat
		device.FindInputDeviceFunc = originalFindInputDevice
	}()

	// Mock os.Stat
	osStat = func(name string) (os.FileInfo, error) {
		if name == "/dev/hidg0" || name == "/dev/hidg1" {
			return nil, nil
		}
		return nil, os.ErrNotExist
	}

	// Mock FindInputDevice
	device.FindInputDeviceFunc = func(deviceType string) (string, error) {
		switch deviceType {
		case "mouse":
			return "/dev/input/event4", nil
		case "keyboard":
			return "/dev/input/event5", nil
		default:
			return "", nil
		}
	}

	err := verifyDevices()
	if err != nil {
		t.Errorf("verifyDevices() returned unexpected error: %v", err)
	}
}

func TestEchoDeviceInputs(t *testing.T) {
	// Disable logging for tests
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	// Save original functions
	originalFindInputDevice := device.FindInputDeviceFunc
	originalOpenFile := osOpenFile
	defer func() {
		device.FindInputDeviceFunc = originalFindInputDevice
		osOpenFile = originalOpenFile
	}()

	// Mock FindInputDevice
	device.FindInputDeviceFunc = func(deviceType string) (string, error) {
		switch deviceType {
		case "mouse":
			return "/dev/input/event4", nil
		case "keyboard":
			return "/dev/input/event5", nil
		default:
			return "", nil
		}
	}

	// Mock file operations
	osOpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		f, _ := os.CreateTemp("", "mock")
		return f, nil
	}

	// Create a channel to signal test completion
	done := make(chan bool)

	go func() {
		err := echoDeviceInputs()
		if err != nil {
			t.Errorf("echoDeviceInputs() returned unexpected error: %v", err)
		}
		done <- true
	}()

	// Wait briefly then cancel
	select {
	case <-done:
		// Test completed normally
	case <-time.After(100 * time.Millisecond):
		// Test timeout is expected since this would normally run indefinitely
	}
}
