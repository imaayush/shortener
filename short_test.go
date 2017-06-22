package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var ts *httptest.Server
var app App

func TestMain(m *testing.M) {

	app = App{}
	app.Initialize()
	app.SlagLength = 1
	ts = httptest.NewServer(app.Router)

	ret := m.Run()
	SlagLength = 1
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
func TestShortUrlEndPointPassCase(t *testing.T) {
	cleanTable()
	TestUrl := "https://goolge.com/home/param=12"
	Input := ShortInput{TestUrl}
	db := app.DB
	var data ShortOut
	var short Short
	resp, _ := MakePostRequest(t, Input, &data)
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
	MakePostRequest(t, Input, &data)
	assert.Equal(t, Input.Url, data.Url)
	url := ts.URL + "/" + data.ShortUrl
	Input = ShortInput{url}
	resp := MakeGetRequest(t, Input)
	assert.Equal(t, resp.StatusCode, 301)

}

func TestWrongInput(t *testing.T) {
	TestCase := "google"

	var data ShortOut
	Input := ShortInput{TestCase}
	resp, _ := MakePostRequest(t, Input, &data)
	fmt.Println(resp)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestShortUrlNotFound(t *testing.T) {
	url := ts.URL + "/" + "ASDFW"
	Input := ShortInput{url}
	resp := MakeGetRequest(t, Input)
	assert.Equal(t, resp.StatusCode, 404)
}

func TestResponseOnCollision(t *testing.T) {
	cleanTable()
	TestUrl := "https://goolge.com/"

	var data ShortOut
	var respon *http.Response
	for i := 0; i < 67; i++ {
		url := TestUrl + strconv.Itoa(i)
		Input := ShortInput{string(url)}
		resp, _ := MakePostRequest(t, Input, &data)
		json.NewDecoder(resp.Body).Decode(&data)
		respon = resp
		if resp.StatusCode == 500 {
			respon = resp

		}
	}

	assert.Equal(t, respon.StatusCode, 500)
}

func TestCheckShorUrlIsAddedToDB(t *testing.T) {
	cleanTable()
	TestUrl := "https://goolge.com/home/param=11"
	Input := ShortInput{TestUrl}

	var data ShortOut

	resp, _ := MakePostRequest(t, Input, &data)
	assert.Equal(t, resp.StatusCode, 200)

	assert.Equal(t, app.CheckIsOnDb(data.ShortUrl), true)
}

func TestSuccessfulSaveUnqiue(t *testing.T) {
	cleanTable()
	ShortUrl := app.GenerateShortUrl()
	fmt.Println(ShortUrl)

	assert.Equal(t, app.CheckIsOnDb(ShortUrl), false)

	LongUrl := "https://google.com"
	resp, err := app.SaveUrl(ShortUrl, LongUrl)

	assert.Equal(t, resp.ShortUrl, ShortUrl)
	assert.Equal(t, resp.Url, LongUrl)
	assert.Equal(t, err, nil)
}
func TestFailTOSaveDuplicate(t *testing.T) {
	cleanTable()
	ShortUrl := app.GenerateShortUrl()

	assert.Equal(t, app.CheckIsOnDb(ShortUrl), false)

	LongUrl := "https://google.com"
	resp, err := app.SaveUrl(ShortUrl, LongUrl)

	assert.Equal(t, resp.ShortUrl, ShortUrl)
	assert.Equal(t, resp.Url, LongUrl)
	assert.Equal(t, err, nil)
	//save dupliact
	LongUrl = "https://facebook.com"
	resp, err = app.SaveUrl(ShortUrl, LongUrl)
	exceptedErr := errors.New("pq: duplicate key value violates unique constraint \"shorts_short_url_key\"")
	assert.Equal(t, err, exceptedErr)

}
