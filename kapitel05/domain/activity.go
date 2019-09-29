package domain

import (
	"fmt"
	"time"
)

type Activity struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	UserID  string    `json:"userId"` // belongs to user
	Updated time.Time `json:"-"`
}

func (a Activity) String() string {
	return fmt.Sprintf("Id: %d Name: %s", a.ID, a.Name)
}
