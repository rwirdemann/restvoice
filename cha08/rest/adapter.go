package rest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/rwirdemann/restvoice/kapitel05/domain"
	"github.com/rwirdemann/restvoice/kapitel06/usecase"

	"github.com/gorilla/mux"
)

type Adapter struct {
	r *mux.Router
}

func NewAdapter() Adapter {
	return Adapter{mux.NewRouter()}
}

func (a Adapter) ListenAndServe() {
	log.Printf("Listening on http://0.0.0.0%s\n", ":8080")
	http.ListenAndServe(":8080", a.r)
}

func (a Adapter) HandleFunc(path string, f func(http.ResponseWriter,
	*http.Request)) *mux.Route {
	return a.r.NewRoute().Path(path).HandlerFunc(f)
}

func (a Adapter) MakeCreateInvoiceHandler(createInvoice usecase.CreateInvoice) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func (a Adapter) MakeGetInvoiceHandler(getInvoice usecase.GetInvoice) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["invoiceId"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		join := ""
		q := r.URL.Query()
		if v, ok := q["expand"]; ok {
			join = v[0]
		}

		if presenter, ok := a.InvoicePresenter(w, r); ok {
			invoice := getInvoice.Run(id, join)
			presenter.Present(NewHALInvoice(invoice))
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
		}
	}
}

func (a Adapter) InvoicePresenter(w http.ResponseWriter, r *http.Request) (InvoicePresenter, bool) {
	switch r.Header.Get("Accept") {
	case "application/json":
		return NewJSONInvoicePresenter(w), true
	case "application/hal+json":
		return NewHALInvoicePresenter(w), true
	case "application/pdf":
		return NewPDFInvoicePresenter(w, r), true
	default:
		return NewJSONInvoicePresenter(w), true
	}
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
