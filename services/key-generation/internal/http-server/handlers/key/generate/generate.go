package generate

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/timeb30/techstreamshop/services/key-generation/internal/lib/api/response"
)

type Request struct {
	UserID      int64  `json:"user_id" validate:"required"`
	SoftwareID  string `json:"software_id" validate:"required,len=32"`
	KeyDuration int64  `json:"key_duration" validate:"required"`
}
type Response struct {
	response.Response
	KeyID int `json:"key_id,omitempty"`
}

type KeySaverProvider interface {
	SaveKey(userID int64, key string, startDate time.Time, endDate time.Time) (int64, error)
}
type KeyGeneratorProvider interface {
	GenerateKey(SoftwareID string, days int64, version string) (string, error)
}

func New(log *slog.Logger, keySaver KeySaverProvider, keyGen KeyGeneratorProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.key.generate.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to parse request", "error", err)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", "error", err)

			render.JSON(w, r, response.ValidationError(validateErr)) // can add custom validation erro 1:22:25
			return

		}
		key, err := keyGen.GenerateKey(req.SoftwareID, req.KeyDuration, "")
		if err != nil || key == "" {
			log.Error("failed to generate key", "error", err)
			render.JSON(w, r, response.Error("failed to generate key"))
			return
		}
		_, err = keySaver.SaveKey(req.UserID, key, time.Now(), time.Now().AddDate(0, 0, 5))
		if err != nil {
			log.Error("failed to save key", "error", err)
			render.JSON(w, r, response.Error("failed to save key"))
			return
		}
		log.Info("key saved", slog.String("key", key))
		render.JSON(w, r, response.OK())
	}
}
