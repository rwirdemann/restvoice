package usecase

import (
	"github.com/rwirdemann/restvoice/cha05/domain"
)

type UpdateInvoiceRepositoryPort interface {
	UpdateInvoice(invoice domain.Invoice) error
	RateByProjectIdAndActivityId(projectId int, activityId int) domain.Rate
	ActivityById(user string, id int) domain.Activity
	GetBookingsByInvoiceId(invoiceId int) []domain.Booking
}

type UpdateInvoice struct {
	repository UpdateInvoiceRepositoryPort
}

func NewUpdateInvoice(repository UpdateInvoiceRepositoryPort) UpdateInvoice {
	return UpdateInvoice{repository: repository}
}

func (u UpdateInvoice) Run(invoice domain.Invoice) error {
	if invoice.Status == "payment expected" {
		bookings := u.repository.GetBookingsByInvoiceId(invoice.Id)
		for _, b := range bookings {
			activity := u.repository.ActivityById("ralf", b.ActivityId)
			rate := u.repository.RateByProjectIdAndActivityId(b.ProjectId, b.ActivityId)
			invoice.AddPosition(b.ProjectId, activity.Name, b.Hours, rate.Price)
		}
	}
	return u.repository.UpdateInvoice(invoice)
}
