package device

import (
	"os"
	"strings"
)

// CheckUSBHostSupportFunc is a type for the host support check function
type CheckUSBHostSupportFunc func() (bool, bool, error)

// CheckUSBHostSupport is the function variable that can be replaced in tests
var CheckUSBHostSupport CheckUSBHostSupportFunc = checkUSBHostSupport

// CheckUSBHostSupport verifies both USB host hardware capability and if it's enabled
func checkUSBHostSupport() (hasCapability bool, isEnabled bool, err error) {
	// Check hardware capability first
	hasCapability = false

	// Check for USB host controller presence
	if _, err := os.Stat("/sys/class/usb_host"); err == nil {
		hasCapability = true
	}

	// Check USB controllers capabilities through kernel info
	controllers, err := os.ReadFile("/sys/kernel/debug/usb/devices")
	if err == nil {
		content := string(controllers)
		// Look for host controller interfaces (EHCI, XHCI, OHCI)
		if strings.Contains(content, "Cls=09") || // USB Hub Class
			strings.Contains(content, "EHCI") ||
			strings.Contains(content, "XHCI") ||
			strings.Contains(content, "OHCI") {
			hasCapability = true
		}
	}

	// If no hardware capability, return early
	if !hasCapability {
		return hasCapability, false, nil
	}

	// Check if it's enabled/configured
	isEnabled = false

	// Check if the device tree has USB OTG support
	dtOverlay, err := os.ReadFile("/boot/config.txt")
	if err == nil && strings.Contains(string(dtOverlay), "dtoverlay=dwc2") {
		isEnabled = true
	}

	// Check if the necessary modules are loaded
	modules, err := os.ReadFile("/proc/modules")
	if err == nil {
		moduleContent := string(modules)
		if strings.Contains(moduleContent, "dwc2") &&
			strings.Contains(moduleContent, "libcomposite") {
			isEnabled = true
		}
	}

	// Check for USB gadget configfs support
	if _, err := os.Stat("/sys/kernel/config/usb_gadget"); err == nil {
		isEnabled = true
	}

	return hasCapability, isEnabled, nil
}
