package service

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	evaluationStore "pickrewardapi/internal/domain/evaluation/store"

	evaluationDTO "pickrewardapi/internal/domain/evaluation/dto"
)

type EvaluationService interface {
	ModifiedEvaluation(ctx context.Context, dto *evaluationDTO.EvaluationDTO) error
	GetEvaluationByOwnerID(ctx context.Context, ID string) (*evaluationDTO.EvaluationDTO, error)
	DeleteEvaluationByID(ctx context.Context, ID string) error
}

type impl struct {
	dig.In

	evaluationStore evaluationStore.EvaluationStore
}

func New(
	evaluationStore evaluationStore.EvaluationStore,
) EvaluationService {
	logPos := "[evaluation.service][New]"

	log.WithFields(log.Fields{
		"pos": logPos,
	}).Info("init")

	impl := &impl{
		evaluationStore: evaluationStore,
	}
	return impl
}

var timeNow = time.Now

func (im *impl) DeleteEvaluationByID(ctx context.Context, id string) error {
	return im.evaluationStore.DeleteEvaluationByID(ctx, id)
}

func (im *impl) GetEvaluationByOwnerID(ctx context.Context, ownerID string) (*evaluationDTO.EvaluationDTO, error) {
	logPos := "[evaluation.service][GetEvaluationByOwnerID]"

	dto, err := im.evaluationStore.GetEvaluationByOwnerID(ctx, ownerID)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("evaluationStore.GetEvaluationByOwnerID failed")
		return nil, err
	}
	return dto, nil
}

func (im *impl) ModifiedEvaluation(ctx context.Context, evaluationDTO *evaluationDTO.EvaluationDTO) error {
	logPos := "[evaluation.service][ModifiedEvaluation]"

	evaluationDTO.CreateDate = timeNow().Unix()
	evaluationDTO.UpdateDate = timeNow().Unix()

	if err := im.evaluationStore.ModifiedEvaluation(ctx, evaluationDTO); err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("evaluationStore.ModifiedEvaluation failed")
		return err
	}

	return nil
}
