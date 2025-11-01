package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"scvis.vaccsedu.org/internal/github"
	"scvis.vaccsedu.org/internal/models"
	"scvis.vaccsedu.org/internal/vaccs"
	"scvis.vaccsedu.org/internal/validator"
)

// setProjectSession sets the current project ID in the session and redirects to target
func (app *application) setProjectSession(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	target := r.FormValue("target")
	if target == "" {
		target = "/scvis/project/view"
	}

	app.sessionManager.Put(r.Context(), "currentProjectID", id)
	http.Redirect(w, r, target, http.StatusSeeOther)
}

// setCFileSession sets the current C file ID in the session and redirects to view
func (app *application) setCFileSession(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	app.sessionManager.Put(r.Context(), "currentCFileID", id)
	http.Redirect(w, r, "/scvis/cfile/view", http.StatusSeeOther)
}

// setInputFileSession sets the current input file ID in the session and redirects to view
func (app *application) setInputFileSession(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	app.sessionManager.Put(r.Context(), "currentInputFileID", id)
	http.Redirect(w, r, "/scvis/input/view", http.StatusSeeOther)
}

// setIssueSession sets the current issue ID in the session and redirects to view
func (app *application) setIssueSession(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	app.sessionManager.Put(r.Context(), "currentIssueID", id)
	http.Redirect(w, r, "/scvis/issue/view", http.StatusSeeOther)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	projects, err := app.projects.GetByUserId(userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Projects = projects

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) projectVisual(w http.ResponseWriter, r *http.Request) {
	// Get project ID from session
	id := app.sessionManager.GetInt(r.Context(), "currentProjectID")
	if id < 1 {
		http.NotFound(w, r)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "projectID",
		Value: strconv.Itoa(id),
		Path:  "/scvis/",
	})

	app.logger.Info("Serving Flutter web application")
	http.Redirect(w, r, "/scvis/visual", http.StatusSeeOther)
}

func (app *application) projectView(w http.ResponseWriter, r *http.Request) {
	// Get project ID from session
	id := app.sessionManager.GetInt(r.Context(), "currentProjectID")
	if id < 1 {
		http.NotFound(w, r)
		return
	}

	// Handle POST request to update project
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		// Update title if changed
		title := r.PostForm.Get("title")
		if title != "" {
			err = app.projects.UpdateTitle(id, title)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		}

		// Handle file deletions
		cfileIds := r.PostForm["cfilesids"]
		if len(cfileIds) > 0 {
			var cfileIdsInt []int
			for _, idStr := range cfileIds {
				fileId, err := strconv.Atoi(idStr)
				if err == nil {
					cfileIdsInt = append(cfileIdsInt, fileId)
				}
			}
			if len(cfileIdsInt) > 0 {
				err = app.projects.DeleteCFiles(id, cfileIdsInt)
				if err != nil {
					app.serverError(w, r, err)
					return
				}
			}
		}

		inputFileIds := r.PostForm["inputfilesids"]
		if len(inputFileIds) > 0 {
			var inputFileIdsInt []int
			for _, idStr := range inputFileIds {
				fileId, err := strconv.Atoi(idStr)
				if err == nil {
					inputFileIdsInt = append(inputFileIdsInt, fileId)
				}
			}
			if len(inputFileIdsInt) > 0 {
				err = app.projects.DeleteInputFiles(id, inputFileIdsInt)
				if err != nil {
					app.serverError(w, r, err)
					return
				}
			}
		}

		app.sessionManager.Put(r.Context(), "flash", "Project updated successfully!")
		http.Redirect(w, r, "/scvis/project/view", http.StatusSeeOther)
		return
	}

	// GET request - display the project view
	project, err := app.projects.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Project = project

	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *application) projectUpdate(w http.ResponseWriter, r *http.Request) {
	// Redirect to the new view page
	http.Redirect(w, r, "/scvis/project/view", http.StatusSeeOther)
}

func (app *application) projectUpdatePost(w http.ResponseWriter, r *http.Request) {
	// Redirect POST requests to view page which now handles updates
	http.Redirect(w, r, "/scvis/project/view", http.StatusSeeOther)
}

type projectCreateForm struct {
	Title               string `form:"title"`
	validator.Validator `form:"-"`
}

func (app *application) projectCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = projectCreateForm{}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) projectCreatePost(w http.ResponseWriter, r *http.Request) {
	var form projectCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	_, err = app.projects.Insert(form.Title, app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Project successfully created!")

	http.Redirect(w, r, "/scvis/", http.StatusSeeOther)
}

type userSignupForm struct {
	Username            string `form:"username"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	Administrator       bool   `form:"administrator"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Username), "username", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Username, form.Email, form.Password, form.Administrator)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "User created successfully.")

	http.Redirect(w, r, "/scvis/", http.StatusSeeOther)
}

type userLoginForm struct {
	EmailOrUsername     string `form:"emailOrUsername"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.EmailOrUsername), "emailOrUsername", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.EmailOrUsername, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email/Username or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/scvis/home", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/scvis/", http.StatusSeeOther)
}

func (app *application) projectAnalyze(w http.ResponseWriter, r *http.Request) {
	// Get project ID from session
	id := app.sessionManager.GetInt(r.Context(), "currentProjectID")
	if id < 1 {
		http.NotFound(w, r)
		return
	}

	project, err := app.projects.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Project = project

	app.render(w, r, http.StatusOK, "analyze.tmpl", data)
}

func (app *application) projectAnalyzePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get project ID from session
	id := app.sessionManager.GetInt(r.Context(), "currentProjectID")
	if id < 1 {
		http.NotFound(w, r)
		return
	}

	stdinFile := r.FormValue("stdin-file")

	paramArray := r.Form["command-params"]
	project, err := app.projects.Get(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	paramStr := strings.Join(paramArray, " ")
	if len(paramStr) > 255 {
		app.sessionManager.Put(r.Context(), "flash", "Command line parameters exceed the limit of 255 charachters")
		app.logger.Error("Error performing analysis: command line parameters more than 255 charachters")
	} else {

		cfiles, err := app.cfiles.GetByProjectid(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		var fileNames []string
		var cSources []string
		for _, cfile := range cfiles {
			fileNames = append(fileNames, cfile.Filename)
			cSources = append(cSources, cfile.CSource)
		}

		vargs := []string{"-u", "-r", "-j"}

		vaccs_comm := vaccs.NewVaccsComm(app.logger)

		err = vaccs_comm.Connect()
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		defer vaccs_comm.Close()

		safeTitle := strings.Map(func(r rune) rune {
			if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
				return '0'
			}
			return r
		}, project.Title)

		err = app.projects.UpdateParameters(id, paramStr)

		if err != nil {
			app.serverError(w, r, err)
			return
		}

		inputfiles, err := app.inputfiles.GetByProjectid(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		var inputFileNames []string
		var inputSources []string
		for _, inputfile := range inputfiles {
			inputFileNames = append(inputFileNames, inputfile.Filename)
			inputSources = append(inputSources, inputfile.Data)
		}
		data, err := vaccs_comm.PerformAnalysis(fileNames, cSources, inputFileNames, inputSources, vargs, safeTitle, paramArray, stdinFile)
		vaccs_comm.Close()

		if err != nil {
			app.sessionManager.Put(r.Context(), "flash", "Analysis was unsuccessful!")
			app.logger.Error(fmt.Sprintf("Error performing analysis: %s", err.Error()))
			err = app.projects.DeleteAnalysis(id)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		} else {

			err = app.projects.UpdateAnalysis(id, data)

			if err != nil {
				app.serverError(w, r, err)
				return
			}

			app.sessionManager.Put(r.Context(), "flash", "Analysis was successful!")
		}
	}

	http.Redirect(w, r, "/scvis/project/analyze", http.StatusSeeOther)

}

func (app *application) projectGetAnalysis(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("projectID")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.NotFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}

	id, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Handle single analysis (unauthenticated) case with projectID = -1
	if id == -1 {
		tempDir := app.sessionManager.GetString(r.Context(), "tempDir")
		analysisFile := app.sessionManager.GetString(r.Context(), "analysisFile")
		if tempDir == "" || analysisFile == "" {
			http.NotFound(w, r)
			return
		}

		jsonPath := filepath.Join(tempDir, analysisFile)
		jsonBytes, err := os.ReadFile(jsonPath)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)

		// Clean up session data
		app.sessionManager.Remove(r.Context(), "tempDir")
		app.sessionManager.Remove(r.Context(), "analysisFile")
		app.sessionManager.Remove(r.Context(), "projectID")

		// Clean up temp directory in background
		go func() {
			os.RemoveAll(tempDir)
		}()
		return

	} else {
		// Handle authenticated project analysis case
		// Check if user is authenticated
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/scvis/user/login", http.StatusSeeOther)
			return
		}

		project, err := app.projects.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				http.NotFound(w, r)
			} else {
				app.serverError(w, r, err)
			}
			return
		}

		jsonData, err := json.Marshal(string(project.Analysis))
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)

	}
}

func (app *application) singleGetAnalysis(w http.ResponseWriter, r *http.Request) {
	// Check for projectID cookie to determine if this is a single analysis
	cookie, err := r.Cookie("projectID")
	if err != nil || cookie.Value != "-1" {
		// Not a single analysis request, return 404
		http.NotFound(w, r)
		return
	}

	// Get temp directory and analysis file from session
	tempDir := app.sessionManager.GetString(r.Context(), "tempDir")
	analysisFile := app.sessionManager.GetString(r.Context(), "analysisFile")

	if tempDir == "" || analysisFile == "" {
		http.NotFound(w, r)
		return
	}

	// Read the analysis JSON file
	jsonPath := filepath.Join(tempDir, analysisFile)
	jsonBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Return the JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)

	// Clean up session data and temp directory after serving
	app.sessionManager.Remove(r.Context(), "tempDir")
	app.sessionManager.Remove(r.Context(), "analysisFile")
	app.sessionManager.Remove(r.Context(), "projectID")

	// Clean up temp directory in background
	go func() {
		os.RemoveAll(tempDir)
	}()
}

func (app *application) projectDelete(w http.ResponseWriter, r *http.Request) {
	// Get project ID from session
	id := app.sessionManager.GetInt(r.Context(), "currentProjectID")
	if id < 1 {
		http.NotFound(w, r)
		return
	}

	err := app.projects.Delete(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Project successfully deleted!")

	http.Redirect(w, r, "/scvis/home", http.StatusSeeOther)
}

func (app *application) cfileView(w http.ResponseWriter, r *http.Request) {
	// Get file ID from session
	id := app.sessionManager.GetInt(r.Context(), "currentCFileID")
	if id < 1 {
		http.NotFound(w, r)
		return
	}

	cfile, err := app.cfiles.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.CFile = cfile

	app.render(w, r, http.StatusOK, "cfile.tmpl", data)
}

func (app *application) cfileViewPost(w http.ResponseWriter, r *http.Request) {
	// Get file ID from session
	id := app.sessionManager.GetInt(r.Context(), "currentCFileID")
	if id < 1 {
		http.NotFound(w, r)
		return
	}

	cfile, err := app.cfiles.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	project, err := app.projects.Get(cfile.ProjectID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Set current project ID in session for analyze page
	app.sessionManager.Put(r.Context(), "currentProjectID", project.ID)

	http.Redirect(w, r, "/scvis/project/analyze", http.StatusSeeOther)
}

type userChangePasswordForm struct {
	CurrentPassword     string `form:"currentPassword"`
	NewPassword         string `form:"newPassword"`
	ConfirmPassword     string `form:"confirmPassword"`
	validator.Validator `form:"-"`
}

func (app *application) userChangePassword(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userChangePasswordForm{}
	app.render(w, r, http.StatusOK, "change_password.tmpl", data)
}

func (app *application) userChangePasswordPost(w http.ResponseWriter, r *http.Request) {
	var form userChangePasswordForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.ConfirmPassword), "confirmPassword", "This field cannot be blank")
	form.CheckField(form.NewPassword == form.ConfirmPassword, "confirmPassword", "Passwords do not match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "change_password.tmpl", data)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	err = app.users.ChangePassword(userID, form.CurrentPassword, form.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("currentPassword", "Current password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "change_password.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your password has been changed successfully.")

	http.Redirect(w, r, "/scvis/", http.StatusSeeOther)
}

// Issue management handlers

type issueCreateForm struct {
	Title               string `form:"title"`
	Description         string `form:"description"`
	IssueType           string `form:"issueType"`
	Priority            string `form:"priority"`
	SubmitToGitHub      bool   `form:"submitToGitHub"`
	validator.Validator `form:"-"`
}

func (app *application) issueCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = issueCreateForm{}

	// Check if GitHub is configured
	data.GitHubConfigured = app.githubClient.IsConfigured()
	owner, repo := app.githubClient.GetRepository()
	data.GitHubRepo = fmt.Sprintf("%s/%s", owner, repo)

	app.render(w, r, http.StatusOK, "issue_create.tmpl", data)
}

func (app *application) issueCreatePost(w http.ResponseWriter, r *http.Request) {
	var form issueCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validation
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 200), "title", "This field cannot be more than 200 characters long")
	form.CheckField(validator.NotBlank(form.Description), "description", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Description, 2000), "description", "This field cannot be more than 2000 characters long")

	validTypes := []string{"bug", "feature", "enhancement", "documentation", "question"}
	form.CheckField(validator.PermittedValue(form.IssueType, validTypes...), "issueType", "Invalid issue type")

	validPriorities := []string{"low", "medium", "high", "critical"}
	form.CheckField(validator.PermittedValue(form.Priority, validPriorities...), "priority", "Invalid priority level")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		data.GitHubConfigured = app.githubClient.IsConfigured()
		owner, repo := app.githubClient.GetRepository()
		data.GitHubRepo = fmt.Sprintf("%s/%s", owner, repo)
		app.render(w, r, http.StatusUnprocessableEntity, "issue_create.tmpl", data)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	// Create the issue in the database first
	issueID, err := app.issues.Insert(userID, form.Title, form.Description, form.IssueType, form.Priority, nil, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// If submit to GitHub is requested and GitHub is configured
	if form.SubmitToGitHub && app.githubClient.IsConfigured() {
		// Create GitHub issue
		githubIssue := &github.IssueRequest{
			Title: form.Title,
			Body:  app.formatIssueBody(form.Description, form.IssueType, form.Priority, userID),
			Labels: []string{
				"scvis-issue",
				form.IssueType,
				fmt.Sprintf("priority-%s", form.Priority),
			},
		}

		githubResponse, err := app.githubClient.CreateIssue(githubIssue)
		if err != nil {
			app.logger.Error("Failed to create GitHub issue", "error", err, "issueID", issueID)
			app.sessionManager.Put(r.Context(), "flash", "Issue created locally but failed to submit to GitHub: "+err.Error())
		} else {
			// Update the database with GitHub information
			err = app.issues.UpdateGitHubInfo(issueID, githubResponse.Number, githubResponse.HTMLURL)
			if err != nil {
				app.logger.Error("Failed to update issue with GitHub info", "error", err, "issueID", issueID)
			}
			app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Issue created successfully and submitted to GitHub as issue #%d!", githubResponse.Number))
		}
	} else {
		app.sessionManager.Put(r.Context(), "flash", "Issue created successfully!")
	}

	http.Redirect(w, r, "/scvis/issue/list", http.StatusSeeOther)
}

func (app *application) issueList(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	isAdmin, err := app.users.IsAdministrator(userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var issues []*models.Issue
	if isAdmin {
		// Administrators can see all issues
		issues, err = app.issues.GetAll()
	} else {
		// Regular users can only see their own issues
		issues, err = app.issues.GetByUser(userID)
	}

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Issues = issues
	data.GitHubConfigured = app.githubClient.IsConfigured()

	app.render(w, r, http.StatusOK, "issue_list.tmpl", data)
}

func (app *application) issueView(w http.ResponseWriter, r *http.Request) {
	// Get issue ID from session
	id := app.sessionManager.GetInt(r.Context(), "currentIssueID")
	if id < 1 {
		http.NotFound(w, r)
		return
	}

	issue, err := app.issues.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Check if user can view this issue
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	isAdmin, err := app.users.IsAdministrator(userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !isAdmin && issue.UserID != userID {
		app.clientError(w, http.StatusForbidden)
		return
	}

	data := app.newTemplateData(r)
	data.Issue = issue

	app.render(w, r, http.StatusOK, "issue_view.tmpl", data)
}

func (app *application) issueStatusUpdate(w http.ResponseWriter, r *http.Request) {
	// Get issue ID from session
	id := app.sessionManager.GetInt(r.Context(), "currentIssueID")
	if id < 1 {
		http.NotFound(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	status := r.FormValue("status")
	validStatuses := []string{"pending", "submitted", "acknowledged", "in_progress", "resolved", "closed"}

	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.issues.UpdateStatus(id, status)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Issue status updated successfully!")
	http.Redirect(w, r, fmt.Sprintf("/scvis/issue/view/%d", id), http.StatusSeeOther)
}

// Helper function to format the GitHub issue body
func (app *application) formatIssueBody(description, issueType, priority string, userID int) string {
	var body strings.Builder

	body.WriteString("**Issue submitted from SCVis application**\n\n")
	body.WriteString(fmt.Sprintf("**Type:** %s\n", strings.Title(issueType)))
	body.WriteString(fmt.Sprintf("**Priority:** %s\n", strings.Title(priority)))
	body.WriteString(fmt.Sprintf("**Submitted by User ID:** %d\n\n", userID))
	body.WriteString("**Description:**\n")
	body.WriteString(description)
	body.WriteString("\n\n---\n")
	body.WriteString("*This issue was automatically created from the SCVis web application.*")

	return body.String()
}

func (app *application) inputFileView(w http.ResponseWriter, r *http.Request) {
	// Get input file ID from session
	id := app.sessionManager.GetInt(r.Context(), "currentInputFileID")
	if id < 1 {
		http.NotFound(w, r)
		return
	}

	inputfile, err := app.inputfiles.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.InputFile = inputfile

	app.render(w, r, http.StatusOK, "input.tmpl", data)
}

func (app *application) inputFileViewPost(w http.ResponseWriter, r *http.Request) {
	// Get input file ID from session
	id := app.sessionManager.GetInt(r.Context(), "currentInputFileID")
	if id < 1 {
		http.NotFound(w, r)
		return
	}

	inputfile, err := app.inputfiles.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	project, err := app.projects.Get(inputfile.ProjectID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Set current project ID in session for analyze page
	app.sessionManager.Put(r.Context(), "currentProjectID", project.ID)

	http.Redirect(w, r, "/scvis/project/analyze", http.StatusSeeOther)
}

// Single analysis handlers - no login required
func (app *application) singleAnalyze(w http.ResponseWriter, r *http.Request) {
	// If user is authenticated, redirect to home page
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/scvis/home", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.Form = map[string]string{
		"Title":      "",
		"Parameters": "",
	}
	app.render(w, r, http.StatusOK, "single.tmpl", data)
}

func (app *application) singleAnalyzePost(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (32MB max)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get form values
	title := r.FormValue("title")
	if title == "" {
		title = "Quick Analysis"
	}
	parameters := r.FormValue("parameters")
	stdinFileIndex := r.FormValue("stdin-file")

	// Create temporary directory for this analysis
	tempDir, err := os.MkdirTemp("", "scvis-single-*")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Note: We keep the temp directory for the visualization step
	// It will be cleaned up after visualization or on server restart

	// Process uploaded C files
	cFiles := r.MultipartForm.File["upload-files"]
	if len(cFiles) == 0 {
		data := app.newTemplateData(r)
		data.Form = map[string]string{
			"Title":      title,
			"Parameters": parameters,
		}
		data.Flash = "At least one C file is required"
		app.render(w, r, http.StatusUnprocessableEntity, "single.tmpl", data)
		os.RemoveAll(tempDir)
		return
	}

	var cFileList []string
	cFilesDir := filepath.Join(tempDir, "cfiles")
	if err := os.MkdirAll(cFilesDir, 0755); err != nil {
		app.serverError(w, r, err)
		os.RemoveAll(tempDir)
		return
	}

	for _, fileHeader := range cFiles {
		file, err := fileHeader.Open()
		if err != nil {
			app.serverError(w, r, err)
			os.RemoveAll(tempDir)
			return
		}
		defer file.Close()

		// Save C file to temp directory
		destPath := filepath.Join(cFilesDir, fileHeader.Filename)
		destFile, err := os.Create(destPath)
		if err != nil {
			app.serverError(w, r, err)
			os.RemoveAll(tempDir)
			return
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, file)
		if err != nil {
			app.serverError(w, r, err)
			os.RemoveAll(tempDir)
			return
		}

		cFileList = append(cFileList, destPath)
	}

	// Process uploaded input files
	var stdinFilename string
	var inputFileNames []string
	var inputSources []string
	inputFiles := r.MultipartForm.File["upload-input-files"]
	if len(inputFiles) > 0 {
		inputFilesDir := filepath.Join(tempDir, "inputfiles")
		if err := os.MkdirAll(inputFilesDir, 0755); err != nil {
			app.serverError(w, r, err)
			os.RemoveAll(tempDir)
			return
		}

		for i, fileHeader := range inputFiles {
			file, err := fileHeader.Open()
			if err != nil {
				app.serverError(w, r, err)
				os.RemoveAll(tempDir)
				return
			}
			defer file.Close()

			// Save input file to temp directory
			destPath := filepath.Join(inputFilesDir, fileHeader.Filename)
			destFile, err := os.Create(destPath)
			if err != nil {
				app.serverError(w, r, err)
				os.RemoveAll(tempDir)
				return
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, file)
			if err != nil {
				app.serverError(w, r, err)
				os.RemoveAll(tempDir)
				return
			}

			// Check if this file should be used as stdin
			if stdinFileIndex != "" {
				selectedIndex, _ := strconv.Atoi(stdinFileIndex)
				if selectedIndex == i {
					stdinFilename = destPath
				}
			}

			inputFileNames = append(inputFileNames, fileHeader.Filename)
			sourceBytes, err := os.ReadFile(destPath)
			if err != nil {
				app.serverError(w, r, err)
				os.RemoveAll(tempDir)
				return
			}
			inputSources = append(inputSources, string(sourceBytes))
		}
	}

	// Perform analysis using vaccs
	vc := vaccs.NewVaccsComm(app.logger)
	err = vc.Connect()
	if err != nil {
		data := app.newTemplateData(r)
		data.Form = map[string]string{
			"Title":      title,
			"Parameters": parameters,
		}
		data.Flash = "Analysis service unavailable. Please try again later."
		app.render(w, r, http.StatusServiceUnavailable, "single.tmpl", data)
		os.RemoveAll(tempDir)
		return
	}
	defer vc.Close()

	// Read C file sources
	var cSources []string
	var cFilenames []string
	for _, cfilePath := range cFileList {
		source, err := os.ReadFile(cfilePath)
		if err != nil {
			app.serverError(w, r, err)
			os.RemoveAll(tempDir)
			return
		}
		cSources = append(cSources, string(source))
		cFilenames = append(cFilenames, filepath.Base(cfilePath))
	}

	// Parse parameters into array
	var paramArray []string
	if parameters != "" {
		paramArray = strings.Fields(parameters)
	}

	vargs := []string{"-u", "-r", "-j"}

	// Call PerformAnalysis with the correct signature
	analysisResult, err := vc.PerformAnalysis(
		cFilenames,
		cSources,
		inputFileNames,
		inputSources,
		vargs,         // vargs (vaccs arguments)
		"analysis",    // executable name
		paramArray,    // params
		stdinFilename, // stdin filename
	)
	if err != nil {
		data := app.newTemplateData(r)
		data.Form = map[string]string{
			"Title":      title,
			"Parameters": parameters,
		}
		data.Flash = fmt.Sprintf("Analysis failed: %v", err)
		app.render(w, r, http.StatusUnprocessableEntity, "single.tmpl", data)
		os.RemoveAll(tempDir)
		return
	}

	// Save analysis result to temp directory
	analysisPath := filepath.Join(tempDir, "analysis.json")
	if err := os.WriteFile(analysisPath, analysisResult, 0644); err != nil {
		app.serverError(w, r, err)
		os.RemoveAll(tempDir)
		return
	}

	// Store temp directory, analysis file name and a projectId of -1 in the session
	app.sessionManager.Put(r.Context(), "tempDir", tempDir)
	app.sessionManager.Put(r.Context(), "analysisFile", "analysis.json")
	app.sessionManager.Put(r.Context(), "projectID", -1)

	// Set cookie for projectID so Flutter app knows this is a single analysis
	http.SetCookie(w, &http.Cookie{
		Name:  "projectID",
		Value: "-1",
		Path:  "/scvis/",
	})

	// Redirect to visualization
	http.Redirect(w, r, "/scvis/visual", http.StatusSeeOther)
}

// projectUploadCFile handles AJAX uploads of C files
func (app *application) projectUploadCFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get project ID from session
	id := app.sessionManager.GetInt(r.Context(), "currentProjectID")
	if id < 1 {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error getting file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Insert into database
	cfileID, err := app.cfiles.Insert(id, fileHeader.Filename, string(fileContent))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving file: %v", err), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"id":       cfileID,
		"filename": fileHeader.Filename,
	})
}

// projectUploadInputFile handles AJAX uploads of input files
func (app *application) projectUploadInputFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get project ID from session
	id := app.sessionManager.GetInt(r.Context(), "currentProjectID")
	if id < 1 {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error getting file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Insert into database
	inputFileID, err := app.inputfiles.Insert(id, fileHeader.Filename, string(fileContent))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving file: %v", err), http.StatusInternalServerError)
		return
	}

	// Add to project
	err = app.projects.AddInputFiles(id, []int{inputFileID})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding file to project: %v", err), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"id":       inputFileID,
		"filename": fileHeader.Filename,
	})
}
