package transfer

import (
	event "pickrewardapi/internal/app/evaluation/domain/event"
	evaluationDTO "pickrewardapi/internal/app/evaluation/dto"
)

func (im *impl) transferContainerEventResult(containerEventResult *event.ContainerEventResult) *evaluationDTO.ContainerEventResultDTO {

	if containerEventResult == nil {
		return &evaluationDTO.ContainerEventResultDTO{}
	}

	containerEventResultDTOs := []*evaluationDTO.ContainerEventResultDTO{}

	for _, c := range containerEventResult.InnerContainerEventResults {
		containerEventResultDTOs = append(containerEventResultDTOs, im.transferContainerEventResult(c))
	}

	containerEventResultDTO := &evaluationDTO.ContainerEventResultDTO{
		ID:                    containerEventResult.ID,
		Pass:                  containerEventResult.Pass,
		Matches:               containerEventResult.Matches,
		MisMatches:            containerEventResult.MisMatches,
		ContainerEventResults: containerEventResultDTOs,
	}

	return containerEventResultDTO
}

func (im *impl) transferFeedbackEventResult(feedbackEventResult *event.FeedbackEventResult) *evaluationDTO.FeedbackEventResultDTO {

	if feedbackEventResult == nil {
		return &evaluationDTO.FeedbackEventResultDTO{}
	}

	feedbackEventResultDTO := &evaluationDTO.FeedbackEventResultDTO{
		RewardType:                feedbackEventResult.RewardType,
		Cost:                      feedbackEventResult.Cost,
		GetReturn:                 feedbackEventResult.GetReturn,
		GetPercentage:             feedbackEventResult.GetPercentage,
		FeedbackEventResultStatus: int32(feedbackEventResult.FeedbackEventResultStatus),
		CalculateType:             feedbackEventResult.CalculateType,
	}

	return feedbackEventResultDTO
}

func (im *impl) transferToPayloadEventResultDTO(payloadEventResult *event.PayloadEventResult) *evaluationDTO.PayloadEventResultDTO {

	if payloadEventResult == nil {
		return &evaluationDTO.PayloadEventResultDTO{}
	}

	payloadEventResultDTOs := []*evaluationDTO.PayloadEventResultDTO{}

	for _, p := range payloadEventResult.PayloadEventResults {
		payloadEventResultDTOs = append(payloadEventResultDTOs, im.transferToPayloadEventResultDTO(p))
	}

	payloadEventResultDTO := &evaluationDTO.PayloadEventResultDTO{
		ID:                   payloadEventResult.ID,
		Pass:                 payloadEventResult.Pass,
		FeedbackEventResult:  im.transferFeedbackEventResult(payloadEventResult.FeedbackEventResult),
		PayloadEventResults:  payloadEventResultDTOs,
		ContainerEventResult: im.transferContainerEventResult(payloadEventResult.ContainerEventResult),
	}

	return payloadEventResultDTO

}
