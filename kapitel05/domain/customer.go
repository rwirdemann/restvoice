package domain

type Customer struct {
	Id     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	UserId int    `json:"userId,omitempty"`
}
