package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrInvalidSign  = errors.New("invalid signing method")
	ErrSignFailed   = errors.New("signing failed")
)

type JWTMaker struct {
	SecretKey string
}

func NewJWTMaker(key string) (Maker, error) {

	if len(key) < minSecretKeySize {
		return nil, fmt.Errorf("secret key does not meet the specifications")
	}

	return &JWTMaker{
		SecretKey: key,
	}, nil
}

func (jwtmaker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte(jwtmaker.SecretKey))

	if err != nil {
		return "", nil, ErrSignFailed
	}

	return tokenString, payload, nil
}

func (jwtmaker *JWTMaker) ValidateToken(token string) (*Payload, error) {

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(jwtmaker.SecretKey), nil
		}
		return nil, ErrInvalidSign
	}

	JWTtoken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		return nil, err
	}

	if !JWTtoken.Valid {
		return nil, ErrInvalidToken
	}

	payload, ok := JWTtoken.Claims.(*Payload)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return payload, nil
}
