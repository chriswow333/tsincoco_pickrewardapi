package event

type ContainerEventResult struct {
	ID                         string                  `json:"id"`
	ContainerType              int32                   `json:"containerType"`
	EventID                    string                  `json:"eventID"`
	Pass                       bool                    `json:"pass"`
	InnerContainerEventResults []*ContainerEventResult `json:"containerEventResults"`
	Matches                    []string                `json:"matches"`
	MisMatches                 []string                `json:"misMatches"`
	// ConstraintEventResult      *ConstraintEventResult
}
