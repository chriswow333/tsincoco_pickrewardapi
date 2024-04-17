package service

import (
	"context"

	commonM "pickrewardapi/internal/shared/common/model"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	payDTO "pickrewardapi/internal/domain/pay/dto"
	payStore "pickrewardapi/internal/domain/pay/store"
)

type PayService interface {
	GetPayByID(ctx context.Context, ID string) (*payDTO.PayDTO, error)
	GetAllPays(ctx context.Context) ([]*payDTO.PayDTO, error)
}

type payImpl struct {
	dig.In

	payStore payStore.PayStore
}

func NewPay(
	payStore payStore.PayStore,
) PayService {
	return &payImpl{
		payStore: payStore,
	}
}

func (im *payImpl) GetPayByID(ctx context.Context, ID string) (*payDTO.PayDTO, error) {
	return im.payStore.GetPayByID(ctx, ID)
}

func (im *payImpl) GetAllPays(ctx context.Context) ([]*payDTO.PayDTO, error) {
	logPos := "[pay.service][GetAllPays]"

	dtos, err := im.payStore.GetAllPays(ctx, commonM.Active)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("payStore.GetAllPays failed")
		return nil, err
	}

	return dtos, nil
}
