package database

import "github.com/rwirdemann/restvoice/kapitel05/domain"

type Repository struct {
	customers  map[int]domain.Customer
	projects   map[int]domain.Project
	invoices   map[int]domain.Invoice
	bookings   map[int]domain.Booking
	activities map[int]domain.Activity
	rates      map[int]map[int]domain.Rate
}

func NewRepository() *Repository {
	r := &Repository{
		customers:  make(map[int]domain.Customer),
		projects:   make(map[int]domain.Project),
		invoices:   make(map[int]domain.Invoice),
		bookings:   make(map[int]domain.Booking),
		activities: make(map[int]domain.Activity),
		rates:      make(map[int]map[int]domain.Rate)}
	_, _ = r.CreateInvoice(domain.Invoice{Month: 6, Year: 2018, CustomerID: 1})
	return r
}

func (r *Repository) GetCustomers() []domain.Customer {
	var customers []domain.Customer
	for _, c := range r.customers {
		customers = append(customers, c)
	}
	return customers
}

func (r *Repository) GetProjects(customerID int) []domain.Project {
	var projects []domain.Project
	for _, p := range r.projects {
		if p.CustomerID == customerID {
			projects = append(projects, p)
		}
	}
	return projects
}

func (r *Repository) GetActivities() []domain.Activity {
	var activities []domain.Activity
	for _, c := range r.activities {
		activities = append(activities, c)
	}
	return activities
}

func (r *Repository) AddActivity(name string) int {
	a := domain.Activity{ID: r.nextActivityID(), Name: name}
	r.activities[a.ID] = a
	return a.ID
}

func (r *Repository) AddCustomer(name string) int {
	c := domain.Customer{ID: r.nextCustomerID(), Name: name}
	r.customers[c.ID] = c
	return c.ID
}

func (r *Repository) AddProject(name string, customerID int) int {
	p := domain.Project{ID: r.nextProjectID(), Name: name, CustomerID: customerID}
	r.projects[p.ID] = p
	return p.ID
}

func (r *Repository) CreateInvoice(invoice domain.Invoice) (domain.Invoice, error) {
	invoice.ID = r.nextInvoiceID()
	invoice.Status = "open"
	invoice.Bookings = []domain.Booking{}
	r.invoices[invoice.ID] = invoice
	return invoice, nil
}

func (r *Repository) CreateBooking(booking domain.Booking) (domain.Booking, error) {
	booking.ID = r.nextBookingID()
	r.bookings[booking.ID] = booking
	return booking, nil
}

func (r *Repository) DeleteBooking(id int) {
	delete(r.bookings, id)
}

func (r *Repository) GetBookingsByInvoiceID(invoiceID int) []domain.Booking {
	var bookings []domain.Booking
	for _, b := range r.bookings {
		if b.InvoiceID == invoiceID {
			bookings = append(bookings, b)
		}
	}
	return bookings
}

func (r *Repository) Update(invoice domain.Invoice) {
	r.invoices[invoice.ID] = invoice
}

func (r *Repository) FindByID(id int) (domain.Invoice, bool) {
	i, ok := r.invoices[id]
	return i, ok
}

func (r *Repository) nextInvoiceID() int {
	nextID := 1
	for _, v := range r.invoices {
		if v.ID >= nextID {
			nextID = v.ID + 1
		}
	}
	return nextID
}

func (r *Repository) nextCustomerID() int {
	nextID := 1
	for _, v := range r.customers {
		if v.ID >= nextID {
			nextID = v.ID + 1
		}
	}
	return nextID
}

func (r *Repository) nextProjectID() int {
	nextID := 1
	for _, v := range r.projects {
		if v.ID >= nextID {
			nextID = v.ID + 1
		}
	}
	return nextID
}

func (r *Repository) nextBookingID() int {
	nextID := 1
	for _, v := range r.bookings {
		if v.ID >= nextID {
			nextID = v.ID + 1
		}
	}
	return nextID
}

func (r *Repository) nextActivityID() int {
	nextID := 1
	for _, v := range r.activities {
		if v.ID >= nextID {
			nextID = v.ID + 1
		}
	}
	return nextID
}

func (r *Repository) ActivityByID(id int) domain.Activity {
	return r.activities[id]
}

func (r *Repository) GetInvoice(id int, join ...string) domain.Invoice {
	return r.invoices[id]
}

func (r *Repository) RateByProjectIDAndActivityID(projectID int, activityID int) domain.Rate {
	return domain.Rate{}
}
