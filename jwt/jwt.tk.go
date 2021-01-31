package tkjwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/proto"
	tkpb "github.com/ikaiguang/srv_toolkit/api"
	tkjwtpb "github.com/ikaiguang/srv_toolkit/api/jwt"
	tke "github.com/ikaiguang/srv_toolkit/error"
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
func (s *jwtToken) Login(param *LoginParam, cacheInfo proto.Message) (token string, err error) {
	//claims := &jwt.StandardClaims{
	//	Audience:  "Audience", // aud 目标收件人(签发给谁)
	//	ExpiresAt: 0,          // exp 过期时间(有效期时间 exp)
	//	Id:        "Id",       // jti 编号
	//	IssuedAt:  0,          // iat 签发时间
	//	Issuer:    "Issuer",   // iss 签发者
	//	NotBefore: 0,          // nbf 生效时间(nbf 时间后生效)
	//	Subject:   "Subject",  // sub 主题
	//}
	if param.Claims.Audience == "" {
		err = tke.New(tke.JwtAudienceEmpty)
		return
	}
	switch param.UserInfo.UserStatus {
	case tkjwtpb.JwtActiveStatus_active_status_valid, tkjwtpb.JwtActiveStatus_active_status_temp, tkjwtpb.JwtActiveStatus_active_status_access:
	default:
		param.UserInfo.UserStatus = tkjwtpb.JwtActiveStatus_active_status_unknown
		err = s.activeStatusError(param.UserInfo.UserStatus)
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
		return nil
	}
}
