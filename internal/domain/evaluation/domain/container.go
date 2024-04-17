package domain

import (
	"strconv"

	log "github.com/sirupsen/logrus"

	"pickrewardapi/internal/domain/evaluation/domain/event"
	"pickrewardapi/internal/domain/evaluation/dto"
	commonM "pickrewardapi/internal/shared/common/model"
)

type ContainerOperator int32

const (
	AndContainer ContainerOperator = iota
	OrContainer
	NotContainer
)

type ContainerType int32

const (
	InnerContainer ContainerType = iota
	ConstraintContainer
	CardRewardTaskLabelContainer
	ChannelContainer
	PayContainer
	ChannelLabelContainer
)

type Container struct {
	ID                string
	ContainerOperator ContainerOperator
	ContainerType     ContainerType
	InnerContainers   []*Container
	Constraints       []*Constraint

	CardRewardTaskLabels []int32

	ChannelEvaluations []*dto.ChannelEvaluationDTO
	PayEvaluations     []string

	ChannelLabelEvaluations []int32
}

func (c *Container) Satisfy(e *commonM.Event) *event.ContainerEventResult {
	logPos := "[evaluation.domain.container][Satisfy]"

	containerEventResult := &event.ContainerEventResult{
		ID:            c.ID,
		ContainerType: int32(c.ContainerType),
		EventID:       e.ID,
	}
	switch c.ContainerType {
	case InnerContainer:
		c.satisfyInnerContainer(e, containerEventResult)
	case CardRewardTaskLabelContainer:
		c.satisfyCardRewardTaskLabel(e, containerEventResult)
	case ChannelContainer:
		c.satisfyChannel(e, containerEventResult)
	case PayContainer:
		c.satisfyPay(e, containerEventResult)
	case ChannelLabelContainer:
		c.satisfyChannelLabel(e, containerEventResult)
	case ConstraintContainer:
		c.satisfyConstraint(e, containerEventResult)
	}

	pass := c.operate(containerEventResult)
	containerEventResult.Pass = pass
	log.WithFields(log.Fields{
		"pos":         logPos,
		"containerID": c.ID,
		"eventID":     e.ID,
		"pass":        pass,
	}).Info("in.Satisfy done")

	return containerEventResult
}

func (c *Container) satisfyConstraint(e *commonM.Event, containerEventResult *event.ContainerEventResult) {
	logPos := "[evaluation.domain.container][satisfyConstraint]"

	defer log.WithFields(log.Fields{
		"pos":                logPos,
		"containerID":        c.ID,
		"eventID":            e.ID,
		"container.match":    containerEventResult.Matches,
		"container.mismatch": containerEventResult.MisMatches,
	}).Info("in.satisfyLabel done")

	for _, cs := range c.Constraints {
		cs.Satisfy(e, containerEventResult)
	}

	log.WithFields(log.Fields{
		"pos":         logPos,
		"containerID": c.ID,
		"eventID":     e.ID,
	}).Info("in.satisfyConstraint done")
}

func (c *Container) satisfyChannelLabel(e *commonM.Event, containerEventResult *event.ContainerEventResult) {
	logPos := "[evaluation.domain.container][satisfyChannelLabel]"

	defer log.WithFields(log.Fields{
		"pos":                logPos,
		"containerID":        c.ID,
		"eventID":            e.ID,
		"container.match":    containerEventResult.Matches,
		"container.mismatch": containerEventResult.MisMatches,
	}).Info("in.satisfyLabel done")

	matchesMapper := make(map[string]bool)
	misMatchesMapper := make(map[string]bool)

	for _, l := range c.ChannelLabelEvaluations {

		match := false
		if _, ok := e.ChannelEvent.ChannelLabels[l]; ok {
			matchesMapper[strconv.FormatInt(int64(l), 10)] = true
			match = true
		} else {
			misMatchesMapper[strconv.FormatInt(int64(l), 10)] = true
		}

		if match {
			continue
		}

		for _, ec := range e.ChannelEvent.ChannelIDs {
			if _, ok := ec.ChannelLabels[l]; ok {
				matchesMapper[strconv.FormatInt(int64(l), 10)] = true
				match = true
				break
			} else {
				misMatchesMapper[strconv.FormatInt(int64(l), 10)] = true
			}
		}
	}

	matches := []string{}
	misMatches := []string{}

	for k := range matchesMapper {
		matches = append(matches, k)
	}

	for k := range misMatchesMapper {
		misMatches = append(misMatches, k)
	}

	containerEventResult.Matches = matches
	containerEventResult.MisMatches = misMatches

}

func (c *Container) satisfyPay(e *commonM.Event, containerEventResult *event.ContainerEventResult) {
	logPos := "[evaluation.domain.container][satisfyPay]"

	defer log.WithFields(log.Fields{
		"pos":                logPos,
		"containerID":        c.ID,
		"eventID":            e.ID,
		"container.match":    containerEventResult.Matches,
		"container.mismatch": containerEventResult.MisMatches,
	}).Info("in.satisfyPay done")

	matches := []string{}
	misMatches := []string{}

	for _, p := range c.PayEvaluations {
		if e.PayEvent == nil {
			misMatches = append(misMatches, p)
			continue
		}

		if e.PayEvent.Status == commonM.Use || e.PayEvent.Status == commonM.Whatever || e.PayEvent.PayIDs[p] {
			matches = append(matches, p)
		} else {
			misMatches = append(misMatches, p)
		}
	}

	containerEventResult.Matches = matches
	containerEventResult.MisMatches = misMatches
}

func (c *Container) satisfyChannel(e *commonM.Event, containerEventResult *event.ContainerEventResult) {
	logPos := "[evaluation.domain.container][satisfyChannel]"
	defer log.WithFields(log.Fields{
		"pos":                logPos,
		"containerID":        c.ID,
		"eventID":            e.ID,
		"container.match":    containerEventResult.Matches,
		"container.mismatch": containerEventResult.MisMatches,
	}).Info("in.satisfyPay done")

	matches := []string{}
	misMatches := []string{}

	for _, containerChannel := range c.ChannelEvaluations {
		if e.ChannelEvent == nil {
			misMatches = append(misMatches, containerChannel.ID)
			continue
		}

		match := false
		for _, eventChannel := range e.ChannelEvent.ChannelIDs {
			if eventChannel.ChannelID == containerChannel.ID {
				matches = append(matches, containerChannel.ID)
				match = true
				break
			}
		}

		if match {
			continue
		}

		for _, containerLabel := range containerChannel.ChannelLabels {
			if _, ok := e.ChannelEvent.ChannelLabels[containerLabel]; ok {
				matches = append(matches, containerChannel.ID)
				match = true
				break
			}
		}

		if !match {
			misMatches = append(misMatches, c.ID)
		}
	}

	containerEventResult.Matches = matches
	containerEventResult.MisMatches = misMatches
}

func (c *Container) satisfyCardRewardTaskLabel(e *commonM.Event, containerEventResult *event.ContainerEventResult) {
	logPos := "[evaluation.domain.container][satisfyCardRewardTaskLabel]"
	defer log.WithFields(log.Fields{
		"pos":                logPos,
		"containerID":        c.ID,
		"eventID":            e.ID,
		"container.match":    containerEventResult.Matches,
		"container.mismatch": containerEventResult.MisMatches,
	}).Info("in.satisfyCardRewardTaskLabel done")

	matches := []string{}
	misMatches := []string{}

	for _, t := range c.CardRewardTaskLabels {

		if e.CardEvent.TaskLabels[t] {
			matches = append(matches, string(t))
		} else {
			misMatches = append(misMatches, string(t))
		}
	}

	containerEventResult.Matches = matches
	containerEventResult.MisMatches = misMatches
}

func (c *Container) satisfyInnerContainer(e *commonM.Event, containerEventResult *event.ContainerEventResult) {
	logPos := "[evaluation.domain.container][satisfyInnerContainer]"

	containerEventResults := []*event.ContainerEventResult{}
	for _, in := range c.InnerContainers {
		innerContainerEvent := in.Satisfy(e)
		containerEventResults = append(containerEventResults, innerContainerEvent)
	}

	containerEventResult.InnerContainerEventResults = containerEventResults

	defer log.WithFields(log.Fields{
		"pos":             logPos,
		"containerID":     c.ID,
		"eventID":         e.ID,
		"innerContainers": containerEventResults,
	}).Info("in.satisfyInnerContainer done")

}

func (c *Container) operate(containerEvent *event.ContainerEventResult) bool {
	logPos := "[evaluation.domain.container][operate]"

	switch c.ContainerType {
	case ConstraintContainer, CardRewardTaskLabelContainer, ChannelContainer, PayContainer, ChannelLabelContainer:
		return c.operateContainer(containerEvent)
	case InnerContainer:
		return c.operateInnerContainer(containerEvent)
	default:
		log.WithFields(log.Fields{
			"pos":         logPos,
			"containerID": c.ID,
		}).Error("Cannot find container type")
		return false
	}

}

func (c *Container) operateContainer(containerEvent *event.ContainerEventResult) bool {
	logPos := "[evaluation.domain.container][operateContainer]"

	matches := containerEvent.Matches
	misMatches := containerEvent.MisMatches
	switch c.ContainerOperator {
	case AndContainer:
		if len(misMatches) > 0 {
			return false
		} else {
			return true
		}
	case OrContainer:
		if len(matches) > 0 {
			return true
		} else {
			return false
		}
	case NotContainer:
		if len(matches) > 0 {
			return false
		} else {
			return true
		}
	default:
		log.WithFields(log.Fields{
			"pos":         logPos,
			"containerID": c.ID,
		}).Error("Others container type has no DEFAULT type")
		return false
	}
}

func (c *Container) operateInnerContainer(containerEvent *event.ContainerEventResult) bool {
	logPos := "[evaluation.domain.container][operateInnerContainer]"

	switch c.ContainerOperator {
	case AndContainer:
		pass := true
		for _, e := range containerEvent.InnerContainerEventResults {
			if !e.Pass {
				pass = pass && e.Pass
			}
		}
		return pass
	case OrContainer:
		pass := false
		for _, e := range containerEvent.InnerContainerEventResults {
			pass = pass || e.Pass
		}
		return pass
	case NotContainer:
		log.WithFields(log.Fields{
			"pos":         logPos,
			"containerID": c.ID,
		}).Error("Inner container type has no NOT type")
		return false
	default:
		log.WithFields(log.Fields{
			"pos":         logPos,
			"containerID": c.ID,
		}).Error("Inner container type has no DEFAULT type")
		return false
	}
}
