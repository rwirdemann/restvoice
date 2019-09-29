package domain

type Project struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customerId"`
	Name       string `json:"name"`
}
