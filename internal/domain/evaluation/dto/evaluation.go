package dto

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type EvaluationDTO struct {
	ID        string `json:"id"`
	Owner     int32  `json:"owner"`
	OwnerID   string `json:"ownerID"`
	StartDate int64  `json:"startDate"`
	EndDate   int64  `json:"endDate"`

	FeedbackID string `json:"feedbackID"`

	CreateDate int64       `json:"createDate"`
	UpdateDate int64       `json:"updateDate"`
	Payload    *PayloadDTO `json:"payload"`
}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a EvaluationDTO) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *EvaluationDTO) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type EvaluationEventResultDTO struct {
	ID                    string                 `json:"id"`
	EvaluationEventStatus int32                  `json:"evaluationEventStatus"`
	PayloadEventResult    *PayloadEventResultDTO `json:"payloadEventResult"`
}

type EvaluationEventResultRespDTO struct {
	ID                         string                  `json:"id"`
	FeedbackEventResultResp    *FeedbackEventResultDTO `json:"feedbackEventResultResp"`
	CardRewardTaskLabelMatched []string                `json:"cardRewardTaskLabelMatched"`
	ChannelMatched             []string                `json:"channelMatched"`
	ChannelLabelMatched        []string                `json:"channelLabelMatched"`
	PayMatched                 []string                `json:"payMatched"`
}

type EvaluationRespDTO struct {
	ID                   string  `json:"id"`
	ChannelCategoryTypes []int32 `json:"channelCategoryTypes"`

	ChannelsEvaluationResp            *ChannelsEvaluationRespDTO            `json:"channelEvaluationResp"`
	PayEvaluationResp                 *PayEvaluationRespDTO                 `json:"payEvaluationResp"`
	CardRewardTaskLabelEvaluationResp *CardRewardTaskLabelEvaluationRespDTO `json:"cardRewardTaskLabelEvaluationResp"`

	ConstraintsEvaluationResp  *ConstraintsEvaluationRespDTO  `json:"constraintEvaluationResp"`
	ChannelLabelEvaluationResp *ChannelLabelEvaluationRespDTO `json:"channelLabelEvaluationResp"`
}

type ChannelsEvaluationRespDTO struct {
	Matches    map[string]bool `json:"matches"`
	MisMatches map[string]bool `json:"misMatches"`
	// ChannelEvaluationRespMapper map[int32]*ChannelEvaluationRespDTO `json:"channelEvaluationRespMapper"`
}

// type ChannelEvaluationRespDTO struct {
// 	ChannelCategoryType int32           `json:"channelCategoryType"`
// 	Matches             map[string]bool `json:"matches"`
// 	MisMatches          map[string]bool `json:"misMatches"`
// }

type ChannelLabelEvaluationRespDTO struct {
	Matches    map[int32]bool `json:"matches"`
	MisMatches map[int32]bool `json:"misMatches"`
}

type CardRewardTaskLabelEvaluationRespDTO struct {
	Matches    map[int32]bool `json:"matches"`
	MisMatches map[int32]bool `json:"misMatches"`
}

type ConstraintsEvaluationRespDTO struct {
	Matches    map[int32]*ConstraintDTO `json:"matches"`
	MisMatches map[int32]*ConstraintDTO `json:"misMatches"`
}

type PayEvaluationRespDTO struct {
	Matches    map[string]bool `json:"matches"`
	MisMatches map[string]bool `json:"misMatches"`
}
