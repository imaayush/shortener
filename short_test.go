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

func MakeRequest(t *testing.T, Input ShortInput, data *ShortOut) *http.Response {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(Input)
	req, _ := http.NewRequest("POST", ts.URL+"/short", b)
	client := &http.Client{}
	resp, _ := client.Do(req)

	json.NewDecoder(resp.Body).Decode(&data)
	return resp
}

func TestShortUrlEndPointPassCase(t *testing.T) {

	TestUrl := "https://goolge.com/home/param=11"
	Input := ShortInput{TestUrl}
	db := app.DB
	var data ShortOut
	var short Short
	resp := MakeRequest(t, Input, &data)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, Input.Url, data.Url)
	db.Where("short_url = ?", data.ShortUrl).Find(&short)
	assert.Equal(t, short.ShortUrl, data.ShortUrl)

}

func TestExpandUrlEndPointPassCase(t *testing.T) {

	TestUrl := "http://goolge.com/"
	var data ShortOut
	Input := ShortInput{TestUrl}
	MakeRequest(t, Input, &data)
	assert.Equal(t, Input.Url, data.Url)
	url := ts.URL + "/" + data.ShortUrl

	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	resp, _ := client.Do(req)
	assert.Equal(t, resp.StatusCode, 301)

}

func TestWrongInput(t *testing.T) {
	TestCase := "google"

	var data ShortOut
	Input := ShortInput{TestCase}
	resp := MakeRequest(t, Input, &data)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestShortUrlNotFound(t *testing.T) {
	var data ShortOut
	url := ts.URL + "/" + "ASDFW"

	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, _ := client.Do(req)

	json.NewDecoder(resp.Body).Decode(&data)
	assert.Equal(t, resp.StatusCode, 404)
}
