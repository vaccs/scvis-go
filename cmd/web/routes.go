package main

import (
	"net/http"

	"scvis.vaccsedu.org/ui"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.Handle("GET /{$}", http.RedirectHandler("/scvis/", http.StatusFound))

	mux.Handle("/scvis/visual/", http.StripPrefix("/scvis/visual/", http.FileServer(http.Dir("./visual/scvis"))))
	// mux.Handle("GET /scvis/visual/", http.FileServerFS(visual.Files))
	mux.Handle("GET /scvis/static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	loginOnly := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate).Append(app.redirectIfAuthenticated)
	protected := dynamic.Append(app.requireAuthentication)
	adminOnly := dynamic.Append(app.requireAdministrator)

	// Root and single analysis routes - no authentication required, default landing page
	mux.Handle("GET /scvis/{$}", dynamic.ThenFunc(app.singleAnalyze))
	mux.Handle("POST /scvis/single/analyze", dynamic.ThenFunc(app.singleAnalyzePost))
	mux.Handle("GET /scvis/single/getanalysis", dynamic.ThenFunc(app.singleGetAnalysis))

	// Authenticated home page for logged-in users
	mux.Handle("GET /scvis/home", protected.ThenFunc(app.home))
	mux.Handle("POST /scvis/home", protected.ThenFunc(app.home))

	// Session setters - these set IDs in session and redirect
	mux.Handle("POST /scvis/project/set", protected.ThenFunc(app.setProjectSession))
	mux.Handle("POST /scvis/cfile/set", protected.ThenFunc(app.setCFileSession))
	mux.Handle("POST /scvis/input/set", protected.ThenFunc(app.setInputFileSession))
	mux.Handle("POST /scvis/issue/set", protected.ThenFunc(app.setIssueSession))

	mux.Handle("GET /scvis/project/view", protected.ThenFunc(app.projectView))
	mux.Handle("POST /scvis/project/view", protected.ThenFunc(app.projectView))
	mux.Handle("GET /scvis/project/create", protected.ThenFunc(app.projectCreate))
	mux.Handle("POST /scvis/project/create", protected.ThenFunc(app.projectCreatePost))
	mux.Handle("GET /scvis/project/update", protected.ThenFunc(app.projectUpdate))
	mux.Handle("POST /scvis/project/update", protected.ThenFunc(app.projectUpdatePost))
	mux.Handle("POST /scvis/project/upload-cfile", protected.ThenFunc(app.projectUploadCFile))
	mux.Handle("POST /scvis/project/upload-inputfile", protected.ThenFunc(app.projectUploadInputFile))
	mux.Handle("GET /scvis/project/analyze", protected.ThenFunc(app.projectAnalyze))
	mux.Handle("POST /scvis/project/analyze", protected.ThenFunc(app.projectAnalyzePost))
	mux.Handle("GET /scvis/project/visual", protected.ThenFunc(app.projectVisual))
	mux.Handle("POST /scvis/project/visual", protected.ThenFunc(app.projectVisual))
	mux.Handle("GET /scvis/project/getanalysis", dynamic.ThenFunc(app.projectGetAnalysis))
	mux.Handle("POST /scvis/project/delete", protected.ThenFunc(app.projectDelete))
	mux.Handle("GET /scvis/user/signup", adminOnly.ThenFunc(app.userSignup))
	mux.Handle("POST /scvis/user/signup", adminOnly.ThenFunc(app.userSignupPost))
	mux.Handle("GET /scvis/user/login", loginOnly.ThenFunc(app.userLogin))
	mux.Handle("POST /scvis/user/login", loginOnly.ThenFunc(app.userLoginPost))
	mux.Handle("GET /scvis/user/change-password", protected.ThenFunc(app.userChangePassword))
	mux.Handle("POST /scvis/user/change-password", protected.ThenFunc(app.userChangePasswordPost))
	mux.Handle("POST /scvis/user/logout", protected.ThenFunc(app.userLogoutPost))
	mux.Handle("GET /scvis/cfile/view", protected.ThenFunc(app.cfileView))
	mux.Handle("POST /scvis/cfile/view", protected.ThenFunc(app.cfileViewPost))
	mux.Handle("GET /scvis/input/view", protected.ThenFunc(app.inputFileView))
	mux.Handle("POST /scvis/input/view", protected.ThenFunc(app.inputFileViewPost))
	mux.Handle("GET /scvis/issue/create", adminOnly.ThenFunc(app.issueCreate))
	mux.Handle("POST /scvis/issue/create", adminOnly.ThenFunc(app.issueCreatePost))
	mux.Handle("GET /scvis/issue/list", protected.ThenFunc(app.issueList))
	mux.Handle("GET /scvis/issue/view", protected.ThenFunc(app.issueView))
	mux.Handle("POST /scvis/issue/view", protected.ThenFunc(app.issueView))
	mux.Handle("POST /scvis/issue/status", adminOnly.ThenFunc(app.issueStatusUpdate))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
