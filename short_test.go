package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var ts *httptest.Server
var app App

func TestMain(m *testing.M) {

	app = App{}
	app.Initialize("", "", TestDb)
	ts = httptest.NewServer(app.Router)
	ret := m.Run()

	os.Exit(ret)
}

func MakeRequest(t *testing.T, Input ShortInput, data *ShortOut, method string) *http.Response {
	if method == "POST" {

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(Input)
		req, _ := http.NewRequest("POST", ts.URL+"/short", b)
		client := &http.Client{}
		resp, _ := client.Do(req)

		json.NewDecoder(resp.Body).Decode(&data)
		return resp
	} else {
		req, _ := http.NewRequest("GET", Input.Url, nil)

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}
		resp, _ := client.Do(req)
		return resp
	}
}

func cleanTable() {
	app.DB.Delete(&Short{})
	app.DB.AutoMigrate(&Short{})
}
func TestShortUrlEndPointPassCase(t *testing.T) {
	cleanTable()
	TestUrl := "https://goolge.com/home/param=11"
	Input := ShortInput{TestUrl}
	db := app.DB
	var data ShortOut
	var short Short
	resp := MakeRequest(t, Input, &data, "POST")
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, Input.Url, data.Url)
	db.Where("short_url = ?", data.ShortUrl).Find(&short)
	assert.Equal(t, short.ShortUrl, data.ShortUrl)

}

func TestExpandUrlEndPointPassCase(t *testing.T) {
	cleanTable()

	TestUrl := "http://goolge.com/"
	var data ShortOut
	Input := ShortInput{TestUrl}
	MakeRequest(t, Input, &data, "POST")
	assert.Equal(t, Input.Url, data.Url)
	url := ts.URL + "/" + data.ShortUrl
	Input = ShortInput{url}
	resp := MakeRequest(t, Input, &data, "GET")
	assert.Equal(t, resp.StatusCode, 301)

}

func TestWrongInput(t *testing.T) {
	TestCase := "google"

	var data ShortOut
	Input := ShortInput{TestCase}
	resp := MakeRequest(t, Input, &data, "POST")
	assert.Equal(t, resp.StatusCode, 400)
}

func TestShortUrlNotFound(t *testing.T) {
	var data ShortOut
	url := ts.URL + "/" + "ASDFW"
	Input := ShortInput{url}
	resp := MakeRequest(t, Input, &data, "GET")
	assert.Equal(t, resp.StatusCode, 404)
}

func TestCollision(t *testing.T) {
	cleanTable()
	TestUrl := "https://goolge.com/home/param=11"
	Input := ShortInput{TestUrl}
	db := app.DB
	var data ShortOut
	var short Short
	resp := MakeRequest(t, Input, &data, "POST")
	assert.Equal(t, resp.StatusCode, 200)
	short = Short{Url: data.Url, ShortUrl: data.ShortUrl}
	err := db.Save(&short).Error
	errStr := "pq: duplicate key value violates unique constraint \"shorts_short_url_key\""
	assert.EqualError(t, err, errStr)

}
