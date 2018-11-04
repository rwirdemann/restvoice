package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/restvoice/cha05/database"
	"github.com/rwirdemann/restvoice/cha05/domain"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/customers", readCustomersHandler).Methods("GET")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices", createInvoiceHandler).Methods("POST")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}/bookings", createBookingHandler).Methods("POST")

	r.HandleFunc("/customers/{customerId:[0-9]+}/projects", createProjectHandler).Methods("POST")
	r.HandleFunc("/customers/{customerId:[0-9]+}/projects", readProjectsHandler).Methods("GET")

	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}/bookings/{bookingId:[0-9]+}", deleteBookingHandler).Methods("DELETE")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", updateInvoiceHandler).Methods("PUT")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", readInvoiceHandler).Methods("GET")
	r.HandleFunc("/activities", readActivitiesHandler).Methods("GET")

	fmt.Println("Restvoice started on http://localhost:8080...")
	http.ListenAndServe(":8080", r)
}

func readActivitiesHandler(writer http.ResponseWriter, request *http.Request) {
	b, _ := json.Marshal(activityRepository.All())
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(b)
}

func readProjectsHandler(writer http.ResponseWriter, request *http.Request) {
	customerId, _ := strconv.Atoi(mux.Vars(request)["customerId"])
	b, _ := json.Marshal(projectRepository.ByCustomer(customerId))
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(b)
}

func readCustomersHandler(writer http.ResponseWriter, _ *http.Request) {
	b, _ := json.Marshal(customerRepository.All())
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(b)
}

var repository = database.NewRepository()
var activityRepository = activity.NewRepository()
var bookingRepository = booking.NewRepository()
var projectRepository = project.NewRepository()
var customerRepository = customer.NewRepository()
var rateRepository = rate.NewRepository()

func createInvoiceHandler(writer http.ResponseWriter, request *http.Request) {
	// Read invoice data from request body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// CreateInvoice invoice and marshal it to JSON
	var i domain.Invoice
	json.Unmarshal(body, &i)

	i.CustomerId, _ = strconv.Atoi(mux.Vars(request)["customerId"])
	created := repository.CreateInvoice(i)
	b, err := json.Marshal(created)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response
	location := fmt.Sprintf("%s/%d", request.URL.String(), created.Id)
	writer.Header().Set("Location", location)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(b)
}

func createBookingHandler(writer http.ResponseWriter, request *http.Request) {
	// Read booking data from request body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create booking booking and marshal it to JSON
	var booking domain.Booking
	json.Unmarshal(body, &booking)
	created := repository.CreateBooking(booking)
	created.InvoiceId, _ = strconv.Atoi(mux.Vars(request)["invoiceId"])
	b, err := json.Marshal(created)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response
	location := fmt.Sprintf("%s/%d", request.URL.String(), created.Id)
	writer.Header().Set("Location", location)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(b)
}

func deleteBookingHandler(writer http.ResponseWriter, request *http.Request) {
	bookingId, _ := strconv.Atoi(mux.Vars(request)["bookingId"])
	repository.DeleteBooking(bookingId)
	writer.WriteHeader(http.StatusNoContent)
}

func updateInvoiceHandler(writer http.ResponseWriter, request *http.Request) {
	// Read invoice data from request body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal and update invoice
	var i domain.Invoice
	json.Unmarshal(body, &i)
	i.Id, _ = strconv.Atoi(mux.Vars(request)["invoiceId"])
	i.CustomerId, _ = strconv.Atoi(mux.Vars(request)["customerId"])

	// Aggregate positions
	if i.Status == "payment expected" {
		bookings := bookingRepository.ByInvoiceId(i.Id)
		for _, b := range bookings {
			activity := repository.GetActivity(b.ActivityId)
			rate := repository.GetRate(b.ProjectId, b.ActivityId)

			i.AddPosition(b.ProjectId, activity.Name, b.Hours, rate.Price)
		}
	}

	repository.Update(i)

	// Write response
	writer.WriteHeader(http.StatusNoContent)
}

func readInvoiceHandler(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(request)["invoiceId"])
	i, _ := repository.FindById(id)
	accept := request.Header.Get("Accept")
	switch accept {
	case "application/pdf":
		content := bytes.NewReader(i.ToPdf())
		http.ServeContent(writer, request, "invoice.pdf", time.Now(), content)
	case "application/json":
		b, _ := json.Marshal(i)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(b)
	default:
		writer.WriteHeader(http.StatusNotAcceptable)
	}
}

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if username, password, ok := r.BasicAuth(); ok {
			if username == os.Getenv("USERNAME") && password == os.Getenv("PASSWORD") {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", "Basic realm=\"example@restvoice.org\"")
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func jwtAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractJwtFromHeader(r.Header)
		if verifyJWT(token) {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", "Bearer realm=\"example@restvoice.org\"")
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func assertCustomer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractJwtFromHeader(r.Header)
		customerId, _ := strconv.Atoi(mux.Vars(r)["customerId"])
		if ownsCustomer(token, customerId) {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", "Bearer realm=\"example@restvoice.org\"")
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func ownsCustomer(token string, customerId int) bool {
	userId := claim(token, "sub")
	customer := customerRepository.ById(customerId)
	return customer.UserId == userId
}

func assertAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractJwtFromHeader(r.Header)
		if isAdmin(token) {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", "Bearer realm=\"example@restvoice.org\"")
		w.WriteHeader(http.StatusUnauthorized)
	}
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

func isAdmin(token string) bool {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
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

var password = "time"

func digestAuth(next http.HandlerFunc) http.HandlerFunc {
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
		auth := fmt.Sprintf("Digest realm=\"%s\" qop=\"auth\" nonce=\"%s\" opaque=\"%s\"", realm(), nonce(), opaque())
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

func nonce() string {
	return "UAZs1dp3wX5BtXEpoCXKO2lHhap564rX"
}

func opaque() string {
	return "xU2Z4FyqwKUBdwTMRYdGtAG1ppaT0bNm"
}

func realm() string {
	return "example@restvoice.org"
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

func createProjectHandler(writer http.ResponseWriter, request *http.Request) {
	// Read invoice data from request body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// CreateInvoice project and marshal it to JSON
	var p domain.Project
	json.Unmarshal(body, &p)

	p.CustomerId, _ = strconv.Atoi(mux.Vars(request)["customerId"])
	created := projectRepository.Create(p)
	b, err := json.Marshal(created)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response
	writer.Header().Set("Location", fmt.Sprintf("%s/%d", request.URL.String(), created.Id))
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(b)
}
