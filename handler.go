package short

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"net/url"
	"strings"

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
	UrlStr := ShortInput.Url
	if u, err := url.Parse(UrlStr); err == nil {
		if u.Host == "" {
			HostName := strings.Split(u.Path, ".")
			if len(HostName) <= 1{
				http.Error(w,"please enter correct url", 402)
				return

			}else if u.Scheme == ""{
				u.Scheme = "https"
				fmt.Println("set https scheme ")
			}
		}

		LongUrl := u.String()
		ShortUrl := GenerateShortUrl()
		var short Short

		db := Database()
		if err := db.Where("url = ?", LongUrl).Find(&short).Error; err != nil {
			short = Short{Url: LongUrl, ShortUrl:ShortUrl}
			db.Save(&short)
		} else {
			short = Short{Url:short.Url, ShortUrl: short.ShortUrl}
		}

		var data = ShortOut{Url:short.Url, ShortUrl: short.ShortUrl}
		if err := json.NewEncoder(w).Encode(data); err != nil {
			panic(err)
		}
	}

}

func ExpandUrl(w http.ResponseWriter, r *http.Request) {
	var short Short
	vars := mux.Vars(r)
	ShortUrl := vars["uuid"]
	fmt.Println(ShortUrl)
	db := Database()
	db.Where("short_url = ?", ShortUrl).Find(&short)
	fmt.Println(short.Url)
	http.Redirect(w, r, short.Url, 200)

}
