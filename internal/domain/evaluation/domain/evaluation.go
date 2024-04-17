package domain

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"pickrewardapi/internal/domain/evaluation/domain/event"
	rewardDTO "pickrewardapi/internal/domain/reward/dto"
	commonM "pickrewardapi/internal/shared/common/model"
)

type Owner int32

const (
	CardReward Owner = iota
	Pay
	Channel
)

type Evaluation struct {
	ID string

	Reward *rewardDTO.RewardDTO

	Owner   Owner
	OwnerID string

	CreateDate int64
	UpdateDate int64
	StartDate  int64
	EndDate    int64

	Payload *Payload
}

type Description struct {
	Name  string
	Order int32
	Desc  []string
}

func (eva *Evaluation) Judge(e *commonM.Event) (*event.EvaluationEventResult, error) {
	logPos := "[evaluation.domain][Evaluation.Judge]"

	log.WithFields(log.Fields{
		"pos":          logPos,
		"EvaluationID": eva.ID,
		"eventID":      e.ID,
	}).Info("Start Evaluation Judge")

	evaluationEventResult := &event.EvaluationEventResult{
		ID: eva.ID,
	}

	payloadEventResult, err := eva.Payload.Judge(e)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":          logPos,
			"EvaluationID": eva.ID,
			"eventID":      e.ID,
		}).Error("Payload.Judge failed ", err)
		return nil, err
	}

	feedbackEventResult, err := eva.calculateFeedbackEventResult(payloadEventResult)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":          logPos,
			"EvaluationID": eva.ID,
			"eventID":      e.ID,
		}).Error("calculateFeedbackEventResult failed ", err)
		return nil, err
	}
	evaluationEventResult.PayloadEventResult = payloadEventResult
	evaluationEventResult.FeedbackEventResult = feedbackEventResult
	return evaluationEventResult, nil
}

func (eva *Evaluation) calculateFeedbackEventResult(payloadEventResult *event.PayloadEventResult) (*event.FeedbackEventResult, error) {
	logPos := "[evaluation.domain][calculateFeedbackEventResult]"

	if eva.Payload.PayloadType == ContainerPayloadType {
		return &event.FeedbackEventResult{
			RewardType:                payloadEventResult.FeedbackEventResult.RewardType,
			CalculateType:             payloadEventResult.FeedbackEventResult.CalculateType,
			Cost:                      payloadEventResult.FeedbackEventResult.Cost,
			GetReturn:                 payloadEventResult.FeedbackEventResult.GetReturn,
			GetPercentage:             payloadEventResult.FeedbackEventResult.GetPercentage,
			FeedbackEventResultStatus: payloadEventResult.FeedbackEventResult.FeedbackEventResultStatus,
		}, nil
	}

	if eva.Payload.PayloadType == SelfPayloadType {
		switch eva.Payload.PayloadOperator {
		case MaxAndPayloadOperator:
			return eva.maxAndCalculator(payloadEventResult), nil
		case MaxOrPayloadOperator:
			return eva.maxOrCalculator(payloadEventResult), nil
		case XorPayloadOperator:
			return eva.xorCalculator(payloadEventResult), nil
		case AddPayloadOperator:
			return eva.addCalculator(payloadEventResult), nil
		default:
			log.WithFields(log.Fields{
				"pos":          logPos,
				"EvaluationID": eva.ID,
			}).Error("Cannot find payload operator")
			return nil, errors.New("Cannot find payload operator")
		}
	}

	log.WithFields(log.Fields{
		"pos":          logPos,
		"EvaluationID": eva.ID,
	}).Error("Cannot find payload type")
	return nil, errors.New("Cannot find payload type")
}

func (eva *Evaluation) maxAndCalculator(payloadEventResult *event.PayloadEventResult) *event.FeedbackEventResult {

	noneFeedbackEventResult := &event.FeedbackEventResult{
		RewardType:                eva.Reward.RewardType,
		FeedbackEventResultStatus: event.GetNone,
	}

	maxFeedbackEventResult := &event.FeedbackEventResult{
		RewardType:                eva.Reward.RewardType,
		FeedbackEventResultStatus: event.GetNone,
	}

	for _, p := range payloadEventResult.PayloadEventResults {

		if p.FeedbackEventResult.FeedbackEventResultStatus == event.GetNone {
			noneFeedbackEventResult.Cost = p.FeedbackEventResult.Cost
			return noneFeedbackEventResult
		}

		if maxFeedbackEventResult.GetReturn < p.FeedbackEventResult.GetReturn {
			maxFeedbackEventResult = &event.FeedbackEventResult{
				RewardType:                eva.Reward.RewardType,
				CalculateType:             p.FeedbackEventResult.CalculateType,
				Cost:                      p.FeedbackEventResult.Cost,
				GetReturn:                 p.FeedbackEventResult.GetReturn,
				GetPercentage:             p.FeedbackEventResult.GetPercentage,
				FeedbackEventResultStatus: p.FeedbackEventResult.FeedbackEventResultStatus,
			}
		}
	}

	return maxFeedbackEventResult
}

func (eva *Evaluation) maxOrCalculator(payloadEventResult *event.PayloadEventResult) *event.FeedbackEventResult {

	maxFeedbackEventResult := &event.FeedbackEventResult{
		RewardType:                eva.Reward.RewardType,
		FeedbackEventResultStatus: event.GetNone,
	}

	for _, p := range payloadEventResult.PayloadEventResults {

		if maxFeedbackEventResult.GetReturn < p.FeedbackEventResult.GetReturn {
			maxFeedbackEventResult = &event.FeedbackEventResult{
				RewardType:                eva.Reward.RewardType,
				CalculateType:             p.FeedbackEventResult.CalculateType,
				Cost:                      p.FeedbackEventResult.Cost,
				GetReturn:                 p.FeedbackEventResult.GetReturn,
				GetPercentage:             p.FeedbackEventResult.GetPercentage,
				FeedbackEventResultStatus: p.FeedbackEventResult.FeedbackEventResultStatus,
			}
		}
	}

	return maxFeedbackEventResult
}

func (eva *Evaluation) xorCalculator(payloadEventResult *event.PayloadEventResult) *event.FeedbackEventResult {

	noneFeedbackEventResult := &event.FeedbackEventResult{
		RewardType:                eva.Reward.RewardType,
		FeedbackEventResultStatus: event.GetNone,
	}
	pass := false

	feedbackEventResult := &event.FeedbackEventResult{
		RewardType:                eva.Reward.RewardType,
		FeedbackEventResultStatus: event.GetNone,
	}

	for _, p := range payloadEventResult.PayloadEventResults {

		if pass && p.FeedbackEventResult.FeedbackEventResultStatus != event.GetNone {
			noneFeedbackEventResult.Cost = p.FeedbackEventResult.Cost
			return noneFeedbackEventResult
		}

		if p.FeedbackEventResult.FeedbackEventResultStatus == event.GetNone {
			continue
		}

		pass = true

		feedbackEventResult = &event.FeedbackEventResult{
			RewardType:                eva.Reward.RewardType,
			CalculateType:             p.FeedbackEventResult.CalculateType,
			Cost:                      p.FeedbackEventResult.Cost,
			GetReturn:                 p.FeedbackEventResult.GetReturn,
			GetPercentage:             p.FeedbackEventResult.GetPercentage,
			FeedbackEventResultStatus: p.FeedbackEventResult.FeedbackEventResultStatus,
		}
	}

	return feedbackEventResult
}

func (eva *Evaluation) addCalculator(payloadEventResult *event.PayloadEventResult) *event.FeedbackEventResult {

	addFeedbackEventResult := &event.FeedbackEventResult{
		RewardType:                eva.Reward.RewardType,
		FeedbackEventResultStatus: event.GetNone,
	}

	for _, p := range payloadEventResult.PayloadEventResults {

		if p.FeedbackEventResult.FeedbackEventResultStatus != event.GetNone {
			addFeedbackEventResult.Cost = p.FeedbackEventResult.Cost
			addFeedbackEventResult.GetReturn += p.FeedbackEventResult.GetReturn
			addFeedbackEventResult.GetPercentage += p.FeedbackEventResult.GetPercentage
			// TODO there have different calculate type.
			addFeedbackEventResult.CalculateType = p.FeedbackEventResult.CalculateType

			if addFeedbackEventResult.FeedbackEventResultStatus != event.GetSome {
				addFeedbackEventResult.FeedbackEventResultStatus = p.FeedbackEventResult.FeedbackEventResultStatus
			}
		}

	}

	return addFeedbackEventResult
}

func (eva *Evaluation) areaCalculator(payloadEventResult *event.PayloadEventResult) *event.FeedbackEventResult {

	areaFeedbackEventResult := &event.FeedbackEventResult{
		RewardType:                eva.Reward.RewardType,
		FeedbackEventResultStatus: event.GetNone,
	}

	// TODO

	return areaFeedbackEventResult
}
