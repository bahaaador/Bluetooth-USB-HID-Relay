package relay

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bahaaador/Bluetooth-USB-HID-Relay/internal/device"
	"github.com/bahaaador/Bluetooth-USB-HID-Relay/internal/logger"
)

type Config struct {
	MouseInput     string
	KeyboardInput  string
	MouseOutput    string
	KeyboardOutput string
}

type Relay struct {
	config  Config
	ctx     context.Context
	cancel  context.CancelFunc
	errChan chan error
	sigChan chan os.Signal
}

func NewRelay(config Config) *Relay {
	ctx, cancel := context.WithCancel(context.Background())
	return &Relay{
		config:  config,
		ctx:     ctx,
		cancel:  cancel,
		errChan: make(chan error, 2),
		sigChan: make(chan os.Signal, 1),
	}
}

func (r *Relay) Start() error {
	logger.Println("Bluetooth HID Relay starting...")

	// Setup signal handling
	signal.Notify(r.sigChan, syscall.SIGINT, syscall.SIGTERM)

	go r.handleSignals()

	// Start device relaying
	go r.handleMouseEvents()
	go r.handleKeyboardEvents()

	// Wait for completion or error
	return r.wait()
}

func (r *Relay) handleSignals() {
	sig := <-r.sigChan
	logger.Printf("Received signal: %v, initiating shutdown...", sig)
	r.Shutdown()
}

func (r *Relay) wait() error {
	select {
	case err := <-r.errChan:
		r.Shutdown()
		return fmt.Errorf("relay error: %v", err)
	case <-r.ctx.Done():
		if r.ctx.Err() == context.Canceled {
			return nil
		}
		return fmt.Errorf("relay error: %v", r.ctx.Err())
	}
}

func (r *Relay) handleMouseEvents() {
	for {
		mouse, err := device.FindInputDevice("mouse")
		if err != nil {
			logger.Printf("Mouse not found: %v, retrying in 1s...", err)
			time.Sleep(time.Second)
			continue
		}

		if err := streamDeviceEvents(r.ctx, mouse, r.config.MouseOutput, &MouseRelay{}); err != nil {
			r.errChan <- fmt.Errorf("mouse relay: %v", err)
		}
	}
}

func (r *Relay) handleKeyboardEvents() {
	for {
		// Find keyboard keyboard
		keyboard, err := device.FindInputDevice("keyboard")
		if err != nil {
			logger.Printf("Keyboard not found: %v, retrying in 1s...", err)
			time.Sleep(time.Second)
			continue
		}

		logger.Printf("Found keyboard at: %s", keyboard)

		// Start relaying until error occurs
		err = streamDeviceEvents(r.ctx, keyboard, r.config.KeyboardOutput, &KeyboardRelay{})
		if err != nil {
			logger.Printf("Keyboard relay error: %v, reconnecting...", err)
		}
	}
}

// Shutdown gracefully stops the relay service
func (r *Relay) Shutdown() {
	logger.Println("Shutting down...")
	r.sendReleaseEvents()
	time.Sleep(100 * time.Millisecond)
	r.cancel()
}

func (r *Relay) sendReleaseEvents() {
	logger.Println("Sending release events...")

	// For keyboard: clear all modifiers and keys
	keyboardRelease := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	if f, err := os.OpenFile(r.config.KeyboardOutput, os.O_WRONLY, 0666); err == nil {
		for i := 0; i < 3; i++ { // Send multiple times to ensure it's received
			f.Write(keyboardRelease)
			time.Sleep(10 * time.Millisecond)
		}
		f.Close()
	}

	// For mouse: clear all buttons and movement
	mouseRelease := []byte{0, 0, 0, 0}
	if f, err := os.OpenFile(r.config.MouseOutput, os.O_WRONLY, 0666); err == nil {
		for i := 0; i < 3; i++ {
			f.Write(mouseRelease)
			time.Sleep(10 * time.Millisecond)
		}
		f.Close()
	}
}
