package domain

type Booking struct {
	Id          int     `json:"-"`
	Day         int     `json:"day"`
	Hours       float32 `json:"hours"`
	Description string  `json:"description"`
	InvoiceId   int     `json:"invoiceId"`            // belongs to invoice
	ProjectId   int     `json:"projectId,omitempty"`  // belongs to project
	ActivityId  int     `json:"activityId,omitempty"` // belongs to activity
}
