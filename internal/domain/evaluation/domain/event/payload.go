package event

type PayloadEventResult struct {
	ID                   string                `json:"id"`
	Pass                 bool                  `json:"pass"`
	FeedbackEventResult  *FeedbackEventResult  `json:"feedbackEventResult"`
	PayloadEventResults  []*PayloadEventResult `json:"payloadEventResults"`
	ContainerEventResult *ContainerEventResult `json:"containerEventResult"`
}
