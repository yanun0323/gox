package usecase

import (
	"gox/internal/payload"
	"gox/internal/payload/internal/domain/repository"
)

type UpdatePhoneEntity struct {
	Phone       string `json:"phone" binding:"required"`
	AreaCode    string `json:"area_code" binding:"required"`
	CaptchaCode string `json:"code" binding:"required"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	CreateAt    int64  `json:"create_at"`
	UpdateAt    int64  `json:"update_at"`
}

func (upe *UpdatePhoneEntity) ToRepository() *repository.UpdatePhoneEntity {
	return &repository.UpdatePhoneEntity{
		Phone:       upe.Phone,
		AreaCode:    upe.AreaCode,
		CaptchaCode: upe.CaptchaCode,
		CreateTime:  upe.CreateTime,
		UpdateTime:  upe.UpdateTime,
		CreateAt:    upe.CreateAt,
		UpdateAt:    upe.UpdateAt,
	}
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

func (upr *UpdatePhoneResp) ToPayload() *payload.UpdatePhoneResp {
	return &payload.UpdatePhoneResp{
		Phone:       upr.Phone,
		AreaCode:    upr.AreaCode,
		CaptchaCode: upr.CaptchaCode,
		CreateTime:  upr.CreateTime,
		UpdateTime:  upr.UpdateTime,
		CreateAt:    upr.CreateAt,
		UpdateAt:    upr.UpdateAt,
	}
}
