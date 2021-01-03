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
	Redis             Code = 10004 // Redis错误
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
)
