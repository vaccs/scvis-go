package main

import (
	"net/http"

	"scvis.vaccsedu.org/scvis"
	"scvis.vaccsedu.org/ui"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /scvis/", http.FileServerFS(scvis.Files))
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /ping", ping)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /project/view/{id}", dynamic.ThenFunc(app.projectView))
	mux.Handle("POST /project/view/{id}", dynamic.ThenFunc(app.projectViewPost))
	mux.Handle("GET /project/create", dynamic.ThenFunc(app.projectCreate))
	mux.Handle("POST /project/create/{$}", dynamic.ThenFunc(app.projectCreatePost))
	mux.Handle("POST /project/analyze/{id}", dynamic.ThenFunc(app.projectAnalyze))
	mux.Handle("GET /project/getanalysis/{id}", dynamic.ThenFunc(app.projectGetAnalysis))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
