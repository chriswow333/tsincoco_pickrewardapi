package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ConstraintType int32

const (
	NewCustomer  ConstraintType = iota // 新戶
	Register                           // 需登錄
	LimitCount                         // 限量
	LimitWeekDay                       // 限定日
)

var (
	constraintTypeMapper = make(map[ConstraintType]*Constraint)
)

func init() {

	constraintTypeMapper = map[ConstraintType]*Constraint{
		NewCustomer: {
			ConstraintType: NewCustomer,
			ConstraintName: "新戶",
		},
		Register: {
			ConstraintType: Register,
			ConstraintName: "登錄",
		},
		LimitCount: {
			ConstraintType: LimitCount,
			ConstraintName: "限量",
		},
		LimitWeekDay: {
			ConstraintType: LimitWeekDay,
			ConstraintName: "限定日",
		},
	}
}

type Constraint struct {
	ConstraintType ConstraintType `json:"constraintType"`
	ConstraintName string         `json:"constraintName"`
	WeekDays       []int32        `json:"weekDays"`
}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a Constraint) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *Constraint) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func GetAllConstraintTypes() []*Constraint {

	constraints := []*Constraint{}

	for _, v := range constraintTypeMapper {
		constraints = append(constraints, v)
	}
	return constraints
}

func GetConstraintType(constraintType int32) (*Constraint, error) {

	constraint, ok := constraintTypeMapper[ConstraintType(constraintType)]

	if !ok {
		return nil, errors.New("Cannot find currency type")
	}
	return constraint, nil
}
