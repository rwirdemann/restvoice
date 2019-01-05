package main

import (
	"github.com/rwirdemann/restvoice/cha05/database"
	"github.com/rwirdemann/restvoice/cha06/usecase"
	"github.com/rwirdemann/restvoice/cha08/rest"
)

func main() {
	repository := database.NewRepository()
	adapter := rest.NewAdapter()

	createInvoice := usecase.NewCreateInvoice(repository)
	createInvoiceHandler := adapter.MakeCreateInvoiceHandler(createInvoice)
	adapter.HandleFunc("/invoice", createInvoiceHandler).Methods("POST")

	getInvoice := usecase.NewGetInvoice(repository)
	adapter.HandleFunc("/invoice/{invoiceId:[0-9]+}", adapter.MakeGetInvoiceHandler(getInvoice)).Methods("GET")

	adapter.ListenAndServe()
}
