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
	Status     string                      `json:"status"`
	CustomerId int                         `json:"customerId"`
	Year       int                         `json:"year"`
	Month      int                         `json:"month"`
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
