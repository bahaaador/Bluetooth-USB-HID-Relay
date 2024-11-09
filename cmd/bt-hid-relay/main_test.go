package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/bahaaador/bluetooth-usb-peripheral-relay/internal/device"
	"github.com/bahaaador/bluetooth-usb-peripheral-relay/internal/logger"
)

var (
	osExit     = os.Exit
	exitCalled = false
)

func init() {
	// Mock FindInputDeviceFunc to immediately return success
	device.FindInputDeviceFunc = func(deviceType string) (string, error) {
		switch deviceType {
		case "mouse":
			return "/dev/input/event4", nil
		case "keyboard":
			return "/dev/input/event5", nil
		default:
			return "", fmt.Errorf("device not found")
		}
	}
}

func TestParseFlags(t *testing.T) {
	// Save original flag values to restore after test
	origArgs := os.Args
	defer func() {
		os.Args = origArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	tests := []struct {
		name     string
		args     []string
		wantConf struct {
			debug          bool
			mouseOutput    string
			keyboardOutput string
		}
	}{
		{
			name: "default values",
			args: []string{"cmd"},
			wantConf: struct {
				debug          bool
				mouseOutput    string
				keyboardOutput string
			}{
				debug:          false,
				mouseOutput:    "/dev/hidg0",
				keyboardOutput: "/dev/hidg1",
			},
		},
		{
			name: "custom values",
			args: []string{
				"cmd",
				"-debug",
				"-mouse-output=/dev/custom0",
				"-keyboard-output=/dev/custom1",
			},
			wantConf: struct {
				debug          bool
				mouseOutput    string
				keyboardOutput string
			}{
				debug:          true,
				mouseOutput:    "/dev/custom0",
				keyboardOutput: "/dev/custom1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags before each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = tt.args

			got := parseFlags()

			if logger.Debug != tt.wantConf.debug {
				t.Errorf("parseFlags() debug = %v, want %v", logger.Debug, tt.wantConf.debug)
			}
			if got.MouseOutput != tt.wantConf.mouseOutput {
				t.Errorf("parseFlags() mouseOutput = %v, want %v", got.MouseOutput, tt.wantConf.mouseOutput)
			}
			if got.KeyboardOutput != tt.wantConf.keyboardOutput {
				t.Errorf("parseFlags() keyboardOutput = %v, want %v", got.KeyboardOutput, tt.wantConf.keyboardOutput)
			}
		})
	}
}

func TestMain(t *testing.T) {
	// Disable logging for tests
	log.SetOutput(io.Discard)

	// Save original values
	originalExit := osExit
	originalFindInputDevice := device.FindInputDeviceFunc

	// Restore after test
	defer func() {
		osExit = originalExit
		device.FindInputDeviceFunc = originalFindInputDevice
		log.SetOutput(os.Stderr)
	}()

	// Mock device detection to succeed immediately
	device.FindInputDeviceFunc = func(deviceType string) (string, error) {
		switch deviceType {
		case "mouse":
			return "/dev/input/event4", nil
		case "keyboard":
			return "/dev/input/event5", nil
		default:
			return "", fmt.Errorf("unknown device type")
		}
	}

	// Mock os.Exit
	osExit = func(code int) {
		exitCalled = true
		panic(fmt.Sprintf("os.Exit(%d)", code))
	}

	// Save original flag values
	origArgs := os.Args
	defer func() {
		os.Args = origArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	// Test with default arguments
	os.Args = []string{"cmd"}

	// Create a channel to capture potential panics
	done := make(chan bool)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Only report error if it's not our expected exit
				if r != "os.Exit(0)" {
					t.Errorf("main() panicked: %v", r)
				}
			}
			done <- true
		}()

		// Run main
		go main()
		time.Sleep(100 * time.Millisecond) // Let main initialize
		osExit(0)
	}()

	// Wait for either completion or timeout
	select {
	case <-done:
		if !exitCalled {
			t.Error("Expected os.Exit to be called")
		}
	case <-time.After(2 * time.Second):
		t.Error("main() timed out")
	}
}
