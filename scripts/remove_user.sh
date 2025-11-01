#!/bin/bash

# Script to remove a user from the SCVis database
# This will delete the user and all associated data (projects, cfiles, inputfiles, sessions)

set -e

echo "SCVis User Removal Script"
echo "========================="
echo ""

# Check if running from correct directory
if [ ! -f ".env" ]; then
    echo "❌ .env file not found"
    echo "   Please run this script from the scvis-go project root"
    exit 1
fi

# Load database configuration from .env file
source .env

DB_HOST="${DB_HOST:-localhost}"
DB_USER="${DB_USER:-scvis}"
DB_NAME="${DB_NAME:-scvis_local}"

echo "Database Configuration:"
echo "  Host: $DB_HOST"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo ""

# Function to execute MySQL command
mysql_exec() {
    local command="$1"
    mysql -h"$DB_HOST" -u"$DB_USER" -p"$DB_PASSWD" "$DB_NAME" -e "$command" 2>/dev/null
}

# Function to execute MySQL query and return results
mysql_query() {
    local query="$1"
    mysql -h"$DB_HOST" -u"$DB_USER" -p"$DB_PASSWD" "$DB_NAME" -N -e "$query" 2>/dev/null
}

# Test database connection
if ! mysql_query "SELECT 1" >/dev/null 2>&1; then
    echo "❌ Cannot connect to database"
    echo "   Check your .env configuration and MySQL service"
    exit 1
fi

# Display available users
echo "Available Users:"
echo "----------------"
mysql -h"$DB_HOST" -u"$DB_USER" -p"$DB_PASSWD" "$DB_NAME" -t -e "SELECT id, name, email, CASE WHEN admin = 1 THEN 'Yes' ELSE 'No' END AS admin, created FROM users ORDER BY id;" 2>/dev/null
echo ""

# Method 1: Remove by user ID
read -p "Enter user ID to remove (or press Enter to search by email): " user_id

if [ -z "$user_id" ]; then
    # Method 2: Search by email
    read -p "Enter email address: " email
    
    if [ -z "$email" ]; then
        echo "❌ No user ID or email provided"
        exit 1
    fi
    
    # Find user by email
    user_id=$(mysql_query "SELECT id FROM users WHERE email = '$email';" | head -n1)
    
    if [ -z "$user_id" ]; then
        echo "❌ User with email '$email' not found"
        exit 1
    fi
    
    echo "Found user ID: $user_id"
fi

# Validate user ID is a number
if ! [[ "$user_id" =~ ^[0-9]+$ ]]; then
    echo "❌ Invalid user ID. Must be a number."
    exit 1
fi

# Get user details
user_name=$(mysql_query "SELECT name FROM users WHERE id = $user_id;" | head -n1)
user_email=$(mysql_query "SELECT email FROM users WHERE id = $user_id;" | head -n1)
user_admin=$(mysql_query "SELECT admin FROM users WHERE id = $user_id;" | head -n1)

if [ -z "$user_name" ]; then
    echo "❌ User with ID $user_id not found"
    exit 1
fi

echo ""
echo "User Details:"
echo "-------------"
echo "  ID: $user_id"
echo "  Name: $user_name"
echo "  Email: $user_email"
echo "  Admin: $([ "$user_admin" = "1" ] && echo "Yes" || echo "No")"
echo ""

# Get statistics
project_count=$(mysql_query "SELECT COUNT(*) FROM projects WHERE userid = $user_id;" | head -n1)
project_count=${project_count:-0}

echo "User Statistics:"
echo "  Projects: $project_count"
echo ""

# Warning about cascading deletes
echo "⚠️  WARNING: This action will DELETE:"
echo "  ✗ The user account for '$user_name' ($user_email)"
echo "  ✗ All $project_count project(s) owned by this user"
echo "  ✗ All C files and input files in those projects"
echo "  ✗ All analysis data for those projects"
echo "  ✗ All sessions for this user"
echo ""
echo "⚠️  This action CANNOT be undone!"
echo ""

# Confirmation step 1: Email verification
read -p "Type the user's email address to confirm: " confirm_email

if [ "$confirm_email" != "$user_email" ]; then
    echo "❌ Email does not match. Operation cancelled."
    exit 1
fi

# Confirmation step 2: Final verification
read -p "Are you absolutely sure? Type 'DELETE' in uppercase to proceed: " final_confirm

if [ "$final_confirm" != "DELETE" ]; then
    echo "❌ Operation cancelled."
    exit 0
fi

echo ""
echo "🗑️  Removing user '$user_name' and all associated data..."
echo ""

# Delete user (database constraints will cascade delete related records)
mysql_exec "DELETE FROM users WHERE id = $user_id;"

if [ $? -eq 0 ]; then
    echo "✅ User successfully removed!"
    echo ""
    echo "Summary of deleted data:"
    echo "  ✓ User account: $user_name ($user_email)"
    echo "  ✓ Projects: $project_count"
    echo "  ✓ All associated C files and input files"
    echo "  ✓ All associated sessions"
    echo ""
    echo "🎉 User removal complete!"
else
    echo "❌ Error removing user from database"
    exit 1
fi
