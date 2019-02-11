package main

import (
	"github.com/rwirdemann/restvoice/kapitel06/database"
	"github.com/rwirdemann/restvoice/kapitel06/usecase"
	"github.com/rwirdemann/restvoice/kapitel09/rest"
	"github.com/rwirdemann/restvoice/kapitel09/roles"
)

func main() {
	repository := database.NewFakeRepository()
	r := rest.NewAdapter()

	createInvoiceHandler := r.MakeCreateInvoiceHandler(usecase.NewCreateInvoice(repository))
	createBookingHandler := r.MakeCreateBookingHandler(usecase.NewCreateBooking(repository))

	r.HandleFunc("/invoice", rest.JWTAuth(roles.AssertAdmin(createInvoiceHandler))).Methods("POST")
	r.HandleFunc("/book/{invoiceId:[0-9]+}", rest.JWTAuth(createBookingHandler)).Methods("POST")

	r.ListenAndServe()
}
