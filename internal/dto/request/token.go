package request

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("secret_key__DO_NOT_POST_IT_TO_GITHUB")

type Claims struct {
	UserId int `json:"userID"`
	jwt.StandardClaims
}

type Token struct {
	JWTToken string `json:"token"`
}

func GenerateToken(userID int) (Token, error) {
	claims := Claims{
		UserId: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var (
		err    error
		answer Token
	)

	answer.JWTToken, err = token.SignedString(jwtKey)
	if err != nil {
		return Token{}, err

	}

	return answer, nil
}

func ParseToken(tokenString string) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		return Claims{}, err
	}

	if !token.Valid {
		return Claims{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return Claims{}, err
	}

	return *claims, nil
}
