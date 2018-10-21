package main

import (
	"github.com/rwirdemann/restvoice/cha05/database"
	"github.com/rwirdemann/restvoice/cha05/rest"
	"github.com/rwirdemann/restvoice/cha05/usecase"
)

func main() {
	repository := database.NewMySQLRepository()
	createInvoice := usecase.NewCreateInvoice(repository)
	createBooking := usecase.NewCreateBooking(repository)
	getInvoice := usecase.NewGetInvoice(repository)

	restAdapter := rest.NewAdapter()
	restAdapter.MakeCreateInvoiceHandler(createInvoice)
	restAdapter.MakeCreateBookingHandler(createBooking)
	restAdapter.MakeGetInvoiceHandler(getInvoice)

	restAdapter.HandleFunc("/customers/{customerId:[0-9]+}/invoices/{invoiceId:[0-9]+}",
		restAdapter.MakeUpdateInvoiceHandler(repository)).Methods("PUT")

	restAdapter.ListenAndServe()
}
