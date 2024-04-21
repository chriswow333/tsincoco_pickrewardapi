package handler

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	domain "pickrewardapi/internal/domain/evaluation/domain"
	evaluationDTO "pickrewardapi/internal/domain/evaluation/dto"

	channelService "pickrewardapi/internal/domain/channel/service"
)

type Transfer struct {
	channelService channelService.ChannelService
}

func NewTransfer(
	channelService channelService.ChannelService,
) *Transfer {
	return &Transfer{
		channelService: channelService,
	}
}

func (t *Transfer) TransferToEvaluation(ctx context.Context, evaluationDTO *evaluationDTO.EvaluationDTO) (*domain.Evaluation, error) {
	logPos := "[evaluation.transfer][TransferToEvaluation]"

	log.WithFields(log.Fields{
		"pos":           logPos,
		"evaluation.ID": evaluationDTO.ID,
	}).Info("loading evaluation")

	payload, err := t.transferToPayload(ctx, evaluationDTO.Payload)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":           logPos,
			"evaluation.ID": evaluationDTO.ID,
		}).Error("transferToPayload failed", err)
		return nil, err
	}

	evaluation := &domain.Evaluation{
		ID:         evaluationDTO.ID,
		CreateDate: evaluationDTO.CreateDate,
		UpdateDate: evaluationDTO.UpdateDate,
		StartDate:  evaluationDTO.StartDate,
		EndDate:    evaluationDTO.EndDate,
		FeedbackID: evaluationDTO.FeedbackID,
		OwnerID:    evaluationDTO.OwnerID,
		Payload:    payload,
	}

	return evaluation, nil

}

func (t *Transfer) transferToPayload(ctx context.Context, payloadDTO *evaluationDTO.PayloadDTO) (*domain.Payload, error) {
	logPos := "[evaluation.transfer][transferToPayload]"

	if payloadDTO == nil {
		return nil, nil
	}

	payloads := []*domain.Payload{}

	if payloadDTO.PayloadType == int32(domain.SelfPayloadType) {
		if payloadDTO.Payloads == nil {
			log.WithFields(log.Fields{
				"pos":        logPos,
				"payload.ID": payloadDTO.ID,
			}).Error("payloads is nil")
			return nil, errors.New("payloads is nil")
		}

		for _, p := range payloadDTO.Payloads {
			payload, err := t.transferToPayload(ctx, p)
			if err != nil {
				log.WithFields(log.Fields{
					"pos":        logPos,
					"payload.ID": payloadDTO.ID,
				}).Error("transferToPayload failed", err)
				return nil, err
			}

			if payload != nil {
				payloads = append(payloads, payload)
			}
		}

		return &domain.Payload{
			ID:              payloadDTO.ID,
			PayloadOperator: domain.PayloadOperator(payloadDTO.PayloadOperator),
			PayloadType:     domain.PayloadType(payloadDTO.PayloadType),
			Payloads:        payloads,
		}, nil
	}

	container, err := t.transferToContainer(ctx, payloadDTO.Container)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":        logPos,
			"payload.ID": payloadDTO.ID,
		}).Error("transferToContainer failed,", err)
		return nil, err
	}

	payload := &domain.Payload{
		ID:              payloadDTO.ID,
		PayloadOperator: domain.PayloadOperator(payloadDTO.PayloadOperator),
		PayloadType:     domain.PayloadType(payloadDTO.PayloadType),
		Feedback:        t.transferToFeedback(payloadDTO.Feedback),
		Container:       container,
	}

	return payload, nil
}

func (t *Transfer) transferToFeedback(dto *evaluationDTO.FeedbackDTO) *domain.Feedback {

	if dto == nil {
		return nil
	}
	return &domain.Feedback{
		FeedbackID:    dto.FeedbackID,
		CalculateType: domain.CalculateType(dto.CalculateType),
		MinCost:       dto.MinCost,
		Fixed:         dto.Fixed,
		Percentage:    dto.Percentage,
		ReturnMax:     dto.ReturnMax,
	}
}

func (t *Transfer) transferToContainer(ctx context.Context, containerDTO *evaluationDTO.ContainerDTO) (*domain.Container, error) {
	logPos := "[evaluation.transfer][transferToContainer]"
	containerID := containerDTO.ID

	containers := []*domain.Container{}

	if containerDTO.ContainerType == int32(domain.InnerContainer) {

		if containerDTO.Containers == nil {
			log.WithFields(log.Fields{
				"pos":          logPos,
				"container.ID": containerID,
			}).Error("containers are nil")
			return nil, errors.New("containers are nil")
		}

		for _, c := range containerDTO.Containers {
			val, err := t.transferToContainer(ctx, c)
			if err != nil {
				log.WithFields(log.Fields{
					"pos":          logPos,
					"container.ID": containerID,
				}).Error("transferToContainer failed,", err)
				return nil, err
			}
			containers = append(containers, val)
		}

		return &domain.Container{
			ID:                containerID,
			ContainerOperator: domain.ContainerOperator(containerDTO.ContainerOperator),
			ContainerType:     domain.ContainerType(containerDTO.ContainerType),
			InnerContainers:   containers,
		}, nil
	}

	var err error

	channelEvaluationDTOs, err := t.transferChannelEvaluationDTOs(ctx, containerDTO)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":          logPos,
			"container.ID": containerDTO.ID,
		}).Error("transferChannelDTOs failed", err)
		return nil, err
	}

	constraints, err := transferToConstraint(containerDTO)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":          logPos,
			"container.ID": containerDTO.ID,
		}).Error("transferToConstraint failed", err)
		return nil, err
	}

	container := &domain.Container{
		ID:                      containerID,
		ContainerOperator:       domain.ContainerOperator(containerDTO.ContainerOperator),
		ContainerType:           domain.ContainerType(containerDTO.ContainerType),
		InnerContainers:         containers,
		Constraints:             constraints,
		CardRewardTaskLabels:    containerDTO.CardRewardTaskLabels,
		ChannelEvaluations:      channelEvaluationDTOs,
		PayEvaluations:          containerDTO.PayIDs,
		ChannelLabelEvaluations: containerDTO.ChannelLabels,
	}
	return container, nil
}

func (t *Transfer) transferChannelEvaluationDTOs(ctx context.Context, containerDTO *evaluationDTO.ContainerDTO) ([]*evaluationDTO.ChannelEvaluationDTO, error) {
	logPos := "[evaluation.app.transfer][transferChannel]"

	if containerDTO.ContainerType == int32(domain.ChannelContainer) {
		channelEvaluationDTOs := []*evaluationDTO.ChannelEvaluationDTO{}

		if containerDTO.ChannelIDs == nil {
			log.WithFields(log.Fields{
				"pos":          logPos,
				"container.ID": containerDTO.ID,
			}).Error("channelIDs is nil")
			return nil, errors.New("channelIDs is nil")
		}

		for _, ch := range containerDTO.ChannelIDs {
			channel, err := t.channelService.GetByChannelID(ctx, ch)
			if err != nil {
				log.WithFields(log.Fields{
					"pos":        logPos,
					"channel.ID": ch,
				}).Error("channelAppService.GetByChannelID failed", err)
				return nil, err
			}

			channelEvaluationDTOs = append(channelEvaluationDTOs, &evaluationDTO.ChannelEvaluationDTO{
				ID:            channel.ID,
				ChannelLabels: channel.ChannelLabels,
			})
		}
		return channelEvaluationDTOs, nil
	}

	return nil, nil
}

func transferToConstraint(containerDTO *evaluationDTO.ContainerDTO) ([]*domain.Constraint, error) {
	logPos := "[evaluation.app.transfer][transferToConstraint]"

	if containerDTO.ContainerType != int32(domain.ConstraintContainer) {
		return nil, nil
	}

	constraints := []*domain.Constraint{}
	if containerDTO.Constraints == nil {
		log.WithFields(log.Fields{
			"pos":          logPos,
			"container.ID": containerDTO.ID,
		}).Error("constraints are nil")
		return nil, errors.New("constraints are nil")
	}

	for _, c := range containerDTO.Constraints {
		constraints = append(constraints, &domain.Constraint{
			ConstraintType: domain.ConstraintType(c.ConstraintType),
			ConstraintName: c.ConstraintName,
			WeekDays:       c.WeekDays,
		})
	}

	return constraints, nil

}
