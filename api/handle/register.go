package handle

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RaniSputnik/ok/api/auth"
	"github.com/RaniSputnik/ok/api/model"
	"github.com/RaniSputnik/ok/api/store"
)

func Register(authSvc auth.Service, s store.Player) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if contentType := r.Header.Get("Content-Type"); !strings.Contains(contentType, "application/json") {
			writeError(w, errUnsupportedContent("application/json"))
			return
		}

		var params model.RegisterParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			writeError(w, errInvalidJSON(err))
			return
		}

		if missingPlayerParams(w, &params) {
			return
		}

		player := model.Player{Username: params.Username}
		if err := s.SavePlayer(r.Context(), &player); err != nil {
			if err == store.ErrUsernameTaken {
				writeError(w, errUsernameTaken(player.Username))
				return
			}
			panic(err) // TODO return internal server error instead?
		}

		token := authSvc.Token(&player)
		res := model.LoginResponse{
			Player: &player,
			Token:  token,
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	}
}

func missingPlayerParams(w http.ResponseWriter, params *model.RegisterParams) bool {
	if params.Username == "" {
		writeError(w, errMissingParameter("username"))
		return true
	}
	return false
}
