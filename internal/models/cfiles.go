package models

import (
	"database/sql"
	"fmt"
	"time"
)

type CFile struct {
	ID        int
	ProjectID int
	Filename  string
	CSource   string
	Created   time.Time
}

type CFileModelInterface interface {
	Insert(projectid int, filename string, csource string) (int, error)
	Get(id int) (CFile, error)
	GetByProjectid(projectid int) ([]CFile, error)
	Delete(id int) error
}

// Define a ProjectModel type which wraps a sql.DB connection pool.
type CFileModel struct {
	DB *sql.DB
}

func (c *CFileModel) Insert(projectid int, filename string, csource string) (int, error) {
	stmt := `INSERT INTO cfiles (projectid, filename, csource)
	VALUES(?, ?, ?)`

	result, err := c.DB.Exec(stmt, projectid, filename, csource)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result to get the ID of our
	// newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

func (c *CFileModel) Get(id int) (CFile, error) {
	stmt := `SELECT id, projectid, filename, csource FROM cfiles WHERE id = ?`

	var cfile CFile
	err := c.DB.QueryRow(stmt, id).Scan(&cfile.ID, &cfile.ProjectID, &cfile.Filename, &cfile.CSource)
	if err != nil {
		if err == sql.ErrNoRows {
			return CFile{}, fmt.Errorf("cfile not found")
		}
		return CFile{}, err
	}

	return cfile, nil
}

func (c *CFileModel) GetByProjectid(projectid int) ([]CFile, error) {
	stmt := `SELECT id, projectid, filename, csource FROM cfiles WHERE projectid = ?`

	rows, err := c.DB.Query(stmt, projectid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cfiles []CFile
	for rows.Next() {
		var cfile CFile
		err := rows.Scan(&cfile.ID, &cfile.ProjectID, &cfile.Filename, &cfile.CSource)
		if err != nil {
			return nil, err
		}
		cfiles = append(cfiles, cfile)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cfiles, nil
}

func (c *CFileModel) Delete(id int) error {
	stmt := `DELETE FROM cfiles WHERE id = ?`

	_, err := c.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
