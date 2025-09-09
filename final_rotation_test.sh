#!/bin/bash

# Final test to verify log rotation functionality
echo "Final log rotation test - targeting exactly 1MB..."

# Create temporary version with even smaller threshold for guaranteed rotation
echo "Creating test version with 500KB threshold..."

# Backup and modify
cp /home/stevecarr/scvis-go/cmd/web/main.go /home/stevecarr/scvis-go/cmd/web/main.go.backup2

# Set very small threshold to guarantee rotation
sed -i 's/maxSizeMB := 10/maxSizeMB := 1/' /home/stevecarr/scvis-go/cmd/web/main.go
sed -i 's/int64(maxSizeMB)\*1024\*1024/int64(500*1024)/' /home/stevecarr/scvis-go/cmd/web/main.go

echo "Modified rotation threshold to 500KB for guaranteed testing"

# Rebuild and redeploy
cd /home/stevecarr/scvis-go
sudo ./deploy.sh > /dev/null 2>&1

echo "Waiting for service to stabilize..."
sleep 2

# Function to check log files
check_logs() {
    echo "Log files:"
    sudo ls -lah /var/log/scvisual/scvisual.log* 2>/dev/null || echo "No log files found"
    
    if [ -f "/var/log/scvisual/scvisual.log" ]; then
        local size=$(stat -c%s "/var/log/scvisual/scvisual.log")
        local size_kb=$((size / 1024))
        echo "Current log: $size bytes ($size_kb KB)"
    fi
}

echo "Initial state:"
check_logs

echo ""
echo "Generating activity to trigger rotation..."

# Generate enough activity to definitely trigger rotation
for i in $(seq 1 500); do
    # Generate requests with longer URLs to create bigger log entries
    curl -k -s "https://localhost:4000/this/is/a/very/long/url/path/that/will/generate/larger/log/entries/when/the/server/returns/404/errors/request-number-$i-with-timestamp-$(date +%s%N)" > /dev/null 2>&1 || true
    
    # Check every 100 requests
    if [ $((i % 100)) -eq 0 ]; then
        echo "Completed $i requests..."
        
        if [ -f "/var/log/scvisual/scvisual.log.1" ]; then
            echo "🎉 ROTATION DETECTED after $i requests!"
            break
        fi
    fi
done

echo ""
echo "Final state:"
check_logs

# Check rotation result
if [ -f "/var/log/scvisual/scvisual.log.1" ]; then
    echo ""
    echo "✅ SUCCESS: Log rotation working correctly!"
    echo ""
    echo "Rotation details:"
    sudo ls -lah /var/log/scvisual/scvisual.log*
    echo ""
    echo "Original log (rotated):"
    sudo ls -lah /var/log/scvisual/scvisual.log.1
    echo ""
    echo "New log file:"
    sudo ls -lah /var/log/scvisual/scvisual.log
else
    echo ""
    echo "❌ Rotation not triggered. This might indicate an issue."
fi

# Restore original configuration
echo ""
echo "Restoring original configuration..."
mv /home/stevecarr/scvis-go/cmd/web/main.go.backup2 /home/stevecarr/scvis-go/cmd/web/main.go
sudo ./deploy.sh > /dev/null 2>&1

echo ""
echo "✅ Test complete! Application restored to production configuration (10MB threshold)."
