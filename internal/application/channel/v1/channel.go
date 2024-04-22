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

	channelService      channelService.ChannelService
	channelLabelService channelService.ChannelLabelService
}

func NewChannelServer(
	s *grpc.Server,

	channelService channelService.ChannelService,
	channelLabelService channelService.ChannelLabelService,
) {
	log.WithFields(log.Fields{
		"pos": "[channel.api][NewChannelServer]",
	}).Info("Init")

	pb.RegisterChannelV1Server(s, &server{
		channelService:      channelService,
		channelLabelService: channelLabelService,
	})
}

func (s *server) GetShowLabels(ctx context.Context, in *pb.EmptyReq) (*pb.ShowLabelsReply, error) {
	logPos := "[channel.api][GetShowLabels]"

	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	showLabelDTOs, err := s.channelLabelService.GetShowLabels(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("channelService.GetShowLabels failed: ", err)

		return &pb.ShowLabelsReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "GetShowLabels failed",
				},
			},
		}, nil
	}

	showLabels := handler.TransferShowLabelsReply(showLabelDTOs)

	showLabelLogs, _ := json.Marshal(showLabels)
	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": string(showLabelLogs),
	}).Info("Response")

	return &pb.ShowLabelsReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		ChannelLabels: showLabels,
	}, nil
}

func (s *server) GetChannelsByShowLabel(ctx context.Context, in *pb.ShowLabelReq) (*pb.ChannelsReply, error) {
	logPos := "[channel.api][GetChannelsByShowLabel]"
	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	channelDTOs, err := s.channelService.GetChannelsByLabel(ctx, in.Label, in.Limit, in.Offset)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("channelService.GetChannelsByLabel failed: ", err)

		return &pb.ChannelsReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "GetChannelsByLabel failed",
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
