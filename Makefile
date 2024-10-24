.PHONY: build clean test install uninstall run

# Default target
all: build

# Build the application
build:
	@echo "Building bt-hid-relay..."
	@go build -o bin/bt-hid-relay ./cmd/bt-hid-relay/

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f bin/bt-hid-relay
	@echo "Done"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Install the service
install: build
	@echo "Installing bt-hid-relay service..."
	cp bin/bt-hid-relay /usr/local/bin/
	cp scripts/bt-hid-relay.service /etc/systemd/system/
	systemctl daemon-reload
	systemctl enable bt-hid-relay.service
	systemctl start bt-hid-relay.service

# Uninstall the service
uninstall:
	@echo "Uninstalling bt-hid-relay service..."
	systemctl stop bt-hid-relay.service
	systemctl disable bt-hid-relay.service
	rm /etc/systemd/system/bt-hid-relay.service
	rm /usr/local/bin/bt-hid-relay
	systemctl daemon-reload

# Run the application
run: build
	@echo "Running bt-hid-relay..."
	@sudo ./bin/bt-hid-relay -debug
