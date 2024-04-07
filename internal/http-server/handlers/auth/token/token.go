package token

import (
	resp "avito-banner-service/internal/utils/response"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func NewTokenHandler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("token")

			if token == "" {
				resp.Send401Error(w, r)
				slog.Error("user has no token!!!")
				return
			}
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
