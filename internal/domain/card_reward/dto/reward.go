package dto

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type RewardDTO struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	RewardType int32  `json:"rewardType"`

	CreateDate int64 `json:"createDate"`
	UpdateDate int64 `json:"updateDate"`
}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a RewardDTO) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *RewardDTO) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
