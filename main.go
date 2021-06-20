package main

import (
	"flag"
	"html/template"

	"github.com/lackerman/shrtnr/handlers"
	"github.com/syndtr/goleveldb/leveldb"

	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
)

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()

	port := flag.Int("port", 8080, "Specify the port to use for the server.")
	flag.Parse()

	logr := klogr.New()

	logr.V(0).Info("starting the database")
	db, err := leveldb.OpenFile("bin/data", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	logr.V(0).Info("loading the templates")
	templates, err := template.ParseGlob("templates/*")
	if err != nil {
		panic(err)
	}

	logr.V(0).Info("starting the server")
	server := handlers.NewServer(*port, templates, db, logr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
