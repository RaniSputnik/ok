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
	auth := handle.Auth(config.AuthSvc, config.Store)

	createGame := auth(handle.CreateGame(config.Store))
	getOneGame := auth(handle.OneGame(fromMuxVars("id"), config.Store))
	playStone := auth(handle.Play(fromMuxVars("id"), config.Store))

	r := mux.NewRouter()
	r.Handle("/register", handle.Register(config.AuthSvc, config.Store))
	r.Handle("/games", createGame).Methods(http.MethodPost)
	r.Handle("/games/{id}", getOneGame).Methods(http.MethodGet)
	r.Handle("/games/{id}/stones", playStone).Methods(http.MethodPost)

	r.NotFoundHandler = handle.NotFound()

	return &http.Server{
		Addr:    config.Addr,
		Handler: handle.WrapGlobalMiddleware(r),

		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 10,
	}
}
