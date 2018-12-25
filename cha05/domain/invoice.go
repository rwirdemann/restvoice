package domain

import (
	"io/ioutil"
	"time"
)

type Position struct {
	Hours float32
	Price float32
}

type Invoice struct {
	Id         int                         `json:"id"`
	Month      int                         `json:"month"`
	Year       int                         `json:"year"`
	Status     string                      `json:"status"`
	CustomerId int                         `json:"customerId"`
	Positions  map[int]map[string]Position `json:"positions,omitempty"`
	Bookings   []Booking                   `json:"-"`
	Updated    time.Time
}

func (i *Invoice) AddBooking() {

}
func (i *Invoice) AddPosition(projectId int, activity string, hours float32,
	rate float32) {
	if i.Positions == nil {
		i.Positions = make(map[int]map[string]Position)
	}

	if i.Positions[projectId] == nil {
		i.Positions[projectId] = make(map[string]Position)
	}

	if p, ok := i.Positions[projectId][activity]; ok {
		p.Hours = p.Hours + hours
		p.Price = p.Price + hours*rate
		i.Positions[projectId][activity] = p
	} else {
		position := Position{Hours: hours, Price: hours * rate}
		i.Positions[projectId][activity] = position
	}
}

func (i *Invoice) ToPDF() []byte {
	dat, _ := ioutil.ReadFile("/tmp/invoice.pdf")
	return dat
}

func (i Invoice) IsReadyForAggregation() bool {
	return i.Status == "ready for aggregation"
}

type Operation string

func (invoice Invoice) GetOperations() []Operation {
	switch invoice.Status {
	case "open":
		return []Operation{"book", "charge", "cancel", "bookings"}
	case "payment expected":
		return []Operation{"payment", "bookings"}
	case "payed":
		return []Operation{"archive"}
	case "archived":
		return []Operation{"revoke"}
	case "revoked":
		return []Operation{"archive"}
	default:
		return []Operation{}
	}
}
