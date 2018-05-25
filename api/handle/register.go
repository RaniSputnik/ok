package handle

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/RaniSputnik/ok/api/model"
	"github.com/RaniSputnik/ok/api/store"
)

func Register(s store.Player) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if contentType := r.Header.Get("Content-Type"); !strings.Contains(contentType, "application/json") {
			writeError(w, errUnsupportedContent("application/json"))
			return
		}

		var player model.Player
		if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
			writeError(w, errInvalidJSON(err))
			return
		}

		if missingPlayerParams(w, &player) {
			return
		}

		_, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(player.PublicKey)))
		if err != nil {
			writeError(w, errInvalidParameter("public_key"))
			return
		}

		player.Verified = false // The player can not be verified at this point
		if err := s.SavePlayer(r.Context(), &player); err != nil {
			if err == store.ErrUsernameTaken {
				writeError(w, errUsernameTaken(player.Username))
				return
			}
			panic(err) // TODO return internal server error instead?
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(player)
	}
}

func missingPlayerParams(w http.ResponseWriter, player *model.Player) bool {
	if player.Username == "" {
		writeError(w, errMissingParameter("username"))
		return true
	}
	if player.PublicKey == "" {
		writeError(w, errMissingParameter("public_key"))
		return true
	}
	return false
}
