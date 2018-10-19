package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/rwirdemann/restvoice/cha05/usecase"

	"github.com/gorilla/mux"
	"github.com/rwirdemann/restvoice/cha05/domain"
)

type Adapter struct {
	r *mux.Router
}

func NewAdapter() *Adapter {
	return &Adapter{mux.NewRouter()}
}

func (a Adapter) ListenAndServe() {
	log.Printf("Listening on http://0.0.0.0%s\n", ":8080")
	http.ListenAndServe(":8080", a.r)
}

func (a Adapter) readInvoice(r *http.Request) (domain.Invoice, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return domain.Invoice{}, err
	}

	var invoice domain.Invoice
	if err := json.Unmarshal(body, &invoice); err != nil {
		return domain.Invoice{}, err
	}

	if invoiceId, ok := mux.Vars(r)["invoiceId"]; ok {
		invoice.Id, _ = strconv.Atoi(invoiceId)
	}
	invoice.CustomerId, _ = strconv.Atoi(mux.Vars(r)["customerId"])
	return invoice, nil
}

func (a Adapter) writeInvoice(invoice domain.Invoice, w http.ResponseWriter) error {
	b, err := json.Marshal(invoice)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	return nil
}

func (a Adapter) MakeCreateInvoiceHandler(createInvoice usecase.CreateInvoice) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		invoice, err := a.readInvoice(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		created, err := createInvoice.Run(invoice)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = a.writeInvoice(created, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	a.r.HandleFunc("/customers/{customerId:[0-9]+}/invoices", handler).Methods("POST")
}

func (a Adapter) MakeCreateBookingHandler(createBooking usecase.CreateBooking) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		booking, err := a.readBooking(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		created, err := createBooking.Run(booking)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = a.writeBooking(created, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	a.r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}/bookings", handler).Methods("POST")
}

func (a Adapter) MakeUpdateInvoiceHandler(updateInvoice usecase.UpdateInvoice) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		invoice, err := a.readInvoice(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := updateInvoice.Run(invoice); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
	a.r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", handler).Methods("PUT")
}

func (a Adapter) MakeGetInvoiceHandler(getInvoice usecase.GetInvoice) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["invoiceId"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		invoice := getInvoice.Run(id)
		accept := r.Header.Get("Accept")
		switch accept {
		case "application/pdf":
			modTime := invoice.Updated
			content := bytes.NewReader(invoice.ToPDF())
			http.ServeContent(w, r, "invoice.pdf", modTime, content)
		case "application/json":
			b, _ := json.Marshal(invoice)
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		default:
			w.WriteHeader(http.StatusNotAcceptable)
		}
	}
	a.r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}", handler).Methods("GET")
}

type InvoicePresenter interface {
	Present(w http.ResponseWriter, i domain.Invoice)
}

type PDFInvoicePresenter struct {
	w http.ResponseWriter
	r *http.Request
}

func NewPDFInvoicePresenter(w http.ResponseWriter, r *http.Request) PDFInvoicePresenter {
	return PDFInvoicePresenter{w: w, r: r}
}

func (p PDFInvoicePresenter) Present(i domain.Invoice) {
	modTime := time.Now()
	content := bytes.NewReader(i.ToPDF())
	http.ServeContent(p.w, p.r, "invoice.pdf", modTime, content)
}

func (a Adapter) readBooking(r *http.Request) (domain.Booking, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return domain.Booking{}, err
	}

	var booking domain.Booking
	if err := json.Unmarshal(body, &booking); err != nil {
		return domain.Booking{}, err
	}

	booking.InvoiceId, _ = strconv.Atoi(mux.Vars(r)["invoiceId"])
	return booking, nil
}

func (a Adapter) writeBooking(booking domain.Booking, w http.ResponseWriter) error {
	b, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	return nil
}
