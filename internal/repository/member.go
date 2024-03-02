package repository

type MemberUsecaseParam struct {
	dig.In

	Config *cfg.Config[configs.Config] `name:"config"`
	Store  *redis.ClusterClient        `name:"redis"`
}

type memberUsecase struct {
	config *cfg.Config[configs.Config]
	store  *redis.ClusterClient
}

func NewMemberUsecase(param MemberUsecaseParam) repository.MemberUsecase {
	return &memberUsecase{
		config: param.Config,
		store:  param.Store,
	}
}

func (repo *memberUsecase) Start(ctx context.Context, req *UpdatePhoneReq) (*UpdatePhoneResp, error) {
	// TODO: implement me
}
