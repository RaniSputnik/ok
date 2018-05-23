package handle

import (
	"net/http"

	"github.com/gorilla/mux"
)

func New() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/games/{id}", oneGameHandler()).Methods(http.MethodGet)

	return wrapGlobalMiddleware(r)
}
