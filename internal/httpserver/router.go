package httpserver

import (
	"net/http"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	appmw "github.com/hihikaAAa/meeting-events/internal/httpserver/middleware"
)

type Handlers struct {
	Create http.Handler
	Get    http.Handler
	Update http.Handler
	Delete http.Handler
}

func NewRouter(h Handlers, log *slog.Logger, basicAuthUser, basicAuthPass string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(appmw.Slog(log)) 
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Route("/v1/meetings", func(r chi.Router) {
		if basicAuthUser != "" {
			r.Use(middleware.BasicAuth("meeting-svc", map[string]string{
				basicAuthUser: basicAuthPass,
			}))
		}
		r.Method(http.MethodPost, "/", h.Create)
		r.Method(http.MethodPatch, "/{id}", h.Update)
		r.Method(http.MethodDelete, "/{id}", h.Delete)
	})

	r.Method(http.MethodGet, "/v1/meetings/{id}", h.Get)

	return r
}
