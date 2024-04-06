package handler

import (
	pb "pickrewardapi/internal/application/bank/v1/proto/generated"

	bankDTO "pickrewardapi/internal/domain/card/dto"
)

func TransferBankDTOsToBank(bankDTOs []*bankDTO.BankDTO) []*pb.BanksReply_Bank {

	banks := []*pb.BanksReply_Bank{}

	for _, b := range bankDTOs {

		banks = append(banks, &pb.BanksReply_Bank{
			Id:         b.ID,
			Name:       b.Name,
			Order:      b.Order,
			BankStatus: int32(b.BankStatus),
		})

	}

	return banks
}
