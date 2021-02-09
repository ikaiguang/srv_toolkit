package tke

import (
	"strconv"
	"sync"
	"sync/atomic"
)

func init() {
	Register(map[Code]string{})
}

// msg
var (
	_messages atomic.Value // NOTE: stored map[Code]string
	_msgMu    sync.Mutex
)

// Register register code message map.
func Register(cm map[Code]string) {
	_msgMu.Lock()
	defer _msgMu.Unlock()

	m, ok := _messages.Load().(map[Code]string)
	if !ok {
		m = make(map[Code]string)
		for c := range msg {
			m[c] = msg[c]
		}
	}
	for c := range cm {
		m[c] = cm[c]
	}
	//for c := range msg {
	//	if _, ok := cm[c]; !ok {
	//		cm[c] = msg[c]
	//	}
	//}
	_messages.Store(m)
}

// Msg .
func Msg(c Code) string {
	if cm, ok := _messages.Load().(map[Code]string); ok {
		if msg, ok := cm[c]; ok {
			return msg
		}
	}
	return strconv.Itoa(int(c.Code()))
}

// msg .
var msg = map[Code]string{
	Init:              "init",
	Success:           "成功",
	Unknown:           "未知错误",
	Err:               "系统错误",
	Panic:             "系统错误:Panic",
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
	JwtAckUnknown:     "登录参数有误",
	JwtAckDeny:        "非法登录",
	JwtAckDeleted:     "无效用户",
	JwtAckInvalid:     "请先激活账户",
	JwtParamErr:       "JWT参数有误",
	JwtAudienceEmpty:  "请填写签发的目标收件人JWT_Audience",
	JwtIdEmpty:        "请填写签发的编号JWT_Id",
	JwtSigned:         "您的账户已在其他设备登录",
	InvalidData:       "无效的数据，请重试",
	InvalidDBData:     "无效的数据库数据",
	InvalidRedisData:  "无效的缓存数据，请重试",
}
