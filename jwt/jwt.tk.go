package tkjwt

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/proto"
	tkpb "github.com/ikaiguang/srv_toolkit/api"
	tkjwtpb "github.com/ikaiguang/srv_toolkit/api/jwt"
	tke "github.com/ikaiguang/srv_toolkit/error"
	tkredisutils "github.com/ikaiguang/srv_toolkit/redis/utils"
	"github.com/pkg/errors"
)

// const
const (
	_defaultLoginType    = tkjwtpb.JwtLoginType_login_type_multiple
	_defaultPlatform     = tkpb.Platform_platform_mobile
	_defaultActiveStatus = tkjwtpb.JwtActiveStatus_active_status_valid
)

// jwtToken .
type jwtToken struct{}

// LoginParam .
type LoginParam struct {
	UserInfo  *tkjwtpb.JwtUserInfo
	Claims    *jwt.StandardClaims // claims.Audience 用于redis缓存(必填)
	Platform  tkpb.Platform
	LoginType tkjwtpb.JwtLoginType
}

// Login .
func (s *jwtToken) Login(ctx context.Context, param *LoginParam) (token string, err error) {
	//claims := &jwt.StandardClaims{
	//	Audience:  "Audience", // aud 目标收件人(签发给谁)
	//	ExpiresAt: 0,          // exp 过期时间(有效期时间 exp)
	//	Id:        "Id",       // jti 编号
	//	IssuedAt:  0,          // iat 签发时间
	//	Issuer:    "Issuer",   // iss 签发者
	//	NotBefore: 0,          // nbf 生效时间(nbf 时间后生效)
	//	Subject:   "Subject",  // sub 主题
	//}

	// cache key
	if param.Claims.Audience == "" {
		err = tke.New(tke.JwtAudienceEmpty)
		return
	}

	// can login ?
	switch param.UserInfo.UserStatus {
	case tkjwtpb.JwtActiveStatus_active_status_valid, tkjwtpb.JwtActiveStatus_active_status_temp, tkjwtpb.JwtActiveStatus_active_status_access:
	default:
		param.UserInfo.UserStatus = tkjwtpb.JwtActiveStatus_active_status_unknown
		err = s.activeStatusError(param.UserInfo.UserStatus)
	}
	return
}

// CacheKey .
// @param @jwtAudience cache key
func (s *jwtToken) CacheKey(jwtAudience string) string {
	return tkredisutils.Key("jwt_token:" + jwtAudience)
}

// GetCache .
func (s *jwtToken) GetCache(ctx context.Context, claims *jwt.StandardClaims) (ac *tkjwtpb.JwtAuthCache, hasCache bool, err error) {
	cacheBytes, err := tkredisutils.Bytes(tkredisutils.Get(ctx, s.CacheKey(claims.Audience)))
	if err != nil {
		if tkredisutils.IsRedisNil(err) {
			err = nil
		} else {
			err = errors.WithStack(err)
		}
		return
	}

	// cache
	hasCache = true
	ac = &tkjwtpb.JwtAuthCache{}
	err = proto.Unmarshal(cacheBytes, ac)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// activeStatusError .
func (s *jwtToken) activeStatusError(status tkjwtpb.JwtActiveStatus) error {
	switch status {
	case tkjwtpb.JwtActiveStatus_active_status_unknown:
		return tke.New(tke.JwtAckUnknown)
	case tkjwtpb.JwtActiveStatus_active_status_deny:
		return tke.New(tke.JwtAckDeny)
	case tkjwtpb.JwtActiveStatus_active_status_deleted:
		return tke.New(tke.JwtAckDeleted)
	case tkjwtpb.JwtActiveStatus_active_status_invalid:
		return tke.New(tke.JwtAckInvalid)
	default:
		return tke.New(tke.JwtAckUnknown)
	}
}
