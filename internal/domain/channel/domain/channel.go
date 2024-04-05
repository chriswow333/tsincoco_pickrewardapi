package domain

type ChannelTypeEnum int32

const (
	Ecommerce ChannelTypeEnum = iota
	Food
	Travel
	Transportation
	Oversea
	Streaming
	PayFee
	Mall
	Insurance
	Supermarket
	Delivery
	Sport
)

type ChannelType struct {
	Type  ChannelTypeEnum
	Name  string
	Order int32
}

var categoryTypeMap = map[ChannelTypeEnum]*ChannelType{
	Ecommerce: {
		Type:  Ecommerce,
		Name:  "網路購物",
		Order: 0,
	},
	Food: {
		Type:  Food,
		Name:  "美食",
		Order: 1,
	},
	Travel: {
		Type:  Travel,
		Name:  "旅遊",
		Order: 2,
	},
	Transportation: {
		Type:  Transportation,
		Name:  "交通",
		Order: 3,
	},
	Oversea: {
		Type:  Oversea,
		Name:  "海外消費",
		Order: 4,
	},
	Streaming: {
		Type:  Streaming,
		Name:  "影音/串流",
		Order: 5,
	},
	PayFee: {
		Type:  PayFee,
		Name:  "繳費",
		Order: 6,
	},
	Mall: {
		Type:  Mall,
		Name:  "百貨/影城",
		Order: 7,
	},
	Insurance: {
		Type:  Insurance,
		Name:  "保險",
		Order: 8,
	},
	Supermarket: {
		Type:  Supermarket,
		Name:  "量販/超市",
		Order: 9,
	},
	Delivery: {
		Type:  Delivery,
		Name:  "外送",
		Order: 10,
	},
	Sport: {
		Type:  Sport,
		Name:  "休閒/運動",
		Order: 11,
	},
}

func GetChannelTypes() map[ChannelTypeEnum]*ChannelType {
	return categoryTypeMap
}

func GetChannelType(c ChannelTypeEnum) *ChannelType {
	return categoryTypeMap[c]
}
