package app

import (
	"net/http"
	"time"

	"github.com/RaniSputnik/ok/api/handle"
	"github.com/gorilla/mux"
)

func fromMuxVars(key string) handle.RequestVarFunc {
	return func(r *http.Request) string {
		return mux.Vars(r)[key]
	}
}

func New(config Config) *http.Server {
	config = config.withSensibleDefaults()

	r := mux.NewRouter()

	r.HandleFunc("/register", handle.Register(config.Store))

	r.HandleFunc("/games/{id}", handle.OneGame(fromMuxVars("id"), config.Store)).Methods(http.MethodGet)
	r.HandleFunc("/games/{id}/stones", handle.Play(fromMuxVars("id"), config.Store)).Methods(http.MethodPost)

	r.NotFoundHandler = handle.NotFound()

	return &http.Server{
		Addr:    config.Addr,
		Handler: handle.WrapGlobalMiddleware(r),

		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 10,
	}
}
