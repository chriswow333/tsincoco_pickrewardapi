package dto

import (
	commonM "pickrewardapi/internal/shared/common/model"
)

type BankDTO struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Order      int32          `json:"order"`
	BankStatus commonM.Status `json:"bankStatus"`
	CreateDate int64          `json:"createDate"`
	UpdateDate int64          `json:"updateDate"`
}
