package main

import (
	"net/http"
	"path"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.logRequest, secureHeaders)

	mux := http.NewServeMux()
	mux.HandleFunc(path.Clean("/"), app.home)
	mux.HandleFunc(path.Clean("/upload"), app.uploadFile)
	mux.HandleFunc(path.Clean("/getcount"), app.GetCount)

	return standardMiddleware.Then(mux)
}
