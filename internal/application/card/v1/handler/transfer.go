package handler

import (
	cardDTO "pickrewardapi/internal/domain/card/dto"

	pb "pickrewardapi/internal/application/card/v1/proto/generated"
)

func TransferCardsDTO2CardsReply(cardDTOs []*cardDTO.CardDTO) []*pb.CardsReply_Card {

	cards := []*pb.CardsReply_Card{}

	for _, c := range cardDTOs {
		cards = append(cards, &pb.CardsReply_Card{
			Id:           c.ID,
			Name:         c.Name,
			Descriptions: c.Descriptions,
			LinkURL:      c.LinkURL,
			BankID:       c.BankID,
			ImageName:    c.ImageName,
			Order:        c.Order,
			CardStatus:   int32(c.CardStatus),
			CreateDate:   c.CreateDate,
			UpdateDate:   c.UpdateDate,
		})
	}

	return cards
}

func TransferCardDTO2CardReply(cardDTO *cardDTO.CardDTO) *pb.CardReply_Card {

	return &pb.CardReply_Card{
		Id:           cardDTO.ID,
		Name:         cardDTO.Name,
		Descriptions: cardDTO.Descriptions,
		LinkURL:      cardDTO.LinkURL,
		BankID:       cardDTO.BankID,
		ImageName:    cardDTO.ImageName,
		Order:        cardDTO.Order,
		CardStatus:   int32(cardDTO.CardStatus),
		CreateDate:   cardDTO.CreateDate,
		UpdateDate:   cardDTO.UpdateDate,
	}
}
