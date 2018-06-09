package handle

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/rs/xid"

	"github.com/gorilla/handlers"
)

func WrapGlobalMiddleware(h http.Handler) http.Handler {
	h = recoveryHandler(h)
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

func recoveryHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errorID := xid.New().String()
				log.Printf("An error occurred, ErrorID='%s', Err='%s' Stack='%s'", errorID, err, debug.Stack())
				writeError(w, errInternalServerError(errorID))
			}
		}()

		h.ServeHTTP(w, r)
	})
}
