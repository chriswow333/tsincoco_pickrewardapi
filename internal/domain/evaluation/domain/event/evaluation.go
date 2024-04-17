package event

type EvaluationEventStatus int32

const (
	EvaluationEventResultStatusALL EvaluationEventStatus = iota
	EvaluationEventResultStatusSOME
	EvaluationEventResultStatusNONE
)

type EvaluationEventResult struct {
	ID                  string               `json:"id"`
	FeedbackEventResult *FeedbackEventResult `json:"feedbackEventResult"`
	PayloadEventResult  *PayloadEventResult  `json:"payloadEventResult"`
}
