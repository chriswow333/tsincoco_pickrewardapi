package application

import (
	"context"
	"encoding/json"
	pb "pickrewardapi/internal/application/pay/v1/proto/generated"

	log "github.com/sirupsen/logrus"

	"go.uber.org/dig"
	"google.golang.org/grpc"

	handler "pickrewardapi/internal/application/pay/v1/handler"

	payService "pickrewardapi/internal/domain/pay/service"
)

type server struct {
	dig.In

	pb.UnimplementedPayV1Server

	payService payService.PayService
}

func NewPayServer(

	s *grpc.Server,

	payService payService.PayService,

) {
	log.WithFields(log.Fields{
		"pos": "[pay.api][NewPayServer]",
	}).Info("Init")
	pb.RegisterPayV1Server(s, &server{
		payService: payService,
	})
}

func (s *server) GetPayByID(ctx context.Context, in *pb.PayIDReq) (*pb.PayReply, error) {
	logPos := "[pay.api][GetPayByID]"

	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	payDTO, err := s.payService.GetPayByID(ctx, in.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("payService.GetPayByID failed: ", err)
		return &pb.PayReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "GetChannelsByType failed",
				},
			},
		}, nil
	}

	pay := handler.TransferPay2PayReply(payDTO)
	channelsLog, _ := json.Marshal(pay)
	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": string(channelsLog),
	}).Info("Response")

	return &pb.PayReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		Pay: pay,
	}, nil
}

func (s *server) GetAllPays(ctx context.Context, in *pb.EmptyReq) (*pb.PaysReply, error) {
	logPos := "[pay.api][GetAllPays]"

	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	payDTOs, err := s.payService.GetAllPays(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("payService.GetAllPays failed: ", err)
		return &pb.PaysReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "GetAllPays failed",
				},
			},
		}, nil
	}

	pays := handler.TransferPays2PaysReply(payDTOs)

	paysLog, _ := json.Marshal(pays)
	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": string(paysLog),
	}).Info("Response")

	return &pb.PaysReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		Pays: pays,
	}, nil
}
