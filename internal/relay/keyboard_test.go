package relay

import (
	"bytes"
	"testing"
)

func TestKeyboardRelay_ConvertEvent(t *testing.T) {
	tests := []struct {
		name       string
		event      InputEvent
		wantReport []byte
		wantErr    bool
	}{
		{
			name: "simple key press A",
			event: InputEvent{
				Type:  1,  // EV_KEY
				Code:  30, // KEY_A
				Value: 1,  // Press
			},
			wantReport: []byte{0, 0, 0x04, 0, 0, 0, 0, 0}, // 0x04 is HID code for 'a'
			wantErr:    false,
		},
		{
			name: "key release A",
			event: InputEvent{
				Type:  1,
				Code:  30,
				Value: 0, // Release
			},
			wantReport: []byte{0, 0, 0, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
		{
			name: "left shift press",
			event: InputEvent{
				Type:  1,
				Code:  42, // Left
				Value: 1,
			},
			wantReport: []byte{0x02, 0, 0, 0, 0, 0, 0, 0},
			wantErr:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			k := new(KeyboardRelay)
			gotReport, err := k.convertEvent(test.event)
			if (err != nil) != test.wantErr {
				t.Errorf("KeyboardRelay.ConvertEvent() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !bytes.Equal(gotReport, test.wantReport) {
				t.Errorf("KeyboardRelay.ConvertEvent() = %v, want %v", gotReport, test.wantReport)
			}
		})
	}
}
