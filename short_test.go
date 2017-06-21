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

func MultMakeRequest(t *testing.T, Input ShortInput, c chan ShortOut) {
	var data ShortOut
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(Input)
	req, _ := http.NewRequest("POST", ts.URL+"/short", b)
	client := &http.Client{}
	resp, _ := client.Do(req)
	json.NewDecoder(resp.Body).Decode(&data)
	c <- data

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

func TestShortUrlUniqueness(t *testing.T) {
	cleanTable()
	letterRunes = rune("abcd")
	TestUrl := "http://goolge.com/"
	//var data ShortOut
	c := make(chan ShortOut)
	output := make([]ShortOut, 5)

	Input := ShortInput{TestUrl}
	for i, _ := range output {
		go MultMakeRequest(t, Input, c)
		output[i] = <-c

	}
	for i, _ := range output {
		output[i] = <-c
	}
	for i, _ := range output {
		assert.Equal(t, output[1], output[i])
	}

}
func TestCollisionPreveation(t *testing.T) {
	cleanTable()
	letterRunes = rune("abcd")
	N := 10
	output := make([]ShortOut, N)
	c := make(chan ShortOut)
	TestUrls := []string{}
	for i := 0; i < N; i++ {
		TestUrls = append(TestUrls, ("http://goolge.com/" + string(i)))
	}
	for i, _ := range TestUrls {
		Input := ShortInput{TestUrls[i]}
		go MultMakeRequest(t, Input, c)
	}
	for i, _ := range TestUrls {
		output[i] = <-c
	}
	for i, _ := range TestUrls {
		for j, _ := range TestUrls {
			if j != i {
				assert.NotEqual(t, output[i].ShortUrl, output[j].ShortUrl)
			}

		}
	}
}
