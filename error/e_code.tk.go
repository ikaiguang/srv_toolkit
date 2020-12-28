package tke

// Code .
type Code int32

// Code .
func (c Code) Code() int32 {
	return int32(c)
}

// code
const (
	Success           Code = 200 // 成功
	Unknown           Code = 1   // 未知错误
	Err               Code = 2   // 系统错误
	Db                Code = 3   // 数据库错误
	Redis             Code = 4   // Redis错误
	Forbidden         Code = 5   // 无权限操作资源，访问被拒绝
	InvalidParameters Code = 6   // 参数错误
	BadRequest        Code = 7   // 无效的请求
	NoneToken         Code = 8   // 请携带登录令牌
	InvalidToken      Code = 9   // 无效的登录令牌
	TokenInvalid      Code = 10  // 登录令牌已失效
	NotFound          Code = 11  // 请求资源不存在
	TooManyRequests   Code = 12  // 请求次数过多，请稍后再试
	PhoneNumberErr    Code = 13  // 手机号码不正确
	SmsSendFail       Code = 14  // 发送短信失败
	SmsCodeExpire     Code = 15  // 验证码已失效，请重新获取
	SmsCodeLimit      Code = 16  // 短信验证码获取已达限制
	SmsCodeEmpty      Code = 17  // 请输入验证码
	SmsCodeErr        Code = 18  // 验证码不正确
)
