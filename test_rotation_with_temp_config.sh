#!/bin/bash

# Test log rotation by temporarily modifying the rotation threshold
echo "Testing log rotation with temporary smaller threshold..."

# Create a temporary version of main.go with smaller rotation threshold
echo "Creating temporary main.go with 1MB rotation threshold for testing..."

# Backup current main.go
cp /home/stevecarr/scvis-go/cmd/web/main.go /home/stevecarr/scvis-go/cmd/web/main.go.backup

# Modify the rotation threshold in main.go for testing
sed -i 's/maxSizeMB := 10/maxSizeMB := 1/' /home/stevecarr/scvis-go/cmd/web/main.go

echo "Modified rotation threshold to 1MB for testing"

# Rebuild and redeploy
echo "Rebuilding application with test configuration..."
cd /home/stevecarr/scvis-go
sudo ./deploy.sh

echo "Waiting for service to start..."
sleep 3

# Check service status
sudo systemctl status scvisual --no-pager

echo ""
echo "Starting log rotation test with 1MB threshold..."

# Function to check log file size
check_log_size() {
    local log_file="/var/log/scvisual/scvisual.log"
    if [ -f "$log_file" ]; then
        local size=$(stat -c%s "$log_file")
        local size_kb=$((size / 1024))
        echo "Current log file size: $size bytes ($size_kb KB)"
        return $size
    else
        echo "Log file not found: $log_file"
        return 0
    fi
}

# Generate log activity
echo "Generating log activity to trigger rotation..."
for batch in $(seq 1 10); do
    echo "Batch $batch: Generating requests..."
    
    for i in $(seq 1 200); do
        curl -k -s "https://localhost:4000/test-rotation-batch-$batch-request-$i-$(date +%s)" > /dev/null 2>&1 || true
    done
    
    # Check current size
    check_log_size
    current_size=$?
    
    # Check if rotation occurred
    if [ -f "/var/log/scvisual/scvisual.log.1" ]; then
        echo "🎉 SUCCESS: Log rotation detected!"
        echo ""
        echo "Log files after rotation:"
        sudo ls -lah /var/log/scvisual/scvisual.log*
        break
    fi
    
    # If we're over 800KB, we're getting close
    if [ $current_size -gt 819200 ]; then
        echo "⚠️  Approaching 1MB threshold..."
    fi
    
    sleep 1
done

# Final check
echo ""
echo "=== Final State ==="
sudo ls -lah /var/log/scvisual/

# Restore original main.go
echo ""
echo "Restoring original configuration..."
mv /home/stevecarr/scvis-go/cmd/web/main.go.backup /home/stevecarr/scvis-go/cmd/web/main.go

# Redeploy with original configuration
echo "Redeploying with original 10MB threshold..."
sudo ./deploy.sh

echo ""
echo "Test complete! Application restored to production configuration."
