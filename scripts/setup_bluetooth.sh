#!/bin/bash

# Install Bluetooth libraries and tools
apt-get update
apt-get install -y bluetooth bluez bluez-tools

# Enable Bluetooth service
systemctl enable bluetooth
systemctl start bluetooth

# Make Bluetooth discoverable
bluetoothctl discoverable on

echo "Bluetooth setup complete"