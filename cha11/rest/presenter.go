package rest

import (
	"encoding/json"
	"time"

	"github.com/rwirdemann/restvoice/kapitel05/domain"
)

type CacheableActivities struct {
	Activities   []byte
	LastModified time.Time
}

type ActivitiesPresenter struct {
}

func (j ActivitiesPresenter) Present(i interface{}) CacheableActivities {
	lastModified := time.Unix(0, 0)
	activities := i.([]domain.Activity)
	for _, a := range activities {
		if a.Updated.After(lastModified) {
			lastModified = a.Updated
		}
	}
	b, _ := json.Marshal(i)
	return CacheableActivities{Activities: b, LastModified: lastModified}
}
