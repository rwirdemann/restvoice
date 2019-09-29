package database

import (
	"strings"
	"time"

	"github.com/rwirdemann/restvoice/kapitel05/domain"
)

type FakeRepository struct {
	nextID     int
	invoices   map[int]domain.Invoice
	bookings   map[int]map[int]domain.Booking
	projects   map[int]domain.Project
	customers  map[int]domain.Customer
	activities map[string]map[int]domain.Activity
	rates      map[int]map[int]domain.Rate
}

func (r *FakeRepository) GetActivities(userID string) []domain.Activity {
	var activities []domain.Activity
	for _, a := range r.activities[userID] {
		activities = append(activities, a)
	}
	return activities
}

func NewFakeRepository() *FakeRepository {
	r := FakeRepository{}
	r.invoices = make(map[int]domain.Invoice)
	r.bookings = make(map[int]map[int]domain.Booking)
	r.projects = make(map[int]domain.Project)
	r.customers = make(map[int]domain.Customer)
	r.rates = make(map[int]map[int]domain.Rate)
	r.activities = make(map[string]map[int]domain.Activity)
	return &r
}

func (r *FakeRepository) GetBookingsByInvoiceID(id int) []domain.Booking {
	var bookings []domain.Booking
	for _, b := range r.bookings[id] {
		bookings = append(bookings, b)
	}
	return bookings
}

func (r *FakeRepository) GetInvoice(id int, join ...string) domain.Invoice {
	i := r.invoices[id]
	if len(join) > 0 {
		if strings.Contains(join[0], "bookings") {
			i.Bookings = r.GetBookingsByInvoiceID(id)
		}
	}
	return i
}

func (r *FakeRepository) GetProject(id int) domain.Project {
	return r.projects[id]
}

func (r *FakeRepository) GetCustomer(id int) domain.Customer {
	return r.customers[id]
}

func (r *FakeRepository) CreateInvoice(invoice domain.Invoice) (domain.Invoice, error) {
	if invoice.ID == 0 {
		invoice.ID = r.getNextID()
	}
	if invoice.Status == "" {
		invoice.Status = "open"
	}
	r.invoices[invoice.ID] = invoice
	return invoice, nil
}

func (r *FakeRepository) UpdateInvoice(invoice domain.Invoice) error {
	r.invoices[invoice.ID] = invoice
	return nil
}

func (r *FakeRepository) CreateBooking(booking domain.Booking) (domain.Booking, error) {
	booking.ID = r.getNextID()
	if bookings, ok := r.bookings[booking.InvoiceID]; ok {
		bookings[booking.ID] = booking
	} else {
		bookings := make(map[int]domain.Booking)
		bookings[booking.ID] = booking
		r.bookings[booking.InvoiceID] = bookings
	}
	return booking, nil
}

func (r *FakeRepository) CreateActivity(activity domain.Activity) {
	activity.ID = r.getNextID()
	activity.Updated = time.Now().UTC()
	if activities, ok := r.activities[activity.UserID]; ok {
		activities[activity.ID] = activity
	} else {
		activities := make(map[int]domain.Activity)
		activities[activity.ID] = activity
		r.activities[activity.UserID] = activities
	}
}

func (r *FakeRepository) ActivityByID(id int) domain.Activity {
	return r.activities[""][id]
}

func (r *FakeRepository) RateByProjectIDAndActivityID(projectID int, activityID int) domain.Rate {
	return r.rates[projectID][activityID]
}

func (r *FakeRepository) getNextID() int {
	r.nextID = r.nextID + 1
	return r.nextID
}

func (r *FakeRepository) CreateRate(rate domain.Rate) {
	if projectRates, ok := r.rates[rate.ProjectID]; ok {
		projectRates[rate.ActivityID] = rate
	} else {
		r.rates[rate.ProjectID] = make(map[int]domain.Rate)
		r.rates[rate.ProjectID][rate.ActivityID] = rate
	}
}

func (r *FakeRepository) CreateProject(p domain.Project) {
	p.ID = r.nextProjectID()
	r.projects[p.ID] = p
}

func (r *FakeRepository) nextProjectID() int {
	nextID := 1
	for _, i := range r.projects {
		if i.ID >= nextID {
			nextID = i.ID + 1
		}
	}
	return nextID
}
