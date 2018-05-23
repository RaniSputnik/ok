package handle

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func WrapGlobalMiddleware(h http.Handler) http.Handler {
	h = handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(h)
	h = handlers.LoggingHandler(os.Stdout, h)
	h = alwaysEmit("Content-Type", "application/json", h)
	h = handlers.CORS()(h)
	return h
}

func alwaysEmit(header, value string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header, value)
		h.ServeHTTP(w, r)
	})
}
