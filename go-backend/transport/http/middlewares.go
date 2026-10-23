package transporthttp

import (
	"context"
	"log/slog"
	"net/http"
	"repository/auth"
	"repository/database"

	"slices"
	"strings"
	"sync"
	"uuid"
)

type Claims struct {
	Subject     string `json:"sub"`
	Typ         string `json:"typ"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAcesss map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
}

type ctxKey struct{}

func ClaimsFrom(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(ctxKey{}).(Claims)
	return c, ok
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || raw == "" {
			w.Header().Set("WWW-Authenticate", `Bearer`)
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		tok, err := auth.Verfier.Verify(r.Context(), raw)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		var claims Claims
		if err := tok.Claims(&claims); err != nil || claims.Typ != "Bearer" {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKey{}, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func EnsureUser() func(http.Handler) http.Handler {
	var seen sync.Map

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFrom(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if _, cached := seen.Load(claims.Subject); !cached {
				if err := database.InsertUser(r.Context(), userID); err != nil {
					slog.Error("EnsureUser: provisioning user %s: %v", userID.String(), err)
					http.Error(w, "internal error", http.StatusInternalServerError)
					return
				}

				seen.Store(claims.Subject, struct{}{})
			}

			next.ServeHTTP(w, r)
		})
	}
}

type respWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *respWriter) WriteHeader(code int) {
	if rw.written {
		return // ignore — header already sent, avoid the superfluous-call log
	}
	rw.written = true
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &respWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		slog.Info("HTTP Request", "method", r.Method, "path", r.URL.Path, "status", wrapped.statusCode)
	})
}

func Chain(h http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
	for _, v := range slices.Backward(m) {
		h = v(h)
	}
	return h
}
