package handle_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/RaniSputnik/ok/api/handle"
	"github.com/RaniSputnik/ok/api/kontext"
	"github.com/RaniSputnik/ok/api/model"
	"github.com/RaniSputnik/ok/game"
	"github.com/stretchr/testify/assert"
)

var testGame = model.Game{
	ID:        "game_12345",
	CreatedBy: "Alice",
	CreatedAt: time.Now().In(time.UTC),
	Black:     "Alice",
	White:     "Bob",
	Board: model.Board{
		Size: game.BoardSizeSmall,
		Stones: []model.Stone{
			model.Stone{
				Colour: game.Black.String(),
				X:      3,
				Y:      5,
			},
		},
	},
	Moves: []model.Move{
		model.Move{
			Kind:    model.MoveKindStone,
			Colour:  game.Black.String(),
			Message: "Black stone played at D6",
			X:       intPointer(3),
			Y:       intPointer(5),
		},
	},
}

func intPointer(v int) *int {
	return &v
}

func TestOneGame(t *testing.T) {

	anyStore := func() *mockGameStore {
		oneStonePlayedAtD6 := []game.Move{
			game.Stone{
				Colour: game.Black,
				Position: game.Position{
					X: 3,
					Y: 5,
				},
			},
		}

		m := &mockGameStore{}
		m.Func.GetGameByID.Returns.Game = &testGame
		m.Func.GetGameByID.Returns.Err = nil
		m.Func.GetGameMoves.Returns.Moves = oneStonePlayedAtD6
		m.Func.GetGameMoves.Returns.Err = nil
		return m
	}

	validGameID := testGame.ID
	givesID := func(id string) handle.RequestVarFunc {
		return func(r *http.Request) string {
			return id
		}
	}

	t.Run("WithGameThatDoesNotExist", func(t *testing.T) {
		w, r := setupOneGame()
		mockStore := anyStore()
		mockStore.Func.GetGameByID.Returns.Game = nil
		mockStore.Func.GetGameByID.Returns.Err = nil
		idDoesNotExist := "does_not_exist"
		handle.OneGame(givesID(idDoesNotExist), mockStore)(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		expectJSON := `{"type":"ResourceNotFound", "message":"Could not find game with id: '` + idDoesNotExist + `'"}`
		assert.JSONEq(t, expectJSON, w.Body.String())
	})

	t.Run("WithValidGameID", func(t *testing.T) {
		testGameBytes, err := json.Marshal(testGame)
		if err != nil {
			t.Fatal(err)
		}
		validGameJSON := string(testGameBytes)

		w, r := setupOneGame()
		handle.OneGame(givesID(validGameID), anyStore())(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, validGameJSON, w.Body.String())
	})
}

func TestPlay(t *testing.T) {
	anyStore := func() *mockGameStore {
		m := &mockGameStore{}
		m.Func.GetGameByID.Returns.Game = &testGame
		m.Func.GetGameByID.Returns.Err = nil
		return m
	}

	blackPlayer := &model.Player{
		Username: testGame.Black,
	}
	whitePlayer := &model.Player{
		Username: testGame.White,
	}
	observingPlayer := &model.Player{
		Username: "Clive",
	}

	validGameID := testGame.ID
	getValidGameID := func(r *http.Request) string {
		return validGameID
	}

	invalidGameID := "game_invalid"
	getInvalidGameID := func(r *http.Request) string {
		return invalidGameID
	}

	t.Run("WithEmptyBody", func(t *testing.T) {
		w, r := setupPlay(nil)
		handle.Play(getValidGameID, anyStore())(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"type":"InvalidJSON", "message":"Invalid JSON body: 'empty input'."}`, w.Body.String())
	})

	t.Run("WithoutX", func(t *testing.T) {
		w, r := setupPlay(strings.NewReader(`{"y":0}`))
		handle.Play(getValidGameID, anyStore())(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"type":"MissingParameter", "message":"Missing required parameter: 'x'."}`, w.Body.String())
	})

	t.Run("WithoutY", func(t *testing.T) {
		w, r := setupPlay(strings.NewReader(`{"x":0}`))
		handle.Play(getValidGameID, anyStore())(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"type":"MissingParameter", "message":"Missing required parameter: 'y'."}`, w.Body.String())
	})

	t.Run("WithoutUserContext", func(t *testing.T) {
		w, r := setupPlay(strings.NewReader(`{"x":2, "y":9}`))
		// TODO rather than panic should we return NotAuthorized?
		// This is expected to be handled by middleware so may not
		// be a valid test case for us...
		assert.Panics(t, func() {
			handle.Play(getValidGameID, anyStore())(w, r)
		})
	})

	t.Run("WithInvalidGameID", func(t *testing.T) {
		mockStore := &mockGameStore{}
		// Not found result
		mockStore.Func.GetGameByID.Returns.Game = nil
		mockStore.Func.GetGameByID.Returns.Err = nil

		w, r := setupPlay(strings.NewReader(`{"x":2, "y":9}`))
		r = r.WithContext(kontext.WithPlayer(r.Context(), blackPlayer))
		handle.Play(getInvalidGameID, mockStore)(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		expectJSON := `{"type":"ResourceNotFound", "message":"Could not find game with id: '` + invalidGameID + `'"}`
		assert.JSONEq(t, expectJSON, w.Body.String())
		// We expect get game by id to have been called
		if assert.Equal(t, 1, mockStore.Func.GetGameByID.Called.Times,
			"Expected 'GetGameByID' to have been called exactly once.") {
			assert.Equal(t, r.Context(), mockStore.Func.GetGameByID.Called.With.Ctx)
			assert.Equal(t, invalidGameID, mockStore.Func.GetGameByID.Called.With.ID)
		}
		// We don't expect save stone to have been called
		assert.Zero(t, mockStore.Func.SaveStone.Called.Times)
	})

	t.Run("WhenUserIsNotParticipating", func(t *testing.T) {
		mockStore := anyStore()
		w, r := setupPlay(strings.NewReader(`{"x":2, "y":9}`))
		r = r.WithContext(kontext.WithPlayer(r.Context(), observingPlayer))
		handle.Play(getValidGameID, mockStore)(w, r)

		assert.Equal(t, http.StatusForbidden, w.Code)
		expectJSON := `{"type":"NotParticipating", "message":"You are not a player of this game."}`
		assert.JSONEq(t, expectJSON, w.Body.String())
		assert.Zero(t, mockStore.Func.SaveStone.Called.Times)
	})

	t.Run("FailsWhenItIsNotYourTurn", func(t *testing.T) {
		mockStore := anyStore()
		oneMoveAlreadyPlayed(mockStore, 0, 0)

		w, r := setupPlay(strings.NewReader(`{"x":2, "y":9}`))
		r = r.WithContext(kontext.WithPlayer(r.Context(), blackPlayer))
		handle.Play(getValidGameID, mockStore)(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		expectJSON := `{"type":"IllegalMove", "message":"It is not your turn."}`
		assert.JSONEq(t, expectJSON, w.Body.String())
	})

	t.Run("FailsWhenMoveIsOutsideBoard", func(t *testing.T) {
		mockStore := anyStore()
		invalidX, invalidY := testGame.Board.Size+1, -1

		body := fmt.Sprintf(`{"x":%d, "y":%d}`, invalidX, invalidY)
		w, r := setupPlay(strings.NewReader(body))
		r = r.WithContext(kontext.WithPlayer(r.Context(), blackPlayer))
		handle.Play(getValidGameID, mockStore)(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		expectJSON := `{"type":"IllegalMove", "message":"Move is outside the board."}`
		assert.JSONEq(t, expectJSON, w.Body.String())
		assert.Zero(t, mockStore.Func.SaveStone.Called)
	})

	t.Run("FailsWhenBoardStateIsInvalid", func(t *testing.T) {
		mockStore := anyStore()
		occupiedX, occupiedY := 0, 0
		oneMoveAlreadyPlayed(mockStore, occupiedX, occupiedY)

		body := fmt.Sprintf(`{"x":%d, "y":%d}`, occupiedX, occupiedY)
		w, r := setupPlay(strings.NewReader(body))
		r = r.WithContext(kontext.WithPlayer(r.Context(), whitePlayer))
		handle.Play(getValidGameID, mockStore)(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		expectJSON := `{"type":"IllegalMove", "message":"The position is already occupied."}`
		assert.JSONEq(t, expectJSON, w.Body.String())
		assert.Zero(t, mockStore.Func.SaveStone.Called)
	})

	t.Run("SavesStone", func(t *testing.T) {
		mockStore := anyStore()
		playX, playY := 5, 3

		input := fmt.Sprintf(`{"x":%d, "y":%d}`, playX, playY)
		w, r := setupPlay(strings.NewReader(input))
		r = r.WithContext(kontext.WithPlayer(r.Context(), blackPlayer))
		handle.Play(getValidGameID, mockStore)(w, r)

		log.Println(w.Body.String())
		if assert.Equal(t, 1, mockStore.Func.SaveStone.Called.Times) {
			expectStone := game.Stone{game.Black, game.Position{playX, playY}}
			assert.Equal(t, r.Context(), mockStore.Func.SaveStone.Called.With.Ctx)
			assert.Equal(t, validGameID, mockStore.Func.SaveStone.Called.With.GameID)
			assert.Equal(t, expectStone, mockStore.Func.SaveStone.Called.With.Stone)
		}
	})

	t.Run("ReturnsBlackStone", func(t *testing.T) {
		w, r := setupPlay(strings.NewReader(`{"x":5, "y":3}`))
		r = r.WithContext(kontext.WithPlayer(r.Context(), blackPlayer))
		handle.Play(getValidGameID, anyStore())(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `{"colour":"Black","x":5,"y":3}`, w.Body.String())
	})

	t.Run("ReturnsWhiteStone", func(t *testing.T) {
		mockStore := anyStore()
		oneMoveAlreadyPlayed(mockStore, 0, 0)
		w, r := setupPlay(strings.NewReader(`{"x":4, "y":2}`))
		r = r.WithContext(kontext.WithPlayer(r.Context(), whitePlayer))
		handle.Play(getValidGameID, mockStore)(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `{"colour":"White","x":4,"y":2}`, w.Body.String())
	})
}

func setupOneGame() (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/games/ID", nil)
	return w, r
}

func setupPlay(body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/games/TODO/stones", body)
	r.Header.Set("Content-Type", "application/json")
	return w, r
}

func oneMoveAlreadyPlayed(mockStore *mockGameStore, x, y int) {
	mockStore.Func.GetGameMoves.Returns.Moves = []game.Move{
		game.Stone{Colour: game.Black, Position: game.Position{X: x, Y: y}},
	}
}
