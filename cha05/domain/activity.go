package domain

import (
	"time"
	"fmt"
)

type Activity struct {
	Id      int
	Name    string    `json:"name"`
	UserId  string    `json:"userId"` // belongs to user
	Updated time.Time `json:"-"`
}

func (activity Activity) String() string {
	return fmt.Sprintf("Id: %d Name: %s", activity.Id, activity.Name)
}
