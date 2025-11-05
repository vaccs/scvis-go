#!/bin/bash

# Deployment script for SCVISUAL to run as www-data user
# This script should be run as root

set -e

echo "Setting up SCVISUAL to run as www-data user..."

# Create application directory
sudo mkdir -p /opt/scvisual
sudo chown www-data:www-data /opt/scvisual

# Create log directory
sudo mkdir -p /var/log/scvisual
sudo chown www-data:www-data /var/log/scvisual
sudo chmod 755 /var/log/scvisual

# Create config directory
sudo mkdir -p /etc/scvisual
sudo chown root:www-data /etc/scvisual
sudo chmod 750 /etc/scvisual

# Copy configuration file if it exists
if [ -f "./.env" ]; then
    sudo cp ./.env /opt/scvisual/
    sudo chown www-data:www-data /opt/scvisual/.env
    sudo chmod 600 /opt/scvisual/.env
fi

# Build the application
echo "Building SCVISUAL application..."
go build -o scvisual ./cmd/web

# Stop the service if it's running
echo "Stopping scvisual service if running..."
if sudo systemctl is-active --quiet scvisual; then
    echo "Service is running, stopping it..."
    sudo systemctl stop scvisual
    SERVICE_WAS_RUNNING=true
else
    echo "Service is not running"
    SERVICE_WAS_RUNNING=false
fi

# Copy binary to application directory
sudo cp scvisual /opt/scvisual/
sudo chown www-data:www-data /opt/scvisual/scvisual
sudo chmod 755 /opt/scvisual/scvisual

# Copy static files if they exist
if [ -d "./visual" ]; then
    sudo cp -r ./visual /opt/scvisual/
    sudo chown -R www-data:www-data /opt/scvisual/visual
fi

if [ -d "./ui" ]; then
    sudo cp -r ./ui /opt/scvisual/
    sudo chown -R www-data:www-data /opt/scvisual/ui
fi

# Ensure www-data can access SSL certificates
echo "Checking SSL certificate access..."
if [ -f "/etc/letsencrypt/live/vaccs.ceas.wmich.edu/fullchain.pem" ] && [ -f "/etc/letsencrypt/live/vaccs.ceas.wmich.edu/privkey.pem" ]; then
    # Ensure www-data is in ssl-cert group
    sudo usermod -a -G ssl-cert www-data
    
    # Let's Encrypt certificates should already have proper permissions for ssl-cert group
    echo "SSL certificates configured for www-data access via ssl-cert group"
    echo "Certificate: /etc/letsencrypt/live/vaccs.ceas.wmich.edu/fullchain.pem"
    echo "Private key: /etc/letsencrypt/live/vaccs.ceas.wmich.edu/privkey.pem"
else
    echo "Warning: Let's Encrypt certificate files not found at expected locations"
    echo "Expected: /etc/letsencrypt/live/vaccs.ceas.wmich.edu/fullchain.pem and /etc/letsencrypt/live/vaccs.ceas.wmich.edu/privkey.pem"
fi

# Install systemd service
sudo cp scvisual.service /etc/systemd/system/
sudo systemctl daemon-reload

# Restart the service if it was running before
if [ "$SERVICE_WAS_RUNNING" = true ]; then
    echo "Restarting scvisual service..."
    sudo systemctl start scvisual
    echo "Service restarted!"
fi

echo "Setup complete!"
echo ""
if [ "$SERVICE_WAS_RUNNING" = false ]; then
    echo "To start the service:"
    echo "  sudo systemctl enable scvisual"
    echo "  sudo systemctl start scvisual"
    echo ""
fi
echo "To check status:"
echo "  sudo systemctl status scvisual"
echo ""
echo "To view logs:"
echo "  sudo journalctl -u scvisual -f"
