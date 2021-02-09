package tke

// Code .
type Code int32

// Code .
func (c Code) Code() int32 {
	return int32(c)
}

// code
const (
	Init              Code = 0     // init
	Success           Code = 200   // 成功
	Unknown           Code = 10000 // 未知错误
	Err               Code = 10001 // 系统错误
	Panic             Code = 10002 // 系统错误:Panic
	Db                Code = 10003 // 数据库错误
	Redis             Code = 10004 // 缓存错误
	Forbidden         Code = 10005 // 无权限操作资源，访问被拒绝
	InvalidParameters Code = 10006 // 参数错误
	BadRequest        Code = 10007 // 无效的请求
	NoneToken         Code = 10008 // 请携带登录令牌
	InvalidToken      Code = 10009 // 无效的登录令牌
	TokenInvalid      Code = 10010 // 登录令牌已失效
	NotFound          Code = 10011 // 请求资源不存在
	TooManyRequests   Code = 10012 // 请求次数过多，请稍后再试
	PhoneNumberErr    Code = 10013 // 手机号码不正确
	SmsSendFail       Code = 10014 // 发送短信失败
	SmsCodeExpire     Code = 10015 // 验证码已失效，请重新获取
	SmsCodeLimit      Code = 10016 // 短信验证码获取已达限制
	SmsCodeEmpty      Code = 10017 // 请输入验证码
	SmsCodeErr        Code = 10018 // 验证码不正确
	JwtAckUnknown     Code = 10019 // 登录参数有误
	JwtAckDeny        Code = 10020 // 非法登录
	JwtAckDeleted     Code = 10021 // 无效用户
	JwtAckInvalid     Code = 10022 // 请先激活账户
	JwtParamErr       Code = 10023 // JWT参数有误
	JwtAudienceEmpty  Code = 10024 // 请填写签发的目标收件人JWT_Audience
	JwtIdEmpty        Code = 10025 // 请填写签发的编号JWT_Id
	JwtSigned         Code = 10026 // 您的账户已在其他设备登录
	InvalidData       Code = 10027 // 无效的数据，请重试
	InvalidDBData     Code = 10028 // 无效的数据库数据，请重试
	InvalidRedisData  Code = 10029 // 无效的缓存数据，请重试
)

/*
// All common ecode
var (
	OK = add(0) // 正确

	NotModified        = add(-304) // 木有改动
	TemporaryRedirect  = add(-307) // 撞车跳转
	RequestErr         = add(-400) // 请求错误
	Unauthorized       = add(-401) // 未认证
	AccessDenied       = add(-403) // 访问权限不足
	NothingFound       = add(-404) // 啥都木有
	MethodNotAllowed   = add(-405) // 不支持该方法
	Conflict           = add(-409) // 冲突
	Canceled           = add(-498) // 客户端取消请求
	ServerErr          = add(-500) // 服务器错误
	ServiceUnavailable = add(-503) // 过载保护,服务暂不可用
	Deadline           = add(-504) // 服务调用超时
	LimitExceed        = add(-509) // 超出限制
)
*/
