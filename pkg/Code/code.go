package Code

type ResponseCode int

const (
	CodeSuccess            ResponseCode = 1000 // 成功
	CodeInvalidParam       ResponseCode = 1001 // 参数错误
	CodeInternalError      ResponseCode = 1002 // 内部错误
	CodeNotFound           ResponseCode = 1003 // 资源不存在
	CodeAlreadyExists      ResponseCode = 1004 // 资源已存在
	CodePermissionDenied   ResponseCode = 1005 // 权限不足
	CodeUnauthorized       ResponseCode = 1006 // 未授权
	CodeTimeout            ResponseCode = 1007 // 请求超时
	CodeTooManyRequests    ResponseCode = 1008 // 请求过于频繁
	CodeServiceUnavailable ResponseCode = 1009 // 服务不可用
	CodeTokenError         ResponseCode = 1010 // 传入Token错误
	CodeNeedLogin          ResponseCode = 1011 // 需要登录
)

var codeMsgMap = map[ResponseCode]string{
	CodeSuccess:            "成功",
	CodeInvalidParam:       "参数错误",
	CodeInternalError:      "内部错误",
	CodeNotFound:           "资源不存在",
	CodeAlreadyExists:      "资源已存在",
	CodePermissionDenied:   "权限不足",
	CodeUnauthorized:       "未授权",
	CodeTimeout:            "请求超时",
	CodeTooManyRequests:    "请求过于频繁",
	CodeServiceUnavailable: "服务不可用",
	CodeTokenError:         "传入Token错误",
	CodeNeedLogin:          "需要登录",
}

func (c ResponseCode) Msg() string {
	if msg, ok := codeMsgMap[c]; ok {
		return msg
	}
	return "未知错误"
}
