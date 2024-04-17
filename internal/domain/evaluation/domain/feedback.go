package domain

type Feedback struct {
	RewardType    int32 // CASH, POINT
	CalculateType CalculateType
	MinCost       int32   // 最少要花多少錢
	Fixed         int32   // 折抵多少元/點
	Percentage    float64 // 幾趴回饋
	ReturnMax     float64 // 最多回饋多少錢
}

type CalculateType int32

const (
	Multiply CalculateType = iota
	Fixed
	Area
)
