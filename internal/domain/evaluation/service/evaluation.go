package service

import (
	"context"
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	channelAppService "pickrewardapi/internal/app/channel/service/channel/app"
	evaluationStore "pickrewardapi/internal/app/evaluation/repo/store"
	evaluationAppTransfer "pickrewardapi/internal/app/evaluation/transfer/app"
	rewardAppService "pickrewardapi/internal/app/reward/service/app"

	domain "pickrewardapi/internal/app/evaluation/domain"
	evaluationDTO "pickrewardapi/internal/app/evaluation/dto"
	commonM "pickrewardapi/internal/shared/common/model"
)

type EvaluationAppService interface {
	// GetEvaluationRespByOwnerID(ctx context.Context, ID string) (*evaluationDTO.EvaluationRespDTO, error)
	// EvaluateResp(ctx context.Context, ID string, event *commonM.Event) (*evaluationDTO.EvaluationEventResultRespDTO, error)
	EvaluateRespByOwnerID(ctx context.Context, ID string, event *commonM.Event) (*evaluationDTO.EvaluationEventResultRespDTO, error)
}

type impl struct {
	dig.In

	evaluations []*domain.Evaluation

	evaluationStore       evaluationStore.EvaluationStore
	channelAppService     channelAppService.ChannelAppService
	rewardAppService      rewardAppService.RewardAppService
	evaluationAppTransfer evaluationAppTransfer.EvaluationAppTransfer
}

func New(
	evaluationStore evaluationStore.EvaluationStore,
	channelAppService channelAppService.ChannelAppService,
	rewardAppService rewardAppService.RewardAppService,
	evaluationAppTransfer evaluationAppTransfer.EvaluationAppTransfer,
) EvaluationAppService {
	logPos := "[evaluation.app.service][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Info("init")

	impl := &impl{
		evaluationStore:       evaluationStore,
		channelAppService:     channelAppService,
		rewardAppService:      rewardAppService,
		evaluationAppTransfer: evaluationAppTransfer,
	}

	impl.init()

	return impl
}

func (im *impl) init() {
	logPos := "[evaluation.app.service][init]"

	ctx := context.Background()

	evaluationDTOs, err := im.evaluationStore.GetAllEvaluations(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("evaluationStore.GetAllRewards failed", err)
		return
	}

	evaluations := []*domain.Evaluation{}
	for _, eva := range evaluationDTOs {
		evaluation, err := im.evaluationAppTransfer.TransferToEvaluation(ctx, eva)
		if err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("evaluationAppTransfer.TransferToEvaluation failed", err)
			return
		}
		evaluations = append(evaluations, evaluation)
	}

	im.evaluations = evaluations

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Info("init")
}

func (im *impl) EvaluateRespByOwnerID(ctx context.Context, ID string, event *commonM.Event) (*evaluationDTO.EvaluationEventResultRespDTO, error) {
	logPos := "[evaluation.app.service][EvaluateRespByOwnerID]"

	if event.CardEvent == nil {
		log.WithFields(log.Fields{
			"pos":           logPos,
			"evaluation.ID": ID,
			"event.date":    event.Date,
		}).Error("CardEvent is nil")
		return nil, errors.New("CardEvent is nil")
	}

	if event.ChannelEvent == nil {
		log.WithFields(log.Fields{
			"pos":           logPos,
			"evaluation.ID": ID,
			"event.date":    event.Date,
		}).Error("ChannelEvent is nil")
		return nil, errors.New("ChannelEvent is nil")
	}

	for _, eva := range im.evaluations {
		if eva.OwnerID == ID {
			if eva.StartDate > event.Date || eva.EndDate < event.Date {
				log.WithFields(log.Fields{
					"pos":           logPos,
					"evaluation.ID": ID,
					"event.date":    event.Date,
					"start.date":    eva.StartDate,
					"end.date":      eva.EndDate,
				}).Error("evaluation is not in the validate range")
				return nil, errors.New("evaluation is not in the validate range")
			}

			if eva.Reward.RewardType != event.CardEvent.RewardType {
				log.WithFields(log.Fields{
					"pos":                          logPos,
					"evaluation.ID":                eva.ID,
					"event.CardEvent.RewardType":   event.CardEvent.RewardType,
					"evaluation.Reward.RewardType": eva.Reward.RewardType,
				}).Error("reward type is not correct")
				return nil, errors.New("reward type is not correct")
			}

			eventResult, err := eva.Judge(event)
			if err != nil {
				log.WithFields(log.Fields{
					"pos":           logPos,
					"evaluation.ID": eva.ID,
				}).Error("im.evaluate failed")
				return nil, err
			}

			eventResultResp := calculateEvaluationResult(eventResult)

			eventResultJsonLog, _ := json.Marshal(eventResult)
			eventResultLog, _ := json.Marshal(eventResultResp)
			log.WithFields(log.Fields{
				"pos":            logPos,
				"evaluation.ID":  eva.ID,
				"event.result":   string(eventResultJsonLog),
				"event.feedback": string(eventResultLog),
			}).Info("evaluation done")

			return eventResultResp, nil

		}
	}

	return nil, errors.New("Cannot find evaluation owner id")
}

// func (im *impl) evaluate(ctx context.Context, ID string, e *commonM.Event) (*evaluationDTO.EvaluationEventResultDTO, error) {
// 	logPos := "[evaluation.app.service][evaluate]"

// 	for _, eva := range im.evaluations {
// 		if eva.ID == ID {
// 			eventResult, err := eva.Judge(e)
// 			if err != nil {
// 				log.WithFields(log.Fields{
// 					"pos":            logPos,
// 					"evaluations.ID": ID,
// 				}).Error("eva.Judge failed", err)
// 				return nil, err
// 			}
// 			return im.evaluationAppTransfer.TransferToEvaluationEventResultDTO(eventResult), nil
// 		}
// 	}

// 	log.WithFields(log.Fields{
// 		"pos":            logPos,
// 		"evaluations.ID": ID,
// 	}).Error("Cannot find evaluation id")

// 	return nil, errors.New("Cannot find evaluation id")
// }

// func (im *impl) GetEvaluationByID(ctx context.Context, ID string) (*evaluationDTO.EvaluationDTO, error) {
// 	logPos := "[evaluation.service][GetEvaluationByID]"

// 	dto, err := im.evaluationStore.GetEvaluationByID(ctx, ID)
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"pos": logPos,
// 		}).Error("evaluationStore.GetEvaluationDTOByID failed")
// 		return nil, err
// 	}

// 	if dto == nil || dto.Reward == nil {
// 		log.WithFields(log.Fields{
// 			"pos": logPos,
// 			"ID":  ID,
// 		}).Error("evaluationStore.GetEvaluationByID is nil")
// 		return nil, errors.New("evaluationStore.GetEvaluationByID is nil")
// 	}

// 	reward, err := im.rewardAppService.GetRewardByID(ctx, dto.Reward.ID)
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"pos":      logPos,
// 			"rewardID": dto.Reward.ID,
// 		}).Error("rewardService.GetRewardByID failed")
// 		return nil, err
// 	}
// 	dto.Reward = reward

// 	return dto, nil
// }

// func (im *impl) GetEvaluationByOwnerID(ctx context.Context, ownerID string) (*evaluationDTO.EvaluationDTO, error) {
// 	logPos := "[evaluation.service][GetEvaluationByOwnerID]"

// 	reward, err := im.rewardAppService.GetRewardByID(ctx, ownerID)
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"pos":      logPos,
// 			"rewardID": ownerID,
// 		}).Error("rewardService.GetRewardByID failed")
// 		return nil, err
// 	}

// 	dto, err := im.evaluationStore.GetEvaluationByOwnerID(ctx, ownerID)
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"pos": logPos,
// 		}).Error("evaluationStore.GetEvaluationDTOByID failed")
// 		return nil, err
// 	}

// 	dto.Reward = reward

// 	return dto, nil

// }
