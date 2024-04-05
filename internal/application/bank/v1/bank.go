package application

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"go.uber.org/dig"
	"google.golang.org/grpc"

	"pickrewardapi/internal/application/bank/v1/handler"
	pb "pickrewardapi/internal/application/bank/v1/proto/generated"

	bankService "pickrewardapi/internal/domain/bank/service"
)

type server struct {
	dig.In

	pb.UnimplementedBankV1Server

	bankService bankService.BankService
}

func NewBankServer(
	s *grpc.Server,

	bankService bankService.BankService,
) {
	log.WithFields(log.Fields{
		"pos": "[app.bank][NewBankServer]",
	}).Info("Init")

	pb.RegisterBankV1Server(s, &server{
		bankService: bankService,
	})
}

func (s *server) GetAllBanks(ctx context.Context, in *pb.AllBanksReq) (*pb.BanksReply, error) {
	logPos := "[bank.api][GetAllBanks]"

	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	banksDTOs, err := s.bankService.GetAllBanks(ctx)

	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("bankService.GetAllBanks failed: ", err)

		return &pb.BanksReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "GetAllBanks failed",
				},
			},
		}, nil
	}

	banks := handler.TransferBankDTOsToBank(banksDTOs)

	banksLog, _ := json.Marshal(banks)
	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": string(banksLog),
	}).Info("Response")

	return &pb.BanksReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		Banks: banks,
	}, nil

}
