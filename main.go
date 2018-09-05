package main

import (
	"flag"
	"html/template"
	"log"
	"os"

	"github.com/lackerman/shrtnr/handlers"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	port := flag.Int("port", 8080, "Specify the port to use for the server.")

	db, err := leveldb.OpenFile("bin/data", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	templates, err := template.ParseGlob("templates/*")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdin, "", log.LstdFlags)
	server := handlers.NewServer(*port, templates, db, logger)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
