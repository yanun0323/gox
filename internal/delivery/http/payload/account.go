package payload

import (
	"gox/internal/domain/entity"
	"gox/internal/domain/repository"
	"gox/internal/domain/usecase"
)

//goo:generate inspector -foo=bar -arg2 -arg3=123 -o ./internal/entity
type UpdateEmailReq struct {
	Email       string `json:"email" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
}

//go:generate esc-gen-model -replace -p2u -p2e -e member.go -u member.go -uu
type UpdatePhoneReq struct {
	Phone       string `json:"phone" binding:"required"`
	AreaCode    string `json:"area_code" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	CreateAt    int64  `json:"create_at"`
	UpdateAt    int64  `json:"update_at"`

	repository.UpdatePhoneEntity
}

func (elem *UpdatePhoneReq) ToEntity() *entity.UpdatePhoneReq { /* generate by esc-gen-model */
	return &entity.UpdatePhoneReq{
		Phone:             elem.Phone,
		AreaCode:          elem.AreaCode,
		CaptchaCode:       elem.CaptchaCode,
		CreateTime:        elem.CreateTime,
		UpdateTime:        elem.UpdateTime,
		CreateAt:          elem.CreateAt,
		UpdateAt:          elem.UpdateAt,
		UpdatePhoneEntity: elem.UpdatePhoneEntity,
	}
}

func (elem *UpdatePhoneReq) ToUseCase() *usecase.UpdatePhoneReq { /* generate by esc-gen-model */
	return &usecase.UpdatePhoneReq{
		Phone:             elem.Phone,
		AreaCode:          elem.AreaCode,
		CaptchaCode:       elem.CaptchaCode,
		CreateTime:        elem.CreateTime,
		UpdateTime:        elem.UpdateTime,
		CreateAt:          elem.CreateAt,
		UpdateAt:          elem.UpdateAt,
		UpdatePhoneEntity: elem.UpdatePhoneEntity,
	}
}

//goo:generate esc-gen-model -u="member.go" -r="member.go"
type UpdatePhoneResp struct {
	Phone       string `json:"phone" binding:"required"`
	AreaCode    string `json:"area_code" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	CreateAt    string `json:"create_at"`
	UpdateAt    string `json:"update_at"`
}

type Timer struct {
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}

type UpdateNicknameReq struct {
	Nickname string `json:"nickname" binding:"required"`
}

type UpdateFiatCurrencyReq struct {
	FiatCurrency string `json:"fiat_currency" binding:"required,uppercase"`
}

type GetPublicProfileReq struct {
	UID string `form:"uid"`
}

type UpdateLanguageReq struct {
	Language string `json:"language"`
}
