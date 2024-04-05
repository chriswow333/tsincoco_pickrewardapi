package handler

import (
	pb "pickrewardapi/internal/application/channel/v1/proto/generated"
	channelDTO "pickrewardapi/internal/domain/channel/dto"
)

func TransferChannelTypeDTO2ChannelTypeReply(channelTypeDTOs []*channelDTO.ChannelTypeDTO) []*pb.ChannelTypesReply_ChannelType {
	channelTypes := []*pb.ChannelTypesReply_ChannelType{}
	for _, c := range channelTypeDTOs {
		channelTypes = append(channelTypes, &pb.ChannelTypesReply_ChannelType{
			ChannelType: c.ChannelType,
			Name:        c.Name,
			Order:       c.Order,
		})
	}
	return channelTypes
}

func TransferChannels2ChannelsReply(channelDTOs []*channelDTO.ChannelDTO) []*pb.ChannelsReply_Channel {
	channels := []*pb.ChannelsReply_Channel{}

	for _, c := range channelDTOs {

		channels = append(channels, &pb.ChannelsReply_Channel{
			Id:            c.ID,
			Name:          c.Name,
			LinkURL:       c.LinkURL,
			ChannelType:   c.ChannelType,
			CreateDate:    c.CreateDate,
			UpdateDate:    c.UpdateDate,
			ChannelLabels: c.ChannelLabels,
			Order:         c.Order,
			ChannelStatus: int32(c.ChannelStatus),
		})
	}
	return channels
}

func TransferSearchChannels2SearchChannelsReply(channelDTOs []*channelDTO.ChannelDTO) []*pb.SearchChannelsReply_Channel {
	channels := []*pb.SearchChannelsReply_Channel{}

	for _, c := range channelDTOs {

		channels = append(channels, &pb.SearchChannelsReply_Channel{
			Id:            c.ID,
			Name:          c.Name,
			LinkURL:       c.LinkURL,
			ChannelType:   c.ChannelType,
			CreateDate:    c.CreateDate,
			UpdateDate:    c.UpdateDate,
			ChannelLabels: c.ChannelLabels,
			Order:         c.Order,
			ChannelStatus: int32(c.ChannelStatus),
		})
	}

	return channels
}
