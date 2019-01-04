package usecase

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rwirdemann/restvoice/cha05/domain"
	"github.com/rwirdemann/restvoice/cha06/database"
	"github.com/rwirdemann/restvoice/cha06/rest"
	"github.com/rwirdemann/restvoice/cha06/usecase"
	"github.com/rwirdemann/restvoice/cha07/usecase/mocks"
	"github.com/stretchr/testify/assert"
)

func TestShouldAggregateAndUpdateInvoice(t *testing.T) {
	// Setup
	repository := database.NewMySQLRepository()
	setupBaseData(repository)
	u := usecase.NewUpdateInvoice(repository)

	// Create Bookings Project 1
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 1,
		ActivityId: 1, Hours: 20, Description: "Steuerung umgestellt"})
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 1,
		ActivityId: 1, Hours: 12, Description: "Rating implementiert"})
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 1,

		ActivityId: 2, Hours: 3, Description: "Ratingtest"})
	// Create Bookings Project 2
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 2,
		ActivityId: 3, Hours: 4, Description: "Retrospektive geplant"})
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 2,
		ActivityId: 3, Hours: 3, Description: "Management Offsite"})
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 2,
		ActivityId: 2, Hours: 8, Description: "Suche getestet"})

	invoice := domain.Invoice{Id: 1, Status: "ready for aggregation"}

	// Run
	u.Run(invoice)

	// Assert
	expected := domain.Invoice{Id: 1, Status: "payment expected"}
	expected.AddPosition(1, "Programmierung", 32, 60)
	expected.AddPosition(1, "Qualitätssicherung", 3, 55)
	expected.AddPosition(2, "Projektmanagement", 7, 50)
	expected.AddPosition(2, "Qualitätssicherung", 8, 55)
	actual := repository.GetInvoice(1)
	assert.Equal(t, expected, actual)
}

func TestShouldUpdateState(t *testing.T) {
	// Setup
	repository := NewFakeRepository()
	u := usecase.NewUpdateInvoice(repository)

	// Run
	i := domain.Invoice{Id: 1, Status: "ready for aggregation"}
	repository.CreateInvoice(i)

	u.Run(i)

	// Assert
	actual := repository.GetInvoice(1)
	assert.Equal(t, "payment expected", actual.Status)
}

func TestShouldAggregateAndUpdateInvoiceWithFake(t *testing.T) {
	// Setup
	repository := NewFakeRepository()
	setupBaseData(repository)
	u := usecase.NewUpdateInvoice(repository)

	// Create Bookings Project 1
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 1,
		ActivityId: 1, Hours: 20, Description: "Steuerung umgestellt"})
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 1,
		ActivityId: 1, Hours: 12, Description: "Rating implementiert"})
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 1,
		ActivityId: 2, Hours: 3, Description: "Ratingtest"})

	// Create Bookings Project 2
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 2,
		ActivityId: 3, Hours: 4, Description: "Retrospektive geplant"})
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 2,
		ActivityId: 3, Hours: 3, Description: "Management Offsite"})
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 2,
		ActivityId: 2, Hours: 8, Description: "Suche getestet"})

	invoice := domain.Invoice{Id: 1, Status: "ready for aggregation"}
	repository.CreateInvoice(invoice)

	// Run
	u.Run(invoice)

	// Assert
	expected := domain.Invoice{Id: 1, Status: "payment expected"}
	expected.AddPosition(1, "Programmierung", 32, 60)
	expected.AddPosition(1, "Qualitätssicherung", 3, 55)
	expected.AddPosition(2, "Projektmanagement", 7, 50)
	expected.AddPosition(2, "Qualitätssicherung", 8, 55)
	actual := repository.GetInvoice(1)
	assert.Equal(t, expected, actual)
}

func TestShouldUpdateStateWithMock(t *testing.T) {
	// Setup
	repository := &mocks.UpdateInvoicePort{}
	u := usecase.NewUpdateInvoice(repository)

	// Setup mock interactions
	repository.On("GetBookingsByInvoiceId", 1).Return(nil)
	invoice := domain.Invoice{Id: 1, Status: "payment expected"}
	repository.On("UpdateInvoice", invoice).Return(nil)

	// Run
	u.Run(invoice)

	// Assert
	repository.AssertCalled(t, "UpdateInvoice", invoice)
}

func TestHttpInvoiceAggregation(t *testing.T) {
	// Setup
	repository := NewFakeRepository()
	setupBaseData(repository)
	u := usecase.NewUpdateInvoice(repository)

	// Create Bookings Project 1
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 1,
		ActivityId: 1, Hours: 20, Description: "Steuerung umgestellt"})
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 1,
		ActivityId: 1, Hours: 12, Description: "Rating implementiert"})
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 1,
		ActivityId: 2, Hours: 3, Description: "Ratingtest"})

	// Create Bookings Project 2
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 2,
		ActivityId: 3, Hours: 4, Description: "Retrospektive geplant"})
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 2,
		ActivityId: 3, Hours: 3, Description: "Management Offsite"})
	repository.CreateBooking(domain.Booking{InvoiceId: 1, ProjectId: 2,
		ActivityId: 2, Hours: 8, Description: "Suche getestet"})

	// Prepare HTTP-Request
	i := domain.Invoice{Id: 1, Status: "ready for aggregation"}
	repository.CreateInvoice(i)
	b, _ := json.Marshal(&i)
	r, _ := http.NewRequest("PUT", "/customers/1/invoices/1", bytes.NewReader(b))

	// Run
	response := httptest.NewRecorder()
	restAdapter := rest.NewAdapter()
	handler := http.HandlerFunc(restAdapter.MakeUpdateInvoiceHandler(u))
	handler.ServeHTTP(response, r)

	// Assert
	assert.Equal(t, http.StatusNoContent, response.Code)
	expected := domain.Invoice{Id: 1, Status: "payment expected"}
	expected.AddPosition(1, "Programmierung", 32, 60)
	expected.AddPosition(1, "Qualitätssicherung", 3, 55)
	expected.AddPosition(2, "Projektmanagement", 7, 50)
	expected.AddPosition(2, "Qualitätssicherung", 8, 55)
	assert.Equal(t, expected, repository.GetInvoice(1))
}

func setupBaseData(repository *database.FakeRepository) {
	repository.CreateProject(domain.Project{Id: 1, Name: "Instanfoo.com"})
	repository.CreateProject(domain.Project{Id: 2, Name: "Wo bleibt Kalle"})

	repository.CreateActivity(domain.Activity{Id: 1, Name: "Programmierung"})
	repository.CreateActivity(domain.Activity{Id: 2, Name: "Qualitätssicherung"})
	repository.CreateActivity(domain.Activity{Id: 3, Name: "Projektmanagement"})

	repository.CreateRate(domain.Rate{ProjectId: 1, ActivityId: 1, Price: 60}) // Programmierung
	repository.CreateRate(domain.Rate{ProjectId: 1, ActivityId: 2, Price: 55}) // Qualitätssicherung
	repository.CreateRate(domain.Rate{ProjectId: 2, ActivityId: 2, Price: 55}) // Qualitätssicherung
	repository.CreateRate(domain.Rate{ProjectId: 2, ActivityId: 3, Price: 50}) // Projektmanagement
}

func NewFakeRepository() *database.FakeRepository {
	return database.NewMySQLRepository()
}
