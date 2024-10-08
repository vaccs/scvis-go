package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"scvis.vaccsedu.org/internal/models"
	"scvis.vaccsedu.org/internal/validator"
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

func (app *application) projectViewPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	var form projectCreateForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "view.tmpl", data)
		return
	}

	err = app.projects.UpdateTitle(id, form.Title)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Remove selected C files
	for _, cfileID := range form.CFileIDs {
		err = app.cfiles.Delete(cfileID)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	// Add new C files
	files := r.MultipartForm.File["cfiles"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		defer file.Close()

		content := make([]byte, fileHeader.Size)
		_, err = file.Read(content)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		_, err = app.cfiles.Insert(id, fileHeader.Filename, string(content))
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	app.sessionManager.Put(r.Context(), "flash", "Project successfully updated!")

	http.Redirect(w, r, fmt.Sprintf("/project/view/%d", id), http.StatusSeeOther)
}

type projectCreateForm struct {
	Title               string `form:"title"`
	CFileIDs            []int  `form:"files"`
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

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		defer file.Close()

		content := make([]byte, fileHeader.Size)
		_, err = file.Read(content)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		_, err = app.cfiles.Insert(projectID, fileHeader.Filename, string(content))
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	app.sessionManager.Put(r.Context(), "flash", "Project successfully created!")

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

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
