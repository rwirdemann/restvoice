package roles

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/rwirdemann/restvoice/cha06/domain"
)

type RoleRepository interface {
	GetCustomer(id int) domain.Customer
	GetInvoice(id int, join ...string) domain.Invoice
}

const publicKeyFilePath = "restvoice.pub"

var publicKey *rsa.PublicKey

func init() {
	if f, err := ioutil.ReadFile(publicKeyFilePath); err == nil {
		if publicKey, err = jwt.ParseRSAPublicKeyFromPEM(f); err != nil {
			log.Fatalf("Could not parse public key from pem file")
		}
	} else {
		log.Fatalf("Could not open public key file: %s", publicKeyFilePath)
	}
}

func AssertOwnsInvoice(next http.HandlerFunc, repository RoleRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractJwtFromHeader(r.Header)
		userId, _ := strconv.Atoi(claim(token, "id"))
		invoiceId, _ := strconv.Atoi(mux.Vars(r)["invoiceId"])
		invoice := repository.GetInvoice(invoiceId)
		customer := repository.GetCustomer(invoice.CustomerId)
		if customer.UserId == userId {
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusForbidden)
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

func claim(token string, key string) string {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err == nil {
		if claims, ok := t.Claims.(jwt.MapClaims); ok {
			if claims[key] != nil {
				return claims[key].(string)
			}
		}
	}

	return ""
}
