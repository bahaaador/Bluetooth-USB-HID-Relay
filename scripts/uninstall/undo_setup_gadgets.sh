#!/bin/bash


 # Clean up gadget directory if it exists

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

cleanup_gadget