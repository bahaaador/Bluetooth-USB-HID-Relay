#!/bin/bash

# Function to scan and list available devices
scan_and_list_devices() {
    echo "Scanning for Bluetooth devices... (15 seconds)"
    bluetoothctl power on
    bluetoothctl scan on > /dev/null 2>&1 &
    scan_pid=$!
    sleep 15
    kill $scan_pid
    bluetoothctl scan off > /dev/null 2>&1

    echo "Scan completed. Available devices:"
    devices=$(bluetoothctl devices)
    echo "$devices" | while read -r line; do
        mac=$(echo "$line" | awk '{print $2}')
        name=$(bluetoothctl info "$mac" | grep "Name" | cut -d ":" -f2 | xargs)
        if [ -z "$name" ]; then
            name="Unknown"
        fi
        echo "$line - $name"
    done | nl
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
remove $device_mac
scan on
pair $device_mac
trust $device_mac
connect $device_mac
scan off
EOF

    echo "Pairing process completed for $device_name"
}

# Main script
echo "Bluetooth Device Pairing Script"
echo "--------------------------------"

# Function to get device selection
get_device_selection() {
    local device_type=$1
    local selection
    local device_mac
    local device_name

    while true; do
        read -p "Enter the number of the $device_type (or 'r' to rescan): " selection
        if [[ $selection == "r" ]]; then
            scan_and_list_devices
        elif [[ $selection =~ ^[0-9]+$ ]]; then
            device_info=$(bluetoothctl devices | sed -n "${selection}p")
            device_mac=$(echo "$device_info" | awk '{print $2}')
            device_name=$(bluetoothctl info "$device_mac" | grep "Name" | cut -d ":" -f2 | xargs)
            if [ -z "$device_name" ]; then
                device_name="Unknown"
            fi
            if [[ -n $device_mac ]]; then
                echo "$device_mac:$device_name"
                return
            else
                echo "Invalid selection. Please try again."
            fi
        else
            echo "Invalid input. Please enter a number or 'r' to rescan."
        fi
    done
}

# Scan and list devices
scan_and_list_devices

# Get keyboard selection
keyboard_info=$(get_device_selection "keyboard")
keyboard_mac=$(echo "$keyboard_info" | cut -d':' -f1)
keyboard_name=$(echo "$keyboard_info" | cut -d':' -f2-)
pair_and_connect "$keyboard_mac" "$keyboard_name"

# Get mouse selection
mouse_info=$(get_device_selection "mouse")
mouse_mac=$(echo "$mouse_info" | cut -d':' -f1)
mouse_name=$(echo "$mouse_info" | cut -d':' -f2-)
pair_and_connect "$mouse_mac" "$mouse_name"

echo "Pairing process completed. Your devices should now be connected."
