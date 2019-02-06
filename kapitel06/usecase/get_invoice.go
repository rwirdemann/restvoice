package usecase

import (
	"github.com/rwirdemann/restvoice/kapitel05/domain"
)

type GetInvoicePort interface {
	GetInvoice(id int, join ...string) domain.Invoice
}

type GetInvoice struct {
	repository GetInvoicePort
}

func NewGetInvoice(repository GetInvoicePort) GetInvoice {
	return GetInvoice{repository: repository}
}

func (u GetInvoice) Run(id int, join string) domain.Invoice {
	return u.repository.GetInvoice(id, join)
}
