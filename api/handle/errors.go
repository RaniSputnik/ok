package handle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type HTTPError struct {
	Status  int    `json:"-"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func errInvalidJSON(err error) HTTPError {
	if err == io.EOF {
		err = errors.New("empty input")
	}
	return HTTPError{http.StatusBadRequest, "InvalidJSON", fmt.Sprintf("Invalid JSON body: '%s'.", err.Error())}
}

func errUnsupportedContent(expected string) HTTPError {
	return HTTPError{http.StatusUnsupportedMediaType, "UnsupportedContent", fmt.Sprintf("Content type must be: '%s'.", expected)}
}

func errMissingParameter(param string) HTTPError {
	return HTTPError{http.StatusBadRequest, "MissingParameter", fmt.Sprintf("Missing required parameter: '%s'.", param)}
}

func errInvalidParameter(param string) HTTPError {
	return HTTPError{http.StatusBadRequest, "InvalidParameter", fmt.Sprintf("A parameter is invalid: '%s'.", param)}
}

func errUsernameTaken(username string) HTTPError {
	return HTTPError{http.StatusConflict, "UsernameTaken", fmt.Sprintf("Username '%s' already in use.", username)}
}

func writeError(w http.ResponseWriter, err HTTPError) {
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}
