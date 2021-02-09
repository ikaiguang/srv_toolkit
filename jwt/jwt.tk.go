package tkjwt

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/proto"
	tkpb "github.com/ikaiguang/srv_toolkit/api"
	tkjwtpb "github.com/ikaiguang/srv_toolkit/api/jwt"
	tkapp "github.com/ikaiguang/srv_toolkit/app"
	tke "github.com/ikaiguang/srv_toolkit/error"
	tkredis "github.com/ikaiguang/srv_toolkit/redis"
	tkru "github.com/ikaiguang/srv_toolkit/redis/utils"
	"time"
)

// const
const (
	_defaultLoginType    = tkjwtpb.JwtLoginLimitType_login_type_unlimited
	_defaultPlatform     = tkpb.Platform_platform_mobile
	_defaultActiveStatus = tkjwtpb.JwtActiveStatus_active_status_valid

	_loginLockKeyPrefix  = "tk_jwt_login_lock:"
	_loginCacheKeyPrefix = "tk_jwt_token:"
	_loginCacheKeyUser   = "user"
)

// var
var (
	// logger
	logger tkapp.LoggerInterface = &tkapp.Log{}

	// token 过期时间
	_loginCacheExpire = time.Hour * 2
)

// SetLogger .
func SetLogger(handler tkapp.LoggerInterface) {
	logger = handler
}

// SetExpire .
func SetExpire(duration time.Duration) {
	_loginCacheExpire = duration
}

// GetExpire .
func GetExpire() time.Duration {
	return _loginCacheExpire
}

// jwtToken .
type jwtToken struct{}

// LoginParam .
type LoginParam struct {
	UserInfo  *tkjwtpb.JwtUserInfo
	Claims    *jwt.StandardClaims // claims.Audience 用于redis缓存(必填)
	Platform  tkpb.Platform
	LimitType tkjwtpb.JwtLoginLimitType
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

	// 验证参数
	err = s.validateParam(param)
	if err != nil {
		return
	}

	// 避免同时登录
	lock, err := s.GetLock(ctx, param.Claims.Audience)
	if err != nil {
		return
	}
	defer func() {
		_, _ = lock.Unlock()
	}()

	// 缓存
	allCache, err := s.GetAllCache(ctx, param.Claims.Audience)
	if err != nil {
		return
	}

	// 没有缓存，直接登录
	if !allCache.HasCache {
		return s.login(ctx, param, allCache)
	}

	// 有缓存，检查限制
	err = s.CheckLogin(ctx, param, allCache)
	if err != nil {
		return
	}
	return s.login(ctx, param, allCache)
}

// CheckLogin .
func (s *jwtToken) CheckLogin(ctx context.Context, param *LoginParam, allCache *JwtCache) (err error) {
	// 删除过期的token
	nowUnix := time.Now().Unix()
	var delFields []string

	// 检查缓存
	var platformM = make(map[tkpb.Platform][]*tkjwtpb.JwtAuthInfo)
	for key := range allCache.Tokens {
		// 删除过期的token
		if allCache.Tokens[key].Et <= nowUnix {
			delFields = append(delFields, key)
			delete(allCache.Tokens, key)
			continue
		}
		platformM[allCache.Tokens[key].Platform] = append(platformM[allCache.Tokens[key].Platform], allCache.Tokens[key])
	}

	// 检查限制
	switch param.LimitType {
	case tkjwtpb.JwtLoginLimitType_login_type_only_one:
		// 同一账户仅允许登录一次(验证码...可强制登录)
		if len(platformM) > 0 {
			err = tke.Newf(tke.JwtSigned, err)
			return err
		}
	case tkjwtpb.JwtLoginLimitType_login_type_platform_one:
		// 同一账户每个平台都可登录一次(验证码...可强制登录)
		if _, ok := platformM[param.Platform]; ok {
			err = tke.Newf(tke.JwtSigned, err)
			return err
		}
	default:
		// 未知 || 无限制 || 其他
	}

	// 删除过期的token
	if len(delFields) > 0 {
		// 删除缓存
		_, err = tkru.HDel(ctx, allCache.Key, delFields...)
		if err != nil {
			err = tke.Newf(tke.Err, err)
			return err
		}
	}
	return
}

// login .
func (s *jwtToken) login(ctx context.Context, param *LoginParam, allCache *JwtCache) (token string, err error) {
	token, err = s.GenToken(param.Claims, param.UserInfo.TokenSecret)
	if err != nil {
		return
	}

	// 缓存
	allCache.User = param.UserInfo
	if allCache.Tokens == nil {
		allCache.Tokens = make(map[string]*tkjwtpb.JwtAuthInfo)
	}
	allCache.Tokens[param.Claims.Id] = &tkjwtpb.JwtAuthInfo{
		TokenId:  param.Claims.Id,
		Platform: param.Platform,
		Lt:       param.LimitType,
		Et:       param.Claims.ExpiresAt,
		Ct:       time.Now().Unix(),
	}
	err = s.SaveCache(ctx, allCache)
	if err != nil {
		return "", err
	}
	return
}

// GenToken 生产token
func (s *jwtToken) GenToken(claims jwt.Claims, secret string) (tokenStr string, err error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = tokenClaims.SignedString([]byte(secret))
	if err != nil {
		err = tke.Newf(tke.Err, err)
		return
	}
	return
}

// validateParam .
func (s *jwtToken) validateParam(param *LoginParam) (err error) {
	// cache key
	if param.Claims.Audience == "" {
		err = tke.New(tke.JwtAudienceEmpty)
		return
	}
	if param.Claims.Id == "" {
		err = tke.New(tke.JwtIdEmpty)
		return
	}

	// 用户状态
	switch param.UserInfo.UserStatus {
	//case tkjwtpb.JwtActiveStatus_active_status_valid, tkjwtpb.JwtActiveStatus_active_status_temp, tkjwtpb.JwtActiveStatus_active_status_access:
	//	// 有效的
	case tkjwtpb.JwtActiveStatus_active_status_deny, tkjwtpb.JwtActiveStatus_active_status_deleted, tkjwtpb.JwtActiveStatus_active_status_invalid:
		// 无效的
		err = s.activeStatusError(param.UserInfo.UserStatus)
		return
	default:
		// 默认未有效用户
		//param.UserInfo.UserStatus = tkjwtpb.JwtActiveStatus_active_status_valid
	}

	// 平台
	//switch param.Platform {
	//case tkpb.Platform_platform_computer, tkpb.Platform_platform_mobile,
	//	tkpb.Platform_platform_desktop, tkpb.Platform_platform_android, tkpb.Platform_platform_iphone:
	//	// 有效平台
	//default:
	//	//param.Platform = tkpb.Platform_platform_unknown
	//}

	// 登录限制
	//switch param.LimitType {
	//case tkjwtpb.JwtLoginLimitType_login_type_only_one, tkjwtpb.JwtLoginLimitType_login_type_platform_one:
	//default:
	//	param.LimitType = tkjwtpb.JwtLoginLimitType_login_type_unlimited
	//}
	return
}

// =====================================================================================================================

// JwtCache .
type JwtCache struct {
	Key      string
	HasCache bool
	User     *tkjwtpb.JwtUserInfo
	Tokens   map[string]*tkjwtpb.JwtAuthInfo
}

// CacheKey .
// @param @jwtAudience cache key
func (s *jwtToken) CacheKey(jwtAudience string) string {
	return tkru.Key(_loginCacheKeyPrefix + jwtAudience)
}

// DelCache .
func (s *jwtToken) DelCache(ctx context.Context, jwtAudience string) (err error) {
	_, err = tkru.Del(ctx, s.CacheKey(jwtAudience))
	if err != nil {
		err = tke.Newf(tke.Redis, err)
		return
	}
	return
}

// SaveCache .
func (s *jwtToken) SaveCache(ctx context.Context, allCache *JwtCache) (err error) {
	var cacheM = make(map[string]interface{})
	// user
	buf, err := proto.Marshal(allCache.User)
	if err != nil {
		err = tke.Newf(tke.Err, err)
		return
	}
	cacheM[_loginCacheKeyUser] = buf

	// auth
	for key := range allCache.Tokens {
		buf, err = proto.Marshal(allCache.User)
		if err != nil {
			err = tke.Newf(tke.Err, err)
			return err
		}
		cacheM[key] = buf
	}

	// save
	_, err = tkru.HMSet(ctx, allCache.Key, cacheM)
	if err != nil {
		err = tke.Newf(tke.Redis, err)
		return
	}

	// expire
	_, err = tkru.Expire(ctx, allCache.Key, GetExpire())
	if err != nil {
		err = tke.Newf(tke.Redis, err)
		return
	}
	return
}

// GetAllCache .
func (s *jwtToken) GetAllCache(ctx context.Context, jwtAudience string) (res *JwtCache, err error) {
	res = &JwtCache{Key: s.CacheKey(jwtAudience)}
	defer func() {
		if err != nil {
			if delErr := s.DelCache(ctx, jwtAudience); err != nil {
				logger.Error(delErr)
			}
		}
	}()

	// get
	cacheM, err := tkru.BytesMap(tkru.HGetAll(ctx, res.Key))
	if err != nil {
		if tkru.IsRedisNil(err) {
			err = nil
		} else {
			err = tke.Newf(tke.Redis, err)
		}
		return
	}

	// cache
	res.HasCache = len(cacheM) > 0
	res.User = &tkjwtpb.JwtUserInfo{}
	res.Tokens = make(map[string]*tkjwtpb.JwtAuthInfo)
	for key := range cacheM {
		// 用户
		if key == _loginCacheKeyUser {
			err = proto.Unmarshal(cacheM[key], res.User)
			if err != nil {
				err = tke.Newf(tke.InvalidRedisData, err)
				return res, err
			}
			continue
		}
		// token
		authInfo := &tkjwtpb.JwtAuthInfo{}
		err = proto.Unmarshal(cacheM[key], authInfo)
		if err != nil {
			err = tke.Newf(tke.InvalidRedisData, err)
			return res, err
		}
		res.Tokens[key] = authInfo
	}
	return
}

// LockKey .
// @param @jwtAudience lock key
func (s *jwtToken) LockKey(jwtAudience string) string {
	return tkru.Key(_loginLockKeyPrefix + jwtAudience)
}

// GetLock .
func (s *jwtToken) GetLock(ctx context.Context, jwtAudience string) (lock *tkredis.DLock, err error) {
	lock, isLockFail, err := tkredis.NewLock(ctx, s.LockKey(jwtAudience))
	if err != nil {
		err = tke.Newf(tke.Redis, err)
		return
	}
	if isLockFail {
		err = tke.New(tke.TooManyRequests)
		return
	}
	return
}

// =====================================================================================================================

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
