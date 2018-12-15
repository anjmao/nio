package nio

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-nio/nio/log"

	"github.com/stretchr/testify/assert"
)

type (
	user struct {
		ID   int    `json:"id" xml:"id" form:"id" query:"id"`
		Name string `json:"name" xml:"name" form:"name" query:"name"`
	}
)

const (
	userJSON                    = `{"id":1,"name":"Jon Snow"}`
	userXML                     = `<user><id>1</id><name>Jon Snow</name></user>`
	userForm                    = `id=1&name=Jon Snow`
	invalidContent              = "invalid content"
	userJSONInvalidType         = `{"id":"1","name":"Jon Snow"}`
	userXMLConvertNumberError   = `<user><id>Number one</id><name>Jon Snow</name></user>`
	userXMLUnsupportedTypeError = `<user><>Number one</><name>Jon Snow</name></user>`
)

const userJSONPretty = `{
  "id": 1,
  "name": "Jon Snow"
}`

const userXMLPretty = `<user>
  <id>1</id>
  <name>Jon Snow</name>
</user>`

func TestNio(t *testing.T) {
	e := New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// DefaultHTTPErrorHandler
	e.DefaultHTTPErrorHandler(errors.New("error"), c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestNioStatic(t *testing.T) {
	e := New()

	assert := assert.New(t)

	// OK
	e.Static("/images", "_fixture/images")
	c, b := request(http.MethodGet, "/images/walle.png", e)
	assert.Equal(http.StatusOK, c)
	assert.NotEmpty(b)

	// No file
	e.Static("/images", "_fixture/scripts")
	c, _ = request(http.MethodGet, "/images/bolt.png", e)
	assert.Equal(http.StatusNotFound, c)

	// Directory
	e.Static("/images", "_fixture/images")
	c, _ = request(http.MethodGet, "/images", e)
	assert.Equal(http.StatusNotFound, c)

	// Directory with index.html
	e.Static("/", "_fixture")
	c, r := request(http.MethodGet, "/", e)
	assert.Equal(http.StatusOK, c)
	assert.Equal(true, strings.HasPrefix(r, "<!doctype html>"))

	// Sub-directory with index.html
	c, r = request(http.MethodGet, "/folder", e)
	assert.Equal(http.StatusOK, c)
	assert.Equal(true, strings.HasPrefix(r, "<!doctype html>"))
}

func TestNioLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	e := New(SetLogger(log.NewLogger(ioutil.Discard, ioutil.Discard, buf)))

	e.Logger().Error("ups")

	assert.Contains(t, buf.String(), "ups")
}

func TestNioFile(t *testing.T) {
	e := New()
	e.File("/walle", "_fixture/images/walle.png")
	c, b := request(http.MethodGet, "/walle", e)
	assert.Equal(t, http.StatusOK, c)
	assert.NotEmpty(t, b)
}

func TestNioMiddleware(t *testing.T) {
	e := New()
	buf := new(bytes.Buffer)

	e.Pre(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			assert.Empty(t, c.Path())
			buf.WriteString("-1")
			return next(c)
		}
	})

	e.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("1")
			return next(c)
		}
	})

	e.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("2")
			return next(c)
		}
	})

	e.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("3")
			return next(c)
		}
	})

	// Route
	e.GET("/", func(c Context) error {
		return c.String(http.StatusOK, "OK")
	})

	c, b := request(http.MethodGet, "/", e)
	assert.Equal(t, "-1123", buf.String())
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)
}

func TestNioMiddlewareError(t *testing.T) {
	e := New()
	e.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			return errors.New("error")
		}
	})
	e.GET("/", NotFoundHandler)
	c, _ := request(http.MethodGet, "/", e)
	assert.Equal(t, http.StatusInternalServerError, c)
}

func TestNioHandler(t *testing.T) {
	e := New()

	// HandlerFunc
	e.GET("/ok", func(c Context) error {
		return c.String(http.StatusOK, "OK")
	})

	c, b := request(http.MethodGet, "/ok", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)
}

func TestNioWrapHandler(t *testing.T) {
	e := New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}))
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
	}
}

func TestNioWrapMiddleware(t *testing.T) {
	e := New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	buf := new(bytes.Buffer)
	mw := WrapMiddleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf.Write([]byte("mw"))
			h.ServeHTTP(w, r)
		})
	})
	h := mw(func(c Context) error {
		return c.String(http.StatusOK, "OK")
	})
	if assert.NoError(t, h(c)) {
		assert.Equal(t, "mw", buf.String())
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "OK", rec.Body.String())
	}
}

func TestNioConnect(t *testing.T) {
	e := New()
	testMethod(t, http.MethodConnect, "/", e)
}

func TestNioDelete(t *testing.T) {
	e := New()
	testMethod(t, http.MethodDelete, "/", e)
}

func TestNioGet(t *testing.T) {
	e := New()
	testMethod(t, http.MethodGet, "/", e)
}

func TestNioHead(t *testing.T) {
	e := New()
	testMethod(t, http.MethodHead, "/", e)
}

func TestNioOptions(t *testing.T) {
	e := New()
	testMethod(t, http.MethodOptions, "/", e)
}

func TestNioPatch(t *testing.T) {
	e := New()
	testMethod(t, http.MethodPatch, "/", e)
}

func TestNioPost(t *testing.T) {
	e := New()
	testMethod(t, http.MethodPost, "/", e)
}

func TestNioPut(t *testing.T) {
	e := New()
	testMethod(t, http.MethodPut, "/", e)
}

func TestNioTrace(t *testing.T) {
	e := New()
	testMethod(t, http.MethodTrace, "/", e)
}

func TestNioAny(t *testing.T) { // JFC
	e := New()
	e.Any("/", func(c Context) error {
		return c.String(http.StatusOK, "Any")
	})
}

func TestNioMatch(t *testing.T) { // JFC
	e := New()
	e.Match([]string{http.MethodGet, http.MethodPost}, "/", func(c Context) error {
		return c.String(http.StatusOK, "Match")
	})
}

func TestNioURL(t *testing.T) {
	e := New()
	static := func(Context) error { return nil }
	getUser := func(Context) error { return nil }
	getFile := func(Context) error { return nil }

	e.GET("/static/file", static)
	e.GET("/users/:id", getUser)
	g := e.Group("/group")
	g.GET("/users/:uid/files/:fid", getFile)

	assert := assert.New(t)

	assert.Equal("/static/file", e.URL(static))
	assert.Equal("/users/:id", e.URL(getUser))
	assert.Equal("/users/1", e.URL(getUser, "1"))
	assert.Equal("/group/users/1/files/:fid", e.URL(getFile, "1"))
	assert.Equal("/group/users/1/files/1", e.URL(getFile, "1", "1"))
}

func TestNioRoutes(t *testing.T) {
	e := New()
	routes := []*Route{
		{http.MethodGet, "/users/:user/events", ""},
		{http.MethodGet, "/users/:user/events/public", ""},
		{http.MethodPost, "/repos/:owner/:repo/git/refs", ""},
		{http.MethodPost, "/repos/:owner/:repo/git/tags", ""},
	}
	for _, r := range routes {
		e.Add(r.Method, r.Path, func(c Context) error {
			return c.String(http.StatusOK, "OK")
		})
	}

	if assert.Equal(t, len(routes), len(e.Routes())) {
		for _, r := range e.Routes() {
			found := false
			for _, rr := range routes {
				if r.Method == rr.Method && r.Path == rr.Path {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Route %s %s not found", r.Method, r.Path)
			}
		}
	}
}

func TestNioEncodedPath(t *testing.T) {
	e := New()
	e.GET("/:id", func(c Context) error {
		return c.NoContent(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/with%2Fslash", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestNioGroup(t *testing.T) {
	e := New()
	buf := new(bytes.Buffer)
	e.Use(MiddlewareFunc(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("0")
			return next(c)
		}
	}))
	h := func(c Context) error {
		return c.NoContent(http.StatusOK)
	}

	//--------
	// Routes
	//--------

	e.GET("/users", h)

	// Group
	g1 := e.Group("/group1")
	g1.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("1")
			return next(c)
		}
	})
	g1.GET("", h)

	// Nested groups with middleware
	g2 := e.Group("/group2")
	g2.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("2")
			return next(c)
		}
	})
	g3 := g2.Group("/group3")
	g3.Use(func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			buf.WriteString("3")
			return next(c)
		}
	})
	g3.GET("", h)

	request(http.MethodGet, "/users", e)
	assert.Equal(t, "0", buf.String())

	buf.Reset()
	request(http.MethodGet, "/group1", e)
	assert.Equal(t, "01", buf.String())

	buf.Reset()
	request(http.MethodGet, "/group2/group3", e)
	assert.Equal(t, "023", buf.String())
}

func TestNioNotFound(t *testing.T) {
	e := New()
	req := httptest.NewRequest(http.MethodGet, "/files", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestNioMethodNotAllowed(t *testing.T) {
	e := New()
	e.GET("/", func(c Context) error {
		return c.String(http.StatusOK, "Nio!")
	})
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestNioContext(t *testing.T) {
	e := New()
	c := e.pool.Get().(*context)
	assert.IsType(t, new(context), c)
	e.pool.Put(c)
}

func TestNioStart(t *testing.T) {
	e := New()
	go func() {
		err := http.ListenAndServe(":0", e)
		assert.NoError(t, err)
	}()
	time.Sleep(200 * time.Millisecond)
}

func TestNioStartTLS(t *testing.T) {
	e := New()
	go func() {
		err := http.ListenAndServeTLS(":0", "_fixture/certs/cert.pem", "_fixture/certs/key.pem", e)
		// Prevent the test to fail after closing the servers
		if err != http.ErrServerClosed {
			assert.NoError(t, err)
		}
	}()
	time.Sleep(200 * time.Millisecond)
}

func testMethod(t *testing.T, method, path string, e *Nio) {
	p := reflect.ValueOf(path)
	h := reflect.ValueOf(func(c Context) error {
		return c.String(http.StatusOK, method)
	})
	i := interface{}(e)
	reflect.ValueOf(i).MethodByName(method).Call([]reflect.Value{p, h})
	_, body := request(method, path, e)
	assert.Equal(t, method, body)
}

func request(method, path string, e *Nio) (int, string) {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func TestHTTPError(t *testing.T) {
	err := NewHTTPError(http.StatusBadRequest, map[string]interface{}{
		"code": 12,
	})
	assert.Equal(t, "code=400, message=map[code:12]", err.Error())
}
