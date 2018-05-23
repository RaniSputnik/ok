package handle_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RaniSputnik/ok/api/handle"
	"github.com/RaniSputnik/ok/api/model"
	"github.com/RaniSputnik/ok/api/store"
)

const invalidPublicKey = "this-public-key-is-invalid"
const validPublicKey = "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFDQVFDdkdHeXE3L29VRTVyN2t5ZEF5RFg3c2Rmc0haenk5cnJyYXlscjl5aHlzZkJGcDlJMG1wclB1eUN4UzhQUXlkT0NHVXhoRjA4SlFmYjAzTGpwcmZ2M3gvRTl2eXhzTE1oMGhnUS9QWU9GYVpvbzVkcjlMUlVwVlJNMWgzSHVNbkJWbExzbkpwYklUWDc3aGdKelJDUEJ4c2FrU1B0YWhyUHJNTXpkK25wb1RaQXhtVTBuTzhrR3lwZTE1clZ4aHJ1amEzWHVWWUdzdjVraGc0VkoyVWVTZE41L2k0d0JsTW02OWhIeXFaM05WQU54cnE4ZlBveGRNRjNSUEpDOHZyRCt3YXo5dURFbzF1T2Y4YkxYbWFLcWRwaGpUZkhrWFE5QUl6NXJJaVVIVktGZXJpSldudlJBUm14bGRDemVvR0h5U3ZBVTFXeWlob2hNOUwrMWNOSklBSXA3dEFJSngzU0ZLNmFtM1lGbXpZODkyNWJuaFo2cEVONkttS0EvRVBwZG1CNURoaVd0TXkrR21HL3YyQjdHOEhod1MzVTBpc0RkOERDTGhGeWlxR3J2VDJtRFNxN2R3V3R4QmRsWEZwa2tHYkZKeG94VFNBdmlmN25uL2Mvcnh1RFhERjBjYXVjQ01kWFRJbysxM2UxZ3lKaThkZWtNOW5qeGZOVXRJT1YwdmEzWGZyUTAzOHFMSEVIVTNKayt3bTZVTWpaeEM1eHI2eHBSRWF1QUtObjJ4NDI3cFFraUNDSzQ2cTRDT2pINmp0bWw5WG40UXNDdWFqeHpOaTM3am5CTEdlYTdMTTBTNTRXQitoUFNUUjlHYi93UW5jeHZINGw5T0s4L3VkVEFGeEljWHRPbFpjWkU3TnRxdThEZ05wZHpaMDNuTFdjL2UxRTNTQUxQeVE9PSB0YWxrQHJ5YW5sb2FkZXIubWUK"

func TestRegister(t *testing.T) {
	inMemoryStore := store.NewInMemory()

	validRegistration := func() io.Reader {
		r := `{
			"username": "RaniSputnik",
			"public_key": "` + validPublicKey + `"
		}`
		return strings.NewReader(r)
	}

	t.Run("WithEmptyBody", func(t *testing.T) {
		w, r := setupRegister(nil)
		handle.Register(inMemoryStore)(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"type":"InvalidJSON", "message":"Invalid JSON body: 'empty input'."}`, w.Body.String())
	})

	t.Run("WithNoContentType", func(t *testing.T) {
		w, r := setupRegister(validRegistration())
		r.Header.Set("Content-Type", "")
		handle.Register(inMemoryStore)(w, r)

		assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
		assert.JSONEq(t, `{"type":"UnsupportedContent", "message":"Content type must be: 'application/json'."}`, w.Body.String())
	})

	t.Run("WithoutPublicKey", func(t *testing.T) {
		w, r := setupRegister(strings.NewReader(`{ "username": "RaniSputnik" }`))
		handle.Register(inMemoryStore)(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"type":"MissingParameter", "message":"Missing required parameter: 'public_key'."}`, w.Body.String())
	})

	t.Run("WithInvalidPublicKey", func(t *testing.T) {
		body := `{ "username": "RaniSputnik", "public_key": "` + invalidPublicKey + `" }`
		w, r := setupRegister(strings.NewReader(body))
		handle.Register(inMemoryStore)(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"type":"InvalidParameter", "message":"A parameter is invalid: 'public_key'."}`, w.Body.String())
	})

	t.Run("WithoutUsername", func(t *testing.T) {
		w, r := setupRegister(strings.NewReader(`{ "public_key": "` + validPublicKey + `" }`))
		handle.Register(inMemoryStore)(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"type":"MissingParameter", "message":"Missing required parameter: 'username'."}`, w.Body.String())
	})

	t.Run("WithUsernameInUse", func(t *testing.T) {
		w, r := setupRegister(validRegistration())
		mockStore := new(mockPlayerStore)
		mockStore.Func.SavePlayer.Returns.Err = store.ErrUsernameTaken

		handle.Register(mockStore)(w, r)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, `{"type":"UsernameTaken", "message":"Username 'RaniSputnik' already in use."}`, w.Body.String())
	})

	t.Run("DoesNotAllowFalseVerification", func(t *testing.T) {
		body := `{ "username": "RaniSputnik", "verified":true, "public_key": "` + validPublicKey + `" }`
		w, r := setupRegister(strings.NewReader(body))
		mockStore := new(mockPlayerStore)

		handle.Register(mockStore)(w, r)

		gotPlayer := mockStore.Func.SavePlayer.Called.With.Player
		assert.False(t, gotPlayer.Verified, "Player should not be verified")
	})

	t.Run("IsSuccessful", func(t *testing.T) {
		w, r := setupRegister(validRegistration())
		mockStore := new(mockPlayerStore)

		handle.Register(mockStore)(w, r)

		assert.Equal(t, 1, mockStore.Func.SavePlayer.Called.Times)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `{"username":"RaniSputnik", "public_key":"`+validPublicKey+`", "verified":false }`, w.Body.String())
	})
}

func setupRegister(body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/register", body)
	r.Header.Set("Content-Type", "application/json")
	return w, r
}

type mockPlayerStore struct {
	Func struct {
		SavePlayer struct {
			Called struct {
				With struct {
					Ctx    context.Context
					Player *model.Player
				}
				Times int
			}
			Returns struct {
				Err error
			}
		}
	}
}

func (m *mockPlayerStore) SavePlayer(ctx context.Context, player *model.Player) error {
	m.Func.SavePlayer.Called.Times++
	m.Func.SavePlayer.Called.With.Ctx = ctx
	m.Func.SavePlayer.Called.With.Player = player
	return m.Func.SavePlayer.Returns.Err
}
