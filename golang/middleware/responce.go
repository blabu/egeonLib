package middleware

type Responce struct {
	Body   []byte              `json:"body"`
	Status int                 `json:"status"`
	Header map[string][]string `json:"header"`
}
