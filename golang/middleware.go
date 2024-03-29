package golang

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
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

// AddServerStatsHandler - Выполняет сбор всей необходимой информации о сервисе, и текущем его состоянии
// не зависит от сервиса и будет переиспользован во всех сервисах приложения
func AddServerStatsHandler(router gin.IRoutes, url string, info *ServerInfo, checkService func() error) {
	var mem runtime.MemStats
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
			case http.MethodGet:
				atomic.AddUint64(&stats.FaileGetCnt, 1)
			case http.MethodPost:
				atomic.AddUint64(&stats.FailePostCnt, 1)
			case http.MethodPut:
				atomic.AddUint64(&stats.FailePutCnt, 1)
			case http.MethodDelete:
				atomic.AddUint64(&stats.FaileDelCnt, 1)
			}
		}
		allCnt := successCnt + failedCnt
		middleReqTime := atomic.LoadInt64(&stats.MiddleReqTime)
		atomic.StoreInt64(&stats.MiddleReqTime, int64((uint64(middleReqTime)*(allCnt-1)+uint64(reqStop))/allCnt))
	})

	router.GET(url, func(c *gin.Context) {
		if err := checkService(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, EgeonError{Code: ServiceWorkError, Description: err.Error()})
			return
		}
		runtime.ReadMemStats(&mem)
		var tempStats = ServerStatus{StartDate: stats.StartDate, Info: stats.Info}
		tempStats.Addition = map[string]interface{}{
			"routine":      runtime.NumGoroutine(),
			"cpu":          runtime.NumCPU(),
			"memmory":      mem.HeapAlloc / (1024 * 1024),
			"allObjects":   mem.Mallocs,
			"freesObject":  mem.Frees,
			"activeObject": mem.Mallocs - mem.Frees,
		}
		tempStats.UpTime = time.Now().Sub(stats.StartDate)
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

func GetServerStatusHandler(status ServerStatus, nowConnected *int32, checkService func() error) http.HandlerFunc {
	status.StartDate = time.Now()
	var mem runtime.MemStats
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checkService(); err != nil {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		runtime.ReadMemStats(&mem)
		status.Addition = map[string]interface{}{
			"memmory":      mem.HeapAlloc / (1024 * 1024),
			"allObjects":   mem.Mallocs,
			"freesObject":  mem.Frees,
			"nowConnected": atomic.LoadInt32(nowConnected),
			"routine":      runtime.NumGoroutine(),
			"cpu":          runtime.NumCPU(),
			"activeObject": mem.Mallocs - mem.Frees,
		}
		status.UpTime = time.Now().Sub(status.StartDate)
		status.UpTimeStr = time.Now().Sub(status.StartDate).String()

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	}
}

// FormRequestID - формирует строку с идентификатором запроса
func FormRequestID(user *User) string {
	if user == nil {
		return ""
	}
	return fmt.Sprintf("%d:%s:%s", time.Now().UnixNano(), user.Email, user.SessionKey)
}

// CreateSignature - подписывает через secret пользователя userJSON
func CreateSignature(secret, userJSON []byte) string {
	temp := make([]byte, len(userJSON), len(userJSON)+len(secret))
	copy(temp, userJSON)
	temp = append(temp, secret...)
	signatureHash := sha256.Sum256(temp)
	return base64.StdEncoding.EncodeToString(signatureHash[:])
}

// CheckSignature - check received signature with origin
func CheckSignature(signature, userJSON, secret string) bool {
	temp := []byte(userJSON + secret)
	signatureHash := sha256.Sum256(temp)
	origin := base64.StdEncoding.EncodeToString(signatureHash[:])
	return signature == origin
}

// ParseHeader - формирует контекст запроса исходя из заголовков HTTP запроса
// формирование заголовка выполняется функцией DoRequest
func ParseHeader(r *http.Request) (context.Context, error) {
	userJSON := r.Header.Get(UserHeaderKey)
	signStr := r.Header.Get(SignatureHeaderKey)
	secret := os.Getenv(EgeonSecretKeyEnviron)
	allowedRole := r.Header.Get(AllowedRoleHeaderKey)
	if !CheckSignature(signStr, userJSON, secret) {
		return r.Context(), EgeonError{Code: NotAuthError, Description: "Signature for user is incorrect"}
	}
	var user User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return r.Context(), EgeonError{Code: NotAuthError, Description: "Error when try parse user in header " + err.Error()}
	}
	ctx := context.WithValue(r.Context(), UserKey, user)
	requestID := r.Header.Get(RequestIDHeaderKey)
	if len(requestID) == 0 {
		requestID = FormRequestID(&user)
	}
	ctx = context.WithValue(ctx, RequestID, requestID)
	ctx = context.WithValue(ctx, SignKey, signStr)
	ctx = context.WithValue(ctx, AllowedRoleKey, allowedRole)
	return ctx, nil
}

// ParseHeaderMiddleware - read standart user header in http request to search them user and requestID parameters and add it to context of request
// Парсинг будет переиспользоватся в выше стоящих слоях приложения (сервисах) если используется gin
func ParseHeaderMiddleware(c *gin.Context) {
	ctx, err := ParseHeader(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, err)
	} else {
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// ParseHTTPHeaderMiddleware - парсит заголовки запроса в поисках пользователя, проверяет подпись пользователя и requestID.
// Парсинг будет переиспользоватся в выше стоящих слоях приложения (сервисах) если используется стандартный http или gorrila.mux
func ParseHTTPHeaderMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := ParseHeader(r)
		if err != nil {
			data, _ := json.Marshal(err)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write(data)
			return
		}
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}
