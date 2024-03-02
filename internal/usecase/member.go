package usecase

import (
	"context"
	"errors"
	"fmt"
	"member/internal/domain"
	"member/internal/domain/entity"
	"member/internal/domain/repository"
	"member/internal/domain/usecase"
	"member/internal/libs/pager"
	"member/internal/utils/encrypt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type MemberUseCaseParam struct {
	dig.In

	MemberRepo      repository.MemberRepository
	UserInfoRepo    repository.UserInfoRepository
	LoginLogRepo    repository.LoginLogRepository
	SecurityRepo    repository.SecurityRepository
	RegisterLogRepo repository.RegisterLogRepository
	JobRepo         repository.JobRepository

	Redis *redis.ClusterClient `name:"redis"`

	DB *gorm.DB `name:"dbM"`
}

type memberUseCase struct {
	memberRepo      repository.MemberRepository
	userInfoRepo    repository.UserInfoRepository
	loginLogRepo    repository.LoginLogRepository
	securityRepo    repository.SecurityRepository
	registerLogRepo repository.RegisterLogRepository
	jobRepo         repository.JobRepository

	redis *redis.ClusterClient

	db *gorm.DB
}

func NewMemberUseCase(param MemberUseCaseParam) usecase.MemberUseCase {
	return &memberUseCase{
		memberRepo:      param.MemberRepo,
		userInfoRepo:    param.UserInfoRepo,
		loginLogRepo:    param.LoginLogRepo,
		securityRepo:    param.SecurityRepo,
		registerLogRepo: param.RegisterLogRepo,
		jobRepo:         param.JobRepo,

		redis: param.Redis,

		db: param.DB,
	}
}

const (
	defaultNicknamePrefix        = "ESC-User"
	defaultNicknameRandomLength  = 5
	defaultGenerateNicknameRetry = 5
)

func (use *memberUseCase) CreateGeneralMember(
	ctx context.Context, param *usecase.CreateGeneralMemberParam,
) (*usecase.CreateGeneralMemberResp, error) {
	if strings.ContainsAny(param.Password, " ") {
		return nil, usecase.PasswordFormatErr{Msg: "wrong format"}
	}

	userEntity, err := use.memberRepo.GetByEmail(ctx, param.Email)
	if err != nil && !errors.Is(err, repository.QueryRecordNotFoundError) {
		return nil, usecase.InternalErr{Err: fmt.Errorf("memberRepo.GetByEmail error: %w", err)}
	}

	if userEntity != nil {
		return nil, usecase.EmailExistsErr{Email: param.Email}
	}

	encryptPwd, err := encrypt.BcryptEncode(param.Password)
	if err != nil {
		return nil, usecase.InternalErr{Err: fmt.Errorf("utils.BcryptEncode err: %w", err)}
	}

	newUser := &entity.User{
		Brand:    param.Brand,
		Email:    param.Email,
		Password: encryptPwd,
		UserType: domain.UserTypeGeneral,
		Status:   domain.UserStatusEnable,
	}

	// 先Commit搶member id(做uid)避免tx中間過久導致uid重覆
	if _, err := use.memberRepo.Create(ctx, newUser); err != nil {
		return nil, usecase.InternalErr{Err: fmt.Errorf("memberRepo.Create error: %w", err)}
	}

	tx := use.db.Begin()
	defer tx.Rollback()

	userInfoRepo := use.userInfoRepo.New(tx)
	securityRepo := use.securityRepo.New(tx)
	registerLogRepo := use.registerLogRepo.New(tx)

	nickname := fmt.Sprintf("%s-%s", defaultNicknamePrefix, newUser.UID)

	for i := 0; i < defaultGenerateNicknameRetry; i++ {
		randomNickname, err := encrypt.GenerateSalt(defaultNicknameRandomLength)
		if err != nil {
			return nil, usecase.InternalErr{Err: fmt.Errorf("encrypt.GenerateSalt error: %w", err)}
		}

		tempNickname := fmt.Sprintf("%s-%s", defaultNicknamePrefix, randomNickname)
		_, err = userInfoRepo.GetUserInfoByNickname(ctx, tempNickname)
		if errors.Is(err, repository.QueryRecordNotFoundError) {
			nickname = tempNickname
			break
		}
	}

	newUserInfo := &entity.UserInfo{
		UID:          newUser.UID,
		NickName:     nickname,
		Language:     param.Language,
		FiatCurrency: domain.DefaultFiatCurrency,
	}

	if _, err := userInfoRepo.Create(ctx, newUserInfo); err != nil {
		return nil, usecase.InternalErr{Err: fmt.Errorf("userInfoRepo.Create error: %w", err)}
	}

	// 由於當前註冊只採用Email方式，代表User 已過email安全驗證
	_, err = securityRepo.UpsertUserSecurity(ctx, &entity.UserSecurity{
		Brand:        param.Brand,
		UID:          newUser.UID,
		SecurityType: domain.SecurityEmail.String(),
		Binding:      domain.BindingTrue.ToINT8(),
		Enable:       domain.EnableTrue.ToINT8(),
	})

	if err != nil {
		return nil, usecase.InternalErr{Err: fmt.Errorf("securityRepo.UpsertUserSecurity error: %w", err)}
	}

	_, err = registerLogRepo.Create(ctx, &entity.RegisterLog{
		UID:           newUser.UID,
		Brand:         param.Brand,
		IPAddress:     param.IPAddress,
		IPLocation:    param.IPLocation,
		Device:        param.Device,
		DeviceType:    param.DeviceType,
		DeviceBrand:   param.DeviceBrand,
		DeviceVersion: param.DeviceVersion,
		DeviceID:      param.DeviceID,
	})
	if err != nil {
		return nil, usecase.InternalErr{Err: fmt.Errorf("registerLogRepo.Create error: %w", err)}
	}

	tx.Commit()

	return &usecase.CreateGeneralMemberResp{
		UID:   newUser.UID,
		Email: newUser.Email,
	}, nil
}

func (use *memberUseCase) GetMemberByUID(ctx context.Context, uid string) (*usecase.GetMemberByUIDResp, error) {
	entityUser, err := use.memberRepo.GetByUID(ctx, uid)
	if err != nil {
		if errors.Is(err, repository.QueryRecordNotFoundError) {
			return nil, usecase.UIDNotFoundErr{UID: uid}
		}
		return nil, usecase.InternalErr{Err: fmt.Errorf("memberRepo.GetByUID error: %w", err)}
	}

	entityUserInfo, err := use.userInfoRepo.GetUserInfoByUID(ctx, uid)
	if err != nil {
		if errors.Is(err, repository.QueryRecordNotFoundError) {
			return nil, usecase.UIDNotFoundErr{UID: uid}
		}
		return nil, usecase.InternalErr{Err: fmt.Errorf("userInfoRepo.GetUserInfoByUID error: %w", err)}
	}

	st := domain.SecurityPhishing
	entityUserSecurity, err := use.securityRepo.GetUserSecurity(ctx, repository.GetUserSecurityOption{
		UID:          &uid,
		SecurityType: &st,
	})
	if err != nil && !errors.Is(err, repository.QueryRecordNotFoundError) {
		return nil, usecase.InternalErr{Err: fmt.Errorf("securityRepo.GetUserSecurityByUID error: %w", err)}
	}

	phishingCode := ""
	if entityUserSecurity != nil {
		phishingCode = entityUserSecurity.SecurityCode
	}

	return &usecase.GetMemberByUIDResp{
		UID:                entityUser.UID,
		Email:              entityUser.Email,
		Nickname:           entityUserInfo.NickName,
		UserType:           entityUser.UserType,
		Status:             entityUser.Status,
		RegisterTime:       entityUser.CreateTime,
		Phone:              entityUserInfo.Phone,
		AreaCode:           entityUserInfo.AreaCode,
		PhishingCode:       phishingCode,
		NicknameUpdateTime: entityUserInfo.NickNameUpdateTime,
		FiatCurrency:       entityUserInfo.FiatCurrency,
		Language:           entityUserInfo.Language,
	}, nil
}

func (use *memberUseCase) GetMemberByEmail(ctx context.Context, email string) (*usecase.GetMemberByEmailResp, error) {
	entityUser, err := use.memberRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.QueryRecordNotFoundError) {
			return nil, usecase.EmailNotFoundErr{Email: email}
		}
		return nil, usecase.InternalErr{Err: fmt.Errorf("memberRepo.GetByEmail error: %w", err)}
	}

	return &usecase.GetMemberByEmailResp{
		UID:   entityUser.UID,
		Email: entityUser.Email,
	}, nil
}

func (use *memberUseCase) ListMembersByQuery(
	ctx context.Context, param *usecase.ListMembersByQueryParam,
) (*usecase.ListMembersByQueryResp, error) {
	repoParam := &repository.ListMembersParam{
		UID:               param.UID,
		Email:             param.Email,
		Phone:             param.Phone,
		SecurityLevel:     param.SecurityLevel,
		Region:            param.Region,
		FiatCoinId:        param.FiatCoinId,
		UserType:          param.UserType,
		Status:            param.Status,
		BeginRegisterTime: param.BeginRegisterTime,
		EndRegisterTime:   param.EndRegisterTime,
		Name:              param.Name,
		IdentityType:      param.IdentityType,
		RegisterDevice:    param.RegisterDevice,
	}
	result, count, err := use.memberRepo.ListMembers(ctx, repoParam, pager.Condition{
		PageIndex: param.PageIndex,
		PageSize:  param.PageSize,
	})

	if err != nil {
		return nil, usecase.InternalErr{Err: fmt.Errorf("memberRepo.ListMembers error: %w", err)}
	}

	if len(result) == 0 {
		return nil, nil
	}

	resp := &usecase.ListMembersByQueryResp{
		List: make([]usecase.MemberByQueryResp, len(result)),
		Page: pager.Response{
			Total: int(count),
			Index: param.PageIndex,
			Size:  param.PageSize,
		},
	}

	for index, userInfo := range result {
		resp.List[index] = usecase.MemberByQueryResp{
			UID:           userInfo.UID,
			Email:         userInfo.Email,
			Name:          userInfo.Name,
			Nickname:      userInfo.NickName,
			UserType:      userInfo.UserType,
			Status:        userInfo.Status,
			RegisterTime:  userInfo.CreateTime,
			Phone:         userInfo.Phone,
			AreaCode:      userInfo.AreaCode,
			IPAddress:     userInfo.IPAddress,
			IPLocation:    userInfo.IPLocation,
			Device:        userInfo.Device,
			DeviceType:    userInfo.DeviceType,
			DeviceBrand:   userInfo.DeviceBrand,
			DeviceVersion: userInfo.DeviceVersion,
			DeviceID:      userInfo.DeviceID,
		}
	}
	return resp, nil
}

func (use *memberUseCase) ResetPasswordByEmail(ctx context.Context, email, password string) error {
	entityUser, err := use.memberRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.QueryRecordNotFoundError) {
			return usecase.EmailNotFoundErr{Email: email}
		}
		return usecase.InternalErr{Err: fmt.Errorf("memberRepo.GetByEmail error: %w", err)}
	}

	encryptPwd, err := encrypt.BcryptEncode(password)
	if err != nil {
		return usecase.InternalErr{Err: fmt.Errorf("utils.BcryptEncode err: %w", err)}
	}

	if err := use.memberRepo.UpdatePasswordByUID(ctx, entityUser.UID, encryptPwd); err != nil {
		return usecase.InternalErr{Err: fmt.Errorf("memberRepo.UpdatePasswordByUID error: %w", err)}
	}

	return nil
}

func (use *memberUseCase) ResetPasswordByUID(ctx context.Context, uid, password string) error {
	_, err := use.memberRepo.GetByUID(ctx, uid)
	if err != nil {
		if errors.Is(err, repository.QueryRecordNotFoundError) {
			return usecase.UIDNotFoundErr{UID: uid}
		}
		return usecase.InternalErr{Err: fmt.Errorf("memberRepo.GetByUID error: %w", err)}
	}

	encryptPwd, err := encrypt.BcryptEncode(password)
	if err != nil {
		return usecase.InternalErr{Err: fmt.Errorf("utils.BcryptEncode err: %w", err)}
	}

	if err := use.memberRepo.UpdatePasswordByUID(ctx, uid, encryptPwd); err != nil {
		return usecase.InternalErr{Err: fmt.Errorf("memberRepo.UpdatePasswordByUID error: %w", err)}
	}

	return nil
}

func (use *memberUseCase) UpdateMemberInfo(ctx context.Context, param *usecase.UpdateMemberInfoParam) error {
	if param.Email != nil {
		if err := use.updateEmail(ctx, param.UID, *param.Email); err != nil {
			return err
		}

		return nil
	}

	entityUserInfo, err := use.userInfoRepo.GetUserInfoByUID(ctx, param.UID)
	if err != nil {
		if errors.Is(err, repository.QueryRecordNotFoundError) {
			return usecase.UIDNotFoundErr{UID: param.UID}
		}
		return usecase.InternalErr{Err: fmt.Errorf("userInfoRepo.GetUserInfoByUID error: %w", err)}
	}

	// 如果不是更改NickName則不需要
	entityUserInfo.NickNameUpdateTime = 0

	if param.NickName != nil {
		result, err := use.userInfoRepo.GetUserInfoByNickname(ctx, *param.NickName)
		if err != nil && !errors.Is(err, repository.QueryRecordNotFoundError) {
			return usecase.InternalErr{Err: fmt.Errorf("userInfoRepo.GetUserInfoByNickname error: %w", err)}
		}

		if result != nil {
			return usecase.NicknameExistsError{Name: *param.NickName}
		}

		entityUserInfo.NickName = *param.NickName
		entityUserInfo.NickNameUpdateTime = time.Now().Unix()
	}

	if param.Phone != nil && param.AreaCode != nil {
		if *param.Phone != "" {
			result, err := use.userInfoRepo.GetUserInfoByPhone(ctx, *param.Phone)
			if err != nil && !errors.Is(err, repository.QueryRecordNotFoundError) {
				return usecase.InternalErr{Err: fmt.Errorf("userInfoRepo.GetUserInfoByPhone error: %w", err)}
			}

			if result != nil {
				return usecase.PhoneExistsError{Phone: *param.Phone}
			}
		}

		entityUserInfo.Phone = *param.Phone
		entityUserInfo.AreaCode = *param.AreaCode
	}

	if err := use.userInfoRepo.UpdateInfo(ctx, entityUserInfo); err != nil {
		return usecase.InternalErr{Err: fmt.Errorf("userInfoRepo.UpdateInfo error: %w", err)}
	}

	return nil
}

// updateEmail is responsible for update user,s email
func (use *memberUseCase) updateEmail(ctx context.Context, uid, email string) error {
	entityUser, err := use.memberRepo.GetByUID(ctx, uid)
	if err != nil {
		if errors.Is(err, repository.QueryRecordNotFoundError) {
			return usecase.UIDNotFoundErr{UID: uid}
		}
		return usecase.InternalErr{Err: fmt.Errorf("memberRepo.GetByUID error: %w", err)}
	}

	// 判斷email 是否重復使用
	userByEmail, err := use.memberRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, repository.QueryRecordNotFoundError) {
		return usecase.InternalErr{Err: fmt.Errorf("memberRepo.GetByEmail error: %w", err)}
	}

	if userByEmail != nil {
		return usecase.EmailExistsErr{Email: email}
	}

	entityUser.Email = email
	if _, err := use.memberRepo.Update(ctx, entityUser); err != nil {
		return fmt.Errorf("memberRepo.Update error: %w", err)
	}

	return nil
}

// SetTemporaryFreeze 暫時凍結用戶
func (use *memberUseCase) SetTemporaryFreeze(ctx context.Context, uid string, expireSec int) error {
	user, err := use.memberRepo.GetByUID(ctx, uid)
	if err != nil {
		if errors.Is(err, repository.QueryRecordNotFoundError) {
			return usecase.UIDNotFoundErr{UID: uid}
		}
		return usecase.InternalErr{Err: fmt.Errorf("memberRepo.GetByUID error: %w", err)}
	}

	user.Status = domain.UserStatusTemporaryFreeze
	if _, err := use.memberRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("memberRepo.Update error: %w", err)
	}

	err = use.jobRepo.Create(ctx, &entity.Job{
		Type:      domain.TemporaryFreezeJob,
		UID:       uid,
		StartTime: time.Now().Add(time.Duration(expireSec) * time.Second).Unix(),
	})
	if err != nil {
		return usecase.InternalErr{fmt.Errorf("jobRepo.Create error: %w", err)}
	}

	return nil
}

func (use *memberUseCase) SetFiatCurrency(ctx context.Context, uid string, fiatCurrency string) error {
	err := use.userInfoRepo.UpdateFiatCurrency(ctx, uid, fiatCurrency)
	if err != nil {
		return fmt.Errorf("userInfoRepo.UpdateFiatCurrency err: %w", err)
	}
	return nil
}

func (use *memberUseCase) SetLanguage(ctx context.Context, uid string, language string) error {
	err := use.userInfoRepo.UpdateLanguage(ctx, uid, language)
	if err != nil {
		return fmt.Errorf("userInfoRepo.UpdateLanguage err: %w", err)
	}

	// Save to Redis
	if err := use.redis.HSet(ctx, domain.UserConfigLanguageHashKey.ToString(), uid, language).Err(); err != nil {
		return fmt.Errorf("redis.HSet err: %w", err)
	}

	return nil
}

func (usecase *memberUseCase) Start(ctx context.Context, req *UpdatePhoneReq) (*UpdatePhoneResp, error) {
	// TODO: implement me
}
