package usecase

import (
	"context"

	"gox/internal/domain/usecase"
)

// go:generate inspector
type memberUsecase struct{}

func NewMemberUsecase() usecase.MemberUseCase {
	return &memberUsecase{}
}

func (use *memberUsecase) Start(ctx context.Context, req *usecase.UpdatePhoneReq) (*usecase.UpdatePhoneResp, error) {
	return nil, nil
}

func (use *memberUsecase) End(ctx context.Context, req *usecase.UpdatePhoneReq) (*usecase.UpdatePhoneResp, error) {
	return nil, nil
}

func (use *memberUsecase) Exit(ctx context.Context) {
}
