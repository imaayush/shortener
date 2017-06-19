
package short

import (
	"net/http"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

func main() {

	//Migrate the schema
	db := Database()
	db.AutoMigrate(&Short{})
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))

}

