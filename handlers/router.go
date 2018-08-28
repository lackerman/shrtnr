package handlers

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/lackerman/shrtnr/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

// NewRouter sets up all the app paths
func NewRouter(db *leveldb.DB) *mux.Router {
	mux := mux.NewRouter()
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/", fs)
	mux.Handle("/url", Post(Creater(db)))
	mux.Handle("/all", Get(Lister(db)))
	mux.Handle("/{key}", Get(Retriever(db)))
	return mux
}

// Request is the struct used for http Request representation
type Request struct {
	reqType string
	h       http.HandlerFunc
}

func (r *Request) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("%v :: %v :: %v", req.Method, req.URL, reflect.ValueOf(r.h))
	switch {
	case r.reqType == req.Method:
		r.h(w, req)
		return
	}
	http.Error(w, "Unsupport Method Type for request", http.StatusUnsupportedMediaType)
}

// Post is part of the DSL for limiting requests to POST method types
func Post(h http.HandlerFunc) *Request {
	return &Request{http.MethodPost, h}
}

// Get is part of the DSL for limiting requests to POST method types
func Get(h http.HandlerFunc) *Request {
	return &Request{http.MethodGet, h}
}

// Creater sets up the handler used for creating URLs
func Creater(db *leveldb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()

		url := req.Form.Get("url")
		encoded := utils.EncodeURL(url)

		go func(encoded string, url string) {
			db.Put([]byte(encoded), []byte(url), nil)
		}(encoded, url)

		w.Write([]byte(encoded))
	}
}

// Retriever sets up the handler for retrieving the URLs
func Retriever(db *leveldb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		w.Write([]byte(fmt.Sprintf("%v, %+v", "Retriever", vars)))
	}
}

// Lister sets up a handler for getting the full list of URLs
func Lister(db *leveldb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Lister"))
	}
}
