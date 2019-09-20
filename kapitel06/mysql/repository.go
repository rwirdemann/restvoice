package mysql

import (
	"database/sql"

	"github.com/rwirdemann/restvoice/kapitel05/domain"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) CreateInvoice(invoice domain.Invoice) (domain.Invoice, error) {
	res, err := r.DB.Exec("insert into invoices (status, customer_id, year, month) VALUES (?, ?, ?, ?)",
		invoice.Status, invoice.CustomerID, invoice.Year, invoice.Month)
	if err != nil {
		return domain.Invoice{}, err
	}
	id, _ := res.LastInsertId()
	return r.GetInvoice(int(id)), nil
}

func (r *Repository) GetInvoice(id int, join ...string) domain.Invoice {
	return domain.Invoice{}
}

func (r *Repository) GetBookingsByInvoiceID(id int) []domain.Booking {
	var bookings []domain.Booking
	return bookings
}

func (r *Repository) GetProject(id int) domain.Project {
	return domain.Project{}
}

func (r *Repository) GetCustomer(id int) domain.Customer {
	return domain.Customer{}
}

func (r *Repository) UpdateInvoice(invoice domain.Invoice) error {
	return nil
}

func (r *Repository) CreateBooking(booking domain.Booking) (domain.Booking, error) {
	return domain.Booking{}, nil
}

func (r *Repository) CreateActivity(activity domain.Activity) {
}

func (r *Repository) ActivityByID(id int) domain.Activity {
	return domain.Activity{}
}

func (r *Repository) RateByProjectIDAndActivityID(projectID int, activityID int) domain.Rate {
	return domain.Rate{}
}

func (r *Repository) CreateRate(rate domain.Rate) {
}

func (r *Repository) CreateProject(p domain.Project) {
}
