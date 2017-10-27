package gzip

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

func Gzip(level int) gin.HandlerFunc {
	var gzPool sync.Pool
	gzPool.New = func() interface{} {
		gz, err := gzip.NewWriterLevel(ioutil.Discard, level)
		if err != nil {
			panic(err)
		}
		return gz
	}
	return func(c *gin.Context) {
		if !shouldCompress(c.Request) {
			return
		}

		if !shouldCompressResponse(c.Writer) {
			return
		}

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)
		gz.Reset(c.Writer)

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")
		c.Writer = &gzipWriter{c.Writer, gz}
		defer func() {
			c.Header("Content-Length", "0")
			gz.Close()
		}()
		c.Next()
	}
}

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

// Add by JianjunXie: filter types for compression
var gzip_types = map[string]bool{
	"text/css":                      true,
	"text/plain":                    true,
	"text/javascript":               true,
	"application/javascript":        true,
	"application/json":              true,
	"application/x-javascript":      true,
	"application/xml":               true,
	"application/xml+rss":           true,
	"application/xhtml+xml":         true,
	"application/x-font-ttf":        true,
	"application/x-font-opentype":   true,
	"application/vnd.ms-fontobject": true,
	"image/svg+xml":                 true,
	"image/x-icon":                  true,
	"application/rss+xml":           true,
	"application/atom_xml":          true,
}

func shouldCompressResponse(w http.ResponseWriter) bool {
	headers := w.Header()
	if t, ok := gzip_types[headers.Get("Content-Type")]; ok && t {
		return true
	}
	return false
}

func shouldCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		return false
	}
	extension := filepath.Ext(req.URL.Path)
	if len(extension) < 4 { // fast path
		return true
	}

	switch extension {
	case ".png", ".gif", ".jpeg", ".jpg":
		return false
	default:
		return true
	}
}
