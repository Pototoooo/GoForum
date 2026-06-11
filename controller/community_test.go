package controller

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"GoForum/pkg/Code"

	"github.com/gin-gonic/gin"
)

func TestCommunityHandler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/community", nil)

	CommunityHandler(c)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeSuccess {
		t.Fatalf("code = %v, want %d", code, Code.CodeSuccess)
	}

	data, ok := resp["data"].([]interface{})
	if !ok {
		t.Fatal("data should be an array")
	}
	if len(data) == 0 {
		t.Error("community list should not be empty")
	}
	// 验证第一个社区的字段
	first := data[0].(map[string]interface{})
	if _, ok := first["id"]; !ok {
		t.Error("community missing 'id' field")
	}
	if _, ok := first["name"]; !ok {
		t.Error("community missing 'name' field")
	}
}

func TestCommunityDetailHandler_InvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/community/abc", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}

	CommunityDetailHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestCommunityDetailHandler_EmptyID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/community/", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

	CommunityDetailHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestCommunityDetailHandler_Valid(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/community/1", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	CommunityDetailHandler(c)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeSuccess {
		t.Fatalf("code = %v, want %d", code, Code.CodeSuccess)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data should be an object")
	}
	if data["id"] != float64(1) {
		t.Errorf("community_id = %v, want 1", data["id"])
	}
	if name, ok := data["name"]; !ok || name == "" {
		t.Error("community name should not be empty")
	}
}
