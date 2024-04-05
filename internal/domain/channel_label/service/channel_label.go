package service

import (
	"context"
	"errors"

	channelDTO "pickrewardapi/internal/domain/channel_label/dto"
	channelLabelStore "pickrewardapi/internal/domain/channel_label/store"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type ChannelLabelAppService interface {
	GetShowChannelLabels(ctx context.Context) ([]*channelDTO.ChannelLabelDTO, error)
}

type impl struct {
	dig.In

	channelLabelStore channelLabelStore.ChannelLabelStore
}

func New(
	channelLabelStore channelLabelStore.ChannelLabelStore,
) ChannelLabelAppService {

	impl := &impl{
		channelLabelStore: channelLabelStore,
	}

	return impl
}

func (im *impl) GetShowChannelLabels(ctx context.Context) ([]*channelDTO.ChannelLabelDTO, error) {
	logPos := "[channel_label.app.service][GetShowChannelLabels]"

	channelLabels, err := im.channelLabelStore.GetAllChannelLabels(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("channelLabelStore.GetAllChannelLabels ", err)
		return nil, errors.New("GetChannelCategoryType is nil")

	}

	showChannelLabels := []*channelDTO.ChannelLabelDTO{}

	for _, c := range channelLabels {
		if c.Show == 1 {
			showChannelLabels = append(showChannelLabels, c)
		}
	}
	return showChannelLabels, nil
}
