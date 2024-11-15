package relay

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bahaaador/bluetooth-usb-peripheral-relay/internal/logger"
)

type EventConverter interface {
	name() string
	validateEvent(event InputEvent) bool
	convertEvent(event InputEvent) ([]byte, error)
}

func streamDeviceEvents(ctx context.Context, inputPath, outputPath string, eventConverter EventConverter) error {
	logger.DebugPrintf("InputEvent struct size: %d bytes", binary.Size(InputEvent{}))
	deviceName := filepath.Base(inputPath)

	for {
		inputFile, outputFile, err := openDeviceFiles(inputPath, outputPath)
		if err != nil {
			return err
		}
		defer inputFile.Close()
		defer outputFile.Close()

		if err := processEvents(ctx, inputFile, outputFile, eventConverter, deviceName); err != nil {
			logger.Printf("Error processing events for %s: %v. Reconnecting...", deviceName, err)
			continue
		}
		
		return nil // Clean shutdown via context cancellation
	}
}

func openDeviceFiles(inputPath, outputPath string) (*os.File, *os.File, error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open input device %s: %v", inputPath, err)
	}

	logger.DebugPrintf("Opening output device %s", outputPath)
	outputFile, err := os.OpenFile(outputPath, os.O_WRONLY, 0666)
	if err != nil {
		inputFile.Close()
		return nil, nil, fmt.Errorf("failed to open output device %s: %v", outputPath, err)
	}

	return inputFile, outputFile, nil
}

func processEvents(ctx context.Context, inputFile, outputFile *os.File, eventConverter EventConverter, deviceName string) error {
	event := InputEvent{}

	for {
		select {
		case <-ctx.Done():
			logger.Printf("Relay shutdown for %s", deviceName)
			sendReleaseEvents(outputFile)
			return nil

		default:
			if err := binary.Read(inputFile, binary.LittleEndian, &event); err != nil {
				return fmt.Errorf("read error: %v", err)
			}

			if !eventConverter.validateEvent(event) {
				continue
			}

			logger.Printf("Read event from %s: Type=%d, Code=%d, Value=%d\n", 
				deviceName, event.Type, event.Code, event.Value)

			if err := handleEvent(outputFile, event, eventConverter, deviceName); err != nil {
				return err
			}
		}
	}
}

func handleEvent(outputFile *os.File, event InputEvent, eventConverter EventConverter, deviceName string) error {
	report, err := eventConverter.convertEvent(event)
	if err != nil {
		logger.DebugPrintf("[DEBUG] Error converting event: %v", err)
		return nil // Non-fatal error, continue processing
	}

	if report != nil {
		if _, err := outputFile.Write(report); err != nil {
			return fmt.Errorf("write error: %v", err)
		}
		logger.DebugPrintf("[DEBUG] %s event relayed", deviceName)
	}
	return nil
}

func sendReleaseEvents(outputFile *os.File) {
	releaseReport := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < 3; i++ {
		outputFile.Write(releaseReport)
		time.Sleep(10 * time.Millisecond)
	}
}
