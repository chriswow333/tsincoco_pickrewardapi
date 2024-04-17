package dto

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type PayloadDTO struct {
	ID              string        `bson:"id" json:"id"`
	PayloadOperator int32         `bson:"payload_operator" json:"payloadOperator"`
	PayloadType     int32         `bson:"payload_type" json:"payloadType"`
	Feedback        *FeedbackDTO  `bson:"feedback" json:"feedback"`
	Payloads        []*PayloadDTO `bson:"payloads" json:"payloads"`
	Container       *ContainerDTO `bson:"container" json:"container"`
}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a PayloadDTO) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *PayloadDTO) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type PayloadEventResultDTO struct {
	ID                   string                   `json:"id"`
	Pass                 bool                     `json:"pass"`
	FeedbackEventResult  *FeedbackEventResultDTO  `json:"feedbackEventResult"`
	PayloadEventResults  []*PayloadEventResultDTO `json:"payloadEventResult"`
	ContainerEventResult *ContainerEventResultDTO `json:"containerEventResult"`
}
