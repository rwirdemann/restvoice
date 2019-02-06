package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rwirdemann/restvoice/kapitel05/domain"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) CreateInvoice(invoice domain.Invoice) (domain.Invoice, error) {
	if res, err := r.DB.Exec("insert into invoices (status, customer_id, year, month) VALUES (?, ?, ?, ?)",
		invoice.Status, invoice.CustomerId, invoice.Year, invoice.Month); err == nil {
		id, _ := res.LastInsertId()
		return r.GetInvoice(int(id)), nil
	} else {
		return domain.Invoice{}, err
	}
}

func (r *Repository) GetInvoice(id int, join ...string) domain.Invoice {
	return domain.Invoice{}
}

func (r *Repository) GetBookingsByInvoiceId(id int) []domain.Booking {
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

func (r *Repository) ActivityById(id int) domain.Activity {
	return domain.Activity{}
}

func (r *Repository) RateByProjectIdAndActivityId(projectId int, activityId int) domain.Rate {
	return domain.Rate{}
}

func (r *Repository) CreateRate(rate domain.Rate) {
}

func (r *Repository) CreateProject(p domain.Project) {
}
