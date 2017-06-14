package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
)

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
