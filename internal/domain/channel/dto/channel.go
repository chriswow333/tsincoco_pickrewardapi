package dto

import (
	commonM "pickrewardapi/internal/shared/common/model"
)

type ChannelDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	CreateDate int64 `json:"createDate"`
	UpdateDate int64 `json:"updateDate"`

	ChannelLabels []string `json:"channelLabels"`
	ShowLabel     string   `json:"showLabel"`

	Order         int32          `json:"order"`
	ImageName     string         `json:"imageName"`
	ChannelStatus commonM.Status `json:"channelStatus"`
}
