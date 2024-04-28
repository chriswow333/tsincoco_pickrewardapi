package service

import (
	"context"

	cardDTO "pickrewardapi/internal/domain/card/dto"
	cardStore "pickrewardapi/internal/domain/card/store"

	"go.uber.org/dig"
)

type FeedbackTypeService interface {
	GetFeedbackTypeByID(ctx context.Context, ID string) (*cardDTO.FeedbackTypeDTO, error)
}

type feedbackImpl struct {
	dig.In

	feedbackTypeStore cardStore.FeedbackTypeStore
}

func NewFeedbackType(
	feedbackTypeStore cardStore.FeedbackTypeStore,

) FeedbackTypeService {

	return &feedbackImpl{
		feedbackTypeStore: feedbackTypeStore,
	}

}

func (im *feedbackImpl) GetFeedbackTypeByID(ctx context.Context, ID string) (*cardDTO.FeedbackTypeDTO, error) {
	return im.feedbackTypeStore.GetFeedbackTypeByID(ctx, ID)
}
