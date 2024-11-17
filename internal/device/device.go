package device

import (
	"fmt"
	"os"
	"strings"
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

var FindInputDeviceFunc = FindInputDevice
var readFile = os.ReadFile

func FindInputDevice(deviceType string) (string, error) {
	data, err := readFile("/proc/bus/input/devices")
	if err != nil {
		return "", fmt.Errorf("failed to read devices: %v", err)
	}

	var deviceName string

	// Look through each line
	for _, line := range strings.Split(string(data), "\n") {
		// Look for device name
		if strings.HasPrefix(line, "N: Name=") {
			if strings.Contains(strings.ToLower(line), deviceType) {
				deviceName = line
			}
			continue
		}

		// If we found a matching device name, look for its event
		if deviceName != "" && strings.HasPrefix(line, "H: Handlers=") {
			for _, word := range strings.Fields(line) {
				if strings.HasPrefix(word, "event") {
					return fmt.Sprintf("/dev/input/%s", word), nil
				}
			}
		}
	}

	return "", fmt.Errorf("%s not found", deviceType)
}
