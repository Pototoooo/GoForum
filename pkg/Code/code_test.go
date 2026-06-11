package Code

import "testing"

func TestResponseCode_Msg(t *testing.T) {
	tests := []struct {
		code ResponseCode
		want string
	}{
		{CodeSuccess, "成功"},
		{CodeInvalidParam, "参数错误"},
		{CodeInternalError, "内部错误"},
		{CodeNotFound, "资源不存在"},
		{CodeAlreadyExists, "资源已存在"},
		{CodePermissionDenied, "权限不足"},
		{CodeUnauthorized, "未授权"},
		{CodeTimeout, "请求超时"},
		{CodeTooManyRequests, "请求过于频繁"},
		{CodeServiceUnavailable, "服务不可用"},
		{CodeTokenError, "传入Token错误"},
		{CodeNeedLogin, "需要登录"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.code.Msg(); got != tt.want {
				t.Errorf("ResponseCode(%d).Msg() = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

func TestResponseCode_Msg_Unknown(t *testing.T) {
	if got := ResponseCode(9999).Msg(); got != "未知错误" {
		t.Errorf("unknown ResponseCode.Msg() = %q, want %q", got, "未知错误")
	}
}

func TestResponseCode_Values(t *testing.T) {
	if CodeSuccess != 1000 {
		t.Errorf("CodeSuccess = %d, want 1000", CodeSuccess)
	}
	if CodeNeedLogin != 1011 {
		t.Errorf("CodeNeedLogin = %d, want 1011", CodeNeedLogin)
	}
}

func TestResponseCode_IntConversion(t *testing.T) {
	code := ResponseCode(1000)
	if int(code) != 1000 {
		t.Errorf("ResponseCode(1000) as int = %d, want 1000", int(code))
	}
}
