package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"errors"
)


func TestGenerateShortUrl(t *testing.T) {
	cleanTable()
	ShortUrl := app.GenerateShortUrl()

	assert.Equal(t, len(ShortUrl), 4)

}
func TestSaveUrlAndCheckUniquie(t *testing.T) {
	cleanTable()
	ShortUrl := app.GenerateShortUrl()

	LongUrl := "https://google.com"
	resp, err := app.SaveUrlAndCheckUniquie(ShortUrl, LongUrl)

	assert.Equal(t, resp.ShortUrl, ShortUrl)
	assert.Equal(t, resp.Url, LongUrl)
	assert.Equal(t, err, nil)
	//save dupliact
	LongUrl = "https://facebook.com"
	resp, err = app.SaveUrlAndCheckUniquie(ShortUrl, LongUrl)
	exceptedErr := errors.New("pq: duplicate key value violates unique constraint \"shorts_short_url_key\"")
	assert.Equal(t, err, exceptedErr)

}
