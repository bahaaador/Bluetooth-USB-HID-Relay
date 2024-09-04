#!/bin/bash

# Backup config.txt
sudo cp /boot/config.txt /boot/config.txt.bak

# Enable dwc2 driver
echo "dtoverlay=dwc2" | sudo tee -a /boot/config.txt

# Backup modules
sudo cp /etc/modules /etc/modules.bak

# Load dwc2 module
echo "dwc2" | sudo tee -a /etc/modules

# Load libcomposite module
echo "libcomposite" | sudo tee -a /etc/modules


# Disable USB mass storage gadget
echo "dtoverlay=dwc2,dr_mode=peripheral" | sudo tee -a /boot/config.txt

echo "USB host setup complete. Please reboot for changes to take effect."