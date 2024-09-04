#!/bin/bash

# Function to scan and list available devices
scan_and_list_devices() {
    echo "Scanning for Bluetooth devices... (10 seconds)"
    bluetoothctl scan on &
    sleep 10
    bluetoothctl scan off

    echo "Available devices:"
    bluetoothctl devices | nl
}

# Function to pair and connect to a Bluetooth device
pair_and_connect() {
    local device_mac=$1
    local device_name=$2

    echo "Attempting to pair and connect to $device_name ($device_mac)"
    
    bluetoothctl << EOF
    power on
    agent on
    default-agent
    pair $device_mac
    trust $device_mac
    connect $device_mac
    quit
EOF

    echo "Pairing process completed for $device_name"
}

# Main script
echo "Bluetooth Device Pairing Script"
echo "--------------------------------"

# Scan and list devices
scan_and_list_devices

# Function to get device selection
get_device_selection() {
    local device_type=$1
    local selection
    local device_mac

    while true; do
        read -p "Enter the number of the $device_type (or 'r' to rescan): " selection
        if [[ $selection == "r" ]]; then
            scan_and_list_devices
        elif [[ $selection =~ ^[0-9]+$ ]]; then
            device_mac=$(bluetoothctl devices | sed -n "${selection}p" | awk '{print $2}')
            if [[ -n $device_mac ]]; then
                echo $device_mac
                return
            else
                echo "Invalid selection. Please try again."
            fi
        else
            echo "Invalid input. Please enter a number or 'r' to rescan."
        fi
    done
}

# Get keyboard selection
keyboard_mac=$(get_device_selection "keyboard")
pair_and_connect "$keyboard_mac" "Keyboard"

# Get mouse selection
mouse_mac=$(get_device_selection "mouse")
pair_and_connect "$mouse_mac" "Mouse"

echo "Pairing process completed. Your devices should now be connected."