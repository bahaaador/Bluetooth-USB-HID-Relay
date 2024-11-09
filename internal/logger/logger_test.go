package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestDebugPrintf(t *testing.T) {
	original := log.Writer()
	defer log.SetOutput(original)

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Test debug off
	Debug = false
	DebugPrintf("test %s", "message")
	if buf.Len() > 0 {
		t.Error("Expected no output when Debug is false, got:", buf.String())
	}

	// Test debug on
	Debug = true
	buf.Reset() // Clear the buffer
	DebugPrintf("test %s", "message")
	if !strings.Contains(buf.String(), "test message") {
		t.Error("Expected output to contain 'test message', got:", buf.String())
	}
}
