package dto

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ConstraintDTO struct {
	ConstraintType int32   `json:"constraintType"`
	ConstraintName string  `json:"constraintName"`
	WeekDays       []int32 `json:"weekDays"`
}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a ConstraintDTO) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *ConstraintDTO) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

// type ConstraintEventResultDTO struct {
// 	Pass           bool     `json:"pass"`
// 	ConstraintType int32    `json:"constraintType"`
// 	Matches        []string `json:"matches"`
// 	MisMatches     []string `json:"misMatches"`
// }
