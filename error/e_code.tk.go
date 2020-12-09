package etk

// Code .
type Code int32

// Code .
func (c Code) Code() int32 {
	return int32(c)
}

// code
const (
	SUCCESS           Code = 0  // 成功
	UNKNOWN           Code = 1  // 未知错误
	ERROR             Code = 2  // 系统错误
	DB                Code = 3  // 数据库错误
	Redis             Code = 4  // Redis错误
	Forbidden         Code = 5  // 没有权限
	InvalidParameters Code = 6  // 参数错误
	BadRequest        Code = 7  // 无效的请求
	NoneToken         Code = 8  // 请携带登录令牌
	InvalidToken      Code = 9  // 无效的登录令牌
	TokenIDErr        Code = 10 // 请使用最新的登录令牌
	NotFound          Code = 11 // 请求资源不存在
	TooManyRequests   Code = 12 // 请求次数过多，请稍后再试
	PhoneNumberERR    Code = 13 // 手机号码不正确
	SmsSendFail       Code = 14 // 发送短信失败
	SmsCodeExpire     Code = 15 // 验证码超出有效期，请重新获取
	SmsCodeLimit      Code = 16 // 短信验证码获取已达限制
	SmsCodeEmpty      Code = 17 // 请输入验证码
	SmsCodeErr        Code = 18 // 验证码不正确
)
