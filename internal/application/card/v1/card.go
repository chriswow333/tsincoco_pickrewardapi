package application

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"google.golang.org/grpc"

	pb "pickrewardapi/internal/application/card/v1/proto/generated"

	handler "pickrewardapi/internal/application/card/v1/handler"

	cardService "pickrewardapi/internal/domain/card/service"
)

type server struct {
	dig.In

	pb.UnimplementedCardV1Server

	cardService cardService.CardService
}

func NewCardServer(
	s *grpc.Server,

	cardService cardService.CardService,
) {
	log.WithFields(log.Fields{
		"pos": "[card.api][NewCardServer]",
	}).Info("Init")

	pb.RegisterCardV1Server(s, &server{
		cardService: cardService,
	})
}

func (s *server) GetCardsByBankID(ctx context.Context, in *pb.CardsByBankIDReq) (*pb.CardsReply, error) {
	logPos := "[card.api][GetCardsByBankID]"

	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	cardDTOs, err := s.cardService.GetCardsByBankID(ctx, in.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("cardService.GetCardsByBankID failed: ", err)

		return &pb.CardsReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "GetCardsByBankID failed",
				},
			},
		}, nil
	}

	cards := handler.TransferCardsDTO2CardsReply(cardDTOs)

	cardsLog, _ := json.Marshal(cards)
	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": string(cardsLog),
	}).Info("Response")

	return &pb.CardsReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		Cards: cards,
	}, nil
}

func (s *server) GetLatestCards(ctx context.Context, in *pb.EmptyReq) (*pb.CardsReply, error) {
	logPos := "[card.api][GetLatestCards]"

	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	cardDTOs, err := s.cardService.GetLatestCards(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("cardService.GetLatestCards failed: ", err)

		return &pb.CardsReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "GetLatestCards failed",
				},
			},
		}, nil
	}

	cards := handler.TransferCardsDTO2CardsReply(cardDTOs)

	cardsLog, _ := json.Marshal(cards)
	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": string(cardsLog),
	}).Info("Response")

	return &pb.CardsReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		Cards: cards,
	}, nil
}

func (s *server) GetCardByID(ctx context.Context, in *pb.CardIDReq) (*pb.CardReply, error) {
	logPos := "[card.api][GetCardByID]"

	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	cardDTO, err := s.cardService.GetCardByID(ctx, in.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("cardService.GetCardByID failed: ", err)

		return &pb.CardReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "GetCardByID failed",
				},
			},
		}, nil
	}

	card := handler.TransferCardDTO2CardReply(cardDTO)
	cardsLog, _ := json.Marshal(card)
	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": string(cardsLog),
	}).Info("Response")

	return &pb.CardReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		Card: card,
	}, nil
}

func (s *server) SearchCard(ctx context.Context, in *pb.SearchCardReq) (*pb.CardsReply, error) {
	logPos := "[card.api][SearchCard]"

	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	cardDTOs, err := s.cardService.SearchCard(ctx, in.Keyword)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("cardService.SearchCard failed: ", err)

		return &pb.CardsReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "SearchCard failed",
				},
			},
		}, nil
	}

	cards := handler.TransferCardsDTO2CardsReply(cardDTOs)
	cardsLog, _ := json.Marshal(cards)
	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": string(cardsLog),
	}).Info("Response")

	return &pb.CardsReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		Cards: cards,
	}, nil
}
