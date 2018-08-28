package main

import (
	"fmt"
	"net/http"

	"github.com/lackerman/shrtnr/handlers"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	fmt.Println("Starting server")
	db, err := leveldb.OpenFile("bin/data", nil)
	if err != nil {
		panic("Failed to create a file handle for the database")
	}
	defer db.Close()

	http.ListenAndServe(":3000", handlers.NewRouter(db))
}
