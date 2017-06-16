package main

import (
	"net/http"
	"testing"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestShort(t *testing.T){

	u := ShortInput{"gmail.com"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	var data ShortOut

	res, _ := http.Post("http://localhost:8080/short", "application/json; charset=utf-8", b)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	json.NewDecoder(res.Body).Decode(&data)
	//assert.Equal(t, u.Url, data.Url)
	db := Database(DevDb)
	var short Short
	db.Where("url = ?", data.Url).Find(&short)
	assert.Equal(t, short.ShortUrl, data.ShortUrl)
}

func TestExpand( t *testing.T){
	u := ShortInput{"fb.com"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	var response ShortOut

	res, _ := http.Post("http://localhost:8080/short", "application/json; charset=utf-8", b)

	assert.Equal(t, 200, res.StatusCode)

	json.NewDecoder(res.Body).Decode(&response)
	fmt.Println(response.Url)
	url := "http://localhost:8080/" + string(response.ShortUrl)
	fmt.Println(url)
	resp, _ := http.Get(url)
	assert.Equal(t, 200, resp.StatusCode)
	
}