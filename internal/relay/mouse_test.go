package relay

import (
	"testing"
)

func TestMouseRelay_ConvertEvent(t *testing.T) {
	tests := []struct {
		name       string
		event      InputEvent
		wantReport []byte
		wantErr    bool
	}{
		{
			name: "left button press",
			event: InputEvent{
				Type:  1,   // EV_KEY
				Code:  272, // BTN_LEFT
				Value: 1,
			},
			wantReport: []byte{0x01, 0, 0, 0}, // First bit set for left button
			wantErr:    false,
		},
		{
			name: "mouse move right",
			event: InputEvent{
				Type:  2, // EV_REL
				Code:  0, // REL_X
				Value: 10,
			},
			wantReport: []byte{0, 10, 0, 0},
			wantErr:    false,
		},
		{
			name: "mouse move down",
			event: InputEvent{
				Type:  2, // EV_REL
				Code:  1, // REL_Y
				Value: 5,
			},
			wantReport: []byte{0, 0, 5, 0},
			wantErr:    false,
		},
		{
			name: "scroll wheel",
			event: InputEvent{
				Type:  2, // EV_REL
				Code:  8, // REL_WHEEL
				Value: 1,
			},
			wantReport: []byte{0, 0, 0, 1},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MouseRelay{}
			report, err := m.convertEvent(tt.event)

			if (err != nil) != tt.wantErr {
				t.Errorf("MouseRelay.convertEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(report) != len(tt.wantReport) {
				t.Errorf("MouseRelay.convertEvent() report length = %d, want %d", len(report), len(tt.wantReport))
				return
			}

			for i := range report {
				if report[i] != tt.wantReport[i] {
					t.Errorf("MouseRelay.convertEvent() report[%d] = %#x, want %#x", i, report[i], tt.wantReport[i])
				}
			}
		})
	}
}

func TestMouseRelay_ValidateEvent(t *testing.T) {
	tests := []struct {
		name  string
		event InputEvent
		want  bool
	}{
		{
			name: "valid button event",
			event: InputEvent{
				Type:  1,
				Code:  272, // BTN_LEFT
				Value: 1,
			},
			want: true,
		},
		{
			name: "valid movement event",
			event: InputEvent{
				Type:  2,
				Code:  0, // REL_X
				Value: 10,
			},
			want: true,
		},
		{
			name: "sync event",
			event: InputEvent{
				Type:  0,
				Code:  0,
				Value: 0,
			},
			want: false,
		},
		{
			name: "invalid event type",
			event: InputEvent{
				Type:  5,
				Code:  0,
				Value: 0,
			},
			want: false,
		},
	}

	m := &MouseRelay{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := m.validateEvent(test.event); got != test.want {
				t.Errorf("MouseRelay.validateEvent() = %v, want %v", got, test.want)
			}
		})
	}
}
