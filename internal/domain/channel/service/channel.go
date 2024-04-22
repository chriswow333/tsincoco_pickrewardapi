package service

import (
	"context"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	channelDTO "pickrewardapi/internal/domain/channel/dto"

	channelStore "pickrewardapi/internal/domain/channel/store"
	commonM "pickrewardapi/internal/shared/common/model"
)

type ChannelService interface {

	// GetChannelTypes(ctx context.Context) []*channelDTO.ChannelTypeDTO
	// GetChannelsByType(ctx context.Context, ctype int32, limit, offset int32) ([]*channelDTO.ChannelDTO, error)

	GetsByChannelIDs(ctx context.Context, IDs []string) ([]*channelDTO.ChannelDTO, error)
	SearchChannel(ctx context.Context, keyword string) ([]*channelDTO.ChannelDTO, error)

	GetByChannelID(ctx context.Context, ID string) (*channelDTO.ChannelDTO, error)

	GetChannelsByLabel(ctx context.Context, channelLabel string, offset, limit int32) ([]*channelDTO.ChannelDTO, error)

	// GetChannelTypeByType(ctx context.Context, ctype int32) (*channelDTO.ChannelTypeDTO, error)

}

type channelImpl struct {
	dig.In

	channelStore channelStore.ChannelStore
}

var (
	timeNow = time.Now
)

func NewChannel(
	channelStore channelStore.ChannelStore,
) ChannelService {

	impl := &channelImpl{
		channelStore: channelStore,
	}

	return impl
}

func (im *channelImpl) GetChannelsByLabel(ctx context.Context, channelLabel string, offset, limit int32) ([]*channelDTO.ChannelDTO, error) {
	return im.channelStore.GetChannelsByShowLabel(ctx, channelLabel, commonM.Active, offset, limit)
}

func (im *channelImpl) GetByChannelID(ctx context.Context, ID string) (*channelDTO.ChannelDTO, error) {
	return im.channelStore.GetChannelByID(ctx, ID)
}

func (im *channelImpl) GetsByChannelIDs(ctx context.Context, IDs []string) ([]*channelDTO.ChannelDTO, error) {
	logPos := "[channel.service][GetsByChannelIDs]"

	channelDTOs, err := im.channelStore.GetChannelByIDs(ctx, IDs)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
			"IDs": IDs,
		}).Error("GetChannelByIDs failed: ", err)
		return nil, err
	}

	sort.SliceStable(channelDTOs, func(i, j int) bool {
		return channelDTOs[i].Order < channelDTOs[j].Order
	})

	return channelDTOs, nil
}

func (im *channelImpl) SearchChannel(ctx context.Context, keyword string) ([]*channelDTO.ChannelDTO, error) {
	logPos := "[channel.service][SearchChannel]"

	channelDTOs, err := im.channelStore.SearchChannel(ctx, keyword, commonM.Active)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"keyword": keyword,
		}).Error("channelStore.SearchChannel: ", err)
		return nil, err
	}

	return channelDTOs, nil

}
