package domain

type Customer struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	UserID int    `json:"userId,omitempty"`
}
