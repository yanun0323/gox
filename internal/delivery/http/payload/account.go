package payload

import "gox/internal/domain/usecase"

//go:generate inspector -foo=bar -arg2 -arg3=123 -o ./internal/entity
type UpdateEmailReq struct {
	Email       string `json:"email" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
}

//go:generate esc-gen-model -v -replace -unix -use="member.go" -repo="member.go"
type UpdatePhoneReq struct {
	Phone       string `json:"phone" binding:"required"`
	AreaCode    string `json:"area_code" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	CreateAt    int64  `json:"create_at"`
	UpdateAt    int64  `json:"update_at"`
}

func (upr *UpdatePhoneReq) ToUseCase() *usecase.UpdatePhoneReq {
	return &usecase.UpdatePhoneReq{
		Phone:       upr.Phone,
		AreaCode:    upr.AreaCode,
		CaptchaCode: upr.CaptchaCode,
		CreateTime:  upr.CreateTime,
		UpdateTime:  upr.UpdateTime,
		CreateAt:    upr.CreateAt,
		UpdateAt:    upr.UpdateAt,
	}
}

//go:generate esc-gen-model -v -use="member.go" -repo="member.go"
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
