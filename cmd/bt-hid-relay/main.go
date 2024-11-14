package main

import (
	"flag"
	"log"
	"os"

	"github.com/bahaaador/bluetooth-usb-peripheral-relay/internal/device"
	"github.com/bahaaador/bluetooth-usb-peripheral-relay/internal/logger"
	"github.com/bahaaador/bluetooth-usb-peripheral-relay/internal/relay"
)

func parseFlags() relay.Config {
	var config relay.Config

	flag.BoolVar(&logger.Debug, "debug", false, "enable debug mode")
	flag.StringVar(&config.MouseOutput, "mouse-output", "/dev/hidg0", "mouse output device")
	flag.StringVar(&config.KeyboardOutput, "keyboard-output", "/dev/hidg1", "keyboard output device")

	if !flag.Parsed() {
		flag.Parse()
	}

	return config
}

func main() {
	config := parseFlags()

	hasHostCapability, isHostEnabled, err := device.CheckUSBHostSupport()
	if err != nil {
		log.Fatal(err)
	}

	if !hasHostCapability {
		log.Fatal("USB Host mode is not supported")
	}

	if !isHostEnabled {
		log.Fatal("USB Host mode is not enabled")
	}

	relay := relay.NewRelay(config)
	if err := relay.Start(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	log.Println("Relay stopped successfully")
}
