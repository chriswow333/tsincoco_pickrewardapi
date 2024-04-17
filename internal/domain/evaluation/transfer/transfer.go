package transfer

import (
	"context"

	domain "pickrewardapi/internal/app/evaluation/domain"
	event "pickrewardapi/internal/app/evaluation/domain/event"
	evaluationDTO "pickrewardapi/internal/app/evaluation/dto"
)

type EvaluationAppTransfer interface {
	TransferToEvaluation(ctx context.Context, dto *evaluationDTO.EvaluationDTO) (*domain.Evaluation, error)
	// TransferToEvaluationRespDTO(ctx context.Context, evaluation *domain.Evaluation) (*evaluationDTO.EvaluationRespDTO, error)
	TransferToEvaluationEventResultDTO(eventResult *event.EvaluationEventResult) *evaluationDTO.EvaluationEventResultDTO
}
