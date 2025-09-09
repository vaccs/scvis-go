#!/bin/bash

# Test script to verify log rotation functionality
echo "Testing log rotation for SCVIS application..."

# Function to check log file size
check_log_size() {
    local log_file="/var/log/scvisual/scvisual.log"
    if [ -f "$log_file" ]; then
        local size=$(stat -c%s "$log_file")
        echo "Current log file size: $size bytes"
        return $size
    else
        echo "Log file not found: $log_file"
        return 0
    fi
}

# Function to list all log files
list_log_files() {
    echo "Current log files in /var/log/scvisual/:"
    sudo ls -la /var/log/scvisual/
}

# Function to generate log entries by making HTTP requests
generate_log_activity() {
    echo "Generating log activity by making HTTP requests..."
    local count=${1:-50}
    
    for i in $(seq 1 $count); do
        # Make requests to generate access logs (these will be caught by error handling and logged)
        curl -k -s "https://localhost:4000/nonexistent-page-$i" > /dev/null 2>&1 || true
        curl -k -s "https://localhost:4000/static/nonexistent-file-$i.css" > /dev/null 2>&1 || true
        
        # Add a small delay to avoid overwhelming the server
        sleep 0.1
    done
    echo "Generated $count log entries"
}

# Initial state
echo "=== Initial State ==="
list_log_files
check_log_size
initial_size=$?

echo ""
echo "=== Generating Log Activity ==="
# Generate some log activity
generate_log_activity 100

# Check size after activity
echo ""
echo "=== After Initial Activity ==="
sleep 2  # Give time for logs to be written
check_log_size
after_activity_size=$?

# If the log file is still small (less than 500KB), generate more activity
if [ $after_activity_size -lt 500000 ]; then
    echo ""
    echo "Log file still small, generating more activity..."
    
    # Generate larger log entries by making more requests
    for batch in $(seq 1 10); do
        echo "Batch $batch of intensive logging..."
        generate_log_activity 200
        
        # Check current size
        check_log_size
        current_size=$?
        
        # If we're getting close to 1MB (our rotation threshold for local), continue
        if [ $current_size -gt 800000 ]; then
            echo "Log file approaching rotation threshold..."
        fi
        
        # If rotation happened, break
        if [ -f "/var/log/scvisual/scvisual.log.1" ]; then
            echo "Log rotation detected!"
            break
        fi
    done
fi

echo ""
echo "=== Final State ==="
list_log_files
check_log_size

# Check if rotation occurred
if [ -f "/var/log/scvisual/scvisual.log.1" ]; then
    echo ""
    echo "✅ SUCCESS: Log rotation occurred!"
    echo "Rotated log file details:"
    sudo ls -la /var/log/scvisual/scvisual.log.1
else
    echo ""
    echo "ℹ️  Log rotation threshold not reached yet."
    echo "This is normal if the log file hasn't reached the rotation size limit."
fi

echo ""
echo "=== Testing Complete ==="
