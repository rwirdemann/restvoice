package rest

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/rwirdemann/restvoice/kapitel09/identityprovider/secret"

	"github.com/dgrijalva/jwt-go"
)

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
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret.Shared), nil
	})

	return err == nil && t.Valid
}

const password = "time"
const realm = "restvoice"
const nonce = "UAZs1dp3wX5BtXEpoCXKO2lHhap564rX"
const opaque = "XF3tAJ3483jUUAUJJQJJAHDQP01MJHD"
const cnonce = "oaSHizKi0RcJXmFE2TMtW8IefL799dWU"
const nc = "00000001"
const qop = "auth"
const method = "POST"

func DigestAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if strings.HasPrefix(authorization, "Digest") {
			authFields := digestParts(authorization)
			h1 := hash(authFields["username"] + ":" + authFields["realm"] + ":" + password)
			h2 := hash(r.Method + ":" + authFields["uri"])
			h3 := hash(h1 + ":" +
				authFields["nonce"] + ":" +
				authFields["nc"] + ":" +
				authFields["cnonce"] + ":" +
				authFields["qop"] + ":" + h2)
			if h3 == authFields["response"] {
				next.ServeHTTP(w, r)
				return
			}
		}
		auth := fmt.Sprintf("Digest realm=\"%s\" qop=\"auth\" nonce=\"%s\" opaque=\"%s\"", realm, nonce, opaque)
		w.Header().Set("WWW-Authenticate", auth)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func hash(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

func digestParts(authorization string) map[string]string {
	result := map[string]string{}
	wantedHeaders := []string{"username", "nonce", "realm", "qop", "uri", "nc", "response", "opaque", "cnonce"}
	requestHeaders := strings.Split(authorization, ",")
	for _, r := range requestHeaders {
		for _, w := range wantedHeaders {
			if strings.Contains(r, " "+w) {
				v := strings.Split(r, "=")[1]
				result[w] = strings.Trim(v, `"`)
			}
		}
	}
	return result
}
