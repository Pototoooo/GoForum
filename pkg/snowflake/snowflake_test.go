package snowflake

import (
	"testing"
)

func TestInit_InvalidTime(t *testing.T) {
	err := Init("invalid-date", 1)
	if err == nil {
		t.Fatal("expected error for invalid date, got nil")
	}
}

func TestInit_InvalidMachineID(t *testing.T) {
	err := Init("2023-01-01", -1)
	if err == nil {
		t.Fatal("expected error for negative machine ID, got nil")
	}
}

func TestGenerateID_SuccessiveCalls(t *testing.T) {
	if err := Init("2023-01-01", 1); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	id1 := GenerateID()
	id2 := GenerateID()

	if id1 == 0 {
		t.Error("GenerateID() returned 0")
	}
	if id2 <= id1 {
		t.Error("GenerateID() should return monotonically increasing values")
	}
}

func TestGenerateID_UniqueAcrossMachines(t *testing.T) {
	if err := Init("2023-01-01", 2); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	ids := make(map[int64]bool)
	const count = 100
	for i := 0; i < count; i++ {
		id := GenerateID()
		if ids[id] {
			t.Fatalf("duplicate ID generated: %d", id)
		}
		ids[id] = true
	}
}

func TestGenerateID_Type(t *testing.T) {
	if err := Init("2023-01-01", 3); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 雪花 ID 应该是 64 位正数
	id := GenerateID()
	if id <= 0 {
		t.Errorf("expected positive int64, got %d", id)
	}
}

func TestInit_TimeParsing(t *testing.T) {
	err := Init("2023-01-01", 1)
	if err != nil {
		t.Errorf("Init() with valid date should succeed, got: %v", err)
	}
	if Node == nil {
		t.Error("Node should not be nil after Init()")
	}
}
