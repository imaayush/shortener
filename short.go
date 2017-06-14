
package main

import (
	"net/http"
	"math/rand"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"io/ioutil"
	"encoding/json"

	"fmt"
)


func Database() *gorm.DB {
	//open a db connection
	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func main() {

	//Migrate the schema
	db := Database()
	db.AutoMigrate(&Short{})

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/short", ShortUrl )
	router.HandleFunc("/{uuid}", ExpandUrl)

	log.Fatal(http.ListenAndServe(":8080", router))

}

type Short struct {
	gorm.Model
	Url       string   `json:"Url"`
	ShortUrl  string   `json:"ShortUrl"`
}

type ShortInput struct {
	Url       string   `json:"url"`
}

func ShortUrl(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	ShortInput := ShortInput{}

	if err := json.Unmarshal(body, &ShortInput); err != nil {
		panic(err)
	}
	url := ShortInput.Url
	ShortUrl := GenerateShortUrl()
	short := Short{Url: url, ShortUrl:ShortUrl}
	db := Database()
	db.Save(&short)
	if err := json.NewEncoder(w).Encode(short); err != nil {
		panic(err)
	}

}


func ExpandUrl(w http.ResponseWriter, r *http.Request) {
	var short Short
	vars := mux.Vars(r)
	ShortUrl := vars["uuid"]
	db := Database()
	db.Where("short_url = ?", ShortUrl).Find(&short)
	fmt.Println(short.Url)
	http.Redirect(w, r, short.Url, 301)
}
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateShortUrl()string{
	n := 16
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}