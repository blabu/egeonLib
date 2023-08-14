package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/blabu/egeonLib/golang"
	"github.com/gin-gonic/gin"
)

// writerWrap - need to store request body when you do request to cache it after success execution
type writerWrap struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (rw *writerWrap) Write(body []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(body)
	if err == nil {
		rw.body.Write(body)
	}
	return n, err
}

func hash(url string) string {
	reqURI := md5.Sum([]byte(url))
	return base64.StdEncoding.EncodeToString(reqURI[:])
}

// BuildRequestMiddleware - create middleware that cached requests by user
// requestPerUser - is a template key that is Sprintf template with to parameters %s and %s
// requestPerUser key forms with userID and md5 hash sum for request URI
func BuildRequestMiddleware(cache Model, log io.StringWriter, requestPerUser string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cache.storage == nil {
			log.WriteString("Cache is nil")
			c.Next()
			return
		}
		if c.Request.Method == http.MethodOptions ||
			c.Request.Method == http.MethodTrace ||
			c.Request.Method == http.MethodHead ||
			c.Request.Method == http.MethodConnect {
			log.WriteString("Methods not cached")
			c.Next()
			return
		}
		user, ok := c.Request.Context().Value(golang.UserKey).(golang.User)
		if !ok {
			log.WriteString("Undefined user work without cache")
			c.Next()
			return
		}
		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut ||
			c.Request.Method == http.MethodPatch ||
			c.Request.Method == http.MethodDelete {
			// clear cache for current user if we got some mutation request and it was success
			c.Next()
			if c.Writer.Status() < http.StatusBadRequest {
				cache.Delete(c.Request.Context(), fmt.Sprintf(requestPerUser, strconv.FormatUint(uint64(user.ID), 10), "*"))
			}
			return
		}
		key := fmt.Sprintf(requestPerUser, strconv.FormatUint(uint64(user.ID), 10), hash(c.Request.RequestURI))
		if resp, err := cache.Get(c.Request.Context(), key); err == nil {
			for k, v := range resp.Header {
				for i := range v {
					log.WriteString(fmt.Sprintf("Add header %s: %s", k, v[i]))
					c.Writer.Header().Add(k, v[i])
				}
			}
			log.WriteString(fmt.Sprintf("Result found in cache. Status: %d", resp.Status))
			c.Writer.WriteHeader(resp.Status)
			c.Writer.Write(resp.Body)
			c.Abort() //stop request execution
			return
		}
		log.WriteString("Undefined result in cache")
		rw := writerWrap{ResponseWriter: c.Writer}
		c.Writer = &rw
		c.Next()
		if rw.Status() < http.StatusBadRequest {
			log.WriteString("Result saved to cache")
			cache.Set(c.Request.Context(), key, &Responce{Body: rw.body.Bytes(), Status: rw.Status(), Header: map[string][]string(rw.Header())})
		}
	}
}
