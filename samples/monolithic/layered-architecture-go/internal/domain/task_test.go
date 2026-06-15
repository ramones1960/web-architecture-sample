package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewTask_Valid(t *testing.T) {
	at := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	task, err := NewTask("id-1", "ミーティングの準備", at)
	if err != nil {
		t.Fatalf("有効なタイトルで失敗: %v", err)
	}
	if task.Done {
		t.Fatal("新規タスクは未完了であるべき")
	}
	if !task.CreatedAt.Equal(at) {
		t.Fatal("CreatedAt が一致しない")
	}
}

func TestNewTask_Validation(t *testing.T) {
	at := time.Now()

	if _, err := NewTask("id", "", at); !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("空タイトルでは ErrEmptyTitle を期待: got %v", err)
	}

	boundary := strings.Repeat("x", MaxTitleLength)
	if _, err := NewTask("id", boundary, at); err != nil {
		t.Fatalf("上限ちょうど(%d文字)は許可されるべき: got %v", MaxTitleLength, err)
	}

	over := strings.Repeat("x", MaxTitleLength+1)
	if _, err := NewTask("id", over, at); !errors.Is(err, ErrTitleTooLong) {
		t.Fatalf("上限超過では ErrTitleTooLong を期待: got %v", err)
	}
}

func TestComplete(t *testing.T) {
	task, _ := NewTask("id", "タイトル", time.Now())
	task.Complete()
	if !task.Done {
		t.Fatal("Complete 後は Done=true であるべき")
	}
}
