package relay

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bahaaador/Bluetooth-USB-HID-Relay/internal/logger"
)

type InputEvent struct {
	Time  [2]uint64
	Type  uint16
	Code  uint16
	Value int32
}

type EventConverter interface {
	name() string
	validateEvent(event InputEvent) bool
	convertEvent(event InputEvent) ([]byte, error)
}

func streamDeviceEvents(ctx context.Context, inputPath, outputPath string, eventConverter EventConverter) error {
	var inputFile *os.File
	var outputFile *os.File
	var err error
	var try int = 1

	for {
		// Attempt to open input and output files
		inputFile, err = os.Open(inputPath)
		if err != nil {
			logger.Printf("Failed to open input device %s: %v. Retrying in %d seconds...", inputPath, err, try)
			try++
			time.Sleep(time.Duration(try) * time.Second)

			if try > 5 {
				return fmt.Errorf("failed to open input device %s: %v", inputPath, err)
			}
			continue
		}

		defer inputFile.Close()

		logger.DebugPrintf("Opening output device %s", outputPath)
		outputFile, err = os.OpenFile(outputPath, os.O_WRONLY, 0666)
		if err != nil {
			logger.Printf("Failed to open output device %s: %v. Retrying...", outputPath, err)
			time.Sleep(1 * time.Second)
			continue
		}
		defer outputFile.Close()

		event := InputEvent{}
		deviceName := filepath.Base(inputPath)

		try = 1

		for {
			select {
			case <-ctx.Done():
				logger.Printf("Relay shutdown for %s", inputPath)
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
					logger.Printf("Error reading from %s: %v. Reconnecting...", deviceName, err)
					goto reconnect
				}

				if !eventConverter.validateEvent(event) {
					logger.DebugPrintf("[DEBUG] Invalid event type for %s: %v", eventConverter.name(), event)
					continue
				}

				logger.Printf("Read event from %s: Type=%d, Code=%d, Value=%d\n", deviceName, event.Type, event.Code, event.Value)

				report, err := eventConverter.convertEvent(event)
				if err != nil {
					logger.DebugPrintf("[DEBUG] Error converting event: %v", err)
					continue
				}

				if report != nil {
					_, err = outputFile.Write(report)
					if err != nil {
						logger.Printf("Error writing to %s: %v. Reconnecting...", deviceName, err)
						goto reconnect
					}

					logger.DebugPrintf("[DEBUG] %s event relayed", deviceName)
				}
			}
		}
	reconnect:
		// Close files before attempting to reconnect
		outputFile.Close()
		inputFile.Close()
	}
}
