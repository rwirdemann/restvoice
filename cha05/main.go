package main

import (
	"github.com/rwirdemann/restvoice/cha05/database"
	"github.com/rwirdemann/restvoice/cha05/rest"
	"github.com/rwirdemann/restvoice/cha05/usecase"
)

func main() {
	repository := database.NewMySQLRepository()
	r := rest.NewAdapter()

	createInvoice := usecase.NewCreateInvoice(repository)
	createInvoiceHandler := r.MakeCreateInvoiceHandler(createInvoice)
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices", createInvoiceHandler).Methods("POST")

	updateInvoice := usecase.NewUpdateInvoice(repository)
	r.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}",
		r.MakeUpdateInvoiceHandler(updateInvoice)).Methods("PUT")

	createBooking := usecase.NewCreateBooking(repository)
	getInvoice := usecase.NewGetInvoice(repository)
	r.MakeCreateBookingHandler(createBooking)
	r.MakeGetInvoiceHandler(getInvoice)

	r.ListenAndServe()
}
