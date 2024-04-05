package dto

import (
	commonM "pickrewardapi/internal/shared/common/model"
)

type ChannelTypeDTO struct {
	ChannelType int32  `json:"channelType"`
	Name        string `json:"name"`
	Order       int32  `json:"order"`
}

type ChannelDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	LinkURL string `json:"linkURL"`

	ChannelType int32 `json:"channelType"`
	CreateDate  int64 `json:"createDate"`
	UpdateDate  int64 `json:"updateDate"`

	ChannelLabels []int32        `json:"channelLabels"`
	Order         int32          `json:"order"`
	ChannelStatus commonM.Status `json:"channelStatus"`
}
