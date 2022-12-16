package handler

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func GzipCompress(level int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldCompress(c.Request) {
			return
		}
		gz, err := gzip.NewWriterLevel(c.Writer, level)
		if err != nil {
			return
		}

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")
		c.Writer = &gzipWriter{c.Writer, gz}
		defer func() {
			gz.Close()
		}()
		c.Next()
	}
}

func GzipUnpack() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldUnpack(c.Request) {
			return
		}
		gz, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		defer gz.Close()

		c.Request.Body = gz
		c.Next()
	}
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func shouldCompress(req *http.Request) bool {
	return strings.Contains(req.Header.Get("Accept-Encoding"), "gzip")

}

func shouldUnpack(req *http.Request) bool {
	return strings.Contains(req.Header.Get("Content-Encoding"), "gzip")

}
