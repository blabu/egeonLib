package golang

// KEYs for request cross services
const (
	RequestIDHeaderKey     = "RequestID"
	UserHeaderKey          = "User"
	SignatureHeaderKey     = "Sign"
	EegeonSecretKeyEnviron = "EGEON_SECRET_KEY" // По этому ключу в env операционки лежит секрет, которым подписывают авторизованного пользователя
)

//TODO package egeonGateway/parseuser are shared as external dependencies at another services
type requestIDType string
type contextKey string
type signatureKey string

// UserKey - ключ, по которому в контексте будет сохранен пользователь
var UserKey contextKey

//RequestID - ключ, по которому в контексте будем искать id запроса
var RequestID requestIDType

//SignKey - ключ, по которому в контексте будет сохранена подпись пользователя
var SignKey signatureKey
