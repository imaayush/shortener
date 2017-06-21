package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (app *App) ShortUrl(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	ShortInput := ShortInput{}

	if err := json.Unmarshal(body, &ShortInput); err != nil {
		panic(err)
	}

	u, err := url.Parse(ShortInput.Url)
	if err != nil {
		http.Error(w, "please enter correct url", 400)
		return
	}

	if u.Host == "" {
		HostName := strings.Split(u.Path, ".")
		if len(HostName) <= 1 {
			http.Error(w, "please enter correct url", 400)
			return

		} else if u.Scheme == "" {
			u.Scheme = "https"
			fmt.Println("set https scheme ")
		}
	}

	LongUrl := u.String()

	var ShortUrl string

	var short Short

	db := app.DB
	defer app.Unlock()

	app.Lock()

	if err := db.Where("url = ?", LongUrl).Find(&short).Error; err != nil {
		UnquieUrl := false

		for UnquieUrl != true {
			ShortUrl = GenerateShortUrl()
			UnquieUrl = app.CheckIsUnqiue(ShortUrl)

		}
		short = Short{Url: LongUrl, ShortUrl: ShortUrl}
		if err := db.Save(&short).Error; err != nil {
			http.Error(w, "UNIQUE constraint failed", 400)
		}
	} else {
		short = Short{Url: short.Url, ShortUrl: short.ShortUrl}
	}

	var data = ShortOut{Url: short.Url, ShortUrl: short.ShortUrl}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}

}

func (app *App) ExpandUrl(w http.ResponseWriter, r *http.Request) {
	var short Short
	vars := mux.Vars(r)
	ShortUrl := vars["uuid"]
	db := app.DB
	if err := db.Where("short_url = ?", ShortUrl).First(&short).Error; err != nil {
		http.Error(w, "Page not found", 404)
		return
	}

	http.Redirect(w, r, short.Url, 301)
	return

}
