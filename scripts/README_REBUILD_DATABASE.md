# Database Rebuild Script

## Overview
The `rebuild_scvis_database.sh` script drops and recreates the `scvis` database based on the structure from `scvis_local`. This is useful for deploying to production or setting up new environments.

## ⚠️ WARNING
This script will **PERMANENTLY DELETE** all data in the `scvis` database! Use with caution!

## Prerequisites

Before running this script, ensure you have:

1. **MySQL root access** - The script requires MySQL root credentials
2. **Database user exists** - The `scvis_user` must be created beforehand
   - If the user doesn't exist, the script will provide instructions to create it
   - Or run `./scripts/install_database.sh` to set up both user and database
3. **Backup of existing data** (if any) - This script will DELETE the existing database

### Quick Setup for New Installation

If you're setting up on a fresh machine without the scvis_user:

```bash
# Option 1: Use the install script (recommended for first-time setup)
./scripts/install_database.sh

# Option 2: Create user manually, then run rebuild
mysql -u root -p
CREATE USER 'scvis_user'@'localhost' IDENTIFIED BY 'your_secure_password';
GRANT SELECT, INSERT, UPDATE, DELETE ON scvis.* TO 'scvis_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;

# Then run the rebuild script
./scripts/rebuild_scvis_database.sh
```

## Usage

### Basic Usage
```bash
./scripts/rebuild_scvis_database.sh
```

### With Custom Database User
```bash
DB_USER=custom_user ./scripts/rebuild_scvis_database.sh
```

## What It Does

1. **Prompts for Confirmation** - Asks user to confirm before proceeding
2. **Drops Existing Database** - Removes the `scvis` database if it exists
3. **Creates New Database** - Creates a fresh `scvis` database with UTF8MB4 encoding
4. **Creates Tables** - Creates all tables with the following structure:
   - `users` - User accounts and authentication
   - `projects` - User projects
   - `cfiles` - C source files linked to projects
   - `inputfiles` - Input files linked to projects
   - `issues` - Issue tracking and GitHub integration
   - `sessions` - User session management
5. **Grants Privileges** - Ensures `scvis_user` has proper permissions
6. **Verifies Structure** - Confirms all tables were created successfully

## Database Structure

### users
- User authentication and profile information
- Administrator flag for access control
- Unique email constraint

### projects
- User projects with title and parameters
- JSON arrays for file IDs
- BLOB storage for analysis results
- Foreign key to users

### cfiles
- C source code files
- Linked to projects via foreign key
- Cascade delete when project is removed

### inputfiles
- Input/test data files
- Linked to projects via foreign key
- Cascade delete when project is removed

### issues
- Issue tracking system
- GitHub integration (issue number and URL)
- Status workflow tracking
- Foreign key to users

### sessions
- Secure session token storage
- Automatic expiry timestamps

## Prerequisites

- MySQL/MariaDB installed and running
- Root password for MySQL
- `scvis_user` already created in MySQL

## Environment Variables

- `DB_USER` - Database username (default: `scvis_user`)
- `DB_HOST` - Database host (default: `localhost`)

## Post-Installation Steps

After running this script:

1. **Verify Database**
   ```bash
   mysql -u root -p -e "SHOW TABLES FROM scvis;"
   ```

2. **Create Admin User**
   ```bash
   ./scripts/create_user.sh
   ```

3. **Update Application Configuration**
   - Ensure `.env` file has correct database credentials
   - Verify `DB_NAME=scvis`

4. **Restart Application**
   ```bash
   ./cmd/web/web
   ```

## Troubleshooting

### "Cannot connect to MySQL"
- Verify MySQL is running: `systemctl status mysql`
- Check root password is correct
- Ensure MySQL socket is accessible

### "Access denied for user 'scvis_user'"
- User may not exist, create with:
  ```sql
  CREATE USER 'scvis_user'@'localhost' IDENTIFIED BY 'your_password';
  GRANT SELECT, INSERT, UPDATE, DELETE ON scvis.* TO 'scvis_user'@'localhost';
  ```

### "Foreign key constraint fails"
- This shouldn't happen with a fresh database
- If it does, check table creation order in script

## Backup Recommendation

Before running this script on a production database, **always create a backup**:

```bash
# Backup entire database
mysqldump -u root -p scvis > scvis_backup_$(date +%Y%m%d_%H%M%S).sql

# Backup structure only
mysqldump -u root -p --no-data scvis > scvis_structure_$(date +%Y%m%d_%H%M%S).sql

# Backup data only
mysqldump -u root -p --no-create-info scvis > scvis_data_$(date +%Y%m%d_%H%M%S).sql
```

## Related Scripts

- `create_user.sh` - Create new application users
- `install_database.sh` - Initial database setup
- `validate_database.sh` - Verify database structure

## Version History

- **v1.0** (2025-11-01) - Initial version based on scvis_local structure
  - Added all 6 tables with proper foreign keys
  - Includes issue tracking system
  - Session management support
