package usecase

import (
	"github.com/rwirdemann/restvoice/kapitel05/domain"
)

type GetActivitiesPort interface {
	GetActivities(userID string) []domain.Activity
}

type GetActivities struct {
	repository GetActivitiesPort
}

func NewGetActivities(repository GetActivitiesPort) GetActivities {
	return GetActivities{repository: repository}
}

func (u GetActivities) Run(userID string) []domain.Activity {
	return u.repository.GetActivities(userID)
}
