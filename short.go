package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

func main() {
	app := App{}
	app.Initialize()

	app.Run()

}

func (app *App) Initialize() {
	app.SlugLength = SlugLength
	app.DB = Database()

	app.NewRouter()
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}
