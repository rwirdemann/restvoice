package domain

import (
	"fmt"
	"time"
)

type Activity struct {
	Id      int       `json:"id"`
	Name    string    `json:"name"`
	UserId  string    `json:"userId"` // belongs to user
	Updated time.Time `json:"-"`
}

func (a Activity) String() string {
	return fmt.Sprintf("Id: %d Name: %s", a.Id, a.Name)
}
