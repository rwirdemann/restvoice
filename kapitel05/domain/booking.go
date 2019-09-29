package domain

type Booking struct {
	ID          int     `json:"-"`
	Day         int     `json:"day"`
	Hours       float32 `json:"hours"`
	Description string  `json:"description"`
	InvoiceID   int     `json:"invoiceId"`            // belongs to invoice
	ProjectID   int     `json:"projectId,omitempty"`  // belongs to project
	ActivityID  int     `json:"activityId,omitempty"` // belongs to activity
}
