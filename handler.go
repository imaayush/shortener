package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *App) ShortUrl(w http.ResponseWriter, r *http.Request) {
	LongUrl, err := GetAndValidateInput(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 412)
	}
	var short Short
	db := app.DB
	if err := db.Where("url = ?", LongUrl).Find(&short).Error; err == nil {
		var data = Short{Url: short.Url, ShortUrl: short.ShortUrl}
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

	}

	app.Lock()
	defer app.Unlock()
	var cnt int

	var data Short
	CollisionErr := errors.New("Start")
	cnt = 0
	MaxUnqiueUrl := len(letterRunes)^app.SlugLength -1
	for CollisionErr != nil && cnt < MaxUnqiueUrl {
		db.Table("shorts").Count(&cnt)
		data, CollisionErr = app.GenerateAndSave(LongUrl)

	}
	if CollisionErr != nil {
		http.Error(w, CollisionErr.Error(), 500)
	}
	data.ShortUrl = ConvetSlugTOUrl(data.ShortUrl, r.Host)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	return
}

func (app *App) ExpandUrl(w http.ResponseWriter, r *http.Request) {
	var short Short
	vars := mux.Vars(r)
	ShortUrl := vars["slug"]
	db := app.DB
	if err := db.Where("short_url = ?", ShortUrl).First(&short).Error; err != nil {
		http.Error(w, "Page not found", 404)
		return
	}

	http.Redirect(w, r, short.Url, 301)
	return

}
