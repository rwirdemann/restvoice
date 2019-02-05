package cli

import (
	"flag"

	"github.com/rwirdemann/restvoice/cha06/usecase"
	"github.com/rwirdemann/restvoice/kapitel05/domain"
)

type Adapter struct {
}

func (a Adapter) MakeCreateInvoiceHandler(createInvoice usecase.CreateInvoice) func() (domain.Invoice, error) {
	return func() (domain.Invoice, error) {
		month := flag.Int("month", 10, "a flag")
		year := flag.Int("year", 2018, "a flag")
		customerId := flag.Int("customerId", 1, "a flag")
		flag.Parse()
		invoice := domain.Invoice{Month: *month, Year: *year, CustomerId: *customerId}
		return createInvoice.Run(invoice)
	}
}
