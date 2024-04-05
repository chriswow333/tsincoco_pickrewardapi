package model

import (
	"errors"
	"strconv"
)

type LabelType int32

const (
	All        LabelType = iota // 不分通路
	Domestic                    // 國內消費
	Oversea                     // 海外消費
	Digit                       // 數位通路
	Physical                    // 實體通路
	Restaurant                  // 全臺餐廳

)

var (
	labelMapper = make(map[LabelType]*Label)
)

func init() {
	labelMapper = map[LabelType]*Label{
		All: {
			LabelType: All,
			LabelName: "不分通路",
		},
		Domestic: {
			LabelType: Domestic,
			LabelName: "國內消費",
		},
		Oversea: {
			LabelType: Oversea,
			LabelName: "海外消費",
		},
		Digit: {
			LabelType: Digit,
			LabelName: "數位通路",
		},
		Physical: {
			LabelType: Physical,
			LabelName: "實體通路",
		},
		Restaurant: {
			LabelType: Restaurant,
			LabelName: "全臺餐廳",
		},
	}
}

func GetLabel(labelType int32) (*Label, error) {

	if labelType > int32(Restaurant) {
		return nil, errors.New("Cannot find Label : " + strconv.Itoa(int(labelType)))
	}

	label, ok := labelMapper[LabelType(labelType)]

	if !ok {
		return nil, errors.New("Cannot find point type")
	}

	return label, nil
}

type Label struct {
	LabelType LabelType `json:"labelType"`
	LabelName string    `json:"labelName"`
}
