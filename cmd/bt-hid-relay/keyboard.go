package main

import (
	"log"
)

type KeyboardRelay struct {
	modifiers byte // Track active modifiers
	lastEvent InputEvent
	lastKeyCode byte
}

// Linux event code to HID usage ID mapping
var keyCodeMap = generateKeyCodeMap()

const (
	KEY_ESC = 1
	KEY_1   = 2
	KEY_2   = 3
	// ... etc
	KEY_A = 30
	KEY_B = 48
	KEY_C = 46
	// ... etc
	KEY_LEFTCTRL   = 29
	KEY_LEFTSHIFT  = 42
	KEY_RIGHTSHIFT = 54
	KEY_LEFTALT    = 56
	KEY_LEFTMETA   = 125
	KEY_RIGHTMETA  = 126
	KEY_F1         = 59
	KEY_BACKSPACE  = 14
	KEY_TAB        = 15
	KEY_ENTER      = 28
	KEY_CAPSLOCK   = 57
	KEY_SPACE      = 58
	KEY_RIGHTCTRL  = 105
	KEY_RIGHTALT   = 108
	KEY_HOME       = 302
	KEY_UP         = 103
	KEY_PAGEUP     = 102
	KEY_LEFT       = 100
	KEY_RIGHT      = 104
	KEY_END        = 301
	KEY_DOWN       = 108
	KEY_PAGEDOWN   = 109
	KEY_INSERT     = 210
	KEY_DELETE     = 111
)

func generateKeyCodeMap() map[uint16]byte {
	m := make(map[uint16]byte)

	// QWERTY Layout mapping
	qwertyMap := map[uint16]byte{
		// Top row
		16: 0x14, // Q
		17: 0x1A, // W
		18: 0x08, // E
		19: 0x15, // R
		20: 0x17, // T
		21: 0x1C, // Y
		22: 0x18, // U
		23: 0x0C, // I
		24: 0x12, // O
		25: 0x13, // P

		// Home row
		30: 0x04, // A
		31: 0x16, // S
		32: 0x07, // D
		33: 0x09, // F
		34: 0x0A, // G
		35: 0x0B, // H
		36: 0x0D, // J
		37: 0x0E, // K
		38: 0x0F, // L

		// Bottom row
		44: 0x1D, // Z
		45: 0x1B, // X
		46: 0x06, // C
		47: 0x19, // V
		48: 0x05, // B
		49: 0x11, // N
		50: 0x10, // M
	}

	// Numbers (1-9,0)
	for i := 1; i <= 10; i++ {
		linuxCode := uint16(i + 1)      // KEY_1 is 2, KEY_2 is 3, etc
		hidCode := byte(0x1E + i - 1)   // 0x1E is USB HID code for '1'
		m[linuxCode] = hidCode
	}

	// Add QWERTY layout
	for code, hid := range qwertyMap {
		m[code] = hid
	}

	// Special keys
	m[1] = 0x29    // ESC
	m[14] = 0x2A   // Backspace
	m[15] = 0x2B   // Tab
	m[28] = 0x28   // Enter
	m[29] = 0xE0   // Left Ctrl
	m[42] = 0xE1   // Left Shift
	m[54] = 0xE5   // Right Shift
	m[56] = 0xE2   // Left Alt
	m[57] = 0x2C   // Space
	m[58] = 0x39   // Caps Lock
	m[97] = 0xE4   // Right Ctrl
	m[100] = 0xE6  // Right Alt
	m[125] = 0xE3  // Left Meta (Windows/Command)
	m[126] = 0xE7  // Right Meta (Windows/Command)

	// Additional special keys
	m[41] = 0x35   // ` (backtick/tilde)
	m[43] = 0x31   // \ (backslash)
	m[26] = 0x2F   // [ (left bracket)
	m[27] = 0x30   // ] (right bracket)
	m[39] = 0x33   // ; (semicolon)
	m[40] = 0x34   // ' (single quote)
	m[51] = 0x36   // , (comma)
	m[52] = 0x37   // . (period)
	m[53] = 0x38   // / (forward slash)
	m[12] = 0x2D   // - (minus)
	m[13] = 0x2E   // = (equals)

	// Function keys
	for i := 0; i < 12; i++ {
		m[uint16(59+i)] = byte(0x3A + i) // F1-F12
	}

	// Navigation cluster
	m[102] = 0x4A  // Home
	m[107] = 0x4B  // End
	m[104] = 0x52  // Page Up
	m[109] = 0x51  // Page Down
	m[110] = 0x49  // Insert
	m[111] = 0x4C  // Delete
	m[103] = 0x52  // Up Arrow
	m[108] = 0x51  // Down Arrow
	m[105] = 0x50  // Left Arrow
	m[106] = 0x4F  // Right Arrow

	return m
}

func (k *KeyboardRelay) convertEvent(event InputEvent) ([]byte, error) {
	report := make([]byte, 8)
	
	// Handle modifier keys
	if isModifier(event.Code) {
		k.updateModifiers(event)
		report[0] = k.modifiers
		return report, nil
	}

	// Regular keys
	hidKeyCode, exists := keyCodeMap[event.Code]
	if !exists {
		if debug {
			log.Printf("No mapping for key code: %d", event.Code)
		}
		return nil, nil
	}

	switch event.Value {
	case 0: // Release
		k.lastKeyCode = 0
		// Clear everything except modifiers
		report[0] = k.modifiers
		return report, nil
	case 1, 2: // Press or Repeat
		k.lastKeyCode = hidKeyCode
		report[0] = k.modifiers
		report[2] = hidKeyCode
		return report, nil
	}

	return nil, nil
}

// Helper functions
func isModifier(code uint16) bool {
	return code == 29 || // Left Ctrl
		   code == 97 || // Right Ctrl
		   code == 42 || // Left Shift
		   code == 54 || // Right Shift
		   code == 56 || // Left Alt
		   code == 100 || // Right Alt
		   code == 125 || // Left Meta
		   code == 126    // Right Meta
}

func (k *KeyboardRelay) updateModifiers(event InputEvent) {
	var mask byte
	switch event.Code {
	case 29:  mask = 0x01 // Left Ctrl
	case 97:  mask = 0x10 // Right Ctrl
	case 42:  mask = 0x02 // Left Shift
	case 54:  mask = 0x20 // Right Shift
	case 56:  mask = 0x04 // Left Alt
	case 100: mask = 0x40 // Right Alt
	case 125: mask = 0x08 // Left Meta
	case 126: mask = 0x80 // Right Meta
	}

	if event.Value > 0 {
		k.modifiers |= mask
	} else {
		k.modifiers &^= mask
	}
}

func (k *KeyboardRelay) validateEvent(event InputEvent) bool {
	switch event.Type {
	case 0: // EV_SYN
		if debug {
			log.Printf("Sync event received - marks end of event batch")
		}
		return false  // Don't need to send to HID device
	case 1: // EV_KEY
		_, exists := keyCodeMap[event.Code]
		return exists
	case 4: // EV_MSC
		if debug {
			log.Printf("Misc event received - scancode: %d", event.Code)
		}
		return false  // These are metadata events, not actual key presses
	default:
		if debug {
			log.Printf("Unexpected event type: %d", event.Type)
		}
		return false
	}
}

func (k *KeyboardRelay) name() string {
	return "keyboard"
}
