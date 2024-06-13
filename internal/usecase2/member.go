package usecase

import (
	"context"

	"gox/internal/domain/usecase"
)

type Setting struct{}

type QueryReq struct{}

func Start() {}

type memberUseCase struct {
	s string
}

func NewMemberUseCase() usecase.MemberUseCase {
	return &memberUseCase{s: ""}
}

func (m *memberUseCase) Start(ctx context.Context, req *usecase.UpdatePhoneReq) (res *usecase.UpdatePhoneResp, err error) {
	return nil, nil
}

func (m *memberUseCase) End(ctx context.Context, req *usecase.UpdatePhoneReq) (*usecase.UpdatePhoneResp, error) {
	return nil, nil
}

func (m *memberUseCase) Exit(ctx context.Context) {
	println("Exit")
}
