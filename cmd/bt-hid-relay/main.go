package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var debug bool

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
	log.Println("Bluetooth HID Relay starting...")

	// if err := r.initializeDevices(); err != nil {
	// 	return fmt.Errorf("device initialization failed: %v", err)
	// }

	// Setup signal handling
	signal.Notify(r.sigChan, syscall.SIGINT, syscall.SIGTERM)
	go r.handleSignals()

	// Start device relaying
	go r.relayMouse()
	go r.relayKeyboard()

	// Wait for completion or error
	return r.wait()
}

func (r *Relay) handleSignals() {
	sig := <-r.sigChan
	log.Printf("Received signal: %v, initiating shutdown...", sig)
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

func (r *Relay) relayMouse() {
	if err := relayInput(r.ctx, r.config.MouseInput, r.config.MouseOutput, &MouseRelay{}); err != nil {
		r.errChan <- fmt.Errorf("mouse relay: %v", err)
	}
}

func (r *Relay) relayKeyboard() {
	if err := relayInput(r.ctx, r.config.KeyboardInput, r.config.KeyboardOutput, &KeyboardRelay{}); err != nil {
		r.errChan <- fmt.Errorf("keyboard relay: %v", err)
	}
}

// Shutdown gracefully stops the relay service
func (r *Relay) Shutdown() {
	log.Println("Shutting down...")
	r.sendReleaseEvents()
	time.Sleep(100 * time.Millisecond)
	r.cancel()
}

func parseFlags() Config {
	var config Config

	flag.BoolVar(&debug, "debug", false, "enable debug mode")
	flag.StringVar(&config.MouseInput, "mouse-input", "/dev/input/event1", "mouse input device")
	flag.StringVar(&config.KeyboardInput, "keyboard-input", "/dev/input/event0", "keyboard input device")
	flag.StringVar(&config.MouseOutput, "mouse-output", "/dev/hidg0", "mouse output device")
	flag.StringVar(&config.KeyboardOutput, "keyboard-output", "/dev/hidg1", "keyboard output device")

	if !flag.Parsed() {
		flag.Parse()
	}

	return config
}

func main() {
	config := parseFlags()

	relay := NewRelay(config)
	if err := relay.Start(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	log.Println("Relay stopped successfully")
}

func (r *Relay) sendReleaseEvents() {
	if debug {
		log.Println("Sending release events...")
	}
	
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
