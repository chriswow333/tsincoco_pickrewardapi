package dto

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type FeedbackDTO struct {
	FeedbackID    string  `bson:"feedback_id" json:"feedbackID"`
	CalculateType int32   `bson:"calculate_type" json:"calculateType"`
	MinCost       int32   `bson:"min_cost" json:"minCost"`
	Fixed         int32   `bson:"fixed" json:"fixed"`
	Percentage    float64 `bson:"percentage" json:"percentage"`
	ReturnMax     float64 `bson:"return_max" json:"returnMax"`
}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a FeedbackDTO) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *FeedbackDTO) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type FeedbackEventResultDTO struct {
	FeedbackID                string  `bson:"feedback_id" json:"feedbackID"`
	CalculateType             int32   `bson:"calculate_type" json:"calculateType"`
	Cost                      int32   `bson:"cost" json:"cost"`
	GetReturn                 float64 `bson:"get_return" json:"getReturn"`
	GetPercentage             float64 `bson:"get_percentage" json:"getPercentage"`
	FeedbackEventResultStatus int32   `bson:"feedback_event_result_status" json:"feedbackEventResultStatus"`
}
