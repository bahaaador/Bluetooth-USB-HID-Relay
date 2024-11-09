package main

import (
	"bytes"
	"os"
	"testing"
	"time"
)

// MockFile implements FileWriter
type MockFile struct {
	buffer *bytes.Buffer
	closed bool
}

func NewMockFile() *MockFile {
	return &MockFile{
		buffer: new(bytes.Buffer),
		closed: false,
	}
}

func (m *MockFile) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}

func (m *MockFile) Close() error {
	m.closed = true
	return nil
}

func (m *MockFile) Bytes() []byte {
	return m.buffer.Bytes()
}

func TestMoveMouseInCircle(t *testing.T) {
	// Save original functions
	originalOpenFile := osOpenFile
	originalDelay := delay
	defer func() {
		osOpenFile = originalOpenFile
		delay = originalDelay
	}()

	// Set minimal delay for testing
	delay = 1 * time.Millisecond

	// Create mock file
	mock := NewMockFile()
	
	// Replace osOpenFile with our mock
	osOpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		// We need to type assert our mock to *os.File
		// This is a hack, but it works for testing
		return os.NewFile(0, "mock"), nil
	}

	// Replace openDevice to return our mock directly
	originalOpenDevice := openDevice
	openDevice = func(path string) (FileWriter, error) {
		return mock, nil
	}
	defer func() {
		openDevice = originalOpenDevice
	}()

	// Run the function
	moveMouseInCircle()

	// Verify writes
	got := mock.Bytes()
	expected := []byte{
		0, 40, 0, 0, // Right
		0, 0, 40, 0, // Down
		0, 216, 0, 0, // Left
		0, 0, 216, 0, // Up
	}

	if !bytes.Equal(got, expected) {
		t.Errorf("Expected %v, got %v", expected, got)
	}

	if !mock.closed {
		t.Error("File was not closed")
	}
}

func TestTypeTestMessage(t *testing.T) {
	// Save original functions
	originalOpenFile := osOpenFile
	originalDelay := delay
	defer func() {
		osOpenFile = originalOpenFile
		delay = originalDelay
	}()

	// Set minimal delay for testing
	delay = 1 * time.Millisecond

	// Create mock file
	mock := NewMockFile()
	
	// Replace openDevice to return our mock directly
	originalOpenDevice := openDevice
	openDevice = func(path string) (FileWriter, error) {
		return mock, nil
	}
	defer func() {
		openDevice = originalOpenDevice
	}()

	// Run the function
	typeTestMessage()

	// Verify the file was closed
	if !mock.closed {
		t.Error("File was not closed")
	}

	// Verify some basic expectations about the writes
	got := mock.Bytes()
	if len(got) == 0 {
		t.Error("No data was written")
	}

	// Each keystroke should be 8 bytes (press) + 8 bytes (release)
	if len(got)%(8*2) != 0 {
		t.Errorf("Unexpected data length: %d", len(got))
	}
}
