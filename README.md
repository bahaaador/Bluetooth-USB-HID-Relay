# Bluetooth USB HID Relay

A delightfully overengineered solution to an unusual modern problem: using Bluetooth peripherals with Bluetooth-disabled computers.

This project creates a Bluetooth USB HID relay using a Raspberry Pi Zero (or similar single-board computer). It allows you to use Bluetooth keyboards and mice with computers that have Bluetooth disabled, by presenting the board as a composite USB HID device.

## Inspiration

This project was born out of a desire to help a friend who couldn't use his favorite Bluetooth mouse and keyboard due to Bluetooth being disabled on his work laptop. As someone who enjoys tinkering and problem-solving, I saw this as an opportunity to create something useful while learning more about Linux internals, USB gadgets, and Go programming. The Raspberry Pi Zero became the perfect bridge, connecting Bluetooth peripherals to computers that wouldn't normally allow it.

```mermaid
graph TD
    B(Keyboard) -- Bluetooth -->  A{"HID <br> (Raspberry Pi Zero)"}
    C(Mouse) -- Bluetooth --> A
    A -- USB--> D(Host Computer)
```

## Features

- Connects to Bluetooth keyboards and mice
- Presents itself as a composite USB HID device (keyboard and mouse) to the host computer
- Works with Windows, Mac, and Linux computers
- Automatically starts the relay service at boot
- Configures the board as a USB OTG device
- Includes a script to help pair Bluetooth devices

## Components

- Raspberry Pi Zero or other single-board computer capable of USB OTG
- DietPi or another lightweight Linux distribution
- Bash scripts for setup and configuration
- Go program for handling Bluetooth input and USB HID output
- Systemd service for automatic startup

## Setup

1. Install DietPi (or your chosen Linux distribution) on your Raspberry Pi Zero or similar board.
2. Clone this repository to your board.
3. Run the setup scripts in the following order:

   ```bash
   sudo ./scripts/setup_usb_host.sh
   sudo reboot
   sudo ./scripts/setup_bluetooth.sh
   sudo ./scripts/setup_gadgets.sh
   ```

4. Install Task runner:
   ```bash
   # On Linux/macOS
   sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin

   # Using Go
   go install github.com/go-task/task/v3/cmd/task@latest
   ```

5. Build the project:
   ```bash
   task build
   ```

6. Set up the systemd service:
   ```bash
   sudo task install
   ```

7. Pair your Bluetooth devices:
   ```bash
   sudo ./scripts/pair_devices.sh
   ```
   You will be prompted to select your keyboard and mouse from the list of available devices.

## Usage

1. Power on your board and connect it to the target computer via USB.
2. If you haven't paired your Bluetooth keyboard and mouse yet, run the pairing script:
   ```bash
   sudo ./scripts/pair_devices.sh
   ```
   Follow the prompts to pair your devices.
3. The board will now relay input from the Bluetooth devices to the target computer.

## Common Tasks

This project uses Task runner for common operations. Here are the available commands:

- Build the project: `task build`
- Clean build artifacts: `task clean`
- Run tests: `task test`
- Build and run the application: `task run`
- Install the service: `sudo task install`
- Uninstall the service: `sudo task uninstall`

For example, to build and run the project:

```bash
task build
task run
```

## Verifying Bluetooth Devices

To verify the connection of Bluetooth devices and echo their inputs:

```bash
task verify
```

## Development

This project serves as a learning opportunity for:
- Linux internals
- Creating USB gadgets
- Using Go for hardware interfacing
- Setting up systemd services
- Configuring single-board computers as USB OTG devices
- Bluetooth device pairing and management

It's been a fun journey of discovery, and I hope others find it useful or inspiring for their own projects!

## Prerequisites

Before building and running the project, ensure you have the following installed:

1. Go (version 1.21 or later)
2. Task runner

On most Linux distributions, including Raspberry Pi OS (formerly Raspbian), you can install Go using:

```bash
sudo apt-get update
sudo apt-get install golang
```

## References

- [Adafruit Guide: Turning your Raspberry Pi Zero into a USB Gadget](https://cdn-learn.adafruit.com/downloads/pdf/turning-your-raspberry-pi-zero-into-a-usb-gadget.pdf)
- [Composite USB Gadgets on the Raspberry Pi Zero](https://www.isticktoit.net/?p=1383)
- [Task: A task runner / simpler Make alternative](https://taskfile.dev)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
