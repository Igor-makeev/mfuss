package handler

import (
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	mrand "math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
	secretKey          = "secret key"
	userIDBytes        = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	userCtx            = "userID"
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

func UserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		cook, err := c.Cookie("UserID")
		if err != nil {
			cook = generateCook()
			c.SetCookie("UserID", cook, 3600*24, "/", "localhost", false, false)

		}
		if checkCook(cook) {
			c.Set(userCtx, cook)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

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

func generateCook() string {
	uuid := []byte(genetareUserID())

	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write(uuid)
	cook := hash.Sum(uuid)
	res := hex.EncodeToString(cook)
	return res
}

func checkCook(id string) bool {
	isvalid := false
	var (
		data []byte
		err  error
	)
	data, err = hex.DecodeString(id)

	if err != nil {

		return isvalid
	}
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(data[:5])
	sign := h.Sum(nil)
	if hmac.Equal(sign, data[5:]) {
		isvalid = true
		return isvalid
	}
	return isvalid
}

func getUserID(c *gin.Context) (string, error) {
	id, ok := c.Get(userCtx)

	if !ok {
		http.Error(c.Writer, "user not found", http.StatusInternalServerError)
		return "", errors.New("user id not found")
	}

	idstring, ok := id.(string)
	if !ok {
		http.Error(c.Writer, "user id is of ivalid type", http.StatusInternalServerError)
		return "", errors.New("user id is of ivalid type")
	}

	return idstring, nil
}

func genetareUserID() string {
	buf := make([]byte, 5)
	for i := range buf {
		buf[i] = userIDBytes[mrand.Intn(len(userIDBytes))]
	}
	res := string(buf)
	return res
}
