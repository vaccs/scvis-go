package models

import (
	"database/sql"
	"time"
)

type IssueModelInterface interface {
	Insert(userID int, title, description, issueType, priority string, githubIssueNumber *int, githubURL *string) (int, error)
	Get(id int) (*Issue, error)
	GetAll() ([]*Issue, error)
	GetByUser(userID int) ([]*Issue, error)
	UpdateGitHubInfo(id int, githubIssueNumber int, githubURL string) error
	UpdateStatus(id int, status string) error
}

type Issue struct {
	ID                int
	UserID            int
	Username          string
	Title             string
	Description       string
	IssueType         string // "bug", "feature", "enhancement", "documentation", "question"
	Priority          string // "low", "medium", "high", "critical"
	Status            string // "pending", "submitted", "acknowledged", "in_progress", "resolved", "closed"
	GitHubIssueNumber *int
	GitHubURL         *string
	Created           time.Time
	Updated           time.Time
}

type IssueModel struct {
	DB *sql.DB
}

func (m *IssueModel) Insert(userID int, title, description, issueType, priority string, githubIssueNumber *int, githubURL *string) (int, error) {
	stmt := `INSERT INTO issues (user_id, title, description, issue_type, priority, status, github_issue_number, github_url, created, updated)
    VALUES(?, ?, ?, ?, ?, 'pending', ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, userID, title, description, issueType, priority, githubIssueNumber, githubURL)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *IssueModel) Get(id int) (*Issue, error) {
	stmt := `SELECT i.id, i.user_id, u.name as username, i.title, i.description, i.issue_type, i.priority, i.status, 
             i.github_issue_number, i.github_url, i.created, i.updated
             FROM issues i
             JOIN users u ON i.user_id = u.id
             WHERE i.id = ?`

	row := m.DB.QueryRow(stmt, id)

	issue := &Issue{}
	err := row.Scan(
		&issue.ID,
		&issue.UserID,
		&issue.Username,
		&issue.Title,
		&issue.Description,
		&issue.IssueType,
		&issue.Priority,
		&issue.Status,
		&issue.GitHubIssueNumber,
		&issue.GitHubURL,
		&issue.Created,
		&issue.Updated,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return issue, nil
}

func (m *IssueModel) GetAll() ([]*Issue, error) {
	stmt := `SELECT i.id, i.user_id, u.name as username, i.title, i.description, i.issue_type, i.priority, i.status,
             i.github_issue_number, i.github_url, i.created, i.updated
             FROM issues i
             JOIN users u ON i.user_id = u.id
             ORDER BY i.created DESC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []*Issue

	for rows.Next() {
		issue := &Issue{}
		err = rows.Scan(
			&issue.ID,
			&issue.UserID,
			&issue.Username,
			&issue.Title,
			&issue.Description,
			&issue.IssueType,
			&issue.Priority,
			&issue.Status,
			&issue.GitHubIssueNumber,
			&issue.GitHubURL,
			&issue.Created,
			&issue.Updated,
		)
		if err != nil {
			return nil, err
		}

		issues = append(issues, issue)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return issues, nil
}

func (m *IssueModel) GetByUser(userID int) ([]*Issue, error) {
	stmt := `SELECT i.id, i.user_id, u.name as username, i.title, i.description, i.issue_type, i.priority, i.status,
             i.github_issue_number, i.github_url, i.created, i.updated
             FROM issues i
             JOIN users u ON i.user_id = u.id
             WHERE i.user_id = ?
             ORDER BY i.created DESC`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []*Issue

	for rows.Next() {
		issue := &Issue{}
		err = rows.Scan(
			&issue.ID,
			&issue.UserID,
			&issue.Username,
			&issue.Title,
			&issue.Description,
			&issue.IssueType,
			&issue.Priority,
			&issue.Status,
			&issue.GitHubIssueNumber,
			&issue.GitHubURL,
			&issue.Created,
			&issue.Updated,
		)
		if err != nil {
			return nil, err
		}

		issues = append(issues, issue)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return issues, nil
}

func (m *IssueModel) UpdateGitHubInfo(id int, githubIssueNumber int, githubURL string) error {
	stmt := `UPDATE issues SET github_issue_number = ?, github_url = ?, status = 'submitted', updated = UTC_TIMESTAMP() 
             WHERE id = ?`

	_, err := m.DB.Exec(stmt, githubIssueNumber, githubURL, id)
	return err
}

func (m *IssueModel) UpdateStatus(id int, status string) error {
	stmt := `UPDATE issues SET status = ?, updated = UTC_TIMESTAMP() WHERE id = ?`
	_, err := m.DB.Exec(stmt, status, id)
	return err
}
