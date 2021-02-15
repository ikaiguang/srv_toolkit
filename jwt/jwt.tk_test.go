package tkjwt

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	tkpb "github.com/ikaiguang/srv_toolkit/api"
	tkjwtpb "github.com/ikaiguang/srv_toolkit/api/jwt"
	tkredis "github.com/ikaiguang/srv_toolkit/redis"
	tkru "github.com/ikaiguang/srv_toolkit/redis/utils"
	"github.com/ikaiguang/srv_toolkit/testdata"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	initOnce sync.Once
)

func initSetup() {
	initOnce.Do(func() {
		testdata.Setup()

		tkredis.Setup("redis.toml", "Client")
		tkru.Init(tkredis.Redis())
	})
}

func testLoginParam(tokenExpire time.Duration) (loginParam *LoginParam) {
	loginParam = &LoginParam{
		UserInfo: &tkjwtpb.JwtUserInfo{
			Id: 1, Uuid: "UUID",
		},
		// // 必填参数： Id, Audience
		Claims: &jwt.StandardClaims{
			Id:        uuid.New().String(),
			Audience:  "Audience",
			ExpiresAt: time.Now().Add(tokenExpire).Unix(),
		},
		Platform:  tkpb.Platform_platform_android,
		LimitType: tkjwtpb.JwtLoginLimitType_login_type_unlimited,
	}
	return
}

func TestJwtToken_Login(t *testing.T) {
	initSetup()

	var (
		token string
		err   error
	)

	// param
	// token 有效期 10s，所以：执行一次完整的测试需要间隔 10s
	tokenExpire := 10 * time.Second
	ctx := context.Background()
	loginParam := testLoginParam(tokenExpire)

	data := []struct {
		canLogin bool
		platform tkpb.Platform
		limit    tkjwtpb.JwtLoginLimitType
	}{
		{true, tkpb.Platform_platform_android, tkjwtpb.JwtLoginLimitType_login_type_unlimited},
		{true, tkpb.Platform_platform_android, tkjwtpb.JwtLoginLimitType_login_type_unlimited},
		{false, tkpb.Platform_platform_iphone, tkjwtpb.JwtLoginLimitType_login_type_only_one},
		{false, tkpb.Platform_platform_iphone, tkjwtpb.JwtLoginLimitType_login_type_only_one},
		{true, tkpb.Platform_platform_iphone, tkjwtpb.JwtLoginLimitType_login_type_platform_one},
		{false, tkpb.Platform_platform_iphone, tkjwtpb.JwtLoginLimitType_login_type_platform_one},
	}

	// uuid
	var uidS []string
	for _, val := range data {
		loginParam.Claims.Id = uuid.New().String()
		uid := loginParam.Claims.Id
		loginParam.Platform = val.platform
		loginParam.LimitType = val.limit
		token, err = Handler.Login(ctx, loginParam)
		msg := "param : " + strconv.Itoa(int(val.platform)) + "-" + strconv.Itoa(int(val.limit))
		if val.canLogin {
			uidS = append(uidS, uid)
			assert.Nil(t, err, msg)
			assert.NotEmpty(t, token, msg)
		} else {
			assert.NotNil(t, err, msg)
			assert.Empty(t, token, msg)
		}
	}

	// cache
	allCache, err := Handler.GetAllCache(ctx, loginParam.Claims.Audience)
	assert.Nil(t, err)
	assert.True(t, allCache.HasCache)
	assert.Equal(t, loginParam.UserInfo.Id, allCache.User.Id)
	assert.Equal(t, loginParam.UserInfo.Uuid, allCache.User.Uuid)
	for i := range uidS {
		assert.Contains(t, allCache.Tokens, uidS[i])
	}
}

func TestJwtToken_IsValid(t *testing.T) {
	initSetup()

	var (
		token string
		err   error
	)

	// param
	tokenExpire := 10 * time.Second
	ctx := context.Background()
	loginParam := testLoginParam(tokenExpire)

	token, err = Handler.Login(ctx, loginParam)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	validRes, err := Handler.IsValid(ctx, token)
	assert.Nil(t, err)
	assert.NotNil(t, validRes)
	assert.Equal(t, loginParam.UserInfo.Id, validRes.UserInfo.Id)
	assert.Equal(t, loginParam.UserInfo.Uuid, validRes.UserInfo.Uuid)
	assert.Equal(t, loginParam.Claims.Id, validRes.Claims.Id)
}

func TestJwtToken_Refresh(t *testing.T) {
	initSetup()

	var (
		token string
		err   error
	)

	// param
	tokenExpire := 10 * time.Second
	ctx := context.Background()
	loginParam := testLoginParam(tokenExpire)

	token, err = Handler.Login(ctx, loginParam)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	validRes, err := Handler.IsValid(ctx, token)
	assert.Nil(t, err)
	assert.NotNil(t, validRes)
	assert.Equal(t, loginParam.UserInfo.Id, validRes.UserInfo.Id)
	assert.Equal(t, loginParam.UserInfo.Uuid, validRes.UserInfo.Uuid)
	assert.Equal(t, loginParam.Claims.Id, validRes.Claims.Id)

	// reset expire time
	validRes.Claims.ExpiresAt = time.Now().Add(tokenExpire + time.Second).Unix()
	newToken, err := Handler.Refresh(ctx, validRes)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	assert.NotEqual(t, token, newToken)

	validRes, err = Handler.IsValid(ctx, newToken)
	assert.Nil(t, err)
	assert.NotNil(t, validRes)
	assert.Equal(t, loginParam.UserInfo.Id, validRes.UserInfo.Id)
	assert.Equal(t, loginParam.UserInfo.Uuid, validRes.UserInfo.Uuid)
	assert.Equal(t, loginParam.Claims.Id, validRes.Claims.Id)
}

func TestJwtToken_Logout(t *testing.T) {
	initSetup()

	var (
		token string
		err   error
	)

	// param
	tokenExpire := 10 * time.Second
	ctx := context.Background()
	loginParam := testLoginParam(tokenExpire)

	token, err = Handler.Login(ctx, loginParam)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	err = Handler.Logout(ctx, loginParam.Claims)
	assert.Nil(t, err)

	tokenCache, err := Handler.GetTokenCache(ctx, loginParam.Claims)
	assert.Nil(t, err)
	assert.False(t, tokenCache.HasCache)
}
