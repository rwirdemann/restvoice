package usecase

import (
	"github.com/rwirdemann/restvoice/kapitel05/domain"
)

type UpdateInvoicePort interface {
	UpdateInvoice(invoice domain.Invoice) error
	RateByProjectIDAndActivityID(projectID int, activityID int) domain.Rate
	ActivityByID(id int) domain.Activity
	GetBookingsByInvoiceID(invoiceID int) []domain.Booking
}

type UpdateInvoice struct {
	repository UpdateInvoicePort
}

func NewUpdateInvoice(repository UpdateInvoicePort) UpdateInvoice {
	return UpdateInvoice{repository: repository}
}

func (u UpdateInvoice) Run(invoice domain.Invoice) error {
	if invoice.IsReadyForAggregation() {
		bookings := u.repository.GetBookingsByInvoiceID(invoice.ID)
		for _, b := range bookings {
			activity := u.repository.ActivityByID(b.ActivityID)
			rate := u.repository.RateByProjectIDAndActivityID(b.ProjectID, b.ActivityID)
			invoice.AddPosition(b.ProjectID, activity.Name, b.Hours, rate.Price)
		}
		invoice.Status = "payment expected"
	}

	return u.repository.UpdateInvoice(invoice)
}
