package main

import (
	"errors"
	"net/http"
	"time"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Auth struct {
	jwtKey []byte
}

type Claims struct {
	Email string
	jwt.StandardClaims
}


var auth *Auth

func (a *Auth) GetToken(email string) (token string, err error) {
	expireTime := time.Now().Add(30 * time.Minute)
	claims := &Claims {
		Email: email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expireTime.Unix(),
			},
		}

	claim := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := claim.SignedString(a.jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *Auth) Verify(req *http.Request) (err error) {
	bearToken := req.Header.Get("Authorization")
	if string(bearToken) == "" {
		return errors.New("Request is missing authorization information")
	}
	token := strings.Split(bearToken, " ")[1]
	claim := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error){
		return a.jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return errors.New("Signature is invalid")
		}
		return errors.New("Failed to parse token")
	}
	if !tkn.Valid {
		return errors.New("Token is not valid");
	}
	return nil
}

func GetAuth() *Auth {
	auth = &Auth{
		jwtKey: []byte("some_random_key"),
	}
	return auth
}
