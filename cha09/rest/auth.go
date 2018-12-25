package rest

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	jwt "github.com/dgrijalva/jwt-go"
)

const publicKeyFilePath = "restvoice.pub"

var publicKey *rsa.PublicKey

func init() {
	var b []byte
	var err error
	if b, err = ioutil.ReadFile(publicKeyFilePath); err != nil {
		log.Fatalf("Could not open public key file: %s", publicKeyFilePath)
	}

	if publicKey, err = jwt.ParseRSAPublicKeyFromPEM(b); err != nil {
		log.Fatalf("Could not parse public key from pem file")
	}
}

func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if username, password, ok := r.BasicAuth(); ok {
			if username == os.Getenv("USERNAME") && password == os.Getenv("PASSWORD") {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", "Basic realm=\"restvoice.org\"")
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractJwtFromHeader(r.Header)
		if verifyJWT(token) {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", "Bearer realm=\"restvoice.org\"")
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func extractJwtFromHeader(header http.Header) (jwt string) {
	var jwtRegex = regexp.MustCompile(`^Bearer (\S+)$`)

	if val, ok := header["Authorization"]; ok {
		for _, value := range val {
			if result := jwtRegex.FindStringSubmatch(value); result != nil {
				jwt = result[1]
				return
			}
		}
	}

	return
}

func verifyJWT(token string) bool {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	return err == nil && t.Valid
}
