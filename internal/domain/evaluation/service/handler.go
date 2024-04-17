package service

import (
	domain "pickrewardapi/internal/app/evaluation/domain"
	event "pickrewardapi/internal/app/evaluation/domain/event"
	evaluationDTO "pickrewardapi/internal/app/evaluation/dto"
)

func calculateEvaluationResult(eventResult *event.EvaluationEventResult) *evaluationDTO.EvaluationEventResultRespDTO {

	evaluationEventResultResp := &evaluationDTO.EvaluationEventResultRespDTO{
		CardRewardTaskLabelMatched: []string{},
		ChannelMatched:             []string{},
		PayMatched:                 []string{},
		ChannelLabelMatched:        []string{},
	}

	containerTypeMatchedMap := make(map[int32]map[string]bool)

	processMatchedPayload(eventResult.PayloadEventResult, containerTypeMatchedMap)

	for k, v := range containerTypeMatchedMap {
		switch k {
		case int32(domain.CardRewardTaskLabelContainer):
			for label := range v {
				evaluationEventResultResp.CardRewardTaskLabelMatched = append(evaluationEventResultResp.CardRewardTaskLabelMatched, label)
			}
		case int32(domain.ChannelContainer):
			for id := range v {
				evaluationEventResultResp.ChannelMatched = append(evaluationEventResultResp.ChannelMatched, id)
			}
		case int32(domain.PayContainer):
			for id := range v {
				evaluationEventResultResp.PayMatched = append(evaluationEventResultResp.PayMatched, id)
			}
		case int32(domain.ChannelLabelContainer):
			for label := range v {
				evaluationEventResultResp.ChannelLabelMatched = append(evaluationEventResultResp.ChannelLabelMatched, label)
			}
		}
	}

	evaluationEventResultResp.ID = eventResult.ID
	evaluationEventResultResp.FeedbackEventResultResp = &evaluationDTO.FeedbackEventResultDTO{
		RewardType:                eventResult.FeedbackEventResult.RewardType,
		CalculateType:             eventResult.FeedbackEventResult.CalculateType,
		Cost:                      eventResult.FeedbackEventResult.Cost,
		GetReturn:                 eventResult.FeedbackEventResult.GetReturn,
		GetPercentage:             eventResult.FeedbackEventResult.GetPercentage,
		FeedbackEventResultStatus: int32(eventResult.FeedbackEventResult.FeedbackEventResultStatus),
	}

	return evaluationEventResultResp
}

func processMatchedPayload(payLoadEventResult *event.PayloadEventResult, matchedContainerTypeMap map[int32]map[string]bool) {
	if payLoadEventResult == nil {
		return
	}

	if payLoadEventResult.PayloadEventResults != nil {
		for _, p := range payLoadEventResult.PayloadEventResults {
			processMatchedPayload(p, matchedContainerTypeMap)
		}
		return
	}

	if payLoadEventResult.ContainerEventResult == nil {
		return
	}

	processMatchedContainer(payLoadEventResult.ContainerEventResult, matchedContainerTypeMap)

}

func processMatchedContainer(containerEventResult *event.ContainerEventResult, matchedContainerTypeMap map[int32]map[string]bool) {

	if containerEventResult == nil {
		return
	}

	if containerEventResult.ContainerType == int32(domain.InnerContainer) {
		for _, c := range containerEventResult.InnerContainerEventResults {
			processMatchedContainer(c, matchedContainerTypeMap)
		}
		return
	}

	if _, ok := matchedContainerTypeMap[containerEventResult.ContainerType]; !ok {
		matchedContainerTypeMap[containerEventResult.ContainerType] = make(map[string]bool)
	}

	for _, m := range containerEventResult.Matches {
		matchedContainerTypeMap[containerEventResult.ContainerType][m] = true
	}

}
