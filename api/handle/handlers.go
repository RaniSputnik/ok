package handle

import (
	"fmt"
	"net/http"
)

func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, HTTPError{http.StatusNotFound, "NoRoute", fmt.Sprintf("The path '%s' is invalid.", r.URL.Path)})
	}
}
