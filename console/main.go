package main

import (
	"fmt"
	"os"

	"github.com/rwirdemann/restvoice/cha05/database"
	"github.com/rwirdemann/restvoice/console/cli"

	"github.com/rwirdemann/restvoice/cha05/usecase"
)

func main() {
	repository := database.NewMySQLRepository()
	createInvoice := usecase.NewCreateInvoice(repository)

	cliAdapter := cli.Adapter{}
	createInvoiceHandler := cliAdapter.MakeCreateInvoiceHandler(createInvoice)

	if contains("create-invoice", os.Args) {
		created, _ := createInvoiceHandler()
		fmt.Println("Created:", created)
	}
}

func contains(s string, a []string) bool {
	for _, e := range a {
		if e == s {
			return true
		}
	}
	return false
}
