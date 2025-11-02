#!/bin/bash

################################################################################
# Script: rebuild_scvis_database.sh
# Description: Drops and rebuilds the scvis database based on scvis_local structure
# Usage: ./rebuild_scvis_database.sh
# Note: This script will DELETE all data in the scvis database!
################################################################################

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
DB_NAME="scvis"
DB_USER="${DB_USER:-scvis}"
DB_HOST="${DB_HOST:-localhost}"

echo -e "${YELLOW}════════════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}  SCVIS Database Rebuild Script${NC}"
echo -e "${YELLOW}════════════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${RED}WARNING: This script will DROP and RECREATE the '$DB_NAME' database!${NC}"
echo -e "${RED}         All existing data will be PERMANENTLY DELETED!${NC}"
echo ""
read -p "Are you sure you want to continue? (yes/no): " -r
echo ""

if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo -e "${YELLOW}Operation cancelled.${NC}"
    exit 0
fi

echo -e "${GREEN}Starting database rebuild...${NC}"
echo ""

# Prompt for MySQL root password
read -sp "Enter MySQL root password: " MYSQL_ROOT_PASSWORD
echo ""

# Test MySQL connection
echo -e "${GREEN}[1/5] Testing MySQL connection...${NC}"
if ! mysql -u root -p"${MYSQL_ROOT_PASSWORD}" -e "SELECT 1;" &>/dev/null; then
    echo -e "${RED}ERROR: Cannot connect to MySQL. Please check your password.${NC}"
    exit 1
fi
echo -e "${GREEN}      ✓ Connected successfully${NC}"
echo ""

# Drop existing scvis database
echo -e "${GREEN}[2/5] Dropping existing '$DB_NAME' database...${NC}"
mysql -u root -p"${MYSQL_ROOT_PASSWORD}" <<EOF
DROP DATABASE IF EXISTS $DB_NAME;
EOF
echo -e "${GREEN}      ✓ Database dropped${NC}"
echo ""

# Create new scvis database
echo -e "${GREEN}[3/5] Creating new '$DB_NAME' database...${NC}"
mysql -u root -p"${MYSQL_ROOT_PASSWORD}" <<EOF
CREATE DATABASE $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
EOF
echo -e "${GREEN}      ✓ Database created${NC}"
echo ""

# Create database structure
echo -e "${GREEN}[4/5] Creating database tables...${NC}"
mysql -u root -p"${MYSQL_ROOT_PASSWORD}" "$DB_NAME" <<'EOF'
-- Users table
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `hashed_password` char(60) NOT NULL,
  `administrator` tinyint(1) NOT NULL DEFAULT '0',
  `logged_in` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_uc_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Projects table
CREATE TABLE `projects` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `userid` int NOT NULL,
  `cfileids` json NOT NULL,
  `analysis` blob,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `parameters` varchar(255) DEFAULT NULL,
  `inputfileids` json NOT NULL DEFAULT (json_array()),
  PRIMARY KEY (`id`),
  KEY `fk_projects_userid` (`userid`),
  CONSTRAINT `fk_projects_userid` FOREIGN KEY (`userid`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- C Files table
CREATE TABLE `cfiles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `projectid` int NOT NULL,
  `filename` varchar(255) NOT NULL,
  `csource` text NOT NULL,
  `created` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_cfiles_projectid` (`projectid`),
  CONSTRAINT `fk_cfiles_projectid` FOREIGN KEY (`projectid`) REFERENCES `projects` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Input Files table
CREATE TABLE `inputfiles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `projectid` int NOT NULL,
  `filename` varchar(255) NOT NULL,
  `data` text NOT NULL,
  `created` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_inputfiles_projectid` (`projectid`),
  CONSTRAINT `fk_inputfiles_projectid` FOREIGN KEY (`projectid`) REFERENCES `projects` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Issues table
CREATE TABLE `issues` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `title` varchar(200) NOT NULL,
  `description` text NOT NULL,
  `issue_type` enum('bug','feature','enhancement','documentation','question') NOT NULL,
  `priority` enum('low','medium','high','critical') NOT NULL,
  `status` enum('pending','submitted','acknowledged','in_progress','resolved','closed') NOT NULL DEFAULT 'pending',
  `github_issue_number` int DEFAULT NULL,
  `github_url` varchar(255) DEFAULT NULL,
  `created` datetime NOT NULL,
  `updated` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created` (`created`),
  KEY `idx_github_issue_number` (`github_issue_number`),
  CONSTRAINT `issues_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Sessions table
CREATE TABLE `sessions` (
  `token` char(43) NOT NULL,
  `data` blob NOT NULL,
  `expiry` timestamp(6) NOT NULL,
  PRIMARY KEY (`token`),
  KEY `sessions_expiry_idx` (`expiry`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
EOF

echo -e "${GREEN}      ✓ Tables created successfully${NC}"
echo ""

# Grant privileges to scvis
echo -e "${GREEN}[5/5] Setting up database user '$DB_USER'...${NC}"

# Check if user exists
USER_EXISTS=$(mysql -u root -p"${MYSQL_ROOT_PASSWORD}" -N -e "SELECT COUNT(*) FROM mysql.user WHERE user='$DB_USER' AND host='localhost';" 2>/dev/null || echo "0")

if [ "$USER_EXISTS" = "0" ]; then
    echo -e "${YELLOW}      User '$DB_USER' does not exist. Please create it first.${NC}"
    echo -e "${YELLOW}      Run the following commands:${NC}"
    echo ""
    echo "      mysql -u root -p"
    echo "      CREATE USER '$DB_USER'@'localhost' IDENTIFIED BY 'your_password';"
    echo "      GRANT SELECT, INSERT, UPDATE, DELETE ON $DB_NAME.* TO '$DB_USER'@'localhost';"
    echo "      FLUSH PRIVILEGES;"
    echo ""
    echo -e "${YELLOW}      Or run: ./scripts/install_database.sh to set up the user${NC}"
else
    # User exists, just grant privileges
    mysql -u root -p"${MYSQL_ROOT_PASSWORD}" <<EOF
GRANT SELECT, INSERT, UPDATE, DELETE ON $DB_NAME.* TO '$DB_USER'@'localhost';
FLUSH PRIVILEGES;
EOF
    echo -e "${GREEN}      ✓ Privileges granted to existing user${NC}"
fi
echo ""

# Verify the structure
echo -e "${GREEN}Verifying database structure...${NC}"
TABLE_COUNT=$(mysql -u root -p"${MYSQL_ROOT_PASSWORD}" -N -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = '$DB_NAME';")
echo -e "${GREEN}      ✓ Created $TABLE_COUNT tables${NC}"
echo ""

echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  Database rebuild completed successfully!${NC}"
echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${YELLOW}Tables created:${NC}"
mysql -u root -p"${MYSQL_ROOT_PASSWORD}" -e "SHOW TABLES FROM $DB_NAME;"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Update your .env file with the correct database credentials"
echo "  2. Create an admin user using: ./scripts/create_user.sh"
echo "  3. Restart the application"
echo ""
