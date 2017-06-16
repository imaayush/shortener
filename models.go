package short

import "github.com/jinzhu/gorm"

type Short struct {
gorm.Model
	Url       string   `sql:"not null`
	ShortUrl  string   `sql:"not null;unique"`
}

type ShortInput struct {
	Url       string   `sql:"not null`
}

type ShortOut struct{
	Url       string   `sql:"not null`
	ShortUrl  string   `sql:"not null;unique"`
}

