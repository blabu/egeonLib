package golang

// KEYs for request cross services
const (
	AllowedRoleHeaderKey  = "AllowedUserRole"
	RequestIDHeaderKey    = "RequestID"
	UserHeaderKey         = "User"
	UserNameHeaderKey     = "X-UserName"
	SignatureHeaderKey    = "Sign"
	TokenQueryKey         = "token"
	EgeonSecretKeyEnviron = "EGEON_SECRET_KEY" // По этому ключу в env операционки лежит секрет, которым подписывают авторизованного пользователя
)

// TODO package egeonGateway/parseuser are shared as external dependencies at another services
type requestIDType string
type contextKey string
type signatureKey string
type allowedRoleType string
type tokenKey string

// UserKey - ключ, по которому в контексте будет сохранен пользователь
var UserKey contextKey

// AllowedRoleKey - ключ по которому в контексте будет сохранена роль по которой пользователю разрешен ответ по запросу
var AllowedRoleKey allowedRoleType

// RequestID - ключ, по которому в контексте будем искать id запроса
var RequestID requestIDType

// SignKey - ключ, по которому в контексте будет сохранена подпись пользователя
var SignKey signatureKey

// TokenKey - the key for token value in the context
var TokenCtxKey tokenKey
