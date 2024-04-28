package dto

type FeedbackType int32

const (
	Currency FeedbackType = iota
	Point
)

type FeedbackTypeDTO struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	FeedbackType FeedbackType `json:"feedbackType"`
	CreateDate   int64        `json:"createDate"`
	UpdateDate   int64        `json:"updateDate"`
}
