package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
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
	signed, _ := token.SignedString(privateKey)
	return signed
}

func login(user string, password string) (User, bool) {
	return User{"Jo Brunner", true}, true
}
