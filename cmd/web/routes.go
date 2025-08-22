package main

import (
	"net/http"

	"scvis.vaccsedu.org/ui"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/scvis/visual/", http.StripPrefix("/scvis/visual/", http.FileServer(http.Dir("./visual/scvis"))))
	// mux.Handle("GET /scvis/visual/", http.FileServerFS(visual.Files))
	mux.Handle("GET /scvis/static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	loginOnly := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate).Append(app.redirectIfAuthenticated)
	protected := dynamic.Append(app.requireAuthentication)
	adminOnly := dynamic.Append(app.requireAdministrator)

	mux.Handle("GET /scvis/{$}", protected.ThenFunc(app.home))
	mux.Handle("POST /scvis/{$}", protected.ThenFunc(app.home))
	mux.Handle("GET /scvis/project/view/{id}", protected.ThenFunc(app.projectView))
	mux.Handle("GET /scvis/project/create", protected.ThenFunc(app.projectCreate))
	mux.Handle("POST /scvis/project/create", protected.ThenFunc(app.projectCreatePost))
	mux.Handle("GET /scvis/project/update", protected.ThenFunc(app.projectUpdate))
	mux.Handle("POST /scvis/project/update", protected.ThenFunc(app.projectUpdatePost))
	mux.Handle("GET /scvis/project/analyze", protected.ThenFunc(app.projectAnalyze))
	mux.Handle("POST /scvis/project/analyze", protected.ThenFunc(app.projectAnalyzePost))
	mux.Handle("GET /scvis/project/visual", protected.ThenFunc(app.projectVisual))
	mux.Handle("GET /scvis/project/getanalysis", protected.ThenFunc(app.projectGetAnalysis))
	mux.Handle("GET /scvis/project/delete", protected.ThenFunc(app.projectDelete))
	mux.Handle("GET /scvis/user/signup", adminOnly.ThenFunc(app.userSignup))
	mux.Handle("POST /scvis/user/signup", adminOnly.ThenFunc(app.userSignupPost))
	mux.Handle("GET /scvis/user/login", loginOnly.ThenFunc(app.userLogin))
	mux.Handle("POST /scvis/user/login", loginOnly.ThenFunc(app.userLoginPost))
	mux.Handle("GET /scvis/user/change-password", protected.ThenFunc(app.userChangePassword))
	mux.Handle("POST /scvis/user/change-password", protected.ThenFunc(app.userChangePasswordPost))
	mux.Handle("POST /scvis/user/logout", protected.ThenFunc(app.userLogoutPost))
	mux.Handle("GET /scvis/cfile/view/{id}", protected.ThenFunc(app.cfileView))
	mux.Handle("POST /scvis/cfile/view/{id}", protected.ThenFunc(app.cfileViewPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
