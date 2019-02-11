package roles

import (
	"fmt"
	"github.com/rwirdemann/restvoice/kapitel09/identityprovider/secret"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/dgrijalva/jwt-go"
	"github.com/rwirdemann/restvoice/kapitel05/domain"
)

type RoleRepository interface {
	GetCustomer(id int) domain.Customer
	GetInvoice(id int, join ...string) domain.Invoice
}

func AssertAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractJwtFromHeader(r.Header)
		if isAdmin(token) {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", "Bearer realm=\"restvoice.org\"")
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func isAdmin(token string) bool {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret.Shared), nil
	})

	if err == nil {
		if claims, ok := t.Claims.(jwt.MapClaims); ok {
			if claims["admin"] != nil {
				return claims["admin"].(bool)
			}
		}
	}

	return false
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
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret.Shared), nil
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
