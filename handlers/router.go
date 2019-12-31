package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
	"gitlab.com/lackerman/shrtnr/utils"
)

const handlerKey = "handler"

// NewServer creates a specific server
func NewServer(port int, t *template.Template, db *leveldb.DB, l *log.Logger) http.Server {
	return http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      newRouter(t, db, l),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func newRouter(t *template.Template, db *leveldb.DB, l *log.Logger) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.Use(middleware(l))
	e.GET("/", wrapper(home(t), l))
	e.POST("/url", wrapper(creater(t, db), l))
	e.Any("/all", wrapper(lister(t, db), l))
	e.POST("/edit", wrapper(edit(db), l))
	e.GET("/u/:key", wrapper(retriever(db), l))
	return e
}

func middleware(l *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		defer func() {
			handler, exists := c.Get(handlerKey)
			if !exists {
				handler = c.HandlerName
			}
			if err := recover(); err != nil {
				l.Printf("%+v", err)
				c.String(500, "Failed to process the request: Error: %v", err)
			}
			l.Printf("Time taken for '%v': %v", handler, time.Since(start))
		}()

		c.Next()
	}
}

func wrapper(h handler, l *log.Logger) gin.HandlerFunc {
	fnName := utils.FuncName(h)
	return func(c *gin.Context) {
		c.Set(handlerKey, fnName)
		err := h(c.Writer, c.Request)
		if err != nil {
			panic(err)
		}
	}
}

type handler func(http.ResponseWriter, *http.Request) error

func home(t *template.Template) handler {
	return func(w http.ResponseWriter, req *http.Request) error {
		return t.ExecuteTemplate(w, "index.tmpl", map[string]string{
			"Title":   "Shrtnr",
			"Heading": "Paste your URL below",
		})
	}
}

// Creater sets up the handler used for creating URLs
func creater(t *template.Template, db *leveldb.DB) handler {
	return func(w http.ResponseWriter, req *http.Request) error {
		req.ParseForm()

		url := req.Form.Get("url")
		encoded, err := utils.EncodeURL(url)
		if err != nil {
			return err
		}

		go func(encoded string, url string) {
			db.Put([]byte(encoded), []byte(url), nil)
		}(encoded, url)

		return t.ExecuteTemplate(w, "url.tmpl", map[string]string{
			"Title":   "Shrtnr",
			"Heading": "Copy your new shrtnd URL",
			"URL":     fmt.Sprintf("http://%s/u/%s", req.Host, encoded),
		})
	}
}

// Retriever sets up the handler for retrieving the URLs
func retriever(db *leveldb.DB) handler {
	return func(w http.ResponseWriter, req *http.Request) error {
		parts := strings.Split(req.RequestURI, "/")
		key := parts[len(parts)-1]
		fmt.Println(key)
		url, err := db.Get([]byte(key), nil)
		if err != nil {
			return err
		}
		http.Redirect(w, req, string(url), http.StatusTemporaryRedirect)
		return nil
	}
}

// Lister sets up a handler for getting the full list of URLs
func lister(t *template.Template, db *leveldb.DB) handler {
	return func(w http.ResponseWriter, req *http.Request) error {
		urls := map[string]string{}

		iter := db.NewIterator(nil, nil)
		for iter.Next() {
			// Remember that the contents of the returned slice should not be modified, and
			// only valid until the next call to Next.
			urls[string(iter.Key())] = string(iter.Value())
		}
		iter.Release()
		err := iter.Error()
		if err != nil {
			return err
		}

		return t.ExecuteTemplate(w, "all.tmpl", map[string]interface{}{
			"Title":   "Shrtnr",
			"Heading": "Select the URL you want to shorten",
			"URLs":    urls,
		})
	}
}

// edit sets up a handler for editing a selection of urls
func edit(db *leveldb.DB) handler {
	return func(w http.ResponseWriter, req *http.Request) error {
		req.ParseForm()

		for k := range req.Form {
			go func(k string) {
				db.Delete([]byte(k), nil)
			}(k)
		}

		http.Redirect(w, req, "/all", http.StatusTemporaryRedirect)
		return nil
	}
}
