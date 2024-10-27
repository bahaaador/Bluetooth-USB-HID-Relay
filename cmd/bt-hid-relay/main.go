package main

import (
	"flag"
	"log"
	"os"

	"github.com/bahaaador/Bluetooth-USB-HID-Relay/internal/relay"
	"github.com/bahaaador/Bluetooth-USB-HID-Relay/internal/logger"
)


func parseFlags() relay.Config {
	var config relay.Config

	flag.BoolVar(&logger.Debug, "debug", false, "enable debug mode")

	if !flag.Parsed() {
		flag.Parse()
	}

	return config
}

func main() {
	config := parseFlags()

	relay := relay.NewRelay(config)
	if err := relay.Start(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	log.Println("Relay stopped successfully")
}
