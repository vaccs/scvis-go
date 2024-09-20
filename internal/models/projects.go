package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Project struct {
	ID       int
	Title    string
	Userid   int
	CFileIds []interface{}            `json:"cfileids"`
	Analysis []map[string]interface{} `json:"analysis"`
	CFiles   []CFile
	Created  time.Time
}

type ProjectModelInterface interface {
	Insert(title string, userid int) (int, error)
	Get(id int) (Project, error)
	GetByUserId(userid int) ([]Project, error)
	UpdateAnalysis(id int, analysis []map[string]interface{}) error
	AddSingleCFile(id int, cfileid int) error
	Latest() ([]Project, error)
}

// Define a ProjectModel type which wraps a sql.DB connection pool.
type ProjectModel struct {
	DB *sql.DB
}

func (p *ProjectModel) Insert(title string, userid int) (int, error) {
	stmt := `INSERT INTO projects (title, userid, cfileids, analysis)
	VALUES(?, ?, ?, ?, UTC_TIMESTAMP())`

	result, err := p.DB.Exec(stmt, title, userid, []int{}, []map[string]interface{}{})
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
	stmt := `SELECT id, title, userid, cfileids, analysis, created_at FROM projects WHERE id = ?`

	var project Project
	var cfileids []byte
	var analysis []byte
	err := p.DB.QueryRow(stmt, id).Scan(&project.ID, &project.Title, &project.Userid, &cfileids, &analysis, &project.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			return Project{}, fmt.Errorf("project not found")
		}
		return Project{}, err
	}
	if err := json.Unmarshal(cfileids, &project.CFileIds); err != nil {
		return Project{}, err
	}
	if err := json.Unmarshal(analysis, &project.Analysis); err != nil {
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
		var analysis []byte
		err := rows.Scan(&project.ID, &project.Title, &project.Userid, &cfileids, &analysis, &project.Created)

		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(cfileids, &project.CFileIds); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(analysis, &project.Analysis); err != nil {
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

func (p *ProjectModel) UpdateAnalysis(id int, json_analysis []map[string]interface{}) error {
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

func (p *ProjectModel) AddSingleCFile(id int, cfileid int) error {
	cfileidJSON, jerr := json.Marshal(cfileid)
	if jerr != nil {
		return jerr
	}
	stmt := `UPDATE projects SET cfileids = JSON_ARRAY_APPEND(cfileids, '$', ?) WHERE id = ?`
	_, err := p.DB.Exec(stmt, cfileidJSON, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectModel) Latest() ([]Project, error) {
	stmt := `SELECT id, title, userid, cfileids, analysis, created_at FROM projects ORDER BY id DESC LIMIT 5`

	rows, err := p.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		var cfileids []byte
		var analysis []byte
		err := rows.Scan(&project.ID, &project.Title, &project.Userid, &cfileids, &analysis, &project.Created)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(cfileids, &project.CFileIds); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(analysis, &project.Analysis); err != nil {
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
