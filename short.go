
package main

import (
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

func main() {

	//Migrate the schema
	db := Database()
	db.AutoMigrate(&Short{})

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/short", ShortUrl )
	router.HandleFunc("/{uuid}", ExpandUrl)

	log.Fatal(http.ListenAndServe(":8080", router))

}

