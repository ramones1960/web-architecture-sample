// Package application はユースケース（業務手続き）を表現する層。
// ドメイン層のインターフェース（domain.TaskRepository）にのみ依存し、
// 具体的な永続化技術や HTTP には依存しない。
package application

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"example.com/layered-architecture-go/internal/domain"
)

// IDGenerator はタスク ID を生成する関数型。テストで差し替え可能にするために抽象化する。
type IDGenerator func() string

// Clock は現在時刻を返す関数型。テストで時刻を固定できるようにする。
type Clock func() time.Time

// TaskService はタスク管理のユースケースを提供する。
type TaskService struct {
	repo  domain.TaskRepository
	newID IDGenerator
	now   Clock
}

// NewTaskService はリポジトリを受け取りサービスを生成する。
// ID 生成と時刻はデフォルト実装を用いる。
func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{
		repo:  repo,
		newID: defaultIDGenerator,
		now:   time.Now,
	}
}

// NewTaskServiceWith は ID 生成・時刻を明示的に注入してサービスを生成する（主にテスト用）。
func NewTaskServiceWith(repo domain.TaskRepository, newID IDGenerator, now Clock) *TaskService {
	return &TaskService{repo: repo, newID: newID, now: now}
}

// Create は新しいタスクを生成して保存し、保存されたタスクを返す。
// タイトルが不正な場合はドメインエラー（ErrEmptyTitle / ErrTitleTooLong）を返す。
func (s *TaskService) Create(title string) (*domain.Task, error) {
	task, err := domain.NewTask(s.newID(), title, s.now())
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(task); err != nil {
		return nil, err
	}
	return task, nil
}

// List は全タスクを返す。
func (s *TaskService) List() ([]*domain.Task, error) {
	return s.repo.FindAll()
}

// Complete は指定 ID のタスクを完了状態にして保存する。
// 対象が存在しない場合は domain.ErrNotFound を返す。
func (s *TaskService) Complete(id string) (*domain.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	task.Complete()
	if err := s.repo.Save(task); err != nil {
		return nil, err
	}
	return task, nil
}

func defaultIDGenerator() string {
	b := make([]byte, 16)
	// crypto/rand.Read は通常失敗しないが、念のためフォールバックする。
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000000")))
	}
	return hex.EncodeToString(b)
}
