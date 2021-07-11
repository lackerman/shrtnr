package main

import (
	"flag"
	"fmt"

	"github.com/lackerman/shrtnr/handlers"
	"github.com/lackerman/shrtnr/tracing"

	"github.com/syndtr/goleveldb/leveldb"

	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
)

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()

	port := flag.Int("port", 8080, "Specify the port to use for the server.")
	zipkin := flag.String("zipkin", "localhost:9411", "Specify the zipkin server url and port")
	flag.Parse()

	logr := klogr.New()

	logr.V(0).Info("starting the database")
	db, err := leveldb.OpenFile("bin/data", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	logr.V(0).Info("starting the tracer")
	shutdown, err := tracing.InitTracer(fmt.Sprintf("http://%s/api/v2/spans", *zipkin))
	if err != nil {
		panic(err)
	}
	defer shutdown()

	listenOn := fmt.Sprintf(":%v", *port)
	logr.V(0).Info("starting the server", "listening_on", listenOn)
	router := handlers.NewRouter(db, logr)
	if err := router.Run(listenOn); err != nil {
		panic(err)
	}
}
