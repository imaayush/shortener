package short

import (
	"github.com/jinzhu/gorm"
	"math/rand"
	"fmt"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var DevDb = "/tmp/dev.db"
var TestDb = "/tmp/testdb.db"


func GenerateShortUrl()string{
	n := 4
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	db := Database(DevDb)
	var short Short
	if err := db.Where("short_url = ?", string(b)).Find(&short).Error; err != nil {
		return string(b)
	}else{
		fmt.Println("collision")
		return GenerateShortUrl()
	}

}

func Database(DbName  string) *gorm.DB {
	//open a db connection
	db, err := gorm.Open("sqlite3", DbName )
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
