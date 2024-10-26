package main

import (
	"fmt"
	"os"
)

// DeviceType represents the type of HID device
type DeviceType int

const (
	Mouse DeviceType = iota
	Keyboard
)

// Device represents a HID device interface
type Device interface {
	Open() error
	Close() error
	Write([]byte) error
	SendRelease() error
}

// DeviceConfig holds device configuration
type DeviceConfig struct {
	InputPath  string
	OutputPath string
	Type       DeviceType
}

func validateDevice(input, output string) error {
	// Check if input device exists and is readable
	if _, err := os.Stat(input); err != nil {
		return fmt.Errorf("input device %s: %v", input, err)
	}

	// Check if output device exists and is writable
	if _, err := os.Stat(output); err != nil {
		return fmt.Errorf("output device %s: %v", output, err)
	}

	return nil
}
