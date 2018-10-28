package main

import (
	"github.com/rwirdemann/restvoice/cha05/database"
	"github.com/rwirdemann/restvoice/cha05/usecase"
	"github.com/rwirdemann/restvoice/cha09/rest"
	"github.com/rwirdemann/restvoice/cha09/roles"
)

func main() {
	repository := database.NewMySQLRepository()
	r := rest.NewAdapter()

	createInvoiceHandler := r.MakeCreateInvoiceHandler(usecase.NewCreateInvoice(repository))
	r.HandleFunc("/invoice", rest.BasicAuth(createInvoiceHandler)).Methods("POST")

	createBookingHandler := r.MakeCreateBookingHandler(usecase.NewCreateBooking(repository))
	r.HandleFunc("/booking/{invoiceId:[0-9]+}", rest.BasicAuth(
		roles.AssertOwnsInvoice(createBookingHandler, repository))).Methods("POST")

	r.ListenAndServe()
}
