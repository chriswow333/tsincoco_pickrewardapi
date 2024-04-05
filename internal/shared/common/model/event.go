package model

type Event struct {
	ID string `json:"id"`

	Date int64 `json:"date"`
	Cost int32 `json:"cost"`

	ChannelEvent *ChannelEvent `json:"channelEvent"`
	PayEvent     *PayEvent     `json:"payEvent"`
	CardEvent    *CardEvent    `json:"cardEvent"`
}

type CardEvent struct {
	RewardType int32          `json:"rewardType"`
	TaskLabels map[int32]bool `json:"taskLabels"`
}

type PayStatus int32

const (
	Whatever PayStatus = iota
	Use
	No
)

type PayEvent struct {
	Status PayStatus       `json:"status"`
	PayIDs map[string]bool `json:"payIDs"`
}

type ChannelEvent struct {
	ChannelIDs    []*ChannelIDEvent `json:"channelIDs"`
	ChannelLabels map[int32]bool    `json:"channelLabels"`
}

type ChannelIDEvent struct {
	ChannelID     string         `json:"channelID"`
	ChannelLabels map[int32]bool `json:"channelLabels"`
}
