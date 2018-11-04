package usecase

import (
	"github.com/rwirdemann/restvoice/cha05/domain"
)

type GetActivitiesPort interface {
	GetActivities(userId string) []domain.Activity
}

type GetActivities struct {
	repository GetActivitiesPort
}

func NewGetActivities(repository GetActivitiesPort) GetActivities {
	return GetActivities{repository: repository}
}

func (u GetActivities) Run(userId string) []domain.Activity {
	return u.repository.GetActivities(userId)
}
