package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	cardService "pickrewardapi/internal/domain/card/service"
	channelService "pickrewardapi/internal/domain/channel/service"
	handler "pickrewardapi/internal/domain/evaluation/service/evaluate/handler"
	evaluationStore "pickrewardapi/internal/domain/evaluation/store"

	domain "pickrewardapi/internal/domain/evaluation/domain"
	evaluationDTO "pickrewardapi/internal/domain/evaluation/dto"

	commonM "pickrewardapi/internal/shared/common/model"
)

type EvaluateService interface {
	EvaluateRespByOwnerID(ctx context.Context, ID string, event *commonM.Event) (*evaluationDTO.EvaluationEventResultRespDTO, error)
}

type impl struct {
	dig.In

	evaluations []*domain.Evaluation

	evaluationStore     evaluationStore.EvaluationStore
	channelService      channelService.ChannelService
	feedbackTypeService cardService.FeedbackTypeService

	transfer handler.Transfer
}

func New(
	evaluationStore evaluationStore.EvaluationStore,
	channelService channelService.ChannelService,
	feedbackTypeService cardService.FeedbackTypeService,
) EvaluateService {
	logPos := "[evaluation.evaluate.service][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Info("init")

	impl := &impl{
		evaluationStore:     evaluationStore,
		channelService:      channelService,
		feedbackTypeService: feedbackTypeService,
		transfer: *handler.NewTransfer(
			channelService,
		),
	}

	impl.init()

	return impl
}

var timeNow = time.Now

func (im *impl) init() {
	logPos := "[evaluation.service][init]"

	ctx := context.Background()

	evaluationDTOs, err := im.evaluationStore.GetAllEvaluations(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("evaluationStore.GetAllEvaluations failed", err)
		return
	}

	evaluations := []*domain.Evaluation{}

	for _, eva := range evaluationDTOs {
		evaluation, err := im.transfer.TransferToEvaluation(ctx, eva)
		if err != nil {
			log.WithFields(log.Fields{
				"pos": logPos,
			}).Error("handler.TransferToEvaluation failed", err)
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
	logPos := "[evaluation.service][EvaluateRespByOwnerID]"

	if event.CardEvent == nil {
		log.WithFields(log.Fields{
			"pos":           logPos,
			"evaluation.ID": ID,
			"event.date":    event.Date,
		}).Error("cardEvent is nil")
		return nil, errors.New("cardEvent is nil")
	}

	if event.ChannelEvent == nil {

		log.WithFields(log.Fields{
			"pos":           logPos,
			"evaluation.ID": ID,
			"event.date":    event.Date,
		}).Error("channelEvent is nil")
		return nil, errors.New("channelEvent is nil")

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

			if eva.FeedbackType != event.CardEvent.FeedbackType {
				log.WithFields(log.Fields{
					"pos":                          logPos,
					"evaluation.ID":                eva.ID,
					"event.CardEvent.FeedbackType": event.CardEvent.FeedbackType,
					"evaluation.FeedbackType":      eva.FeedbackType,
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

			eventResultResp := handler.CalculateEvaluationResult(eventResult)

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
