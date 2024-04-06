package dto

type TaskLabelDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Show  int32  `json:"show"`
	Order int32  `json:"order"`
}
