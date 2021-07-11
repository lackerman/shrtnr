package handlers

import (
	"fmt"
	"net/http"

	"github.com/lackerman/shrtnr/tracing"
	"github.com/lackerman/shrtnr/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/syndtr/goleveldb/leveldb"
)

type handler struct {
	db  *leveldb.DB
	log logr.Logger
}

func NewRouter(db *leveldb.DB, l logr.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.LoadHTMLGlob("templates/*")
	r.Use(tracing.Middleware())
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	h := handler{db: db, log: l}
	r.GET("/", h.home)
	r.GET("/spantest", h.spanTest)
	r.POST("/url", h.creater)
	r.Any("/all", h.lister)
	r.POST("/edit", h.edit)
	r.GET("/u/:key", h.retriever)
	return r
}

func (h *handler) home(c *gin.Context) {
	ctx, span := tracing.NewSpan(h.log.V(0).Info, c.Request.Context(), c.HandlerName())
	defer span.End()

	res, err := tracing.HttpRequest(ctx, http.MethodGet, "http://localhost:8080/spantest", nil)
	if err != nil {
		panic(err)
	}
	h.log.V(0).Info("Received a response", "status_code", res.StatusCode)

	c.HTML(http.StatusOK, "index.tmpl", map[string]string{
		"Title":   "Shrtnr",
		"Heading": "Paste your URL below",
	})
}

// Creater sets up the handler used for creating URLs
func (h *handler) creater(c *gin.Context) {
	_, span := tracing.NewSpan(h.log.V(0).Info, c.Request.Context(), c.HandlerName())
	defer span.End()

	c.Request.ParseForm()

	url := c.Request.Form.Get("url")
	encoded, err := utils.EncodeURL(url)
	if err != nil {
		panic(err)
	}

	go func(encoded string, url string) {
		h.db.Put([]byte(encoded), []byte(url), nil)
	}(encoded, url)

	c.HTML(http.StatusOK, "url.tmpl", map[string]string{
		"Title":   "Shrtnr",
		"Heading": "Copy your new shrtnd URL",
		"URL":     fmt.Sprintf("http://%s/u/%s", c.Request.Host, encoded),
	})
}

// Retriever sets up the handler for retrieving the URLs
func (h *handler) retriever(c *gin.Context) {
	_, span := tracing.NewSpan(h.log.V(0).Info, c.Request.Context(), c.HandlerName())
	defer span.End()

	key := c.Query("key")
	url, err := h.db.Get([]byte(key), nil)
	if err != nil {
		panic(err)
	}
	c.Redirect(http.StatusTemporaryRedirect, string(url))
}

// Lister sets up a handler for getting the full list of URLs
func (h *handler) lister(c *gin.Context) {
	_, span := tracing.NewSpan(h.log.V(0).Info, c.Request.Context(), c.HandlerName())
	defer span.End()

	urls := map[string]string{}
	iter := h.db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		urls[string(iter.Key())] = string(iter.Value())
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "all.tmpl", map[string]interface{}{
		"Title":   "Shrtnr",
		"Heading": "Select the URL you want to shorten",
		"URLs":    urls,
	})
}

// edit sets up a handler for editing a selection of urls
func (h *handler) edit(c *gin.Context) {
	_, span := tracing.NewSpan(h.log.V(0).Info, c.Request.Context(), c.HandlerName())
	defer span.End()

	url := c.PostForm("url")
	go func(url string) {
		h.db.Delete([]byte(url), nil)
	}(url)

	c.Redirect(http.StatusTemporaryRedirect, string("/all"))
}

func (h *handler) spanTest(c *gin.Context) {
	_, span := tracing.NewSpan(h.log.V(0).Info, c.Request.Context(), c.HandlerName())
	defer span.End()

	c.JSON(http.StatusOK, "{\"hello\":\"there\"}")
}
