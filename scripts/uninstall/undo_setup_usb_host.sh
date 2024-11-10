#!/bin/bash

# Restore config.txt
sudo cp /boot/config.txt.bak /boot/config.txt

# Remove dwc2 driver
sudo sed -i '/dtoverlay=dwc2/d' /boot/config.txt