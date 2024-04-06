package service

import (
	"context"
	"errors"

	channelDTO "pickrewardapi/internal/domain/channel/dto"
	channelLabelStore "pickrewardapi/internal/domain/channel/store"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type ChannelLabelService interface {
	GetShowChannelLabels(ctx context.Context) ([]*channelDTO.ChannelLabelDTO, error)
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

func (im *channelLabelImpl) GetShowChannelLabels(ctx context.Context) ([]*channelDTO.ChannelLabelDTO, error) {
	logPos := "[channellabel.service][GetShowChannelLabels]"

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
