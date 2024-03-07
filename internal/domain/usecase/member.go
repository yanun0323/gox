package usecase

import (
	"context"
	"gox/internal/domain/repository"
)

//go:generate esc-domain-gen -f member.go
type MemberUseCase interface {
	Start(ctx context.Context, req *UpdatePhoneReq) (*UpdatePhoneResp, error)
}
