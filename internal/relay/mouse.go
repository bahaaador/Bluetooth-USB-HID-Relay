package relay

import (
	"github.com/bahaaador/Bluetooth-USB-HID-Relay/internal/logger"
)

type MouseRelay struct {
	lastState byte
}

func (m *MouseRelay) convertEvent(event InputEvent) ([]byte, error) {
	var report [4]byte

	// Preserve the last button state
	report[0] = m.lastState

	logger.Printf("Mouse event: type=%d, code=%d, value=%d, time=%v, lastMouseState=%d", event.Type, event.Code, event.Value, event.Time, m.lastState)

	switch event.Type {
	case 1: // EV_KEY
		if event.Code >= 0x110 && event.Code <= 0x112 { // BTN_MOUSE, BTN_RIGHT, BTN_MIDDLE
			if event.Value == 1 { // Button press
				m.lastState |= 1 << (event.Code - 0x110)
			} else if event.Value == 0 { // Button release
				m.lastState &^= 1 << (event.Code - 0x110)
			}
			report[0] = m.lastState
		}
	case 2: // EV_REL
		switch event.Code {
		case 0: // REL_X
			report[1] = byte(event.Value)
		case 1: // REL_Y
			report[2] = byte(event.Value)
		case 8: // REL_WHEEL
			report[3] = byte(event.Value)
		}
	default:
		return nil, nil // Ignore other event types
	}

	// Only send report if there's actual movement or button state change
	if report[0] != 0 || report[1] != 0 || report[2] != 0 || report[3] != 0 {
		return report[:], nil
	}
	return nil, nil

}

func (m *MouseRelay) validateEvent(event InputEvent) bool {
	switch event.Type {
	case 0: // EV_SYN - synchronization events
		return false // Skip sync events silently
	case 1: // EV_KEY - button events
		return event.Code >= 0x110 && event.Code <= 0x112 // Only accept mouse buttons
	case 2: // EV_REL - movement events
		return event.Code <= 8 // X, Y, and wheel movements
	case 4: // EV_MSC - miscellaneous events
		return false // Skip misc events silently
	default:
		logger.DebugPrintf("Unknown event type: %d", event.Type)
		return false
	}
}

func (m *MouseRelay) name() string {
	return "mouse"
}
