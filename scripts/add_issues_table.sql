-- Add issues table for GitHub issue submission functionality

CREATE TABLE IF NOT EXISTS issues (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id INTEGER NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    issue_type ENUM('bug', 'feature', 'enhancement', 'documentation', 'question') NOT NULL,
    priority ENUM('low', 'medium', 'high', 'critical') NOT NULL,
    status ENUM('pending', 'submitted', 'acknowledged', 'in_progress', 'resolved', 'closed') NOT NULL DEFAULT 'pending',
    github_issue_number INTEGER NULL,
    github_url VARCHAR(255) NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_created (created),
    INDEX idx_github_issue_number (github_issue_number)
);
