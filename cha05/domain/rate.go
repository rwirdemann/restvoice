package domain

type Rate struct {
	ProjectId  int     `json:"projectId"`
	ActivityId int     `json:"activityId"`
	Price      float32 `json:"price"`
}
