package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rwirdemann/restvoice/cha05/database"
	"github.com/rwirdemann/restvoice/cha05/domain"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices", createInvoiceHandler).Methods("POST")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}/bookings", createBookingHandler).Methods("POST")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}/bookings/{bookingId:[0-9]+}", deleteBookingHandler).Methods("DELETE")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", updateInvoiceHandler).Methods("PUT")
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", readInvoiceHandler).Methods("GET")

	fmt.Println("Restvoice started on http://localhost:8080...")
	http.ListenAndServe(":8080", r)
}

var repository = database.NewRepository()

func createInvoiceHandler(writer http.ResponseWriter, request *http.Request) {
	// Read invoice data from request body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// CreateInvoice invoice and marshal it from JSON
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
	if i.Status == "ready for aggregation" {
		bookings := repository.BookingsByInvoiceId(i.Id)
		for _, b := range bookings {
			activity := repository.GetActivity(b.ActivityId)
			rate := repository.GetRate(b.ProjectId, b.ActivityId)

			i.AddPosition(b.ProjectId, activity.Name, b.Hours, rate.Price)
		}
		i.Status = "payment expected"
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
		content := bytes.NewReader(i.ToPDF())
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
