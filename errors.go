package coinstats

import (
	"errors"
	"fmt"
)

var (
	ErrBadRequest            = errors.New("bad Request")             // 400
	ErrUnauthorized          = errors.New("unauthorized")            // 401
	ErrForbidden             = errors.New("forbidden")               // 403
	ErrNotFound              = errors.New("not found")               // 404
	ErrTransactionsNotSynced = errors.New("transactions not synced") // 409
	ErrRateLimit             = errors.New("rate limit exceeded")     // 429
	ErrServiceUnavailable    = errors.New("service unavailable")     // 503
	ErrBadStatusCode         = errors.New("bad status code")         // any, except 200

	ErrBadWalletData = errors.New("bad wallet data")
)

type Error struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	RequestID  string `json:"requestId"`
	Path       string `json:"path"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d %s", e.StatusCode, e.Message)
}
