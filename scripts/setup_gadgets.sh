#!/bin/bash
set -e

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 
   exit 1
fi

# check if modules are loaded
MODULES_LOADED=0
if lsmod | grep -E "g_ether|usb_f_rndis|usb_f_ecm|u_ether" > /dev/null; then
    MODULES_LOADED=1
fi

# check if gadget exists
GADGET_EXISTS=0
if [ -d /sys/kernel/config/usb_gadget/hid_gadget ]; then
    GADGET_EXISTS=1
fi

# if modules are loaded or gadget exists, ask for confirmation
if [ $MODULES_LOADED -eq 1 ] || [ $GADGET_EXISTS -eq 1 ]; then
    echo "WARNING: Found existing USB configuration:"
    [ $MODULES_LOADED -eq 1 ] && echo "- USB modules are currently loaded"
    [ $GADGET_EXISTS -eq 1 ] && echo "- HID gadget configuration exists"
    echo "Proceeding will remove these configurations and may disrupt existing USB devices."
    read -p "Do you want to proceed? (y/N) " -n 1 -r
    echo    # Move to a new line
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Operation cancelled."
        exit 0
    fi
fi

cleanup_gadget() {
    echo "Cleaning up existing USB configurations..."
    
    # Only unload modules if they exist
    if [ $MODULES_LOADED -eq 1 ]; then
        echo "Unloading existing USB modules..."
        modprobe -r g_ether usb_f_rndis usb_f_ecm u_ether || true
    fi
    
    # Clean up gadget directory if it exists
    if [ $GADGET_EXISTS -eq 1 ]; then
        echo "Removing existing HID gadget configuration..."
        cd /sys/kernel/config/usb_gadget/hid_gadget
        if [ -f UDC ]; then
            echo "" > UDC
        fi
        
        # Remove symbolic links
        rm -f configs/c.1/hid.usb0
        rm -f configs/c.1/hid.usb1
        
        # Remove directories
        rm -rf functions/hid.usb0
        rm -rf functions/hid.usb1
        rm -rf configs/c.1/strings/0x409
        rm -rf configs/c.1
        rm -rf strings/0x409
        
        cd ..
        rmdir hid_gadget 2>/dev/null || true
    fi
}

# Only proceed with cleanup if user confirmed when needed
if [ $MODULES_LOADED -eq 0 ] && [ $GADGET_EXISTS -eq 0 ] || [[ $REPLY =~ ^[Yy]$ ]]; then
    cleanup_gadget
else
    echo "Operation cancelled."
    exit 0
fi

echo "Setting up HID gadget..."

# Create gadget
mkdir -p /sys/kernel/config/usb_gadget/hid_gadget
cd /sys/kernel/config/usb_gadget/hid_gadget

# Set up USB device descriptor
echo 0x1d6b > idVendor  # Linux Foundation
echo 0x0104 > idProduct # Multifunction Composite Gadget
echo 0x0100 > bcdDevice # v1.0.0
echo 0x0200 > bcdUSB    # USB2

# Set up strings
mkdir -p strings/0x409
echo "fedcba9876543210" > strings/0x409/serialnumber
echo "Your Name" > strings/0x409/manufacturer
echo "BT HID Relay" > strings/0x409/product

# Set up configuration
mkdir -p configs/c.1/strings/0x409
echo "Config 1: HID" > configs/c.1/strings/0x409/configuration
echo 250 > configs/c.1/MaxPower

# Set up HID function
# Refer to https://www.usb.org/sites/default/files/documents/hid1_11.pdf for full report descriptor reference (pages 69 and 71 contain the keyboard and mouse descriptor examples)

# Set up Mouse HID function
mkdir -p functions/hid.usb0
echo 0 > functions/hid.usb0/protocol
echo 0 > functions/hid.usb0/subclass
echo 4 > functions/hid.usb0/report_length
echo -ne \\x05\\x01\\x09\\x02\\xa1\\x01\\x09\\x01\\xa1\\x00\\x05\\x09\\x19\\x01\\x29\\x03\\x15\\x00\\x25\\x01\\x95\\x03\\x75\\x01\\x81\\x02\\x95\\x01\\x75\\x05\\x81\\x03\\x05\\x01\\x09\\x30\\x09\\x31\\x09\\x38\\x15\\x81\\x25\\x7f\\x75\\x08\\x95\\x03\\x81\\x06\\xc0\\xc0 > functions/hid.usb0/report_desc

# Set up Keyboard HID function
mkdir -p functions/hid.usb1
echo 1 > functions/hid.usb1/protocol
echo 1 > functions/hid.usb1/subclass
echo 8 > functions/hid.usb1/report_length
echo -ne \\x05\\x01\\x09\\x06\\xa1\\x01\\x05\\x07\\x19\\xe0\\x29\\xe7\\x15\\x00\\x25\\x01\\x75\\x01\\x95\\x08\\x81\\x02\\x95\\x01\\x75\\x08\\x81\\x03\\x95\\x06\\x75\\x08\\x15\\x00\\x25\\x65\\x05\\x07\\x19\\x00\\x29\\x65\\x81\\x00\\xc0 > functions/hid.usb1/report_desc

ln -s functions/hid.usb0 configs/c.1/
ln -s functions/hid.usb1 configs/c.1/

# Enable gadget
UDC=$(ls /sys/class/udc)
echo $UDC > UDC

echo "HID gadget setup complete."

# Wait for a moment to ensure the device is recognized
sleep 2

# Check if the HID devices was created
if [ -e /dev/hidg0 ]; then
    echo "HID device mouse /dev/hidg0 created successfully."
else
    echo "Error: HID device mouse /dev/hidg0 not created."
    exit 1
fi

if [ -e /dev/hidg1 ]; then
    echo "HID device keyboard /dev/hidg1 created successfully."
else
    echo "Error: HID device keyboard /dev/hidg1 not created."
    exit 1
fi

echo "Setup completed successfully."