// Package infrastructure はドメイン層が定義するインターフェースの具体実装を提供する。
// ここではインメモリの TaskRepository を実装する（RDB 実装に差し替え可能）。
package infrastructure

import (
	"sync"

	"example.com/layered-architecture-go/internal/domain"
)

// MemoryTaskRepository は domain.TaskRepository のインメモリ実装。
// 複数の HTTP ハンドラから並行アクセスされるため mutex で保護する。
type MemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

// NewMemoryTaskRepository は空のリポジトリを生成する。
func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

// Save はタスクを保存する。外部から保持中の参照を書き換えられないようコピーを格納する。
func (r *MemoryTaskRepository) Save(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	stored := *task
	r.tasks[task.ID] = &stored
	return nil
}

// FindAll は全タスクのコピーを返す。
func (r *MemoryTaskRepository) FindAll() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		copied := *t
		result = append(result, &copied)
	}
	return result, nil
}

// FindByID は ID でタスクを取得する。存在しなければ domain.ErrNotFound を返す。
func (r *MemoryTaskRepository) FindByID(id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	copied := *t
	return &copied, nil
}
