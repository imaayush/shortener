package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

func main() {
	app := App{}
	app.Initialize("", "", DbName)

	app.Run()

}

func (app *App) Initialize(user, password, DbName string) {

	app.DB = Database(DbName)

	app.NewRouter()
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}
