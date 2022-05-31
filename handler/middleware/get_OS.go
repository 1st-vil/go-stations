package middleware

import (
	"context"
	"errors"
	"net/http"

	ua "github.com/mileusna/useragent"
)

type contextKey struct{}

var OSKey contextKey

func ContextWithOS(parent context.Context, r *http.Request) context.Context {
	return context.WithValue(r.Context(), OSKey, ua.Parse(r.UserAgent()).OS)
}

func OSFromContext(ctx context.Context) (string, error) {
	v := ctx.Value(OSKey)
	OS, ok := v.(string)
	if !ok {
		return "", errors.New("OS not found")
	}
	return OS, nil
}

func GetOS(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(ContextWithOS(r.Context(), r))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
