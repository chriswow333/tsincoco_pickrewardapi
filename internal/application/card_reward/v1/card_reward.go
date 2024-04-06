package application

import (
	"context"
	"encoding/json"

	"go.uber.org/dig"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"pickrewardapi/internal/application/card_reward/v1/handler"
	pb "pickrewardapi/internal/application/card_reward/v1/proto/generated"

	cardService "pickrewardapi/internal/domain/card/service"
)

type server struct {
	dig.In

	pb.UnimplementedCardRewardV1Server

	cardRewardService cardService.CardRewardService
}

func NewCardRewardServer(
	s *grpc.Server,

	cardRewardService cardService.CardRewardService,
) {
	log.WithFields(log.Fields{
		"pos": "[card.reward.api][NewCardRewardServer]",
	}).Info("Init")

	pb.RegisterCardRewardV1Server(s, &server{
		cardRewardService: cardRewardService,
	})
}

func (s *server) GetCardRewardsByCardID(ctx context.Context, in *pb.CardRewardIDReq) (*pb.CardRewardsReply, error) {
	logPos := "[card.reward.api][GetCardRewardsByCardID]"
	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	cardRewardDTOs, err := s.cardRewardService.GetCardRewardsByCardID(ctx, in.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("cardRewardService.GetCardRewardsByCardID failed: ", err)

		return &pb.CardRewardsReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "GetCardRewardsByCardID failed",
				},
			},
		}, nil
	}

	cardRewards := handler.TransferCardRewards2CardRewardsReply(cardRewardDTOs)

	cardRewardsLog, _ := json.Marshal(cardRewards)
	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": string(cardRewardsLog),
	}).Info("Response")

	return &pb.CardRewardsReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		CardRewards: cardRewards,
	}, nil
}
