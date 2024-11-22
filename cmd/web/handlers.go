package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"scvis.vaccsedu.org/internal/models"
	"scvis.vaccsedu.org/internal/validator"

	"scvis.vaccsedu.org/internal/vaccs"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
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

func (app *application) projectView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	fmt.Printf("Project ID: %d\n", id)
	if err != nil || id < 1 {
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

	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *application) projectUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id_str := r.Form.Get("id")
	id, err := strconv.Atoi(id_str)
	fmt.Printf("Project ID: %d\n", id)
	if err != nil || id < 1 {
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

	app.render(w, r, http.StatusOK, "update.tmpl", data)
}

func (app *application) projectUpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Println("projectUpdatePost: get project")
	project, err := app.projects.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	fmt.Println("projectUpdatePost: parse form")
	// Parse the multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	fmt.Println("projectUpdatePost: get files")
	// Retrieve the files from form data
	files := r.MultipartForm.File["upload-files"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error opening file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileContent, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		_, err = app.cfiles.Insert(id, fileHeader.Filename, string(fileContent))
		if err != nil {
			app.serverError(w, r, err)
			return
		}

	}

	title := r.MultipartForm.Value["title"][0]

	fmt.Println("projectUpdatePost: update project title")

	if title != project.Title {
		err = app.projects.UpdateTitle(id, title)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	fmt.Println("projectUpdatePost: delete selected cfiles")
	// get the list of cfileids to delete
	cfileids := []int{}
	for _, idStr := range r.MultipartForm.Value["cfileids"] {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		cfileids = append(cfileids, id)
	}

	// delete the cfiles
	for _, cfileID := range cfileids {
		err = app.cfiles.Delete(cfileID)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	app.sessionManager.Put(r.Context(), "flash", "Project successfully updated!")

	http.Redirect(w, r, fmt.Sprintf("/project/update?id=%d", id), http.StatusSeeOther)
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

	projectID, err := app.projects.Insert(form.Title, app.sessionManager.GetInt(r.Context(), "authenticatedUserID"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Project successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/project/view/%d", projectID), http.StatusSeeOther)
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
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

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
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

	err = app.users.Insert(form.Name, form.Email, form.Password)
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

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
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

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) projectAnalyze(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id_str := r.Form.Get("id")
	id, err := strconv.Atoi(id_str)
	fmt.Printf("Project ID: %d\n", id)
	if err != nil || id < 1 {
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
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	params := r.PathValue("cmd-params")
	paramArray := strings.Fields(params)

	project, err := app.projects.Get(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

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

	vaccs_comm := vaccs.NewVaccsComm()

	err = vaccs_comm.Connect()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	defer vaccs_comm.Close()

	data, err := vaccs_comm.PerformAnalysis(fileNames, cSources, vargs, project.Title, paramArray)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.projects.UpdateAnalysis(id, data)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/project/view?id=%d", id), http.StatusSeeOther)

}

func (app *application) projectGetAnalysis(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	project, err := app.projects.Get(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	jsonData, err := json.Marshal(project.Analysis)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
