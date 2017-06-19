package main


import (
	"github.com/jinzhu/gorm"
	"math/rand"
	"github.com/gorilla/mux"
	"sync"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")
type App struct{
	Router *mux.Router
	DB     *gorm.DB
	sync.Mutex
}

var DbName = "/tmp/dev.db"
var TestDb = "/tmp/dev.db"


func GenerateShortUrl()string{
	n := 4
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)


}

func (app *App)CheckIsUnqiue(url string)bool{
	db := app.DB
	var short Short
	if err := db.Where("short_url = ?", string(url)).Find(&short).Error; err != nil {
		return true
	}else{
		return false
	}
}



func Database(DbName string) *gorm.DB {
	//open a db connection
	db, err := gorm.Open("sqlite3", DbName )
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
