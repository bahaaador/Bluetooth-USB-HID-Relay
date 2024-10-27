package relay

import (
	"github.com/bahaaador/Bluetooth-USB-HID-Relay/internal/logger"
)

type MouseRelay struct {
	lastState byte
}

func (m *MouseRelay) convertEvent(event InputEvent) ([]byte, error) {
	var report [4]byte

	logger.DebugPrintf("Mouse event: type=%d, code=%d, value=%d, time=%v", event.Type, event.Code, event.Value, event.Time)

	switch event.Type {
	case 1: // EV_KEY
		if event.Code >= 272 && event.Code <= 276 {
			buttonBit := event.Code - 272 // Convert to 0-based index

			if event.Value == 1 { // Button press
				m.lastState |= 1 << buttonBit
			} else if event.Value == 0 { // Button release
				m.lastState &^= 1 << buttonBit
			}
			report[0] = m.lastState
			return report[:], nil
		}
	case 2: // EV_REL
		report[0] = m.lastState
		switch event.Code {
		case 0: // X axis
			report[1] = byte(event.Value)
		case 1: // Y axis
			report[2] = byte(event.Value)
		case 8: // Wheel
			report[3] = byte(event.Value)
		}
		return report[:], nil
	}

	// For any other event type, still return the current state
	report[0] = m.lastState
	return report[:], nil
}

func (m *MouseRelay) validateEvent(event InputEvent) bool {
	switch event.Type {
	case 0: // EV_SYN
		return false
	case 1: // EV_KEY
		return event.Code >= 272 && event.Code <= 276
	case 2: // EV_REL
		return event.Code <= 8
	case 4: // EV_MSC
		return false // Explicitly ignore these events
	default:
		logger.DebugPrintf("Unknown event type: %d", event.Type)
		return false
	}
}

func (m *MouseRelay) name() string {
	return "mouse"
}
