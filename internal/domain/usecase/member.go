package usecase

import "gox/internal/domain/repository"

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
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	CreateAt    string `json:"create_at"`
	UpdateAt    string `json:"update_at"`
}

type UpdatePhoneReq struct {
	Phone       string
	AreaCode    string
	CaptchaCode string
	CreateTime  int64
	UpdateTime  int64
	CreateAt    int64
	UpdateAt    int64
	repository.UpdatePhoneEntity
}
