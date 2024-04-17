package transfer

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	domain "pickrewardapi/internal/app/evaluation/domain"
	evaluationDTO "pickrewardapi/internal/app/evaluation/dto"
)

func (im *impl) transferToPayload(ctx context.Context, payloadDTO *evaluationDTO.PayloadDTO) (*domain.Payload, error) {
	logPos := "[evaluation.app.transfer][transferToPayload]"

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
			payload, err := im.transferToPayload(ctx, p)
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

	container, err := im.transferToContainer(ctx, payloadDTO.Container)
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
		Feedback:        im.transferToFeedback(payloadDTO.Feedback),
		Container:       container,
	}

	return payload, nil
}

func (im *impl) transferToFeedback(dto *evaluationDTO.FeedbackDTO) *domain.Feedback {

	if dto == nil {
		return nil
	}
	return &domain.Feedback{
		RewardType:    dto.RewardType,
		CalculateType: domain.CalculateType(dto.CalculateType),
		MinCost:       dto.MinCost,
		Fixed:         dto.Fixed,
		Percentage:    dto.Percentage,
		ReturnMax:     dto.ReturnMax,
	}
}

func (im *impl) transferToContainer(ctx context.Context, containerDTO *evaluationDTO.ContainerDTO) (*domain.Container, error) {
	logPos := "[evaluation.app.transfer][transferToContainer]"
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
			val, err := im.transferToContainer(ctx, c)
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

	channelEvaluationDTOs, err := im.transferChannelEvaluationDTOs(ctx, containerDTO)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":          logPos,
			"container.ID": containerDTO.ID,
		}).Error("transferChannelDTOs failed", err)
		return nil, err
	}

	constraints, err := im.transferToConstraint(containerDTO)
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

func (im *impl) transferChannelEvaluationDTOs(ctx context.Context, containerDTO *evaluationDTO.ContainerDTO) ([]*evaluationDTO.ChannelEvaluationDTO, error) {
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
			channel, err := im.channelAppService.GetByChannelID(ctx, ch)
			if err != nil {
				log.WithFields(log.Fields{
					"pos":        logPos,
					"channel.ID": ch,
				}).Error("channelAppService.GetByChannelID failed", err)
				return nil, err
			}
			channelLabels := []int32{}
			for _, cl := range channel.ChannelLabels {
				channelLabels = append(channelLabels, cl.Label)
			}

			channelEvaluationDTOs = append(channelEvaluationDTOs, &evaluationDTO.ChannelEvaluationDTO{
				ID:            channel.ID,
				ChannelLabels: channelLabels,
			})
		}
		return channelEvaluationDTOs, nil
	}

	return nil, nil
}

func (im *impl) transferToConstraint(containerDTO *evaluationDTO.ContainerDTO) ([]*domain.Constraint, error) {
	logPos := "[evaluation.app.transfer][transferToConstraint]"

	if containerDTO.ContainerType != int32(domain.ConstraintContainer) {
		return nil, nil
	}

	constraints := []*domain.Constraint{}
	if containerDTO.Constraints == nil {
		log.WithFields(log.Fields{
			"pos":          logPos,
			"container.ID": containerDTO.ID,
		}).Error("constraints is nil")
		return nil, errors.New("constraints is nil")
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
