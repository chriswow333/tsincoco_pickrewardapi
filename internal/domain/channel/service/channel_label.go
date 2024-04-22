package service

import (
	"context"

	channelDTO "pickrewardapi/internal/domain/channel/dto"
	channelLabelStore "pickrewardapi/internal/domain/channel/store"

	"go.uber.org/dig"
)

type ChannelLabelService interface {
	GetShowLabels(ctx context.Context) ([]*channelDTO.ChannelLabelDTO, error)
}

type channelLabelImpl struct {
	dig.In

	channelLabelStore channelLabelStore.ChannelLabelStore
}

func NewChannelLabel(
	channelLabelStore channelLabelStore.ChannelLabelStore,
) ChannelLabelService {

	impl := &channelLabelImpl{
		channelLabelStore: channelLabelStore,
	}

	return impl
}

func (im *channelLabelImpl) GetShowLabels(ctx context.Context) ([]*channelDTO.ChannelLabelDTO, error) {
	return im.channelLabelStore.GetChannelLabelsByShow(ctx, 1)
}
