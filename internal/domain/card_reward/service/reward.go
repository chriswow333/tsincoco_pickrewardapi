package service

import (
	"context"
	"time"

	cardRewardDTO "pickrewardapi/internal/domain/card_reward/dto"

	cardRewardStore "pickrewardapi/internal/domain/card_reward/store"

	"go.uber.org/dig"
)

type RewardAppService interface {
	GetRewardByID(ctx context.Context, ID string) (*cardRewardDTO.RewardDTO, error)
}

type impl struct {
	dig.In

	rewardStore cardRewardStore.RewardStore
}

var (
	timeNow = time.Now
)

func New(
	rewardStore cardRewardStore.RewardStore,
) RewardAppService {

	impl := &impl{
		rewardStore: rewardStore,
	}

	return impl
}

func (im *impl) GetRewardByID(ctx context.Context, ID string) (*cardRewardDTO.RewardDTO, error) {
	return im.rewardStore.GetRewardByID(ctx, ID)
}
