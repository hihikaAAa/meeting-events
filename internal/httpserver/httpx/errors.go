package httpx

import(
	"errors"
	"net/http"

	app "github.com/hihikaAAa/meeting-events/internal/app/app_errors"
)

func HttpStatusFromErr(err error) (status int, code string, msg string) {
	switch {
	case errors.Is(err, app.ErrValidation):
		return http.StatusBadRequest, "validation_error", err.Error()
	case errors.Is(err, app.ErrNotFound):
		return http.StatusNotFound, "not_found", err.Error()
	case errors.Is(err, app.ErrConflict):
		return http.StatusConflict, "conflict", err.Error()
	default:
		return http.StatusInternalServerError, "internal", "internal error"
	}
}
