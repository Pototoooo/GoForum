package response

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"GoForum/pkg/Code"

	"github.com/gin-gonic/gin"
)

func TestResponseSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"key": "value"}
	ResponseSuccess(c, data)

	var resp ResponseData
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != Code.CodeSuccess {
		t.Errorf("Code = %d, want %d", resp.Code, Code.CodeSuccess)
	}
	if resp.Msg != "成功" {
		t.Errorf("Msg = %q, want %q", resp.Msg, "成功")
	}
	if resp.Data == nil {
		t.Fatal("Data should not be nil")
	}
}

func TestResponseSuccess_NilData(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ResponseSuccess(c, nil)

	var resp ResponseData
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != Code.CodeSuccess {
		t.Errorf("Code = %d, want %d", resp.Code, Code.CodeSuccess)
	}
}

func TestResponseError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ResponseError(c, Code.CodeNotFound)

	var resp ResponseData
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != Code.CodeNotFound {
		t.Errorf("Code = %d, want %d", resp.Code, Code.CodeNotFound)
	}
	if resp.Msg != "资源不存在" {
		t.Errorf("Msg = %q, want %q", resp.Msg, "资源不存在")
	}
	if resp.Data != nil {
		t.Error("Data should be nil for error response")
	}
}

func TestResponseErrorWithMsg(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	customMsg := "自定义错误信息"
	ResponseErrorWithMsg(c, Code.CodeInvalidParam, customMsg)

	var resp ResponseData
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != Code.CodeInvalidParam {
		t.Errorf("Code = %d, want %d", resp.Code, Code.CodeInvalidParam)
	}
	if resp.Msg != customMsg {
		t.Errorf("Msg = %q, want %q", resp.Msg, customMsg)
	}
	if resp.Data != nil {
		t.Error("Data should be nil for error response")
	}
}

func TestResponseStatusCode(t *testing.T) {
	// 所有响应都是 HTTP 200
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ResponseError(c, Code.CodeNeedLogin)

	if w.Code != 200 {
		t.Errorf("HTTP status = %d, want 200", w.Code)
	}
}
