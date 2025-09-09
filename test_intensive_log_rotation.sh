#!/bin/bash

# Enhanced test script to force log rotation by generating more intensive logging
echo "Enhanced log rotation test - generating more intensive activity..."

# Function to check log file size
check_log_size() {
    local log_file="/var/log/scvisual/scvisual.log"
    if [ -f "$log_file" ]; then
        local size=$(stat -c%s "$log_file")
        local size_mb=$((size / 1024 / 1024))
        echo "Current log file size: $size bytes ($size_mb MB)"
        return $size
    else
        echo "Log file not found: $log_file"
        return 0
    fi
}

# Function to list all log files
list_log_files() {
    echo "Current log files in /var/log/scvisual/:"
    sudo ls -lah /var/log/scvisual/
}

# Function to generate intensive log activity
generate_intensive_log_activity() {
    echo "Generating intensive log activity..."
    local batches=${1:-20}
    
    for batch in $(seq 1 $batches); do
        echo "Processing batch $batch/$batches..."
        
        # Generate many concurrent requests to create more log entries
        for i in $(seq 1 100); do
            {
                # Multiple types of requests to generate different log entries
                curl -k -s "https://localhost:4000/nonexistent/path/very/long/url/that/will/generate/error/logs/batch-$batch-request-$i" > /dev/null 2>&1 || true
                curl -k -s "https://localhost:4000/static/css/nonexistent-very-long-filename-that-should-generate-detailed-error-logs-$batch-$i.css" > /dev/null 2>&1 || true
                curl -k -s -X POST "https://localhost:4000/some/post/endpoint/that/does/not/exist/batch-$batch-$i" > /dev/null 2>&1 || true
            } &
        done
        
        # Wait for all background requests to complete
        wait
        
        # Check size every few batches
        if [ $((batch % 5)) -eq 0 ]; then
            check_log_size
            current_size=$?
            
            # Check if rotation occurred
            if [ -f "/var/log/scvisual/scvisual.log.1" ]; then
                echo "🎉 Log rotation detected after batch $batch!"
                return 0
            fi
            
            # If we're over 8MB, we're getting close to the 10MB threshold
            if [ $current_size -gt 8388608 ]; then
                echo "⚠️  Log file is getting large (>8MB), rotation should happen soon..."
            fi
        fi
        
        # Small delay between batches
        sleep 1
    done
}

# Initial state
echo "=== Initial State ==="
list_log_files
check_log_size

echo ""
echo "=== Starting Intensive Log Generation ==="
echo "Note: The application is configured for 10MB rotation in production mode."
echo "This may take a while to reach the threshold..."

# Generate intensive activity
generate_intensive_log_activity 50

echo ""
echo "=== Final Check ==="
list_log_files
check_log_size

# Final check for rotation
if [ -f "/var/log/scvisual/scvisual.log.1" ]; then
    echo ""
    echo "✅ SUCCESS: Log rotation occurred!"
    echo ""
    echo "Rotation details:"
    sudo ls -lah /var/log/scvisual/scvisual.log*
else
    echo ""
    echo "ℹ️  Log rotation threshold (10MB) not reached yet."
    echo "Current log size vs threshold:"
    current_size=$(stat -c%s "/var/log/scvisual/scvisual.log" 2>/dev/null || echo 0)
    threshold=$((10 * 1024 * 1024))
    percentage=$((current_size * 100 / threshold))
    echo "Progress: $percentage% of 10MB threshold"
fi

echo ""
echo "=== Test Complete ==="
