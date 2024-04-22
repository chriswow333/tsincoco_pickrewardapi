package dto

type ChannelLabelDTO struct {
	Label string `json:"label"`
	Name  string `json:"name"`
	Show  int32  `json:"show"`
	Order int32  `json:"order"`
}
