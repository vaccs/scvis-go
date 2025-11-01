package models

import (
	"database/sql"
	"fmt"
	"time"
)

type InputFile struct {
	ID        int
	ProjectID int
	Filename  string
	Data      string
	Created   time.Time
}

type InputFileModelInterface interface {
	Insert(projectid int, filename string, data string) (int, error)
	Get(id int) (InputFile, error)
	GetByProjectid(projectid int) ([]InputFile, error)
	Delete(id int) error
}

// Define a InputFileModel type which wraps a sql.DB connection pool.
type InputFileModel struct {
	DB *sql.DB
}

func (i *InputFileModel) Insert(projectid int, filename string, data string) (int, error) {
	stmt := `INSERT INTO inputfiles (projectid, filename, data)
	VALUES(?, ?, ?)`

	result, err := i.DB.Exec(stmt, projectid, filename, data)
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

func (i *InputFileModel) Get(id int) (InputFile, error) {
	stmt := `SELECT id, projectid, filename, data FROM inputfiles WHERE id = ?`

	var inputfile InputFile
	err := i.DB.QueryRow(stmt, id).Scan(&inputfile.ID, &inputfile.ProjectID, &inputfile.Filename, &inputfile.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return InputFile{}, fmt.Errorf("input file not found")
		}
		return InputFile{}, err
	}

	return inputfile, nil
}

func (i *InputFileModel) GetByProjectid(projectid int) ([]InputFile, error) {
	stmt := `SELECT id, projectid, filename, data FROM inputfiles WHERE projectid = ?`

	rows, err := i.DB.Query(stmt, projectid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inputfiles []InputFile
	for rows.Next() {
		var inputfile InputFile
		err := rows.Scan(&inputfile.ID, &inputfile.ProjectID, &inputfile.Filename, &inputfile.Data)
		if err != nil {
			return nil, err
		}
		inputfiles = append(inputfiles, inputfile)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return inputfiles, nil
}

func (i *InputFileModel) Delete(id int) error {
	stmt := `DELETE FROM inputfiles WHERE id = ?`

	_, err := i.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
