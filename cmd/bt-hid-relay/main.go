package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Bluetooth HID Relay")

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// TODO: Implement Bluetooth to USB HID relay logic
	return nil
}
