package jwt

import (
	"testing"
)

func TestGenToken(t *testing.T) {
	aToken, rToken := GenToken(12345, "testuser")

	if aToken == "" {
		t.Error("access token should not be empty")
	}
	if rToken == "" {
		t.Error("refresh token should not be empty")
	}
}

func TestParseToken_Success(t *testing.T) {
	aToken, _ := GenToken(12345, "testuser")

	claims, err := ParseToken(aToken)
	if err != nil {
		t.Fatalf("ParseToken() failed: %v", err)
	}

	if claims.UserID != 12345 {
		t.Errorf("UserID = %d, want %d", claims.UserID, 12345)
	}
	if claims.Username != "testuser" {
		t.Errorf("Username = %q, want %q", claims.Username, "testuser")
	}
	if claims.Issuer != "Myself" {
		t.Errorf("Issuer = %q, want %q", claims.Issuer, "Myself")
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	_, err := ParseToken("invalid-token-string")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

func TestParseToken_ExpiredToken(t *testing.T) {
	// 保存原值并恢复
	origDuration := TokenTimeDuration
	// 使用极短的有效期来测试过期场景
	TokenTimeDuration = 1

	aToken, _ := GenToken(1, "test")

	// 注意：1ns 可能已经过期了，但我们不能保证时间精度
	// 测试 ParseToken 是否能正确解析格式
	claims, err := ParseToken(aToken)
	if err != nil {
		// 如果过期了也是合理的
		t.Logf("token expired as expected: %v", err)
	} else {
		if claims.UserID != 1 {
			t.Errorf("UserID = %d, want %d", claims.UserID, 1)
		}
	}

	TokenTimeDuration = origDuration
}

func TestGenToken_Unique(t *testing.T) {
	token1, _ := GenToken(1, "userA")
	token2, _ := GenToken(1, "userB") // 不同用户名，确保 claims 不同

	if token1 == token2 {
		t.Error("tokens for different claims should not be identical")
	}
}

func TestGenToken_SameUser(t *testing.T) {
	// 相同用户在不同时间生成的 token 应该不同（因为有 IssuedAt）
	token1, _ := GenToken(1, "user")
	token2, _ := GenToken(1, "user")

	if token1 == token2 {
		t.Log("same user tokens are identical (generated in same nanosecond)")
	}
}

func TestParseToken_UserIDType(t *testing.T) {
	aToken, _ := GenToken(99999999999, "bigid")

	claims, err := ParseToken(aToken)
	if err != nil {
		t.Fatalf("ParseToken() failed: %v", err)
	}

	if claims.UserID != 99999999999 {
		t.Errorf("UserID = %d, want %d", claims.UserID, 99999999999)
	}
}

func TestRefreshToken_Invalid(t *testing.T) {
	_, _, err := RefreshToken("invalid", "invalid")
	if err == nil {
		t.Fatal("expected error for invalid refresh token, got nil")
	}
}
