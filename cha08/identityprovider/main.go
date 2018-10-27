package main

import (
	"fmt"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
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

var privateKey []byte

func init() {
	var err error
	privateKey, err = ioutil.ReadFile("restvoice.key")
	if err != nil {
		panic(err)
	}

}

func main() {
	if user, ok := login("jo", "secret"); ok {
		claims := CustomClaims{user.name, user.admin, jwt.StandardClaims{
			Subject: "restvoice.org",
		}}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString(privateKey)
		fmt.Printf("%v %v", signed, err)
	}
}

func login(user string, password string) (User, bool) {
	return User{"Jo Brunner", true}, true
}
