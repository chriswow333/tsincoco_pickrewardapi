package application

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"google.golang.org/grpc"

	pb "pickrewardapi/internal/application/channel/v1/proto/generated"

	handler "pickrewardapi/internal/application/channel/v1/handler"

	channelService "pickrewardapi/internal/domain/channel/service"
)

type server struct {
	dig.In

	pb.UnimplementedChannelV1Server

	channelService channelService.ChannelService
}

func NewChannelServer(
	s *grpc.Server,

	channelService channelService.ChannelService,
) {
	log.WithFields(log.Fields{
		"pos": "[channel.api][NewChannelServer]",
	}).Info("Init")

	pb.RegisterChannelV1Server(s, &server{
		channelService: channelService,
	})
}

func (s *server) GetChannelTypes(ctx context.Context, in *pb.EmptyReq) (*pb.ChannelTypesReply, error) {
	logPos := "[channel.api][GetChannelTypes]"

	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	channelTypeDTOs := s.channelService.GetChannelTypes(ctx)

	channelTypes := handler.TransferChannelTypeDTO2ChannelTypeReply(channelTypeDTOs)

	channelTypesLog, _ := json.Marshal(channelTypes)
	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": string(channelTypesLog),
	}).Info("Response")

	return &pb.ChannelTypesReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		ChannelTypes: channelTypes,
	}, nil
}

func (s *server) GetChannelsByType(ctx context.Context, in *pb.ChannelTypeReq) (*pb.ChannelsReply, error) {
	logPos := "[channel.api][GetChannelsByType]"
	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	channelDTOs, err := s.channelService.GetChannelsByType(ctx, in.Ctype, in.Limit, in.Offset)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("channelService.GetChannelsByType failed: ", err)

		return &pb.ChannelsReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "GetChannelsByType failed",
				},
			},
		}, nil
	}

	channels := handler.TransferChannels2ChannelsReply(channelDTOs)
	channelsLog, _ := json.Marshal(channels)
	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": string(channelsLog),
	}).Info("Response")

	return &pb.ChannelsReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		Channels: channels,
	}, nil
}

func (s *server) GetsByChannelIDs(ctx context.Context, in *pb.ChannelIDsReq) (*pb.ChannelsReply, error) {
	logPos := "[channel.api][GetsByChannelIDs]"
	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	channelDTOs, err := s.channelService.GetsByChannelIDs(ctx, in.ChannelIDs)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("channelService.GetsByChannelIDs failed: ", err)

		return &pb.ChannelsReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "GetsByChannelIDs failed",
				},
			},
		}, nil
	}

	channels := handler.TransferChannels2ChannelsReply(channelDTOs)
	channelsLog, _ := json.Marshal(channels)
	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": string(channelsLog),
	}).Info("Response")

	return &pb.ChannelsReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		Channels: channels,
	}, nil
}

func (s *server) SearchChannel(ctx context.Context, in *pb.SearchChannelReq) (*pb.SearchChannelsReply, error) {
	logPos := "[channel.api][SearchChannel]"
	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	channelDTOs, err := s.channelService.SearchChannel(ctx, in.Keyword)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("channelService.SearchChannel failed: ", err)

		return &pb.SearchChannelsReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "SearchChannel failed",
				},
			},
		}, nil
	}

	channels := handler.TransferSearchChannels2SearchChannelsReply(channelDTOs)
	channelsLog, _ := json.Marshal(channels)
	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": string(channelsLog),
	}).Info("Response")

	return &pb.SearchChannelsReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		Channels: channels,
	}, nil

}
