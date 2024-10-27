package utils

import (
	"fmt"
	"os"
	"strings"
)

func FindInputDevice(deviceType string) (string, error) {
	content, err := os.ReadFile("/proc/bus/input/devices")
	if err != nil {
		return "", fmt.Errorf("failed to read devices: %v", err)
	}

	var deviceName string

	// Look through each line
	for _, line := range strings.Split(string(content), "\n") {
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

	return "", fmt.Errorf("%s device not found", deviceType)
}
