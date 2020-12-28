package tke

// Msg .
func Msg(c Code) string {
	return msg[c]
}

// msg .
var msg = map[Code]string{
	Success:           "成功",
	Unknown:           "未知错误",
	Err:               "系统错误",
	Db:                "数据库错误",
	Redis:             "Redis错误",
	Forbidden:         "无权限操作资源，访问被拒绝",
	InvalidParameters: "参数错误",
	BadRequest:        "无效的请求",
	NoneToken:         "请携带登录令牌",
	InvalidToken:      "无效的登录令牌",
	TokenInvalid:      "登录令牌已失效",
	NotFound:          "请求资源不存在",
	TooManyRequests:   "请求次数过多，请稍后再试",
	PhoneNumberErr:    "手机号码不正确",
	SmsSendFail:       "发送短信失败",
	SmsCodeExpire:     "验证码已失效，请重新获取",
	SmsCodeLimit:      "短信验证码获取已达限制",
	SmsCodeEmpty:      "请输入验证码",
	SmsCodeErr:        "验证码不正确",
}
