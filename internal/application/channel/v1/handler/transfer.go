package handler

import (
	pb "pickrewardapi/internal/application/channel/v1/proto/generated"
	channelDTO "pickrewardapi/internal/domain/channel/dto"
)

func TransferShowLabelsReply(showLabelDTOs []*channelDTO.ChannelLabelDTO) []*pb.ShowLabelsReply_ChannelLabel {
	showLabels := []*pb.ShowLabelsReply_ChannelLabel{}
	for _, c := range showLabelDTOs {
		showLabels = append(showLabels, &pb.ShowLabelsReply_ChannelLabel{
			Label: c.Label,
			Name:  c.Name,
			Order: c.Order,
		})
	}
	return showLabels
}

func TransferChannels2ChannelsReply(channelDTOs []*channelDTO.ChannelDTO) []*pb.ChannelsReply_Channel {
	channels := []*pb.ChannelsReply_Channel{}

	for _, c := range channelDTOs {

		channels = append(channels, &pb.ChannelsReply_Channel{
			Id:            c.ID,
			Name:          c.Name,
			CreateDate:    c.CreateDate,
			UpdateDate:    c.UpdateDate,
			ChannelLabels: c.ChannelLabels,
			ShowLabel:     c.ShowLabel,
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
			CreateDate:    c.CreateDate,
			UpdateDate:    c.UpdateDate,
			ChannelLabels: c.ChannelLabels,
			ShowLabel:     c.ShowLabel,
			Order:         c.Order,
			ChannelStatus: int32(c.ChannelStatus),
		})
	}

	return channels
}
