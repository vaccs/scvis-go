#!/bin/bash

# Script to create a new user in the scvis_local database
# This script builds and runs the Go create_user program

echo "SCVIS User Creation Tool"
echo "======================="
echo ""

# Check if we're in the scripts directory
if [ ! -f "create_user.go" ]; then
    echo "Error: create_user.go not found. Please run this script from the scripts directory."
    exit 1
fi

# Build the Go program
echo "Building user creation tool..."
go build -o create_user create_user.go

if [ $? -ne 0 ]; then
    echo "Error: Failed to build the user creation tool."
    exit 1
fi

# Run the program
echo "Starting user creation process..."
echo ""
./create_user

# Clean up the binary
rm -f create_user

echo ""
echo "User creation tool finished."
