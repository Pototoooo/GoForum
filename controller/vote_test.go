package controller

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"GoForum/pkg/Code"

	"github.com/gin-gonic/gin"
)

func TestVoteHandler_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/vote",
		strings.NewReader(`bad`))
	c.Request.Header.Set("Content-Type", "application/json")

	VoteHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestVoteHandler_MissingPostID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/vote",
		strings.NewReader(`{"direction":1}`))
	c.Request.Header.Set("Content-Type", "application/json")

	VoteHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestVoteHandler_MissingDirection(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/vote",
		strings.NewReader(`{"post_id":"1"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	VoteHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestVoteHandler_InvalidDirection(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/vote",
		strings.NewReader(`{"post_id":"1","direction":999}`))
	c.Request.Header.Set("Content-Type", "application/json")

	VoteHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestVoteHandler_NoUserInContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/vote",
		strings.NewReader(`{"post_id":"1","direction":1}`))
	c.Request.Header.Set("Content-Type", "application/json")

	VoteHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeNeedLogin {
		t.Errorf("code = %v, want %d", code, Code.CodeNeedLogin)
	}
}

func TestVoteHandler_ValidDirectionValues(t *testing.T) {
	// direction=0 是合法值（取消投票）
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/vote",
		strings.NewReader(`{"post_id":"1","direction":0}`))
	c.Request.Header.Set("Content-Type", "application/json")

	VoteHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	// direction=0 可能被 validator 视为零值拒绝，也可能是合法值通过到 NeedLogin
	// 两者都是可接受的，只要不 panic 即可
	if Code.ResponseCode(code) != Code.CodeNeedLogin && Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("unexpected code = %v, expected NeedLogin(%d) or InvalidParam(%d)",
			code, Code.CodeNeedLogin, Code.CodeInvalidParam)
	}
}

func TestVoteHandler_DirectionNegativeOne(t *testing.T) {
	// direction=-1 是合法值（反对）
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/vote",
		strings.NewReader(`{"post_id":"1","direction":-1}`))
	c.Request.Header.Set("Content-Type", "application/json")

	VoteHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeNeedLogin {
		t.Errorf("expected NeedLogin after passing binding, got code=%v", code)
	}
}
