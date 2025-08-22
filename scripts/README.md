# User Creation Script for SCVIS

This directory contains scripts to create new users in the scvis_local database.

## Files

- `create_user.go` - The main Go program that creates users
- `create_user.sh` - Shell script wrapper for easy execution
- `create_user` - Compiled binary (created when script runs)

## Usage

### Method 1: Using the Shell Script (Recommended)

```bash
cd scripts
./create_user.sh
```

### Method 2: Running the Go Program Directly

```bash
cd scripts
go run create_user.go
```

## What the Script Does

1. **Connects to Database**: Uses credentials from `../.env` to connect to the scvis_local database
2. **Collects User Information**:
   - Name (visible input)
   - Email address (visible input)
   - Password (hidden input, minimum 8 characters)
   - Password confirmation (hidden input)
   - Administrator status (y/N)

3. **Validates Input**:
   - Ensures all required fields are filled
   - Validates password length (minimum 8 characters)
   - Confirms passwords match
   - Checks for duplicate email addresses

4. **Secure Password Handling**:
   - Uses bcrypt with cost 12 (same as the main application)
   - Passwords are not visible when typing
   - Passwords are hashed before storage

5. **Creates User**: Inserts the new user into the database with all provided information

## Security Features

- **Hidden Password Input**: Passwords are not displayed on screen while typing
- **Password Confirmation**: Requires password to be entered twice
- **Bcrypt Hashing**: Uses the same hashing method as the main application (cost 12)
- **Duplicate Prevention**: Prevents creating users with duplicate email addresses
- **Input Validation**: Validates all inputs before database insertion

## Example Session

```
SCVIS User Creation Tool
=======================

Building user creation tool...
Starting user creation process...

Connected to database successfully!
Creating new user for scvis_local database
=========================================
Enter user name: John Doe
Enter email address: john.doe@example.com
Enter password (min 8 characters): [hidden]
Confirm password: [hidden]
Is this user an administrator? (y/N): n

User created successfully!
Name: John Doe
Email: john.doe@example.com
Administrator: false

The user can now log in to the application.
```

## Requirements

- Go 1.19+ (for golang.org/x/term package)
- Access to scvis_local database
- Valid `.env.debug` file in parent directory

## Troubleshooting

- **Database Connection Error**: Check that MySQL is running and credentials in `.env.debug` are correct
- **Duplicate Email Error**: The email address is already in use by another user
- **Permission Error**: Ensure the database user has INSERT permissions on the users table
