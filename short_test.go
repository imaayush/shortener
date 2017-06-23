package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"math/rand"
	"testing"
)

var ts *httptest.Server
var app App

func TestMain(m *testing.M) {

	app = App{}
	app.Initialize()

	ts = httptest.NewServer(app.Router)

	ret := m.Run()

	os.Exit(ret)

}

func MakePostRequest(t *testing.T, Input ShortInput, data *ShortOut) (*http.Response, error) {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(Input)
	req, _ := http.NewRequest("POST", ts.URL+"/short", b)
	client := &http.Client{}
	resp, err := client.Do(req)
	json.NewDecoder(resp.Body).Decode(&data)
	return resp, err

}

func MakeGetRequest(t *testing.T, Input ShortInput) *http.Response {
	req, _ := http.NewRequest("GET", Input.Url, nil)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	resp, _ := client.Do(req)
	return resp

}

func cleanTable() {
	app.DB.Delete(&Short{})
	app.DB.AutoMigrate(&Short{})
}
func TestShortUrlEndPointVaildUrls(t *testing.T) {
	cleanTable()
	rndInt := rand.Int()
	TestUrl := []string{"https://goolge.com/home/param=" +strconv.Itoa(rndInt), "127.0.0.1", "google.com"}
	for i  :=  range TestUrl {
		Input := ShortInput{TestUrl[i]}
		db := app.DB
		var data ShortOut
		var short Short
		resp, _ := MakePostRequest(t, Input, &data)
		assert.Equal(t, resp.StatusCode, 200)
		db.Where("short_url = ?", data.ShortUrl).Find(&short)
		assert.Equal(t, short.ShortUrl, data.ShortUrl)
	}

}

func TestExpandUrlEndPointVaildSlug(t *testing.T) {
	cleanTable()
	rndInt := rand.Int()
	TestUrl := []string{"https://goolge.com/home/param="+strconv.Itoa(rndInt), "127.0.0.1", "google.com"}
	for i  :=  range TestUrl {

		var data ShortOut
		Input := ShortInput{TestUrl[i]}
		resp, _ :=MakePostRequest(t, Input, &data)
		assert.Equal(t, resp.StatusCode, 200)
		url := ts.URL + "/" + data.ShortUrl
		Input = ShortInput{url}
		resp = MakeGetRequest(t, Input)
		assert.Equal(t, resp.StatusCode, 301)
	}

}

func TestShortUrlEndPointInvaildUrls(t *testing.T) {
	TestUrl := []string{"google", "127.0.0.1:8000", ""}
	for i  :=  range TestUrl {
		var data ShortOut
		Input := ShortInput{TestUrl[i]}
		resp, _ := MakePostRequest(t, Input, &data)
		assert.Equal(t, resp.StatusCode, 400)
	}
}

func TestSExpandUrlEndPointInvaildSlug(t *testing.T) {
	app.SlugLength = 5
	slug := app.GenerateShortUrl()
	app.SlugLength = 4
	url := ts.URL + "/" + slug
	Input := ShortInput{url}
	resp := MakeGetRequest(t, Input)
	assert.Equal(t, resp.StatusCode, 404)
}

func TestResponseOnCollision(t *testing.T) {
	cleanTable()
	TestUrl := "https://goolge.com/"

	var data ShortOut
	app.SlugLength =1
	for i := 0; i < 62; i++ {
		url := TestUrl + strconv.Itoa(i)
		Input := ShortInput{string(url)}
		resp, _ := MakePostRequest(t, Input, &data)
		json.NewDecoder(resp.Body).Decode(&data)
		assert.Equal(t, resp.StatusCode, 200)

	}

	url := TestUrl + strconv.Itoa(62)
	Input := ShortInput{string(url)}
	resp, _ := MakePostRequest(t, Input, &data)
	json.NewDecoder(resp.Body).Decode(&data)

	assert.Equal(t, resp.StatusCode, 500)
	app.SlugLength = 4
}

