package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/otakakot/sample-go-gqlgen/internal/domain"
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		token := strings.Split(auth, " ")

		if len(token) != 2 || token[0] != "Bearer" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		slog.Info("auth", slog.String("auth", token[1]))

		slog.Info("Authorize middleware begin")
		defer slog.Info("Authorize middleware end")

		ctx := domain.CtxWithUserID(r.Context(), token[1])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
