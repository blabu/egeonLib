package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/blabu/egeonLib/golang"
	"github.com/gin-gonic/gin"
)

type LogWriterWrapper func(string)

func (l LogWriterWrapper) WriteString(data string) (int, error) {
	l(data)
	return len(data), nil
}

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

func BuildRequestMiddleware(cache Model, log io.StringWriter, requestPerUser string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			log.WriteString("Method not cached")
			c.Next()
			return
		}
		if cache.storage == nil {
			log.WriteString("Cache is nil")
			c.Next()
			return
		}
		if user, ok := c.Request.Context().Value(golang.UserKey).(golang.User); ok {
			start := time.Now()
			key := fmt.Sprintf(requestPerUser, user.ID, hash(c.Request.RequestURI))
			if resp, err := cache.Get(c.Request.Context(), key); err == nil {
				for k, v := range resp.Header {
					for i := range v {
						log.WriteString(fmt.Sprintf("Add header %s: %s", k, v[i]))
						c.Writer.Header().Add(k, v[i])
					}
				}
				c.Writer.Header().Add("X-Cache", fmt.Sprintf("%d ms", time.Now().Sub(start).Milliseconds()))
				log.WriteString(fmt.Sprintf("Result finded in cache. Status: %d", resp.Status))
				c.Writer.WriteHeader(resp.Status)
				c.Writer.Write(resp.Body)
				c.Abort()
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
			return
		}
		log.WriteString("Undefined user work without cache")
		c.Next()
	}
}
