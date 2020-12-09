package etk

// Msg .
func Msg(c Code) string {
	return msg[c]
}

// msg .
var msg = map[Code]string{
	SUCCESS:           "成功",
	UNKNOWN:           "未知错误",
	ERROR:             "系统错误",
	DB:                "数据库错误",
	Redis:             "Redis错误",
	Forbidden:         "没有权限",
	InvalidParameters: "参数错误",
	BadRequest:        "无效的请求",
	NoneToken:         "请携带登录令牌",
	InvalidToken:      "无效的登录令牌",
	TokenIDErr:        "请使用最新的登录令牌",
	NotFound:          "请求资源不存在",
	TooManyRequests:   "请求次数过多，请稍后再试",
	PhoneNumberERR:    "手机号码不正确",
	SmsSendFail:       "发送短信失败",
	SmsCodeExpire:     "验证码超出有效期，请重新获取",
	SmsCodeLimit:      "短信验证码获取已达限制",
	SmsCodeEmpty:      "请输入验证码",
	SmsCodeErr:        "验证码不正确",
}
