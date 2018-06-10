package handle

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/RaniSputnik/ok/api/auth"
	"github.com/RaniSputnik/ok/api/kontext"
	"github.com/RaniSputnik/ok/api/store"
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

type Middleware func(h http.Handler) http.Handler

func Auth(authSvc auth.Service, db store.Player) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authorization := r.Header.Get("Authorization")
			tokenString := strings.TrimPrefix(authorization, "Bearer ")
			username, ok := authSvc.Verify(tokenString)
			if !ok {
				writeError(w, errUnauthorized())
				return
			}

			player, err := db.GetPlayer(ctx, username)
			panicIf(err)
			if player == nil {
				writeError(w, errUnauthorized())
				return
			}

			ctxWithPlayer := kontext.WithPlayer(ctx, player)
			r = r.WithContext(ctxWithPlayer)

			h.ServeHTTP(w, r)
		})
	}
}
