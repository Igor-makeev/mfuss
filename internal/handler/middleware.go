package handler

import (
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	urlIDSliceCtx      = "input_id"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

// GzipCompress — мидлваре архивирующий ответ в формате gzip.
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

// GzipUnpack — мидлваре разархивирующий ответ в формате gzip.
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

// userCheck — мидлваре проверяющий прова доступа у пользователя.
func (h *Handler) userCheck(c *gin.Context) {

	cook, err := c.Cookie(userCtx)

	if err != nil {

		cook = generateCook()
		c.SetCookie(userCtx, cook, 3600*24, "/", "localhost", false, false)

	}
	if checkCook(cook) {
		c.Set(userCtx, cook)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.Next()
}

// checkURLSID — мидлваре проверяющий валидность переданных на удаление ID.
func (h *Handler) checkURLSID(c *gin.Context) {

	var input []string

	err := json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)

		return
	}

	for _, id := range input {
		if len([]rune(id)) != 5 {
			http.Error(c.Writer, fmt.Sprintf("invalid url id: %v", id), http.StatusBadRequest)

			return
		}
	}
	c.Set(urlIDSliceCtx, input)

	c.Next()

}

// shouldCompress — функция проверяющая загловок.
func shouldCompress(req *http.Request) bool {
	return strings.Contains(req.Header.Get("Accept-Encoding"), "gzip")

}

// shouldUnpack — функция проверяющая загловок.
func shouldUnpack(req *http.Request) bool {
	return strings.Contains(req.Header.Get("Content-Encoding"), "gzip")

}

// generateCook — функция генерирующая Cook.
func generateCook() string {
	uuid := []byte(genetareUserID())

	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write(uuid)
	cook := hash.Sum(uuid)
	res := hex.EncodeToString(cook)
	return res
}

// checkCook — функция проверяющая Cook.
func checkCook(cook string) bool {

	var (
		data []byte
		err  error
	)
	data, err = hex.DecodeString(cook)

	if err != nil {

		return false
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(data[:5])
	sign := h.Sum(nil)

	return hmac.Equal(sign, data[5:])
}

// getUserID — функция плучающая id  пользователя из контекста.
func getUserID(c *gin.Context) (string, error) {
	id, ok := c.Get(userCtx)

	if !ok {

		return "", ErrNoUserID
	}

	idstring, ok := id.(string)
	if !ok {

		return "", ErrInvalidUserID
	}

	return idstring, nil
}

// getUrlsArray — функция плучающая массив url идентификаторов пользователя из контекста.
func getUrlsArray(c *gin.Context) ([]string, error) {
	ids, ok := c.Get(urlIDSliceCtx)

	if !ok {

		return nil, ErrNoDataArray
	}

	idsArray, ok := ids.([]string)
	if !ok {

		return nil, ErrInvalidDataArray
	}

	return idsArray, nil
}

// genetareUserID — функция генерирующая идентификатор пользователя.
func genetareUserID() string {
	buf := make([]byte, 5)
	for i := range buf {
		buf[i] = userIDBytes[mrand.Intn(len(userIDBytes))]
	}

	return string(buf)
}
