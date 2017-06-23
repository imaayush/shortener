package main

import (
	"github.com/gorilla/mux"
)

func (app *App) NewRouter() *mux.Router {

	app.Router = mux.NewRouter().StrictSlash(true)
	app.Router.HandleFunc("/{slug}", app.ExpandUrl).Methods("GET")
	app.Router.HandleFunc("/short", app.ShortUrl).Methods("POST")
	return app.Router
}
