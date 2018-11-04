package database

import "github.com/rwirdemann/restvoice/cha05/domain"

type Repository struct {
	invoices   map[int]domain.Invoice
	bookings   map[int]domain.Booking
	activities map[int]domain.Activity
	rates      map[int]map[int]domain.Rate
}

func NewRepository() *Repository {
	r := &Repository{invoices: make(map[int]domain.Invoice)}
	r.CreateInvoice(domain.Invoice{Month: 6, Year: 2018, CustomerId: 1})
	return r
}

func (r *Repository) CreateInvoice(invoice domain.Invoice) domain.Invoice {
	invoice.Id = r.nextId()
	invoice.Status = "open"
	invoice.Bookings = []domain.Booking{}
	r.invoices[invoice.Id] = invoice
	return invoice
}

func (r *Repository) CreateBooking(booking domain.Booking) domain.Booking {
	booking.Id = r.nextBookingId()
	r.bookings[booking.Id] = booking
	return booking
}

func (r *Repository) DeleteBooking(id int) {
	delete(r.bookings, id)
}

func (r *Repository) Update(invoice domain.Invoice) {
	r.invoices[invoice.Id] = invoice
}

func (r *Repository) FindById(id int) (domain.Invoice, bool) {
	i, ok := r.invoices[id]
	return i, ok
}

func (r *Repository) nextId() int {
	nextId := 1
	for _, v := range r.invoices {
		if v.Id >= nextId {
			nextId = v.Id + 1
		}
	}
	return nextId
}

func (r *Repository) nextBookingId() int {
	nextId := 1
	for _, v := range r.bookings {
		if v.Id >= nextId {
			nextId = v.Id + 1
		}
	}
	return nextId
}

func (r *Repository) GetActivity(id int) domain.Activity {
	return r.activities[id]
}

func (r *Repository) GetRate(projectId int, activityId int) domain.Rate {
	return r.rates[projectId][activityId]
}
