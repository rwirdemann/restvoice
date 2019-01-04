package usecase

import (
	"github.com/rwirdemann/restvoice/cha05/domain"
)

type CreateBookingPort interface {
	CreateBooking(booking domain.Booking) (domain.Booking, error)
}

type CreateBooking struct {
	repository CreateBookingPort
}

func NewCreateBooking(repository CreateBookingPort) CreateBooking {
	return CreateBooking{repository: repository}
}

func (u CreateBooking) Run(booking domain.Booking) (domain.Booking, error) {
	return u.repository.CreateBooking(booking)
}
