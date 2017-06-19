package short

import (
	"net/http"
	"testing"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"os"
)

var ts *httptest.Server

func TestMain(m *testing.M){

	routes := NewRouter()
	ts = httptest.NewServer(routes)

	ret := m.Run()
	os.Exit(ret)
}

func TestShortPass(t *testing.T){
	TestCase := "https://goolge.com/home/param=11"
	db := Database()
	var data ShortOut
	var short Short

	u := ShortInput{TestCase}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	req, _ := http.NewRequest("POST", ts.URL + "/short", b)
	client := &http.Client{}
	resp, _ := client.Do(req)
	json.NewDecoder(resp.Body).Decode(&data)
	assert.Equal(t, u.Url, data.Url)
	assert.Equal(t, resp.StatusCode, 200)
	db.Where("short_url = ?", data.ShortUrl).Find(&short)

	assert.Equal(t, short.ShortUrl, data.ShortUrl)

}


func TestExpandPass( t *testing.T){

	TestCase := "https://goolge.com/"
	var data ShortOut

	u := ShortInput{TestCase}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	req, _ := http.NewRequest("POST", ts.URL + "/short", b)
	client := &http.Client{}
	resp, _ := client.Do(req)

	json.NewDecoder(resp.Body).Decode(&data)

	assert.Equal(t, u.Url, data.Url)
	url := ts.URL + "/" + data.ShortUrl

	req, _ = http.NewRequest("GET", url, nil)
	client = &http.Client{}
	resp, _ = client.Do(req)

	json.NewDecoder(resp.Body).Decode(&data)
	assert.Equal(t, resp.StatusCode, 200)

}


func TestShortFail(t *testing.T) {
	TestCase := "google"

	var data ShortOut
	u := ShortInput{TestCase}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	req, _ := http.NewRequest("POST", ts.URL + "/short", b)
	client := &http.Client{}
	resp, _ := client.Do(req)
	json.NewDecoder(resp.Body).Decode(&data)
	assert.Equal(t, resp.StatusCode, 402)
}