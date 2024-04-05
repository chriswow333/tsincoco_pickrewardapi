package service

import (
	"context"
	"errors"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"

	channelDTO "pickrewardapi/internal/domain/channel/dto"

	channelDomain "pickrewardapi/internal/domain/channel/domain"
	channelStore "pickrewardapi/internal/domain/channel/store"
	commonM "pickrewardapi/internal/shared/common/model"
)

type ChannelService interface {
	GetChannelTypes(ctx context.Context) []*channelDTO.ChannelTypeDTO
	GetChannelsByType(ctx context.Context, ctype int32, limit, offset int32) ([]*channelDTO.ChannelDTO, error)
	GetsByChannelIDs(ctx context.Context, IDs []string) ([]*channelDTO.ChannelDTO, error)
	SearchChannel(ctx context.Context, keyword string) ([]*channelDTO.ChannelDTO, error)

	GetChannelTypeByType(ctx context.Context, ctype int32) (*channelDTO.ChannelTypeDTO, error)
}

type impl struct {
	dig.In

	channelStore channelStore.ChannelStore
}

var (
	timeNow = time.Now
)

func New(
	channelStore channelStore.ChannelStore,
) ChannelService {

	impl := &impl{
		channelStore: channelStore,
	}

	return impl
}

func (im *impl) GetChannelTypeByType(ctx context.Context, ctype int32) (*channelDTO.ChannelTypeDTO, error) {
	logPos := "[channel.service][GetChannelTypeByType]"

	channelType := channelDomain.GetChannelType(channelDomain.ChannelTypeEnum(ctype))
	if channelType == nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("GetChannelType is nil")
		return nil, errors.New("GetChannelType is nil")
	}

	return &channelDTO.ChannelTypeDTO{
		ChannelType: int32(channelType.Type),
		Name:        channelType.Name,
		Order:       channelType.Order,
	}, nil
}

func (im *impl) GetChannelTypes(ctx context.Context) []*channelDTO.ChannelTypeDTO {

	channelCategoryDTOs := []*channelDTO.ChannelTypeDTO{}

	for _, c := range channelDomain.GetChannelTypes() {
		channelCategoryDTOs = append(channelCategoryDTOs, &channelDTO.ChannelTypeDTO{
			ChannelType: int32(c.Type),
			Name:        c.Name,
			Order:       c.Order,
		})
	}
	sort.SliceStable(channelCategoryDTOs, func(i, j int) bool {
		return channelCategoryDTOs[i].Order < channelCategoryDTOs[j].Order
	})

	return channelCategoryDTOs
}

func (im *impl) GetChannelsByType(ctx context.Context, channelType int32, limit, offset int32) ([]*channelDTO.ChannelDTO, error) {
	logPos := "[channel.service][GetChannelsByType]"

	channelDTOs, err := im.channelStore.GetChannelsByType(ctx, channelType, commonM.Active, limit, offset)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":          logPos,
			"channel.type": channelType,
		}).Error("channelStore.GetChannelsByType: ", err)
		return nil, err
	}

	sort.SliceStable(channelDTOs, func(i, j int) bool {
		return channelDTOs[i].Order < channelDTOs[j].Order
	})

	return channelDTOs, nil
}

func (im *impl) GetByChannelID(ctx context.Context, ID string) (*channelDTO.ChannelDTO, error) {
	return im.channelStore.GetChannelByID(ctx, ID)
}

func (im *impl) GetsByChannelIDs(ctx context.Context, IDs []string) ([]*channelDTO.ChannelDTO, error) {
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

func (im *impl) SearchChannel(ctx context.Context, keyword string) ([]*channelDTO.ChannelDTO, error) {
	logPos := "[channel.service][SearchChannel]"

	channelDTOs, err := im.channelStore.SearchChannel(ctx, keyword, commonM.Active)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":     logPos,
			"keyword": keyword,
		}).Error("channelStore.SearchChannel: ", err)
		return nil, err
	}

	sort.SliceStable(channelDTOs, func(i, j int) bool {
		return channelDTOs[i].ChannelType < channelDTOs[j].ChannelType
	})

	return channelDTOs, nil

}
