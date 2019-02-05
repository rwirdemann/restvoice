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
	_, _ = r.CreateInvoice(domain.Invoice{Month: 6, Year: 2018, CustomerId: 1})
	return r
}

func (r *Repository) GetCustomers() []domain.Customer {
	var customers []domain.Customer
	for _, c := range r.customers {
		customers = append(customers, c)
	}
	return customers
}

func (r *Repository) GetProjects(customerId int) []domain.Project {
	var projects []domain.Project
	for _, p := range r.projects {
		if p.CustomerId == customerId {
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
	a := domain.Activity{Id: r.nextActivityId(), Name: name}
	r.activities[a.Id] = a
	return a.Id
}

func (r *Repository) AddCustomer(name string) int {
	c := domain.Customer{Id: r.nextCustomerId(), Name: name}
	r.customers[c.Id] = c
	return c.Id
}

func (r *Repository) AddProject(name string, customerId int) int {
	p := domain.Project{Id: r.nextProjectId(), Name: name, CustomerId: customerId}
	r.projects[p.Id] = p
	return p.Id
}

func (r *Repository) CreateInvoice(invoice domain.Invoice) (domain.Invoice, error) {
	invoice.Id = r.nextInvoiceId()
	invoice.Status = "open"
	invoice.Bookings = []domain.Booking{}
	r.invoices[invoice.Id] = invoice
	return invoice, nil
}

func (r *Repository) CreateBooking(booking domain.Booking) (domain.Booking, error) {
	booking.Id = r.nextBookingId()
	r.bookings[booking.Id] = booking
	return booking, nil
}

func (r *Repository) DeleteBooking(id int) {
	delete(r.bookings, id)
}

func (r *Repository) GetBookingsByInvoiceId(invoiceId int) []domain.Booking {
	var bookings []domain.Booking
	for _, b := range r.bookings {
		if b.InvoiceId == invoiceId {
			bookings = append(bookings, b)
		}
	}
	return bookings
}

func (r *Repository) Update(invoice domain.Invoice) {
	r.invoices[invoice.Id] = invoice
}

func (r *Repository) FindById(id int) (domain.Invoice, bool) {
	i, ok := r.invoices[id]
	return i, ok
}

func (r *Repository) nextInvoiceId() int {
	nextId := 1
	for _, v := range r.invoices {
		if v.Id >= nextId {
			nextId = v.Id + 1
		}
	}
	return nextId
}

func (r *Repository) nextCustomerId() int {
	nextId := 1
	for _, v := range r.customers {
		if v.Id >= nextId {
			nextId = v.Id + 1
		}
	}
	return nextId
}

func (r *Repository) nextProjectId() int {
	nextId := 1
	for _, v := range r.projects {
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

func (r *Repository) nextActivityId() int {
	nextId := 1
	for _, v := range r.activities {
		if v.Id >= nextId {
			nextId = v.Id + 1
		}
	}
	return nextId
}

func (r *Repository) ActivityById(id int) domain.Activity {
	return r.activities[id]
}

func (r *Repository) GetInvoice(id int, join ...string) domain.Invoice {
	return r.invoices[id]
}

func (r *Repository) RateByProjectIdAndActivityId(projectId int, activityId int) domain.Rate {
	return domain.Rate{}
}

func (r *Repository) UpdateInvoice(invoice domain.Invoice) error {
	return nil
}
