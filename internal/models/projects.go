package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Project struct {
	ID         int
	Title      string
	Userid     int
	CFileIds   []interface{} `json:"cfileids"`
	Analysis   []byte        `json:"analysis"`
	CFiles     []CFile
	Created    time.Time
	Parameters string
}

type ProjectModelInterface interface {
	Insert(title string, userid int) (int, error)
	Get(id int) (Project, error)
	GetByUserId(userid int) ([]Project, error)
	UpdateTitle(id int, title string) error
	UpdateAnalysis(id int, analysis []byte) error
	AddCFiles(id int, cfileids []int) error
	DeleteCFiles(id int, cfileids []int) error
	UpdateParameters(id int, parameters string) error
}

// Define a ProjectModel type which wraps a sql.DB connection pool.
type ProjectModel struct {
	DB *sql.DB
}

func (p *ProjectModel) Insert(title string, userid int) (int, error) {
	stmt := `INSERT INTO projects (title, userid, cfileids, analysis, created_at)
	VALUES(?, ?, ?, ?, UTC_TIMESTAMP())`

	// Define default JSON values for cfileids and analysis
	defaultCFileIds := "[]"
	defaultAnalysis := "[]"

	// Execute the database command with the provided values and default JSON values
	result, err := p.DB.Exec(stmt, title, userid, defaultCFileIds, defaultAnalysis)
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

func (p *ProjectModel) Get(id int) (Project, error) {
	stmt := `SELECT id, title, userid, cfileids, analysis, created_at, parameters FROM projects WHERE id = ?`

	var project Project
	var cfileids []byte
	err := p.DB.QueryRow(stmt, id).Scan(&project.ID, &project.Title, &project.Userid, &cfileids, &project.Analysis, &project.Created, &project.Parameters)
	if err != nil {
		if err == sql.ErrNoRows {
			return Project{}, fmt.Errorf("project not found")
		}
		return Project{}, err
	}
	if err := json.Unmarshal(cfileids, &project.CFileIds); err != nil {
		return Project{}, err
	}

	var CFiles CFileModel = CFileModel{DB: p.DB}
	cfiles, err := CFiles.GetByProjectid(id)

	if err != nil {
		return Project{}, err
	}
	project.CFiles = cfiles
	return project, nil
}

func (p *ProjectModel) GetByUserId(userid int) ([]Project, error) {
	stmt := `SELECT id, title, userid, cfileids, analysis, created_at FROM projects WHERE userid = ?`

	rows, err := p.DB.Query(stmt, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		var cfileids []byte
		err := rows.Scan(&project.ID, &project.Title, &project.Userid, &cfileids, &project.Analysis, &project.Created)

		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(cfileids, &project.CFileIds); err != nil {
			return nil, err
		}

		var CFiles CFileModel = CFileModel{DB: p.DB}
		cfiles, err := CFiles.GetByProjectid(project.ID)

		if err != nil {
			return nil, err
		}
		project.CFiles = cfiles
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *ProjectModel) UpdateParameters(id int, parameters string) error {
	stmt := `UPDATE projects SET parameters = ? WHERE id = ?`
	_, err := p.DB.Exec(stmt, parameters, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectModel) UpdateAnalysis(id int, json_analysis []byte) error {
	analysis, jerr := json.Marshal(json_analysis)
	if jerr != nil {
		return jerr
	}
	stmt := `UPDATE projects SET analysis = ? WHERE id = ?`
	_, err := p.DB.Exec(stmt, analysis, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectModel) UpdateTitle(id int, title string) error {
	stmt := `UPDATE projects SET title = ? WHERE id = ?`
	_, err := p.DB.Exec(stmt, title, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectModel) AddCFiles(id int, cfileids []int) error {
	// Retrieve the existing cfileids
	stmt := `SELECT cfileids FROM projects WHERE id = ?`
	var existingCFileIds []byte
	err := p.DB.QueryRow(stmt, id).Scan(&existingCFileIds)
	if err != nil {
		return err
	}

	var currentCFileIds []int
	if err := json.Unmarshal(existingCFileIds, &currentCFileIds); err != nil {
		return err
	}

	// Append the new cfileids to the existing ones
	currentCFileIds = append(currentCFileIds, cfileids...)

	// Marshal the updated cfileids back to JSON
	updatedCFileIds, err := json.Marshal(currentCFileIds)
	if err != nil {
		return err
	}

	// Update the project with the new cfileids
	stmt = `UPDATE projects SET cfileids = ? WHERE id = ?`
	_, err = p.DB.Exec(stmt, updatedCFileIds, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectModel) DeleteCFiles(id int, cfileids []int) error {
	// Retrieve the existing cfileids
	stmt := `SELECT cfileids FROM projects WHERE id = ?`
	var existingCFileIds []byte
	err := p.DB.QueryRow(stmt, id).Scan(&existingCFileIds)
	if err != nil {
		return err
	}

	var currentCFileIds []int
	if err := json.Unmarshal(existingCFileIds, &currentCFileIds); err != nil {
		return err
	}

	// Create a map for quick lookup of cfileids to delete
	cfileidsToDelete := make(map[int]struct{})
	for _, cfileid := range cfileids {
		cfileidsToDelete[cfileid] = struct{}{}
	}

	// Filter out the cfileids to delete
	var updatedCFileIds []int
	for _, cfileid := range currentCFileIds {
		if _, found := cfileidsToDelete[cfileid]; !found {
			updatedCFileIds = append(updatedCFileIds, cfileid)
		}
	}

	// Marshal the updated cfileids back to JSON
	updatedCFileIdsJSON, err := json.Marshal(updatedCFileIds)
	if err != nil {
		return err
	}

	// Update the project with the new cfileids
	stmt = `UPDATE projects SET cfileids = ? WHERE id = ?`
	_, err = p.DB.Exec(stmt, updatedCFileIdsJSON, id)
	if err != nil {
		return err
	}

	return nil
}
