package dto

import commonM "pickrewardapi/internal/shared/common/model"

type DescriptionDTO struct {
	Name  string   `json:"name"`
	Order int32    `json:"order"`
	Desc  []string `json:"desc"`
}

type CardRewardType int32

const (
	Activity   CardRewardType = iota // 一般活動
	Evaluation                       // 可試算優惠
)

type CardRewardDTO struct {
	ID           string            `json:"id"`
	CardID       string            `json:"cardID"`
	Name         string            `json:"name"`
	Descriptions []*DescriptionDTO `json:"descriptions"`

	StartDate int64 `json:"startDate"`
	EndDate   int64 `json:"endDate"`

	CardRewardType CardRewardType `json:"cardRewardType"`

	FeedbackType *FeedbackTypeDTO `json:"feedbackType"`

	TaskLabelDTOs []*TaskLabelDTO `json:"taskLabels"`

	Order int32 `json:"order"`

	CardRewardStatus commonM.Status `json:"cardRewardStatus"`

	CreateDate int64 `json:"createDate"`
	UpdateDate int64 `json:"updateDate"`
}
