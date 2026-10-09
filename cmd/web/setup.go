package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"scvis.vaccsedu.org/internal/models"
)

var validDBName = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

// ensureDatabase creates the database if it does not exist (used on hosted
// MySQL servers shared by several apps, e.g. Railway)
func ensureDatabase(user, passwd, host, port, name string) error {
	if !validDBName.MatchString(name) {
		return fmt.Errorf("invalid DB_NAME %q", name)
	}
	server, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, passwd, host, port))
	if err != nil {
		return err
	}
	defer server.Close()
	_, err = server.Exec("CREATE DATABASE IF NOT EXISTS `" + name + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci")
	return err
}

// schema matches scripts/rebuild_scvis_database.sh
var schema = []string{
	"CREATE TABLE IF NOT EXISTS `users` (" +
		"`id` int NOT NULL AUTO_INCREMENT," +
		"`name` varchar(255) NOT NULL," +
		"`email` varchar(255) NOT NULL," +
		"`hashed_password` char(60) NOT NULL," +
		"`administrator` tinyint(1) NOT NULL DEFAULT '0'," +
		"`logged_in` tinyint(1) DEFAULT '0'," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `users_uc_email` (`email`))",
	"CREATE TABLE IF NOT EXISTS `projects` (" +
		"`id` int NOT NULL AUTO_INCREMENT," +
		"`title` varchar(255) NOT NULL," +
		"`userid` int NOT NULL," +
		"`cfileids` json NOT NULL," +
		"`analysis` blob," +
		"`created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP," +
		"`parameters` varchar(255) DEFAULT NULL," +
		"`inputfileids` json NOT NULL DEFAULT (json_array())," +
		"PRIMARY KEY (`id`)," +
		"KEY `fk_projects_userid` (`userid`)," +
		"CONSTRAINT `fk_projects_userid` FOREIGN KEY (`userid`) REFERENCES `users` (`id`) ON DELETE RESTRICT)",
	"CREATE TABLE IF NOT EXISTS `cfiles` (" +
		"`id` int NOT NULL AUTO_INCREMENT," +
		"`projectid` int NOT NULL," +
		"`filename` varchar(255) NOT NULL," +
		"`csource` text NOT NULL," +
		"`created` timestamp NULL DEFAULT CURRENT_TIMESTAMP," +
		"PRIMARY KEY (`id`)," +
		"KEY `fk_cfiles_projectid` (`projectid`)," +
		"CONSTRAINT `fk_cfiles_projectid` FOREIGN KEY (`projectid`) REFERENCES `projects` (`id`) ON DELETE CASCADE)",
	"CREATE TABLE IF NOT EXISTS `inputfiles` (" +
		"`id` int NOT NULL AUTO_INCREMENT," +
		"`projectid` int NOT NULL," +
		"`filename` varchar(255) NOT NULL," +
		"`data` text NOT NULL," +
		"`created` timestamp NULL DEFAULT CURRENT_TIMESTAMP," +
		"PRIMARY KEY (`id`)," +
		"KEY `fk_inputfiles_projectid` (`projectid`)," +
		"CONSTRAINT `fk_inputfiles_projectid` FOREIGN KEY (`projectid`) REFERENCES `projects` (`id`) ON DELETE CASCADE)",
	"CREATE TABLE IF NOT EXISTS `issues` (" +
		"`id` int NOT NULL AUTO_INCREMENT," +
		"`user_id` int NOT NULL," +
		"`title` varchar(200) NOT NULL," +
		"`description` text NOT NULL," +
		"`issue_type` enum('bug','feature','enhancement','documentation','question') NOT NULL," +
		"`priority` enum('low','medium','high','critical') NOT NULL," +
		"`status` enum('pending','submitted','acknowledged','in_progress','resolved','closed') NOT NULL DEFAULT 'pending'," +
		"`github_issue_number` int DEFAULT NULL," +
		"`github_url` varchar(255) DEFAULT NULL," +
		"`created` datetime NOT NULL," +
		"`updated` datetime NOT NULL," +
		"PRIMARY KEY (`id`)," +
		"KEY `idx_user_id` (`user_id`)," +
		"KEY `idx_status` (`status`)," +
		"KEY `idx_created` (`created`)," +
		"KEY `idx_github_issue_number` (`github_issue_number`)," +
		"CONSTRAINT `issues_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE)",
	"CREATE TABLE IF NOT EXISTS `sessions` (" +
		"`token` char(43) NOT NULL," +
		"`data` blob NOT NULL," +
		"`expiry` timestamp(6) NOT NULL," +
		"PRIMARY KEY (`token`)," +
		"KEY `sessions_expiry_idx` (`expiry`))",
}

// ensureSchema creates any missing tables
func ensureSchema(db *sql.DB) error {
	for _, stmt := range schema {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

// seedAdmin creates an initial administrator from ADMIN_* env vars when the
// users table is empty
func seedAdmin(db *sql.DB, users *models.UserModel, logger *slog.Logger) error {
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	if email == "" || password == "" {
		return nil
	}
	name := os.Getenv("ADMIN_USERNAME")
	if name == "" {
		name = "admin"
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	if err := users.Insert(name, email, password, true); err != nil {
		return err
	}
	logger.Info("created initial admin user", "name", name, "email", email)
	return nil
}
