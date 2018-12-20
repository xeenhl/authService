package jwt

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/xeenhl/spendlog/backend/authService/model"
)

type Tokenazer interface {
	GenerateToken(u *model.User) (string, error)
	ParceAndVerifyToken(s string) (*jwt.Token, error)
}

type Jwt struct {
	keyLoader KeyLoader
}

const (
	tokenDuration = 1
	expireOffset  = 3600
)

var tokenizer *Jwt = nil

func NewTokenizer(k KeyLoader) *Jwt {
	return &Jwt{
		keyLoader: k,
	}

}

func (j *Jwt) GenerateToken(u *model.User) (string, error) {
	keys := j.keyLoader.InitializeKeysChain()

	token := jwt.NewWithClaims(jwt.SigningMethodPS512, jwt.MapClaims{
		"sub": u.Id,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * time.Duration(tokenDuration)).Unix(),
	})

	tokenString, err := token.SignedString(keys.PrivateKey)

	if err != nil {
		return "", errors.New("No Token has been generated")
	}

	return tokenString, nil

}

func (j *Jwt) ParceAndVerifyToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSAPSS); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return j.keyLoader.InitializeKeysChain().PublicKey, nil
	})
}
