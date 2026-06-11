package controller

import (
	"errors"
	"net/http/httptest"
	"testing"

	"GoForum/middlewire"

	"github.com/gin-gonic/gin"
)

func TestGetCurrentUserID_Success(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middlewire.UserID, int64(12345))

	userID, err := GetCurrentUserID(c)
	if err != nil {
		t.Fatalf("GetCurrentUserID() failed: %v", err)
	}
	if userID != 12345 {
		t.Errorf("userID = %d, want %d", userID, 12345)
	}
}

func TestGetCurrentUserID_NotSet(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	_, err := GetCurrentUserID(c)
	if !errors.Is(err, ErrorRounterNoUserInfo) {
		t.Errorf("expected ErrorRounterNoUserInfo, got %v", err)
	}
}

func TestGetCurrentUserID_WrongType(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middlewire.UserID, "not-an-int64")

	_, err := GetCurrentUserID(c)
	if err == nil {
		t.Fatal("expected error for wrong type, got nil")
	}
}

func TestGetPageNSize_Default(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	page, size := GetPageNSize(c)
	if page != 0 {
		t.Errorf("page = %d, want 0", page)
	}
	if size != 0 {
		t.Errorf("size = %d, want 0", size)
	}
}

func TestGetPageNSize_WithParams(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=3&size=20", nil)

	page, size := GetPageNSize(c)
	if page != 3 {
		t.Errorf("page = %d, want 3", page)
	}
	if size != 20 {
		t.Errorf("size = %d, want 20", size)
	}
}

func TestGetPageNSize_InvalidParams(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=abc&size=xyz", nil)

	// 非法字符串应该被解析为 0（strconv.Atoi 默认值）
	page, size := GetPageNSize(c)
	if page != 0 || size != 0 {
		t.Errorf("expected (0,0) for invalid params, got (%d,%d)", page, size)
	}
}

func TestGetPageNSize_PartialParams(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?page=5", nil)

	page, size := GetPageNSize(c)
	if page != 5 {
		t.Errorf("page = %d, want 5", page)
	}
	if size != 0 {
		t.Errorf("size = %d, want 0 (default)", size)
	}
}
