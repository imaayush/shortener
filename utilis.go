package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io"
	"io/ioutil"
	"math/rand"
	"net/url"
	"strings"
	"sync"
)

type App struct {
	Router     *mux.Router
	DB         *gorm.DB
	SlugLength int
	sync.Mutex
}

func (app *App) GenerateShortUrl() string {
	b := make([]rune, app.SlugLength)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)

}

func (app *App) CheckIsUniquie(url string) bool {
	db := app.DB
	var short Short
	if err := db.Where("short_url = ?", string(url)).Find(&short).Error; err != nil {
		return true

	} else {
		return false
	}
}
func (app *App) GenerateAndSave(LongUrl string) (Short, error) {

	ShortUrl := app.GenerateShortUrl()
	data, err := app.SaveUrlAndCheckUniquie(ShortUrl, LongUrl)
	return data, err
}

func (app *App) SaveUrlAndCheckUniquie(ShortUrl string, LongUrl string) (Short, error) {
	db := app.DB
	var short Short
	IsUniquie := app.CheckIsUniquie(ShortUrl)
	err := errors.New("pq: duplicate key value violates unique constraint \"shorts_short_url_key\"")
	if !IsUniquie {
		return Short{}, err
	}
	short = Short{Url: LongUrl, ShortUrl: ShortUrl}
	if err := db.Save(&short).Error; err != nil {
		return short, err
	}
	var data = Short{Url: short.Url, ShortUrl: short.ShortUrl}
	return data, nil
}

func Database() *gorm.DB {
	//open a db connection
	db, err := gorm.Open(Db, DbConfig)

	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Short{})
	return db
}
func ValidateUrl(Url string) (string, error) {
	u, err := url.Parse(Url)
	if err != nil {

		return "", err
	}

	if u.Host == "" {
		HostName := strings.Split(u.Path, ".")
		if len(HostName) <= 1 {
			err = errors.New("Host Name is not difine")
			return "", err

		} else if u.Scheme == "" {
			u.Scheme = "http"

		}
	}

	return u.String(), nil

}

func GetAndValidateInput(r io.Reader) (string, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	ShortInput := ShortInput{}

	if err := json.Unmarshal(body, &ShortInput); err != nil {
		return "", err

	}
	LongUrl, err := ValidateUrl(ShortInput.Url)

	if err != nil {
		return LongUrl, err

	}
	return LongUrl, nil
}
