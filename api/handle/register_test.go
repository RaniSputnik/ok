package handle_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RaniSputnik/ok/api/handle"
	"github.com/RaniSputnik/ok/api/store"
	"github.com/RaniSputnik/ok/api/store/inmemory"
)

func TestRegister(t *testing.T) {
	authSvc := &mockAuthService{}
	inMemoryStore := inmemory.New()

	validRegistration := func() io.Reader {
		r := `{
			"username": "RaniSputnik"
		}`
		return strings.NewReader(r)
	}

	t.Run("WithEmptyBody", func(t *testing.T) {
		w, r := setupRegister(nil)
		handle.Register(authSvc, inMemoryStore)(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"type":"InvalidJSON", "message":"Invalid JSON body: 'empty input'."}`, w.Body.String())
	})

	t.Run("WithNoContentType", func(t *testing.T) {
		w, r := setupRegister(validRegistration())
		r.Header.Set("Content-Type", "")
		handle.Register(authSvc, inMemoryStore)(w, r)

		assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
		assert.JSONEq(t, `{"type":"UnsupportedContent", "message":"Content type must be: 'application/json'."}`, w.Body.String())
	})

	t.Run("WithNoParams", func(t *testing.T) {
		w, r := setupRegister(strings.NewReader(`{ }`))
		handle.Register(authSvc, inMemoryStore)(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"type":"MissingParameter", "message":"Missing required parameter: 'username'."}`, w.Body.String())
	})

	t.Run("WithUsernameInUse", func(t *testing.T) {
		w, r := setupRegister(validRegistration())
		mockStore := new(mockPlayerStore)
		mockStore.Func.SavePlayer.Returns.Err = store.ErrUsernameTaken

		handle.Register(authSvc, mockStore)(w, r)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, `{"type":"UsernameTaken", "message":"Username 'RaniSputnik' already in use."}`, w.Body.String())
	})

	t.Run("IsSuccessful", func(t *testing.T) {
		const expectedToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c3IiOiJSYW5pU3B1dG5payJ9.valid"

		w, r := setupRegister(validRegistration())
		mockStore := new(mockPlayerStore)
		mockAuthSvc := &mockAuthService{}
		mockAuthSvc.Func.Token.Returns.Token = expectedToken
		expectedResponse := `{ "player": { "username":"RaniSputnik" }, "token": "` + expectedToken + `" }`

		handle.Register(mockAuthSvc, mockStore)(w, r)

		assert.Equal(t, 1, mockStore.Func.SavePlayer.Called.Times)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})
}

func setupRegister(body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/register", body)
	r.Header.Set("Content-Type", "application/json")
	return w, r
}
