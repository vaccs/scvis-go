# Database Rebuild Script - Summary

## Created Files

### 1. `/scripts/rebuild_scvis_database.sh`
**Purpose**: Automated script to drop and recreate the `scvis` database based on `scvis_local` structure

**Features**:
- ✅ Interactive confirmation prompt
- ✅ Color-coded output for better visibility
- ✅ Step-by-step progress indicators
- ✅ Error handling with exit on failure
- ✅ MySQL connection testing
- ✅ Automatic privilege granting
- ✅ Structure verification
- ✅ Summary of created tables

**Usage**:
```bash
./scripts/rebuild_scvis_database.sh
```

### 2. `/scripts/README_REBUILD_DATABASE.md`
**Purpose**: Comprehensive documentation for the rebuild script

**Includes**:
- Detailed usage instructions
- Database structure documentation
- Environment variables reference
- Post-installation steps
- Troubleshooting guide
- Backup recommendations

## Database Structure

The script creates 6 tables based on your `scvis_local` database:

| Table | Purpose | Key Features |
|-------|---------|--------------|
| `users` | User accounts | Unique email, admin flag, login tracking |
| `projects` | User projects | JSON file arrays, blob analysis, parameters |
| `cfiles` | C source files | Linked to projects, cascade delete |
| `inputfiles` | Input/test files | Linked to projects, cascade delete |
| `issues` | Issue tracking | GitHub integration, workflow status |
| `sessions` | Session tokens | Secure storage, auto-expiry |

## Script Workflow

```
┌─────────────────────────────────────┐
│  1. Display Warning & Get Consent  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  2. Test MySQL Connection           │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  3. DROP DATABASE IF EXISTS scvis   │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  4. CREATE DATABASE scvis           │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  5. Create All 6 Tables             │
│     - users                         │
│     - projects                      │
│     - cfiles                        │
│     - inputfiles                    │
│     - issues                        │
│     - sessions                      │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  6. Grant Privileges to scvis_user  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  7. Verify & Display Summary        │
└─────────────────────────────────────┘
```

## Safety Features

1. **Confirmation Prompt**: Requires explicit "yes" to proceed
2. **Password Validation**: Tests MySQL connection before making changes
3. **Error Handling**: Exits immediately on any MySQL error
4. **Visual Feedback**: Color-coded messages (red=warning, green=success, yellow=info)
5. **Progress Tracking**: Shows 5-step progress with checkmarks

## Use Cases

### Development
- Reset database to clean state
- Apply schema changes from scvis_local
- Test fresh installation

### Production Deployment
- Deploy to new server
- Update schema to match development
- Recover from corruption

### Testing
- Create isolated test databases
- Verify migration scripts
- Benchmark fresh installation

## Before Running

### Recommended Backup Commands

```bash
# Full backup with timestamp
mysqldump -u root -p scvis > backups/scvis_backup_$(date +%Y%m%d_%H%M%S).sql

# Quick verification
mysql -u root -p -e "SELECT COUNT(*) FROM scvis.users;"
mysql -u root -p -e "SELECT COUNT(*) FROM scvis.projects;"
```

### Environment Check

```bash
# Verify MySQL is running
systemctl status mysql

# Check scvis_user exists
mysql -u root -p -e "SELECT User, Host FROM mysql.user WHERE User='scvis_user';"

# Test scvis_user can connect
mysql -u scvis_user -p -e "SHOW DATABASES;"
```

## After Running

### Verification Steps

```bash
# 1. Check all tables created
mysql -u root -p -e "SHOW TABLES FROM scvis;"

# 2. Verify table structures
mysql -u root -p -e "DESCRIBE scvis.users;"
mysql -u root -p -e "DESCRIBE scvis.projects;"

# 3. Check foreign keys
mysql -u root -p -e "SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME 
FROM information_schema.KEY_COLUMN_USAGE 
WHERE TABLE_SCHEMA='scvis' AND REFERENCED_TABLE_NAME IS NOT NULL;"

# 4. Verify privileges
mysql -u root -p -e "SHOW GRANTS FOR 'scvis_user'@'localhost';"
```

### Create Admin User

```bash
# Use the create_user script
./scripts/create_user.sh

# Or manually create via application
# (login with any credentials after first user creation)
```

### Start Application

```bash
# Build if needed
go build ./cmd/web

# Run application
./cmd/web/web

# Or use systemd if configured
sudo systemctl start scvisual
```

## Key Differences from Original

This script improves upon manual database creation by:

1. **Automation**: One command instead of multiple SQL files
2. **Safety**: Confirmation prompts and connection testing
3. **Completeness**: Includes all 6 tables with proper relationships
4. **Verification**: Automatically checks table creation
5. **Documentation**: Built-in help and post-install instructions
6. **Portability**: Works on any machine with MySQL installed

## Example Output

```
════════════════════════════════════════════════════════════════
  SCVIS Database Rebuild Script
════════════════════════════════════════════════════════════════

WARNING: This script will DROP and RECREATE the 'scvis' database!
         All existing data will be PERMANENTLY DELETED!

Are you sure you want to continue? (yes/no): yes

Starting database rebuild...

Enter MySQL root password: 
[1/5] Testing MySQL connection...
      ✓ Connected successfully

[2/5] Dropping existing 'scvis' database...
      ✓ Database dropped

[3/5] Creating new 'scvis' database...
      ✓ Database created

[4/5] Creating database tables...
      ✓ Tables created successfully

[5/5] Granting privileges to 'scvis_user'...
      ✓ Privileges granted

Verifying database structure...
      ✓ Created 6 tables

════════════════════════════════════════════════════════════════
  Database rebuild completed successfully!
════════════════════════════════════════════════════════════════

Tables created:
+------------+
| Tables_in_scvis |
+------------+
| cfiles     |
| inputfiles |
| issues     |
| projects   |
| sessions   |
| users      |
+------------+

Next steps:
  1. Update your .env file with the correct database credentials
  2. Create an admin user using: ./scripts/create_user.sh
  3. Restart the application
```

## Files Location

```
scvis-go/
├── scripts/
│   ├── rebuild_scvis_database.sh          ← Main rebuild script
│   └── README_REBUILD_DATABASE.md         ← Detailed documentation
└── DATABASE_REBUILD_SUMMARY.md            ← This file
```

## Quick Reference

```bash
# Make executable (if needed)
chmod +x scripts/rebuild_scvis_database.sh

# Run with defaults
./scripts/rebuild_scvis_database.sh

# Run with custom user
DB_USER=myuser ./scripts/rebuild_scvis_database.sh

# View help
cat scripts/README_REBUILD_DATABASE.md
```
