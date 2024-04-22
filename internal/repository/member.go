package repository

type MemberRepositoryParam struct {
	dig.In

	DB *gorm.DB `name:"dbM"`
}

type memberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(param MemberRepositoryParam) repository.MemberRepository {
	return &memberRepository{
		db: param.DB,
	}
}

func (repo *memberRepository) Start(ctx context.Context, req *UpdatePhoneReq) (*UpdatePhoneResp, error) {
	// TODO: implement me
}
