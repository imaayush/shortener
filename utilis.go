package short

import (
	"github.com/jinzhu/gorm"
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var DbName = "/tmp/dev.db"
//var DbName = "/tmp/testdb.db"


func GenerateShortUrl()string{
	n := 4
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	db := Database()
	var short Short
	if err := db.Where("short_url = ?", string(b)).Find(&short).Error; err != nil {
		return string(b)
	}else{
		return GenerateShortUrl()
	}

}

func Database() *gorm.DB {
	//open a db connection
	db, err := gorm.Open("sqlite3", DbName )
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
