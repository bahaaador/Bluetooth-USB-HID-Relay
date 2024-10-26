package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type InputEvent struct {
	Time  [2]uint64
	Type  uint16
	Code  uint16
	Value int32
}

type InputDevice interface {
	convertEvent(event InputEvent) ([]byte, error)
	validateEvent(event InputEvent) bool
	name() string
}

func relayInput(ctx context.Context, inputPath, outputPath string, deviceType InputDevice) error {
	var inputFile *os.File
	var outputFile *os.File
	var err error
	var try int = 1

	for {
		// Attempt to open input and output files
		inputFile, err = os.Open(inputPath)
		if err != nil {
			if debug {
				log.Printf("Failed to open input device %s: %v. Retrying in %d seconds...", inputPath, err, try)
			}
			try++
			time.Sleep(time.Duration(try) * time.Second)
			continue
		}
		try = 1

		defer inputFile.Close()

		outputFile, err = os.OpenFile(outputPath, os.O_WRONLY, 0666)
		if err != nil {
			if debug {
				log.Printf("Failed to open output device %s: %v. Retrying...", outputPath, err)
			}
			time.Sleep(1 * time.Second)
			continue
		}
		defer outputFile.Close()

		event := InputEvent{}
		deviceName := filepath.Base(inputPath)

		for {
			select {
			case <-ctx.Done():
				log.Printf("Relay shutdown for %s", inputPath)
				// Send multiple release events during shutdown
				for i := 0; i < 3; i++ {
					releaseReport := []byte{0, 0, 0, 0, 0, 0, 0, 0}
					outputFile.Write(releaseReport)
					time.Sleep(10 * time.Millisecond)
				}
				return nil
			default:
				err := binary.Read(inputFile, binary.LittleEndian, &event)
				if err != nil {
					if debug {
						log.Printf("Error reading from %s: %v. Reconnecting...", deviceName, err)
					}
					goto reconnect
				}

				if !deviceType.validateEvent(event) {
					if debug {
						log.Printf("[DEBUG] Invalid event type for %s: %v", deviceType.name(), event)
					}
					continue
				}

				if debug {
					fmt.Printf("Read event from %s: Type=%d, Code=%d, Value=%d\n", deviceName, event.Type, event.Code, event.Value)
				}

				report, err := deviceType.convertEvent(event)
				if err != nil {
					if debug {
						log.Printf("[DEBUG] Error converting event: %v", err)
					}
					continue
				}

				if report != nil {
					_, err = outputFile.Write(report)
					if err != nil {
						if debug {
							log.Printf("Error writing to %s: %v. Reconnecting...", deviceName, err)
						}
						goto reconnect
					}

					if debug {
						log.Printf("[DEBUG] %s event relayed", deviceName)
					}
				}
			}
		}
	reconnect:
		// Close files before attempting to reconnect
		outputFile.Close()
		inputFile.Close()
	}
}
