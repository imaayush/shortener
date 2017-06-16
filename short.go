
package short

import (
	"net/http"
	//"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"github.com/gorilla/mux"
)

func main() {

	//Migrate the schema
	db := Database(DevDb)
	db.AutoMigrate(&Short{})
	router := mux.NewRouter().StrictSlash(true)


	router.Handle("/short", ShortUrl(DevDb) )
	router.Handle("/{uuid}", ExpandUrl(DevDb))

	log.Fatal(http.ListenAndServe(":8080", nil))

}

