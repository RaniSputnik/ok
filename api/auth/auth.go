package auth

import (
	"fmt"

	"github.com/RaniSputnik/ok/api/model"
	jwt "github.com/dgrijalva/jwt-go"
)

func NewHMAC(key []byte) Service {
	return hmac{key}
}

type Service interface {
	Token(player *model.Player) string
	Verify(tokenString string) (player string, ok bool)
}

type hmac struct {
	key []byte
}

type Claims struct {
	jwt.StandardClaims
	Username string `json:"usr"`
}

func (s hmac) Token(player *model.Player) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: player.Username,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(s.key)
	if err != nil {
		panic(err) // TODO
	}

	return tokenString
}

func (s hmac) Verify(tokenString string) (player string, ok bool) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return s.key, nil
	})
	if err != nil {
		return "", false
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Username, true
	}

	return "", false
}
