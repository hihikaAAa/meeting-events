package httpx

import(
	"errors"
	"net/http"

	app "github.com/hihikaAAa/meeting-events/internal/app/app_errors"
)

func HttpStatusFromErr(err error) int {
	switch {
	case errors.Is(err, app.ErrValidation):
		return http.StatusBadRequest
	case errors.Is(err, app.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, app.ErrConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
