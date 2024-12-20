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

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("POST /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /project/view/{id}", protected.ThenFunc(app.projectView))
	mux.Handle("GET /project/create", protected.ThenFunc(app.projectCreate))
	mux.Handle("POST /project/create", protected.ThenFunc(app.projectCreatePost))
	mux.Handle("GET /project/update", protected.ThenFunc(app.projectUpdate))
	mux.Handle("POST /project/update", protected.ThenFunc(app.projectUpdatePost))
	mux.Handle("GET /project/analyze", protected.ThenFunc(app.projectAnalyze))
	mux.Handle("POST /project/analyze", protected.ThenFunc(app.projectAnalyzePost))
	mux.Handle("GET /project/visual", protected.ThenFunc(app.projectVisual))
	mux.Handle("GET /project/getanalysis", protected.ThenFunc(app.projectGetAnalysis))
	mux.Handle("GET /project/delete", protected.ThenFunc(app.projectDelete))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	mux.Handle("GET /cfile/view/{id}", protected.ThenFunc(app.cfileView))
	mux.Handle("POST /cfile/view/{id}", protected.ThenFunc(app.cfileViewPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
