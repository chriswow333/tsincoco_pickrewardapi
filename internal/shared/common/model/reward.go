package model

// type RewardOwner int32

// const (
// 	CARD RewardOwner = iota
// 	CHANNEL
// 	PAY
// 	TASK
// )

// type RewardType int32

// const (
// 	NONE_REWARD RewardType = iota
// 	CASH
// 	POINT
// )

// var (
// 	rewardMapper = make(map[RewardType]*Reward)
// )

// func init() {
// 	rewardMapper = map[RewardType]*Reward{
// 		NONE_REWARD: {
// 			RewardType: NONE_REWARD,
// 			RewardName: "無",
// 		},
// 		CASH: {
// 			RewardType: CASH,
// 			RewardName: "現金回饋",
// 		},
// 		POINT: {
// 			RewardType: POINT,
// 			RewardName: "點數回饋",
// 		},
// 	}
// }

// type Reward struct {
// 	RewardType RewardType `json:"rewardType"`
// 	RewardName string     `json:"rewardName"`
// }

// func GetAllRewardTypes() []*Reward {

// 	rewards := []*Reward{}

// 	for _, v := range rewardMapper {
// 		rewards = append(rewards, v)
// 	}
// 	return rewards
// }

// func GetRewardType(rewardType int32) (*Reward, error) {

// 	reward, ok := rewardMapper[RewardType(rewardType)]

// 	if !ok {
// 		return nil, errors.New("Cannot find reward type")
// 	}
// 	return reward, nil
// }
