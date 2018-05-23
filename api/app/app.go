package app

import (
	"net/http"
	"time"

	"github.com/RaniSputnik/ok/api/handle"
	"github.com/gorilla/mux"
)

func New(config Config) *http.Server {
	config = config.withSensibleDefaults()

	r := mux.NewRouter()

	r.HandleFunc("/register", handle.Register(config.Store))

	r.HandleFunc("/games/{id}", handle.OneGame()).Methods(http.MethodGet)

	r.NotFoundHandler = handle.NotFound()

	return &http.Server{
		Addr:    config.Addr,
		Handler: handle.WrapGlobalMiddleware(r),

		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 10,
	}
}
