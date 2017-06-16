package short

import (
	"net/http"
	"testing"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"fmt"
	"net/http/httptest"
)

func TestShort(t *testing.T){
	db := Database(TestDb)
	db.AutoMigrate(&Short{})
	u := ShortInput{"https://goolge.com"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	var data ShortOut
	res, _ := http.NewRequest("POST", "/short", b)
	w := httptest.NewRecorder()
	ShortUrl := ShortUrl(TestDb)
	ShortUrl.ServeHTTP(w, res)

	json.NewDecoder(w.Body).Decode(&data)
	assert.Equal(t, u.Url, data.Url)

	var short Short
	db.Where("url = ?", data.Url).Find(&short)
	fmt.Println(data)
	assert.Equal(t, short.ShortUrl, data.ShortUrl)
}

func TestExpand( t *testing.T){
	u := ShortInput{"https://www.goolge.com"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	var response ShortOut

	res, _ := http.NewRequest("POST", "/short", b)
	w := httptest.NewRecorder()
	ShortUrl := ShortUrl(TestDb)
	ShortUrl.ServeHTTP(w, res)

	assert.Equal(t, 200, w.Code)

	json.NewDecoder(w.Body).Decode(&response)
	fmt.Println(response.Url)
	url := "/" + string(response.ShortUrl)
	fmt.Println(url)
	res, _ = http.NewRequest("POST", "/short", b)
	w = httptest.NewRecorder()
	ExpandUrl := ExpandUrl(TestDb)
	ExpandUrl.ServeHTTP(w, res)
	fmt.Println(w)
	fmt.Println(res)
	assert.Equal(t, 301, w.Code)

}