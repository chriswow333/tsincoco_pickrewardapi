package handler

import (
	pb "pickrewardapi/internal/application/card_reward/v1/proto/generated"
	cardDTO "pickrewardapi/internal/domain/card/dto"
)

func TransferCardRewards2CardRewardsReply(cardRewardDTOs []*cardDTO.CardRewardDTO) []*pb.CardRewardsReply_CardReward {

	cardRewards := []*pb.CardRewardsReply_CardReward{}

	for _, c := range cardRewardDTOs {

		descriptions := []*pb.CardRewardsReply_Description{}
		for _, d := range c.Descriptions {
			descriptions = append(descriptions, &pb.CardRewardsReply_Description{
				Name:  d.Name,
				Order: d.Order,
				Desc:  d.Desc,
			})
		}

		feedbackType := &pb.CardRewardsReply_FeedbackType{
			Id:           c.FeedbackType.ID,
			Name:         c.FeedbackType.Name,
			FeedbackType: int32(c.FeedbackType.FeedbackType),
			CreateDate:   c.FeedbackType.CreateDate,
			Updatedate:   c.FeedbackType.UpdateDate,
		}

		taskLabels := []*pb.CardRewardsReply_TaskLabel{}

		for _, t := range c.TaskLabelDTOs {
			taskLabels = append(taskLabels, &pb.CardRewardsReply_TaskLabel{
				Label: t.Label,
				Name:  t.Name,
				Show:  t.Show,
				Order: t.Order,
			})
		}

		cardRewards = append(cardRewards, &pb.CardRewardsReply_CardReward{
			Id:               c.ID,
			CardID:           c.CardID,
			Name:             c.Name,
			Descriptions:     descriptions,
			StartDate:        c.StartDate,
			EndDate:          c.EndDate,
			CardRewardType:   int32(c.CardRewardType),
			FeedbackType:     feedbackType,
			TaskLabels:       taskLabels,
			Order:            c.Order,
			CardRewardStatus: int32(c.CardRewardStatus),
			CreateDate:       c.CreateDate,
			UpdateDate:       c.UpdateDate,
		})
	}

	return cardRewards
}
