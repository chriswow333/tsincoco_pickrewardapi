package domain

import (
	"time"

	log "github.com/sirupsen/logrus"

	"pickrewardapi/internal/domain/evaluation/domain/event"
	commonM "pickrewardapi/internal/shared/common/model"
)

type ConstraintType int32

const (
	NewCustomer  ConstraintType = iota // 新戶
	Register                           // 需登錄
	LimitCount                         // 限量
	LimitWeekDay                       // 限定日
)

type Constraint struct {
	ConstraintType ConstraintType `json:"constraintType"`
	ConstraintName string         `json:"constraintName"`
	WeekDays       []int32        `json:"weekDays"`
}

func (c *Constraint) Satisfy(e *commonM.Event, containerEventResult *event.ContainerEventResult) {
	logPos := "[reward.domain][Constraint.Satisfy]"

	logFields := log.Fields{
		"pos":     logPos,
		"eventID": e.ID,
	}

	switch c.ConstraintType {
	case NewCustomer, Register, LimitCount:
		// default is matched for now
		containerEventResult.Matches = append(containerEventResult.Matches, string(c.ConstraintType))
		// if _, ok := e.Constraints[c.ConstraintType]; ok {
		// 	containerEventResult.Matches = append(containerEventResult.Matches, string(c.ConstraintType))
		// } else {
		// 	containerEventResult.MisMatches = append(containerEventResult.MisMatches, string(c.ConstraintType))
		// }
	case LimitWeekDay:
		c.weekDayConstraint(e, containerEventResult)
	default:
		log.WithFields(logFields).Error("Cannot find match ConstraintType")
	}
}

func (c *Constraint) weekDayConstraint(e *commonM.Event, containerEventResult *event.ContainerEventResult) {
	logPos := "[reward.domain][Constraint.weekDayConstraint]"

	weekDay := time.Unix(e.Date, 0).Weekday()
	matches := []string{}
	misMatches := []string{}

	for _, w := range c.WeekDays {
		if w == int32(weekDay) {
			matches = append(matches, string(w))
		} else {
			misMatches = append(misMatches, string(w))
		}
	}

	if len(matches) > 0 {
		containerEventResult.Matches = append(containerEventResult.Matches, string(commonM.LimitWeekDay))
	} else {
		containerEventResult.MisMatches = append(containerEventResult.MisMatches, string(commonM.LimitWeekDay))
	}

	log.WithFields(log.Fields{
		"pos":        logPos,
		"eventID":    e.ID,
		"matches":    containerEventResult.Matches,
		"misMatches": containerEventResult.MisMatches,
	}).Infof("week day constraint finished")

}
