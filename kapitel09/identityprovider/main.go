package main

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/rwirdemann/restvoice/kapitel09/identityprovider/secret"
)

type User struct {
	name  string
	admin bool
}

type CustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func main() {
	fmt.Printf("%v", token("jo", "secret"))
}

func token(name string, password string) string {
	var user User
	var ok bool
	if user, ok = login(name, password); !ok {
		return ""
	}

	claims := CustomClaims{user.name, user.admin, jwt.StandardClaims{
		Subject: "restvoice.org",
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret.Shared))
	if err != nil {
		panic(err)
	}
	return signed
}

func login(_ string, _ string) (User, bool) {
	return User{"Jo Brunner", true}, true
}
