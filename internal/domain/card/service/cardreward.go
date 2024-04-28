package service

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	cardStore "pickrewardapi/internal/domain/card/store"

	cardDTO "pickrewardapi/internal/domain/card/dto"
	commonM "pickrewardapi/internal/shared/common/model"
)

type CardRewardService interface {
	GetCardRewardsByCardID(ctx context.Context, cardID string) ([]*cardDTO.CardRewardDTO, error)
}

type cardRewardImpl struct {
	dig.In

	cardRewardStore   cardStore.CardRewardStore
	feedbackTypeStore cardStore.FeedbackTypeStore
}

func NewCardReward(
	cardRewardStore cardStore.CardRewardStore,
	feedbackTypeStore cardStore.FeedbackTypeStore,
) CardRewardService {

	impl := &cardRewardImpl{
		cardRewardStore:   cardRewardStore,
		feedbackTypeStore: feedbackTypeStore,
	}

	return impl
}

func (im *cardRewardImpl) GetCardRewardsByCardID(ctx context.Context, cardID string) ([]*cardDTO.CardRewardDTO, error) {
	logPos := "[card.reward.service][GetCardRewardsByCardID]"

	cardRewardDTOs, err := im.cardRewardStore.GetCardRewardsByCardID(ctx, cardID, commonM.Active)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"card.id": cardID,
		}).Error("cardRewardStore.GetCardRewardsByCardID failed", err)
		return nil, err
	}

	for _, c := range cardRewardDTOs {
		if c.FeedbackType != nil {

			feedbackTypeDTO, err := im.feedbackTypeStore.GetFeedbackTypeByID(ctx, c.FeedbackType.ID)
			if err != nil {
				log.WithFields(log.Fields{
					"pos":             logPos,
					"feedbacktype.id": c.FeedbackType.ID,
					"card.reward.id":  c.ID,
				}).Error("feedbackTypeStore.GetFeedbackTypeByID failed", err)
				return nil, err
			}
			c.FeedbackType = feedbackTypeDTO
		}
	}
	return cardRewardDTOs, nil
}
