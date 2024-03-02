package repository

type MemberUsecaseParam struct { /* generate by esc-gen-domain */
	dig.In

	Config *cfg.Config[configs.Config] `name:"config"`
	Store  *redis.ClusterClient        `name:"redis"`
}

type memberUsecase struct { /* generate by esc-gen-domain */
	config *cfg.Config[configs.Config]
	store  *redis.ClusterClient
}

func NewMemberUsecase(param MemberUsecaseParam) repository.MemberUsecase { /* generate by esc-gen-domain */
	return &memberUsecase{
		config: param.Config,
		store:  param.Store,
	}
}

func (repo *memberUsecase) Start(ctx context.Context, req *UpdatePhoneReq) (*UpdatePhoneResp, error) { /* generate by esc-gen-domain */
	// TODO: implement me
}
