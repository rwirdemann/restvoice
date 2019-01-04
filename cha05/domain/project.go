package domain

type Project struct {
	Id         int    `json:"id"`
	CustomerId int    `json:"customerId"`
	Name       string `json:"name"`
}
