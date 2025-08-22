package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"scvis.vaccsedu.org/internal/models"
	"scvis.vaccsedu.org/internal/validator"

	"scvis.vaccsedu.org/internal/vaccs"
)

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

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "projectID",
		Value: idStr,
		Path:  "/scvis/",
	})

	app.logger.Info("Serving Flutter web application")
	http.Redirect(w, r, "/scvis/visual", http.StatusSeeOther)
}

func (app *application) projectView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
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
	project, err := app.projects.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Parse the multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

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

	if title != project.Title {
		err = app.projects.UpdateTitle(id, title)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

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

	http.Redirect(w, r, fmt.Sprintf("/scvis/project/update?id=%d", id), http.StatusSeeOther)
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

	http.Redirect(w, r, "/scvis/", http.StatusSeeOther)
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
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id_str := r.Form.Get("id")
	id, err := strconv.Atoi(id_str)
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
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
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
		data, err := vaccs_comm.PerformAnalysis(fileNames, cSources, vargs, safeTitle, paramArray)
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

	http.Redirect(w, r, fmt.Sprintf("/scvis/project/analyze?id=%d", id), http.StatusSeeOther)

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

	jsonData, err := json.Marshal(string(project.Analysis))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func (app *application) projectDelete(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id_str := r.Form.Get("id")
	id, err := strconv.Atoi(id_str)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	err = app.projects.Delete(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Project successfully deleted!")

	http.Redirect(w, r, "/scvis/", http.StatusSeeOther)
}

func (app *application) cfileView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/scvis/cfile/view/"))
	if err != nil || id < 1 {
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
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/scvis/cfile/view/"))
	if err != nil || id < 1 {
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

	http.Redirect(w, r, fmt.Sprintf("/scvis/project/analyze?id=%d", project.ID), http.StatusSeeOther)
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
