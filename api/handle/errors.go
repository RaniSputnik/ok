package handle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/RaniSputnik/ok/game"
)

type HTTPError struct {
	Status  int    `json:"-"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func ErrIllegalMove(err error) HTTPError {
	res := HTTPError{
		Status: http.StatusBadRequest,
		Type:   "IllegalMove",
	}

	switch err {
	case game.ErrNotYourTurn:
		res.Message = "It is not your turn."
	case game.ErrOutsideBoard:
		res.Message = "Move is outside the board."
	case game.ErrPositionOccupied:
		res.Message = "The position is already occupied."
	case game.ErrSuicidalMove:
		res.Message = "Move would result in suicide."
	case game.ErrViolatesKo:
		res.Message = "Move violates ko."
	default:
		res.Message = "The move is illegal."
	}

	return res
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

func errUnauthorized() HTTPError {
	return HTTPError{http.StatusUnauthorized, "Unauthorized", "User credentials are missing or invalid."}
}

func errNotParticipating() HTTPError {
	return HTTPError{http.StatusForbidden, "NotParticipating", "You are not a player of this game."}
}

func errResourceNotFound(resourceKind, resourceID string) HTTPError {
	return HTTPError{http.StatusNotFound, "ResourceNotFound",
		fmt.Sprintf("Could not find %s with id: '%s'", resourceKind, resourceID)}
}

func errUsernameTaken(username string) HTTPError {
	return HTTPError{http.StatusConflict, "UsernameTaken", fmt.Sprintf("Username '%s' already in use.", username)}
}

func errInternalServerError(errorID string) HTTPError {
	msg := fmt.Sprintf("An error occurred. Error id: '%s'", errorID)
	return HTTPError{http.StatusInternalServerError, "InternalServerError", msg}
}

func writeError(w http.ResponseWriter, err HTTPError) {
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}

// Used to report an Internal Server Error
// when unrecoverable errors are encountered
func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
