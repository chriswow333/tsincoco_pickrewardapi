package event

type FeedbackEventResultStatus int32

const (
	GetNone FeedbackEventResultStatus = iota // 遇到feedback限制 拿不到
	GetSome                                  // 遇到feedback限制 拿一點
	GetAll                                   // 回饋全拿
)

type FeedbackEventResult struct {
	FeedbackID                string                    `json:"feedbackID"`
	CalculateType             int32                     `json:"calculateType"`
	Cost                      int32                     `json:"cost"`          // 花費多少
	GetReturn                 float64                   `json:"getReturn"`     // 回饋多少
	GetPercentage             float64                   `json:"getPercentage"` // 如果calculate type is multiply, 則拿到的趴數是多少, 後續作加總用
	FeedbackEventResultStatus FeedbackEventResultStatus `json:"feedbackEventResultStatus"`
}
