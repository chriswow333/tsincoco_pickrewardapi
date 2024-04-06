package service

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	cardRewardStore "pickrewardapi/internal/domain/card/store"

	cardDTO "pickrewardapi/internal/domain/card/dto"
	commonM "pickrewardapi/internal/shared/common/model"
)

type CardRewardService interface {
	GetCardRewardsByCardID(ctx context.Context, cardID string) ([]*cardDTO.CardRewardDTO, error)
}

type cardRewardImpl struct {
	dig.In

	cardRewardStore cardRewardStore.CardRewardStore
}

func NewCardReward(
	cardRewardStore cardRewardStore.CardRewardStore,
) CardRewardService {

	impl := &cardRewardImpl{
		cardRewardStore: cardRewardStore,
	}

	return impl
}

func (im *cardRewardImpl) GetCardRewardsByCardID(ctx context.Context, cardID string) ([]*cardDTO.CardRewardDTO, error) {
	logPos := "[card.reward.service][GetCardRewardsByCardID]"

	cardRewardDTOs, err := im.cardRewardStore.GetCardRewardsByCardID(ctx, cardID, commonM.Active)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("cardRewardStore.GetCardRewardsByCardID failed", err)
		return nil, err
	}
	return cardRewardDTOs, nil
}
