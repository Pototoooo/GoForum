package controller

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"GoForum/dao/mysql"
	"GoForum/middlewire"
	"GoForum/pkg/Code"

	"github.com/gin-gonic/gin"
)

// 工具：设置用户到上下文
func setUser(c *gin.Context, userID int64) {
	c.Set(middlewire.UserID, userID)
}

func TestPostHandler_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/post",
		strings.NewReader(`not json`))
	c.Request.Header.Set("Content-Type", "application/json")

	PostHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestPostHandler_CommunityIDZero(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/post",
		strings.NewReader(`{"title":"test","content":"test","community_id":0}`))
	c.Request.Header.Set("Content-Type", "application/json")

	PostHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestPostHandler_CommunityIDNegative(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/post",
		strings.NewReader(`{"title":"test","content":"test","community_id":-1}`))
	c.Request.Header.Set("Content-Type", "application/json")

	PostHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestPostHandler_MissingCommunityID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/post",
		strings.NewReader(`{"title":"test","content":"test"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	PostHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestPostHandler_NoUserInContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/post",
		strings.NewReader(`{"title":"test","content":"test","community_id":1}`))
	c.Request.Header.Set("Content-Type", "application/json")

	PostHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeNeedLogin {
		t.Errorf("code = %v, want %d", code, Code.CodeNeedLogin)
	}
}

func TestGetDetailPostHandler_EmptyID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/post/", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: ""}}

	GetDetailPostHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestGetDetailPostHandler_InvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/post/abc", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}

	GetDetailPostHandler(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeInvalidParam {
		t.Errorf("code = %v, want %d", code, Code.CodeInvalidParam)
	}
}

func TestGetOrderedPosts_DefaultParams(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/posts2", nil)

	GetOrderedPosts(c)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if _, ok := resp["code"]; !ok {
		t.Error("response missing 'code' field")
	}
}

func TestGetOrderedPosts_WithPageParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/posts2?page=1&size=5&order=time", nil)

	GetOrderedPosts(c)

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
	if len(data) > 5 {
		t.Errorf("data length = %d, want ≤ 5", len(data))
	}
	// 验证帖子结构
	if len(data) > 0 {
		post := data[0].(map[string]interface{})
		expectedFields := []string{"id", "title", "author_name", "community_name", "content"}
		for _, field := range expectedFields {
			if _, ok := post[field]; !ok {
				t.Errorf("post missing '%s' field", field)
			}
		}
	}
}

func TestGetOrderedPosts_ByScore(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/posts2?page=1&size=5&order=score", nil)

	GetOrderedPosts(c)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if _, ok := resp["code"]; !ok {
		t.Error("response missing 'code' field")
	}
}

func TestGetOrderedPosts_ByCommunity(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/posts2?page=1&size=10&order=time&community_id=1", nil)

	GetOrderedPosts(c)

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
	// 验证所有帖子都属于 Go 社区（community_id=1）
	for _, item := range data {
		post := item.(map[string]interface{})
		if post["community_name"] != "Go" {
			t.Errorf("expected community_name 'Go', got '%v'", post["community_name"])
		}
	}
}

// TestCreateAndQueryPostFlow 测试完整的发帖+查询帖子详情流程
func TestCreateAndQueryPostFlow(t *testing.T) {
	db := mysql.GetDB()

	// 使用已有的 alice 用户
	aliceID := mysql.QueryUserIdByName("alice")
	if aliceID == 0 {
		t.Fatal("缺少测试用户 alice，请先注册")
	}

	// Step 1: 创建帖子
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/post",
		strings.NewReader(`{"title":"测试帖子标题","content":"测试帖子内容","community_id":1}`))
	c.Request.Header.Set("Content-Type", "application/json")
	setUser(c, aliceID)

	PostHandler(c)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	code := resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeSuccess {
		t.Fatalf("创建帖子失败: code=%v msg=%v", code, resp["msg"])
	}

	// 查询刚创建的帖子 ID
	var postID int64
	err := db.Get(&postID, "SELECT post_id FROM post WHERE author_id = ? ORDER BY create_time DESC LIMIT 1", aliceID)
	if err != nil {
		t.Fatalf("查询帖子 ID 失败: %v", err)
	}
	if postID == 0 {
		t.Fatal("创建的帖子 ID 为 0")
	}

	// Step 2: 查询帖子详情
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/post/1", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: fmt.Sprintf("%d", postID)}}

	GetDetailPostHandler(c)

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("detail response is not valid JSON: %v", err)
	}
	code = resp["code"].(float64)
	if Code.ResponseCode(code) != Code.CodeSuccess {
		t.Fatalf("查询帖子详情失败: code=%v msg=%v", code, resp["msg"])
	}

	data := resp["data"].(map[string]interface{})
	if data["title"] != "测试帖子标题" {
		t.Errorf("title = %v, want '测试帖子标题'", data["title"])
	}
	if data["author_name"] != "alice" {
		t.Errorf("author_name = %v, want 'alice'", data["author_name"])
	}
	if data["community_name"] != "Go" {
		t.Errorf("community_name = %v, want 'Go'", data["community_name"])
	}

	// 清理测试帖子
	db.Exec("DELETE FROM post WHERE post_id = ?", postID)
}
