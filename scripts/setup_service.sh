#!/bin/bash

# Create a systemd service file
cat << EOF > /etc/systemd/system/bt-hid-relay.service
[Unit]
Description=Bluetooth HID Relay Service
After=network.target bluetooth.target

[Service]
ExecStart=/usr/local/bin/bt-hid-relay
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOF

# Assuming the compiled Go binary is named 'bt-hid-relay'
# and is located in the project directory
cp bt-hid-relay /usr/local/bin/

# Set correct permissions
chmod 755 /usr/local/bin/bt-hid-relay

# Reload systemd, enable and start the service
systemctl daemon-reload
systemctl enable bt-hid-relay.service
systemctl start bt-hid-relay.service

echo "Bluetooth HID Relay service has been set up and started"