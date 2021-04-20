package golang

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// ContextMiddleware - wrap current request context to context with cancel
// Не зависит от сервиса и будет переиспользован во всех сервисах приложения
func ContextMiddleware(c *gin.Context) {
	ctx, cncl := context.WithCancel(c.Request.Context())
	defer cncl()
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

//AddServerStatsHandler - Выполняет сбор всей необходимой информации о сервисе, и текущем его состоянии
// не зависит от сервиса и будет переиспользован во всех сервисах приложения
func AddServerStatsHandler(router gin.IRoutes, url string, info *ServerInfo) {
	var stats = ServerStatus{StartDate: time.Now(), Info: *info}
	router.Use(func(c *gin.Context) {
		method := c.Request.Method
		reqStart := time.Now()
		c.Next()
		var successCnt, failedCnt uint64
		reqStop := time.Now().Sub(reqStart).Milliseconds()
		if !c.IsAborted() {
			successCnt = atomic.AddUint64(&stats.SuccesReqCnt, 1)
		} else {
			failedCnt = atomic.AddUint64(&stats.FaileReqCnt, 1)
			switch method {
			case "GET":
				atomic.AddUint64(&stats.FaileGetCnt, 1)
			case "POST":
				atomic.AddUint64(&stats.FailePostCnt, 1)
			case "PUT":
				atomic.AddUint64(&stats.FailePutCnt, 1)
			case "DELETE":
				atomic.AddUint64(&stats.FaileDelCnt, 1)
			}
		}
		allCnt := successCnt + failedCnt
		middleReqTime := atomic.LoadInt64(&stats.MiddleReqTime)
		atomic.StoreInt64(&stats.MiddleReqTime, int64((uint64(middleReqTime)*(allCnt-1)+uint64(reqStop))/allCnt))
	})

	router.GET(url, func(c *gin.Context) {
		var tempStats = ServerStatus{StartDate: stats.StartDate, Info: stats.Info}
		tempStats.UpTime = time.Now().Sub(tempStats.StartDate)
		tempStats.UpTimeStr = tempStats.UpTime.String()
		tempStats.SuccesReqCnt = atomic.LoadUint64(&stats.SuccesReqCnt)
		tempStats.FaileReqCnt = atomic.LoadUint64(&stats.FaileReqCnt)
		tempStats.FaileGetCnt = atomic.LoadUint64(&stats.FaileGetCnt)
		tempStats.FailePostCnt = atomic.LoadUint64(&stats.FailePostCnt)
		tempStats.FailePutCnt = atomic.LoadUint64(&stats.FailePutCnt)
		tempStats.FaileDelCnt = atomic.LoadUint64(&stats.FaileDelCnt)
		tempStats.MiddleReqTime = atomic.LoadInt64(&stats.MiddleReqTime)
		c.JSON(http.StatusOK, tempStats)
	})
}

//FormRequestID - формирует строку с идентификатором запроса
func FormRequestID(user *User) string {
	if user == nil {
		return ""
	}
	return fmt.Sprintf("%d:%s:%s", time.Now().UnixNano(), user.Email, user.SessionKey)
}

//CreateSignature - подписывает через secret пользователя userJSON
func CreateSignature(secret, usrJSON []byte) string {
	temp := make([]byte, len(usrJSON), len(usrJSON)+len(secret))
	copy(temp, usrJSON)
	temp = append(temp, secret...)
	signatureHash := sha256.Sum256(temp)
	return base64.StdEncoding.EncodeToString(signatureHash[:])
}

//CheckSignature - check received signature with origin
func CheckSignature(signature, userJSON, secret string) bool {
	temp := []byte(userJSON + secret)
	signatureHash := sha256.Sum256(temp)
	origin := base64.StdEncoding.EncodeToString(signatureHash[:])
	return strings.EqualFold(signature, origin)
}

//ParseHeaderMiddleware - read standart user header in http request to search them user and requestID parameters and add it to context of request
// Парсинг будет переиспользоватся в выше стоящих слоях приложения (сервисах)
func ParseHeaderMiddleware(c *gin.Context) {
	ParseHttpHeaderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Request = r
		c.Next()
	}))
	// userJSON := c.Request.Header.Get(UserHeaderKey)
	// signStr := c.Request.Header.Get(SignatureHeaderKey)
	// secret := os.Getenv(EegeonSecretKeyEnviron)
	// if !CheckSignature(signStr, userJSON, secret) {
	// 	c.AbortWithStatusJSON(http.StatusForbidden, EgeonError{Code: NotAuthError, Description: "Signature for user is incorrect"})
	// 	return
	// }
	// var user User
	// if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
	// 	c.AbortWithStatusJSON(http.StatusForbidden, EgeonError{Code: NotAuthError, Description: "Error when try parse user in header " + err.Error()})
	// 	return
	// }
	// ctx := context.WithValue(c.Request.Context(), UserKey, user)
	// requestID := c.Request.Header.Get(RequestIDHeaderKey)
	// if len(requestID) == 0 {
	// 	requestID = FormRequestID(&user)
	// }
	// ctx = context.WithValue(ctx, RequestID, requestID)
	// ctx = context.WithValue(ctx, SignKey, signStr)
	// c.Request = c.Request.WithContext(ctx)
	// c.Next()
}

//ParseHttpHeaderMiddleware - парсит заголовки запроса в поисках пользователя, проверяет подпись пользователя и requestID.
func ParseHttpHeaderMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userJSON := r.Header.Get(UserHeaderKey)
		signStr := r.Header.Get(SignatureHeaderKey)
		secret := os.Getenv(EegeonSecretKeyEnviron)
		if !CheckSignature(signStr, userJSON, secret) {
			data, _ := json.Marshal(EgeonError{Code: NotAuthError, Description: "Signature for user is incorrect"})
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write(data)
			return
		}
		var user User
		if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
			data, _ := json.Marshal(EgeonError{Code: NotAuthError, Description: "Error when try parse user in header " + err.Error()})
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write(data)
			return
		}
		ctx := context.WithValue(r.Context(), UserKey, user)
		requestID := r.Header.Get(RequestIDHeaderKey)
		if len(requestID) == 0 {
			requestID = FormRequestID(&user)
		}
		ctx = context.WithValue(ctx, RequestID, requestID)
		ctx = context.WithValue(ctx, SignKey, signStr)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}
