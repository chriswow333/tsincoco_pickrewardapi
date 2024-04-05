package apis

type Reply struct {
	Status int32  `json:"status"`
	Error  string `json:"error"`
}
