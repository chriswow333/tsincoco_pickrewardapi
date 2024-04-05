package service

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	commonM "pickrewardapi/internal/shared/common/model"

	bankDTO "pickrewardapi/internal/domain/bank/dto"
	bankStore "pickrewardapi/internal/domain/bank/store"
)

type BankService interface {
	GetBankByID(ctx context.Context, ID string) (*bankDTO.BankDTO, error)
	GetAllBanks(ctx context.Context) ([]*bankDTO.BankDTO, error)
}

type impl struct {
	dig.In

	bankStore bankStore.BankStore
}

func New(
	bankStore bankStore.BankStore,
) BankService {
	return &impl{
		bankStore: bankStore,
	}
}

func (im *impl) GetBankByID(ctx context.Context, ID string) (*bankDTO.BankDTO, error) {
	return im.bankStore.GetBankByID(ctx, ID)
}

func (im *impl) GetAllBanks(ctx context.Context) ([]*bankDTO.BankDTO, error) {
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
