package service

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	cardStore "pickrewardapi/internal/domain/card/store"

	cardDTO "pickrewardapi/internal/domain/card/dto"
	commonM "pickrewardapi/internal/shared/common/model"
)

type CardService interface {
	GetCardsByBankID(ctx context.Context, bankID string) ([]*cardDTO.CardDTO, error)
	GetLatestCards(ctx context.Context) ([]*cardDTO.CardDTO, error)
	GetCardByID(ctx context.Context, cardID string) (*cardDTO.CardDTO, error)

	SearchCard(ctx context.Context, keyword string) ([]*cardDTO.CardDTO, error)
}

var (
	timeNow = time.Now
)

type cardImpl struct {
	dig.In

	cardStore cardStore.CardStore
}

func NewCard(
	cardStore cardStore.CardStore,
) CardService {

	impl := &cardImpl{
		cardStore: cardStore,
	}

	return impl
}

func (im *cardImpl) GetCardsByBankID(ctx context.Context, bankID string) ([]*cardDTO.CardDTO, error) {
	logPos := "[card.service][GetCardsByBankID]"

	dtos, err := im.cardStore.GetCardsByBankID(ctx, bankID, commonM.Active)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("cardStore.GetAllBanks failed")
		return nil, err
	}

	return dtos, nil
}

func (im *cardImpl) GetCardByID(ctx context.Context, cardID string) (*cardDTO.CardDTO, error) {
	logPos := "[card.service][GetCardByID]"

	card, err := im.cardStore.GetByCardID(ctx, cardID)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("cardStore.GetByCardID failed")
	}

	if card == nil {
		log.WithFields(log.Fields{
			"pos":    logPos,
			"cardID": cardID,
		}).Error("Cannot find cardID")
		return nil, errors.New("Cannot find cardID")
	}

	return card, nil
}

func (im *cardImpl) GetLatestCards(ctx context.Context) ([]*cardDTO.CardDTO, error) {
	logPos := "[card.service][GetLatestCards]"

	cards, err := im.cardStore.GetLatestCards(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("cardStore.GetLatestCards failed")
		return nil, err
	}

	return cards, nil
}

func (im *cardImpl) SearchCard(ctx context.Context, keyword string) ([]*cardDTO.CardDTO, error) {
	logPos := "[card.service][SearchCard]"

	cards, err := im.cardStore.SearchCard(ctx, keyword, commonM.Active)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("cardStore.SearchCard failed")
		return nil, err
	}

	return cards, nil
}
