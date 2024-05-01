package dto

import (
	commonM "pickrewardapi/internal/shared/common/model"
)

type CardDTO struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Descriptions []string       `json:"descriptions"`
	LinkURL      string         `json:"linkURL"`
	BankID       string         `json:"bankID"`
	ImageName    string         `json:"imageName"`
	Order        int32          `json:"order"`
	CardStatus   commonM.Status `json:"cardStatus"`
	CreateDate   int64          `json:"createDate"`
	UpdateDate   int64          `json:"updateDate"`
}
