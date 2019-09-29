package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rwirdemann/restvoice/kapitel05/domain"
)

type Link struct {
	Href string `json:"href"`
}

type Embedded struct {
	Bookings []domain.Booking `json:"bookings,omitempty"`
}

// HALInvoice dekoriert eine Invoice mit HAL-konformen _link-Elementen.
type HALInvoice struct {
	domain.Invoice
	Links    map[domain.Operation]Link `json:"_links"`
	Embedded *Embedded                 `json:"_embedded,omitempty"`
}

func NewHALInvoice(invoice domain.Invoice) HALInvoice {
	var links = make(map[domain.Operation]Link)
	links["self"] = Link{fmt.Sprintf("/invoice/%d", invoice.ID)}
	for _, o := range invoice.GetOperations() {
		if l, err := translate(o, invoice); err == nil {
			links[o] = l
		} else {
			log.Print(err)
		}
	}
	return HALInvoice{Invoice: invoice, Links: links}
}

func translate(operation domain.Operation, invoice domain.Invoice) (Link, error) {
	switch operation {
	case "book":
		return Link{fmt.Sprintf("/book/%d", invoice.ID)}, nil
	case "charge":
		return Link{fmt.Sprintf("/charge/%d", invoice.ID)}, nil
	case "cancel":
		return Link{fmt.Sprintf("/invoice/%d", invoice.ID)}, nil
	case "payment":
		return Link{fmt.Sprintf("/payment/%d", invoice.ID)}, nil
	case "archive":
		return Link{fmt.Sprintf("/payment/%d", invoice.ID)}, nil
	default:
		return Link{}, fmt.Errorf(fmt.Sprintf("No translation found for operation %s", operation))
	}
}

type HALInvoicePresenter struct {
	writer http.ResponseWriter
}

func NewHALInvoicePresenter(w http.ResponseWriter) HALInvoicePresenter {
	return HALInvoicePresenter{writer: w}
}

func (p HALInvoicePresenter) Present(i interface{}) {
	invoice := i.(HALInvoice)
	if len(invoice.Bookings) > 0 {
		invoice.Embedded = &Embedded{
			Bookings: invoice.Bookings,
		}
	}

	if b, err := json.Marshal(invoice); err == nil {
		p.writer.Write(b)
	}
}
