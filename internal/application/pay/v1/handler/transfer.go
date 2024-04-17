package handler

import (
	pb "pickrewardapi/internal/application/pay/v1/proto/generated"

	payDTO "pickrewardapi/internal/domain/pay/dto"
)

func TransferPay2PayReply(payDTO *payDTO.PayDTO) *pb.PayReply_Pay {

	return &pb.PayReply_Pay{
		Id:         payDTO.ID,
		Name:       payDTO.Name,
		Order:      payDTO.Order,
		PayStatus:  int32(payDTO.PayStatus),
		CreateDate: payDTO.CreateDate,
		UpdateDate: payDTO.UpdateDate,
	}
}

func TransferPays2PaysReply(payDTOs []*payDTO.PayDTO) []*pb.PaysReply_Pay {

	pays := []*pb.PaysReply_Pay{}
	for _, p := range payDTOs {
		pays = append(pays, &pb.PaysReply_Pay{
			Id:         p.ID,
			Name:       p.Name,
			Order:      p.Order,
			PayStatus:  int32(p.PayStatus),
			CreateDate: p.CreateDate,
			UpdateDate: p.UpdateDate,
		})
	}
	return pays
}
