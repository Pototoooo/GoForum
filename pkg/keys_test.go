package pkg

import (
	"testing"
)

func TestGetRedisKey(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		id     string
		want   string
	}{
		{"post prefix", KeyPostInfoPrefix, "123", "bluebell:post:123"},
		{"vote prefix", KeyPostVotedPrefix, "456", "bluebell:post:voted:456"},
		{"community prefix", KeyCommunityPostPrefix, "1", "bluebell:community:1"},
		{"empty id", KeyPostInfoPrefix, "", "bluebell:post:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRedisKey(tt.prefix, tt.id); got != tt.want {
				t.Errorf("GetRedisKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetPostInfoKey(t *testing.T) {
	if got := GetPostInfoKey("789"); got != "bluebell:post:789" {
		t.Errorf("GetPostInfoKey() = %q, want %q", got, "bluebell:post:789")
	}
}

func TestGetPostVotedKey(t *testing.T) {
	if got := GetPostVotedKey("789"); got != "bluebell:post:voted:789" {
		t.Errorf("GetPostVotedKey() = %q, want %q", got, "bluebell:post:voted:789")
	}
}

func TestGetCommunityPostKey(t *testing.T) {
	if got := GetCommunityPostKey("2"); got != "bluebell:community:2" {
		t.Errorf("GetCommunityPostKey() = %q, want %q", got, "bluebell:community:2")
	}
}

func TestConstants(t *testing.T) {
	if KeyPostTime != "bluebell:post:time" {
		t.Errorf("KeyPostTime = %q, want %q", KeyPostTime, "bluebell:post:time")
	}
	if KeyPostScore != "bluebell:post:score" {
		t.Errorf("KeyPostScore = %q, want %q", KeyPostScore, "bluebell:post:score")
	}
	if OrderTime != "time" {
		t.Errorf("OrderTime = %q, want %q", OrderTime, "time")
	}
	if OrderScore != "score" {
		t.Errorf("OrderScore = %q, want %q", OrderScore, "score")
	}
}
