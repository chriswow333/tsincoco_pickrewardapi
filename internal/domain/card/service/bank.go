package service

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	commonM "pickrewardapi/internal/shared/common/model"

	cardDTO "pickrewardapi/internal/domain/card/dto"
	bankStore "pickrewardapi/internal/domain/card/store"
)

type BankService interface {
	GetBankByID(ctx context.Context, ID string) (*cardDTO.BankDTO, error)
	GetAllBanks(ctx context.Context) ([]*cardDTO.BankDTO, error)
}

type bankImpl struct {
	dig.In

	bankStore bankStore.BankStore
}

func NewBank(
	bankStore bankStore.BankStore,
) BankService {
	return &bankImpl{
		bankStore: bankStore,
	}
}

func (im *bankImpl) GetBankByID(ctx context.Context, ID string) (*cardDTO.BankDTO, error) {
	return im.bankStore.GetBankByID(ctx, ID)
}

func (im *bankImpl) GetAllBanks(ctx context.Context) ([]*cardDTO.BankDTO, error) {
	logPos := "[card.app.service][GetAllBanks]"

	dtos, err := im.bankStore.GetAllBanks(ctx, commonM.Active)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("cardStore.GetAllBanks failed")
		return nil, err
	}

	return dtos, nil
}
