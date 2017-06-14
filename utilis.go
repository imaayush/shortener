package main

import (
	"github.com/jinzhu/gorm"
	"math/rand"
	"fmt"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateShortUrl()string{
	n := 32
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Database() *gorm.DB {
	//open a db connection
	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
