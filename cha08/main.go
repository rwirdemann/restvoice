package main

import (
	"github.com/rwirdemann/restvoice/cha05/database"
	"github.com/rwirdemann/restvoice/cha05/usecase"
	"github.com/rwirdemann/restvoice/cha08/rest"
)

func main() {
	repository := database.NewMySQLRepository()
	r := rest.NewAdapter()

	createInvoiceHandler := r.MakeCreateInvoiceHandler(usecase.NewCreateInvoice(repository))
	r.HandleFunc("/invoice", rest.BasicAuth(createInvoiceHandler)).Methods("POST")

	r.ListenAndServe()
}
