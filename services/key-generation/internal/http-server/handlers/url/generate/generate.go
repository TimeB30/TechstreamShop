package generate

import (
	"log/slog"
	"net/http"
	"time"
)

type Request struct {
	UserId int64 `json:"user_id" validate:"required"`
}
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	KeyID  int64  `json:"key_id"`
}

type KeySaver interface {
	SaveKey(userID int64, key string, startDate time.Time, endDate time.Time) (int64, error)
}

func New(log *slog.Logger, keySaver KeySaver) http.Handler {

}
