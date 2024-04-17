package dto

import (
	commonM "pickrewardapi/internal/shared/common/model"
)

type PayDTO struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Order      int32          `json:"order"`
	PayStatus  commonM.Status `json:"payStatus"`
	CreateDate int64          `json:"createDate"`
	UpdateDate int64          `json:"updateDate"`
}
