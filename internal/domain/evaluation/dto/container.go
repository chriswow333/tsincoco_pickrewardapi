package dto

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ContainerDTO struct {
	ID                string `bson:"id" json:"id"`
	ContainerOperator int32  `bson:"container_operator" json:"containerOperator"`
	ContainerType     int32  `bson:"container_type" json:"containerType"`

	Containers []*ContainerDTO `bson:"containers" json:"containers"`

	TaskLabels []string `bson:"task_labels" json:"taskLabels"`

	ChannelLabels []string `bson:"channel_labels" json:"channelLabels"`

	ChannelIDs []string `bson:"channel_ids" json:"channelIDs"`

	PayIDs []string `bson:"pay_ids" json:"payIDs"`

	Constraints []*ConstraintDTO `bson:"constraints" json:"constraints"`
}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a ContainerDTO) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *ContainerDTO) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type ContainerEventResultDTO struct {
	ID                    string                     `json:"id"`
	Pass                  bool                       `json:"pass"`
	Matches               []string                   `json:"matches"`
	MisMatches            []string                   `json:"misMatches"`
	ContainerEventResults []*ContainerEventResultDTO `json:"containerEventResults"`
	// ConstraintEventResult *ConstraintEventResultDTO  `json:"constraintEventResult"`
}
