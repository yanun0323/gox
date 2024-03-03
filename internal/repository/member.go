package repository

type MemberRepositoryParam struct {
	dig.In

	Config *cfg.Config[configs.Config] `name:"config"`
	Store  *redis.ClusterClient        `name:"redis"`
}

type memberRepository struct {
	config *cfg.Config[configs.Config]
	store  *redis.ClusterClient
}

func NewMemberRepository(param MemberRepositoryParam) repository.MemberRepository {
	return &memberRepository{
		config: param.Config,
		store:  param.Store,
	}
}

func (repo *memberRepository) Start(ctx context.Context, req *UpdatePhoneReq) (*UpdatePhoneResp, error) {
	// TODO: implement me
}
