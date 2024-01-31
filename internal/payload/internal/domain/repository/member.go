package repository

import "gox/internal/payload/internal/domain/usecase"

type UpdatePhoneEntity struct {
	Phone       string `json:"phone" binding:"required"`
	AreaCode    string `json:"area_code" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	CreateAt    int64  `json:"create_at"`
	UpdateAt    int64  `json:"update_at"`
}

type UpdatePhoneResp struct {
	Phone       string `json:"phone" binding:"required"`
	AreaCode    string `json:"area_code" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	CreateAt    int64  `json:"create_at"`
	UpdateAt    int64  `json:"update_at"`
}

func (upr *UpdatePhoneResp) ToUseCase() *usecase.UpdatePhoneResp {
	return &usecase.UpdatePhoneResp{
		Phone:       upr.Phone,
		AreaCode:    upr.AreaCode,
		CaptchaCode: upr.CaptchaCode,
		CreateTime:  upr.CreateTime,
		UpdateTime:  upr.UpdateTime,
		CreateAt:    upr.CreateAt,
		UpdateAt:    upr.UpdateAt,
	}
}
