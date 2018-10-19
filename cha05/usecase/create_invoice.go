package usecase

import (
	"github.com/rwirdemann/restvoice/cha05/domain"
)

type CreateInvoicePort interface {
	CreateInvoice(invoice domain.Invoice) (domain.Invoice, error)
}

type CreateInvoice struct {
	repository CreateInvoicePort
}

func NewCreateInvoice(repository CreateInvoicePort) CreateInvoice {
	return CreateInvoice{repository: repository}
}

func (u CreateInvoice) Run(invoice domain.Invoice) (domain.Invoice, error) {
	return u.repository.CreateInvoice(invoice)
}
