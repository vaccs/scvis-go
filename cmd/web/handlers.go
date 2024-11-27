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

func (app *application) projectVisual(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "projectID",
		Value: idStr,
		Path:  "/",
	})

	http.Redirect(w, r, "/scvis", http.StatusSeeOther)
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

	const jsonString = `[{"VACCS_ARCH":{"Event_num":-1,"Arch_type":1,"Heap_start":94400437915648,"Heap_end":140728349331448,"Stack_start":140728349470720,"Stack_end":140728349331456}},{"VACCS_CMD_LINE":{"Event_num":-1,"Cmd_line":"./t1 "}},{"VACCS_BINARY":{"Event_num":-1,"Path":"./t1"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":1,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"#include <cstdio.h>"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":2,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"int a;"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":3,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"int c = 3;"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":4,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"int main() {"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":5,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":""}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":6,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"   int b;"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":7,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":""}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":8,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"   int *p = NULL;"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":9,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":""}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":10,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"   p = b; "}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":11,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"   a = 1;"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":12,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"   b = 2;"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"   printf(\"a = %d, b = %d, *p = %d, c = %d\\n\",a, b, *p, c);"}},{"VACCS_CCODE":{"Event_num":-1,"Line_num":14,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","C_start_pos":0,"Source":"}"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|Type|Address|Size|Section Name"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb8318|28|.interp"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb8338|48|.note.gnu.property"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb8368|36|.note.gnu.build-id"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb838c|32|.note.ABI-tag"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb83b0|36|.gnu.hash"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|16|0x55db4eeb83d8|168|.dynsym"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|20|0x55db4eeb8480|143|.dynstr"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb8510|14|.gnu.version"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeb8520|48|.gnu.version_r"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|2|0x55db4eeb8550|192|.rela.dyn"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|2|0x55db4eeb8610|24|.rela.plt"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|4|0x55db4eeb9000|27|.init"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|4|0x55db4eeb9020|32|.plt"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|4|0x55db4eeb9040|16|.plt.got"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|4|0x55db4eeb9050|16|.plt.sec"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|4|0x55db4eeb9060|331|.text"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|4|0x55db4eeb91ac|13|.fini"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeba000|41|.rodata"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeba02c|52|.eh_frame_hdr"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eeba060|172|.eh_frame"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eebbdb8|8|.init_array"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eebbdc0|8|.fini_array"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|6|0x55db4eebbdc8|496|.dynamic"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|8|0x55db4eebbfb8|72|.got"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|5|0x55db4eebc000|20|.data"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|18|0x55db4eebc014|12|.bss"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|43|.comment"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|48|.debug_aranges"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|244|.debug_info"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|160|.debug_abbrev"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|98|.debug_line"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|220|.debug_str"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|17|0x0|84|.debug_line_str"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|15|0x0|912|.symtab"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|19|0x0|477|.strtab"}},{"VACCS_SECTION":{"Event_num":-1,"Info":"|21|0x0|362|.shstrtab"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37193,"C_line_num":4,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"nop edx, edi"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37197,"C_line_num":4,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"push rbp"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37198,"C_line_num":4,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov rbp, rsp"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37201,"C_line_num":4,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"sub rsp, 0x10"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37205,"C_line_num":8,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov qword ptr [rbp-0x8], 0x0"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37213,"C_line_num":10,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"lea rax, ptr [rbp-0xc]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37217,"C_line_num":10,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov qword ptr [rbp-0x8], rax"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37221,"C_line_num":11,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov dword ptr [rip+0x2ea9], 0x1"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37231,"C_line_num":12,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov dword ptr [rbp-0xc], 0x2"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37238,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov esi, dword ptr [rip+0x2e94]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37244,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov rax, qword ptr [rbp-0x8]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37248,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov ecx, dword ptr [rax]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37250,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov edx, dword ptr [rbp-0xc]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37253,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov eax, dword ptr [rip+0x2e8d]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37259,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov r8d, esi"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37262,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov esi, eax"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37264,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"lea rax, ptr [rip+0xe71]"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37271,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov rdi, rax"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37274,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov eax, 0x0"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37279,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"call 0x55db4eeb9050"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37284,"C_line_num":13,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"mov eax, 0x0"}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37289,"C_line_num":14,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"leave "}},{"VACCS_ASM":{"Event_num":-1,"Asm_line_num":37290,"C_line_num":14,"Asm_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.s","C_file_name":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Inst":"ret "}},{"VACCS_FUNCTION_INV":{"Event_num":1,"Line_num":0,"Path":"__NOCSOURCE__","Inv_line_num":0,"C_inv_path":"__NOCSOURCE__","Func_name":"_start","Address":140064791934592}},{"VACCS_FUNCTION_INV":{"Event_num":3,"Line_num":4,"Path":"","Inv_line_num":0,"C_inv_path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Func_name":"main","Address":140728349461256}},{"VACCS_RETURN_ADDR":{"Event_num":3,"Dynamic_link":140064477568400,"Return_address":1,"Function_name":"main"}},{"VACCS_LINE":{"Event_num":4,"Line_num":4,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c"}},{"VACCS_REGISTER":{"Event_num":4,"Line_num":4,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rsp","Value":140728349461248}},{"VACCS_REGISTER":{"Event_num":5,"Line_num":4,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rbp","Value":140728349461248}},{"VACCS_REGISTER":{"Event_num":6,"Line_num":4,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rsp","Value":140728349461232}},{"VACCS_REGISTER":{"Event_num":7,"Line_num":4,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rflags","Value":518}},{"VACCS_VAR_ACCESS":{"Event_num":8,"Line_num":8,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Scope":"main","Address":140728349461240,"Var_type":"<pointer>","Name":"p","Points_to":"","Value":"<NULL>"}},{"VACCS_RETURN_ADDR":{"Event_num":8,"Dynamic_link":140064477568400,"Return_address":1,"Function_name":"main"}},{"VACCS_LINE":{"Event_num":9,"Line_num":10,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c"}},{"VACCS_REGISTER":{"Event_num":9,"Line_num":10,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rax","Value":140728349461236}},{"VACCS_VAR_ACCESS":{"Event_num":10,"Line_num":10,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Scope":"main","Address":140728349461240,"Var_type":"<pointer>","Name":"p","Points_to":2,"Value":"0x7ffddf478af4"}},{"VACCS_VAR_ACCESS":{"Event_num":11,"Line_num":11,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Scope":"*G*","Address":94400410271768,"Var_type":"int","Name":"a","Points_to":"","Value":"1"}},{"VACCS_VAR_ACCESS":{"Event_num":12,"Line_num":12,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Scope":"main","Address":140728349461236,"Var_type":"int","Name":"b","Points_to":"","Value":"2"}},{"VACCS_VAR_ACCESS":{"Event_num":12,"Line_num":12,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Scope":"main","Address":140728349461240,"Var_type":"<pointer>","Name":"p","Points_to":2,"Value":"0x7ffddf478af4"}},{"VACCS_LINE":{"Event_num":13,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c"}},{"VACCS_REGISTER":{"Event_num":13,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"esi","Value":3}},{"VACCS_REGISTER":{"Event_num":14,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"ecx","Value":2}},{"VACCS_REGISTER":{"Event_num":15,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"edx","Value":2}},{"VACCS_REGISTER":{"Event_num":16,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"eax","Value":1}},{"VACCS_REGISTER":{"Event_num":17,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"r8d","Value":3}},{"VACCS_REGISTER":{"Event_num":18,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"esi","Value":1}},{"VACCS_REGISTER":{"Event_num":19,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rax","Value":94400410263560}},{"VACCS_REGISTER":{"Event_num":20,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rdi","Value":94400410263560}},{"VACCS_REGISTER":{"Event_num":21,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"eax","Value":0}},{"VACCS_REGISTER":{"Event_num":23,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rip","Value":94400410259536}},{"VACCS_REGISTER":{"Event_num":24,"Line_num":13,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rsp","Value":140728349461224}},{"VACCS_LINE":{"Event_num":26,"Line_num":14,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c"}},{"VACCS_REGISTER":{"Event_num":26,"Line_num":14,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rbp","Value":1}},{"VACCS_REGISTER":{"Event_num":27,"Line_num":14,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rsp","Value":140728349461256}},{"VACCS_RETURN":{"Event_num":28}},{"VACCS_REGISTER":{"Event_num":30,"Line_num":14,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rip","Value":140064477568400}},{"VACCS_REGISTER":{"Event_num":31,"Line_num":14,"Path":"/home/carr/dynamic_analysis/pin/source/tools/PAS/tests/test1.c","Register":"rsp","Value":140728349461264}},{"VACCS_RETURN":{"Event_num":35}}]`

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

	// project, err := app.projects.Get(id)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	jsonData, err := json.Marshal(jsonString)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}
