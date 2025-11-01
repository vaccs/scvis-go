package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Project struct {
	ID           int
	Title        string
	Userid       int
	CFileIds     []interface{} `json:"cfileids"`
	InputFileIds []interface{} `json:"inputfileids"`
	Analysis     []byte        `json:"analysis"`
	CFiles       []CFile
	InputFiles   []InputFile
	Created      time.Time
	Parameters   string
}

type ProjectModelInterface interface {
	Insert(title string, userid int) (int, error)
	Get(id int) (Project, error)
	GetByUserId(userid int) ([]Project, error)
	UpdateTitle(id int, title string) error
	UpdateAnalysis(id int, analysis []byte) error
	AddCFiles(id int, cfileids []int) error
	AddInputFiles(id int, inputfileids []int) error
	Delete(id int) error
	DeleteAnalysis(id int) error
	DeleteCFiles(id int, cfileids []int) error
	DeleteInputFiles(id int, inputfileids []int) error
	UpdateParameters(id int, parameters string) error
}

// Define a ProjectModel type which wraps a sql.DB connection pool.
type ProjectModel struct {
	DB *sql.DB
}

func (p *ProjectModel) Insert(title string, userid int) (int, error) {
	stmt := `INSERT INTO projects (title, userid, cfileids, inputfileids, analysis, created_at, parameters)
	VALUES(?, ?, ?, ?, ?, UTC_TIMESTAMP(), ?)`

	// Define default JSON values for cfileids, inputfileids and analysis
	defaultCFileIds := "[]"
	defaultInputFileIds := "[]"
	defaultAnalysis := "[]"
	defaultParameters := ""

	// Execute the database command with the provided values and default JSON values
	result, err := p.DB.Exec(stmt, title, userid, defaultCFileIds, defaultInputFileIds, defaultAnalysis, defaultParameters)
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
	stmt := `SELECT id, title, userid, cfileids, inputfileids, analysis, created_at, parameters FROM projects WHERE id = ?`

	var project Project
	var cfileids []byte
	var inputfileids []byte
	err := p.DB.QueryRow(stmt, id).Scan(&project.ID, &project.Title, &project.Userid, &cfileids, &inputfileids, &project.Analysis, &project.Created, &project.Parameters)
	if err != nil {
		if err == sql.ErrNoRows {
			return Project{}, fmt.Errorf("project not found")
		}
		return Project{}, err
	}
	if err := json.Unmarshal(cfileids, &project.CFileIds); err != nil {
		return Project{}, err
	}
	if err := json.Unmarshal(inputfileids, &project.InputFileIds); err != nil {
		return Project{}, err
	}

	var CFiles CFileModel = CFileModel{DB: p.DB}
	cfiles, err := CFiles.GetByProjectid(id)

	if err != nil {
		return Project{}, err
	}
	project.CFiles = cfiles

	var InputFiles InputFileModel = InputFileModel{DB: p.DB}
	inputfiles, err := InputFiles.GetByProjectid(id)

	if err != nil {
		return Project{}, err
	}
	project.InputFiles = inputfiles

	return project, nil
}

func (p *ProjectModel) GetByUserId(userid int) ([]Project, error) {
	stmt := `SELECT id, title, userid, cfileids, inputfileids, analysis, created_at, parameters FROM projects WHERE userid = ?`

	rows, err := p.DB.Query(stmt, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		var cfileids []byte
		var inputfileids []byte
		err := rows.Scan(&project.ID, &project.Title, &project.Userid, &cfileids, &inputfileids, &project.Analysis, &project.Created, &project.Parameters)

		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(cfileids, &project.CFileIds); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(inputfileids, &project.InputFileIds); err != nil {
			return nil, err
		}

		var CFiles CFileModel = CFileModel{DB: p.DB}
		cfiles, err := CFiles.GetByProjectid(project.ID)

		if err != nil {
			return nil, err
		}
		project.CFiles = cfiles

		var InputFiles InputFileModel = InputFileModel{DB: p.DB}
		inputfiles, err := InputFiles.GetByProjectid(project.ID)

		if err != nil {
			return nil, err
		}
		project.InputFiles = inputfiles

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
	stmt := `UPDATE projects SET analysis = ? WHERE id = ?`
	_, err := p.DB.Exec(stmt, json_analysis, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectModel) DeleteAnalysis(id int) error {
	stmt := `UPDATE projects SET analysis = ? WHERE id = ?`
	_, err := p.DB.Exec(stmt, "[]", id)
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

func (p *ProjectModel) Delete(id int) error {
	stmt := `DELETE FROM projects WHERE id = ?`
	_, err := p.DB.Exec(stmt, id)
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

func (p *ProjectModel) AddInputFiles(id int, inputfileids []int) error {
	// Retrieve the existing inputfileids
	stmt := `SELECT inputfileids FROM projects WHERE id = ?`
	var existingInputFileIds []byte
	err := p.DB.QueryRow(stmt, id).Scan(&existingInputFileIds)
	if err != nil {
		return err
	}

	var currentInputFileIds []int
	if err := json.Unmarshal(existingInputFileIds, &currentInputFileIds); err != nil {
		return err
	}

	// Append the new inputfileids to the existing ones
	currentInputFileIds = append(currentInputFileIds, inputfileids...)

	// Marshal the updated inputfileids back to JSON
	updatedInputFileIds, err := json.Marshal(currentInputFileIds)
	if err != nil {
		return err
	}

	// Update the project with the new inputfileids
	stmt = `UPDATE projects SET inputfileids = ? WHERE id = ?`
	_, err = p.DB.Exec(stmt, updatedInputFileIds, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectModel) DeleteInputFiles(id int, inputfileids []int) error {
	// Retrieve the existing inputfileids
	stmt := `SELECT inputfileids FROM projects WHERE id = ?`
	var existingInputFileIds []byte
	err := p.DB.QueryRow(stmt, id).Scan(&existingInputFileIds)
	if err != nil {
		return err
	}

	var currentInputFileIds []int
	if err := json.Unmarshal(existingInputFileIds, &currentInputFileIds); err != nil {
		return err
	}

	// Create a map for quick lookup of inputfileids to delete
	inputfileidsToDelete := make(map[int]struct{})
	for _, inputfileid := range inputfileids {
		inputfileidsToDelete[inputfileid] = struct{}{}
	}

	// Filter out the inputfileids to delete
	var updatedInputFileIds []int
	for _, inputfileid := range currentInputFileIds {
		if _, found := inputfileidsToDelete[inputfileid]; !found {
			updatedInputFileIds = append(updatedInputFileIds, inputfileid)
		}
	}

	// Marshal the updated inputfileids back to JSON
	updatedInputFileIdsJSON, err := json.Marshal(updatedInputFileIds)
	if err != nil {
		return err
	}

	// Update the project with the new inputfileids
	stmt = `UPDATE projects SET inputfileids = ? WHERE id = ?`
	_, err = p.DB.Exec(stmt, updatedInputFileIdsJSON, id)
	if err != nil {
		return err
	}

	return nil
}
