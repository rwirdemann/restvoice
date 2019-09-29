package domain

type Rate struct {
	ProjectID  int     `json:"projectId"`
	ActivityID int     `json:"activityId"`
	Price      float32 `json:"price"`
}
