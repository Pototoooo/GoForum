package controller

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"GoForum/dao/mysql"
	"GoForum/pkg/Code"

	"github.com/gin-gonic/gin"
)

// 工具函数：解析响应 JSON
func parseResponse(t *testing.T, body []byte) map[string]interface{} {
	t.Helper()
	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	return resp
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/signup",
		strings.NewReader(`{bad json`))
	c.Request.Header.Set("Content-Type", "application/json")

	RegisterHandler(c)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestRegisterHandler_MissingField(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// 缺少 re_password
	c.Request = httptest.NewRequest("POST", "/api/v1/signup",
		strings.NewReader(`{"username":"test","password":"123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	RegisterHandler(c)

	resp := parseResponse(t, w.Body.Bytes())
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestRegisterHandler_EmptyBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/signup",
		strings.NewReader(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")

	RegisterHandler(c)

	resp := parseResponse(t, w.Body.Bytes())
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/login",
		strings.NewReader(`not json`))
	c.Request.Header.Set("Content-Type", "application/json")

	LoginHandler(c)

	resp := parseResponse(t, w.Body.Bytes())
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestLoginHandler_MissingPassword(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/login",
		strings.NewReader(`{"username":"test"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	LoginHandler(c)

	resp := parseResponse(t, w.Body.Bytes())
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestLoginHandler_MissingUsername(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/login",
		strings.NewReader(`{"password":"123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	LoginHandler(c)

	resp := parseResponse(t, w.Body.Bytes())
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

// TestRegisterAndLoginFlow 测试完整的注册+登录流程（依赖 MySQL）
func TestRegisterAndLoginFlow(t *testing.T) {
	testUsername := "testuser_unittest"

	// 清理之前的测试用户
	mysql.QueryUserIdByName(testUsername)
	db := mysql.GetDB()
	db.Exec("DELETE FROM user WHERE username = ?", testUsername)

	// Step 1: 注册新用户
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/signup",
		strings.NewReader(`{"username":"`+testUsername+`","password":"test123","re_password":"test123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	RegisterHandler(c)

	resp := parseResponse(t, w.Body.Bytes())
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeSuccess {
		t.Fatalf("注册失败: code=%v msg=%v", code, resp["msg"])
	}

	// 验证用户已写入数据库
	userID := mysql.QueryUserIdByName(testUsername)
	if userID == 0 {
		t.Fatal("注册后数据库未找到用户")
	}

	// Step 2: 登录
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/login",
		strings.NewReader(`{"username":"`+testUsername+`","password":"test123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	LoginHandler(c)

	resp = parseResponse(t, w.Body.Bytes())
	code = resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeSuccess {
		t.Fatalf("登录失败: code=%v msg=%v", code, resp["msg"])
	}

	// 验证返回了 token
	data := resp["data"].(map[string]interface{})
	if token, ok := data["token"]; !ok || token == "" {
		t.Error("登录响应缺少 token")
	}
	if username, ok := data["username"]; !ok || username != testUsername {
		t.Errorf("登录响应 username = %v, want %s", username, testUsername)
	}

	// 清理测试用户
	db.Exec("DELETE FROM user WHERE username = ?", testUsername)
}

// TestRegisterHandler_DuplicateUsername 测试重复注册
func TestRegisterHandler_DuplicateUsername(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/signup",
		strings.NewReader(`{"username":"alice","password":"test123","re_password":"test123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	RegisterHandler(c)

	resp := parseResponse(t, w.Body.Bytes())
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("重复注册应该返回 InvalidParam, 实际 code=%v", code)
	}
}

// TestLoginHandler_WrongPassword 测试错误密码登录
func TestLoginHandler_WrongPassword(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/login",
		strings.NewReader(`{"username":"alice","password":"wrongpassword"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	LoginHandler(c)

	resp := parseResponse(t, w.Body.Bytes())
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("密码错误应返回 InvalidParam, 实际 code=%v", code)
	}
}

// TestLoginHandler_NotExistUser 测试不存在的用户登录
func TestLoginHandler_NotExistUser(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/login",
		strings.NewReader(`{"username":"user_not_exist_12345","password":"test123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	LoginHandler(c)

	resp := parseResponse(t, w.Body.Bytes())
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("不存在的用户应返回 InvalidParam, 实际 code=%v", code)
	}
}
