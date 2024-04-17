package transfer

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	channelAppService "pickrewardapi/internal/app/channel/service/channel/app"
	payAppService "pickrewardapi/internal/app/pay/service/app"
	rewardAppService "pickrewardapi/internal/app/reward/service/app"

	domain "pickrewardapi/internal/app/evaluation/domain"
	event "pickrewardapi/internal/app/evaluation/domain/event"
	evaluationDTO "pickrewardapi/internal/app/evaluation/dto"
)

type impl struct {
	dig.In

	channelAppService channelAppService.ChannelAppService
	payAppService     payAppService.PayAppService
	rewardAppService  rewardAppService.RewardAppService
}

func New(
	channelAppService channelAppService.ChannelAppService,
	payAppService payAppService.PayAppService,
	rewardAppService rewardAppService.RewardAppService,
) EvaluationAppTransfer {
	logPos := "[evaluation.app.transfer][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Info("init")

	impl := &impl{
		channelAppService: channelAppService,
		payAppService:     payAppService,
		rewardAppService:  rewardAppService,
	}

	return impl
}

func (im *impl) TransferToEvaluation(ctx context.Context, evaluationDTO *evaluationDTO.EvaluationDTO) (*domain.Evaluation, error) {
	logPos := "[evaluation.app.transfer][TransferToEvaluation]"

	rewardDTO, err := im.rewardAppService.GetRewardByID(ctx, evaluationDTO.RewardID)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":                  logPos,
			"evaluation.ID":        evaluationDTO.ID,
			"evaluation.reward.ID": evaluationDTO.RewardID,
		}).Error("rewardAppService.GetRewardByID failed", err)
		return nil, err
	}

	log.WithFields(log.Fields{
		"pos":           logPos,
		"evaluation.ID": evaluationDTO.ID,
	}).Info("Loading evaluation")

	payload, err := im.transferToPayload(ctx, evaluationDTO.Payload)
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
		Reward:     rewardDTO,
		Owner:      domain.Owner(evaluationDTO.Owner),
		OwnerID:    evaluationDTO.OwnerID,
		Payload:    payload,
	}

	return evaluation, nil

}

// func (im *impl) TransferToEvaluationRespDTO(ctx context.Context, evaluation *domain.Evaluation) (*evaluationDTO.EvaluationRespDTO, error) {
// 	return nil, nil
// }

// func (im *impl) TransferToEvaluationRespDTO(ctx context.Context, evaluation *domain.Evaluation) (*evaluationDTO.EvaluationRespDTO, error) {
// 	logPos := "[evaluation.service][transferToEvaluationRespDTO]"

// 	evaluationRespDTO := &evaluationDTO.EvaluationRespDTO{
// 		ID: evaluation.ID,
// 	}

// 	if err := im.transferToPayloadRespDTO(ctx, evaluation.Payload, evaluationRespDTO); err != nil {
// 		log.WithFields(log.Fields{
// 			"pos":           logPos,
// 			"evaluation.ID": evaluation.ID,
// 		}).Error("transferToPayloadRespDTO failed", err)
// 		return nil, err
// 	}

// 	if err := im.transferChannelCategoryTypes(ctx, evaluationRespDTO); err != nil {
// 		log.WithFields(log.Fields{
// 			"pos":           logPos,
// 			"evaluation.ID": evaluation.ID,
// 		}).Error("transferChannelCategoryTypes failed", err)
// 	}

// 	return evaluationRespDTO, nil
// }

func (im *impl) TransferToEvaluationEventResultDTO(eventResult *event.EvaluationEventResult) *evaluationDTO.EvaluationEventResultDTO {

	evaluationEventResultDTO := &evaluationDTO.EvaluationEventResultDTO{
		ID:                    eventResult.ID,
		EvaluationEventStatus: int32(eventResult.FeedbackEventResult.FeedbackEventResultStatus),
		PayloadEventResult:    im.transferToPayloadEventResultDTO(eventResult.PayloadEventResult),
	}

	return evaluationEventResultDTO
}
