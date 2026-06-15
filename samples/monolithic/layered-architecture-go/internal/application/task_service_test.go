package application

import (
	"errors"
	"strings"
	"testing"
	"time"

	"example.com/layered-architecture-go/internal/domain"
	"example.com/layered-architecture-go/internal/infrastructure"
)

// newTestService は決定的な ID と固定時刻を注入したサービスを返す。
func newTestService() *TaskService {
	repo := infrastructure.NewMemoryTaskRepository()
	var seq int
	idGen := func() string {
		seq++
		return "id-" + string(rune('0'+seq))
	}
	clock := func() time.Time { return time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC) }
	return NewTaskServiceWith(repo, idGen, clock)
}

func TestCreate_ValidationFailure(t *testing.T) {
	svc := newTestService()

	if _, err := svc.Create(""); !errors.Is(err, domain.ErrEmptyTitle) {
		t.Fatalf("空タイトルでは ErrEmptyTitle を期待: got %v", err)
	}

	long := strings.Repeat("あ", domain.MaxTitleLength+1)
	if _, err := svc.Create(long); !errors.Is(err, domain.ErrTitleTooLong) {
		t.Fatalf("長すぎるタイトルでは ErrTitleTooLong を期待: got %v", err)
	}
}

func TestCreateAndList(t *testing.T) {
	svc := newTestService()

	created, err := svc.Create("買い物に行く")
	if err != nil {
		t.Fatalf("Create が失敗: %v", err)
	}
	if created.ID == "" {
		t.Fatal("ID が採番されていない")
	}
	if created.Done {
		t.Fatal("新規タスクは未完了であるべき")
	}
	if created.Title != "買い物に行く" {
		t.Fatalf("タイトルが一致しない: %q", created.Title)
	}

	tasks, err := svc.List()
	if err != nil {
		t.Fatalf("List が失敗: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("タスク数は 1 を期待: got %d", len(tasks))
	}
	if tasks[0].ID != created.ID {
		t.Fatalf("保存されたタスクの ID が一致しない: %q != %q", tasks[0].ID, created.ID)
	}
}

func TestComplete(t *testing.T) {
	svc := newTestService()

	created, err := svc.Create("レポートを書く")
	if err != nil {
		t.Fatalf("Create が失敗: %v", err)
	}

	completed, err := svc.Complete(created.ID)
	if err != nil {
		t.Fatalf("Complete が失敗: %v", err)
	}
	if !completed.Done {
		t.Fatal("Complete 後は Done=true であるべき")
	}

	// 永続化されたタスクも完了になっていることを確認。
	tasks, _ := svc.List()
	if len(tasks) != 1 || !tasks[0].Done {
		t.Fatal("リポジトリ上のタスクも完了状態であるべき")
	}
}

func TestComplete_NotFound(t *testing.T) {
	svc := newTestService()

	if _, err := svc.Complete("does-not-exist"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("存在しない ID では ErrNotFound を期待: got %v", err)
	}
}
