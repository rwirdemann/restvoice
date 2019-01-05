package main

import (
	"github.com/rwirdemann/restvoice/cha06/database"
	"github.com/rwirdemann/restvoice/cha06/usecase"
	"github.com/rwirdemann/restvoice/cha09/rest"
	"github.com/rwirdemann/restvoice/cha09/roles"
)

func main() {
	repository := database.NewMySQLRepository()
	r := rest.NewAdapter()

	createInvoiceHandler := r.MakeCreateInvoiceHandler(usecase.NewCreateInvoice(repository))
	createBookingHandler := r.MakeCreateBookingHandler(usecase.NewCreateBooking(repository))

	r.HandleFunc("/invoice", rest.DigestAuth(roles.AssertAdmin(createInvoiceHandler))).Methods("POST")
	r.HandleFunc("/book/{invoiceId:[0-9]+}", rest.JWTAuth(createBookingHandler)).Methods("POST")

	r.ListenAndServe()
}
