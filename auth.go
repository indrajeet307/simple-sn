package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Auth struct {
	jwtKey []byte
}

type Claims struct {
	Email string
	jwt.StandardClaims
}


var auth Auth

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
	return nil
}

func GetAuth() (Auth) {
	auth = Auth{
		jwtKey: []byte("some_random_key"),
	}
	return auth
}
