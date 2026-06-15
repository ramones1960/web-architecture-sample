package domain

// TaskRepository はタスクの永続化に関する契約（インターフェース）を定義する。
//
// インターフェースをドメイン層に置くことで、ドメイン／アプリケーション層は
// 具体的な永続化技術（メモリ、RDB など）に依存しない。実装はインフラ層が提供し、
// 依存方向は外側（インフラ）から内側（ドメイン）へ向かう（依存性逆転）。
type TaskRepository interface {
	// Save はタスクを保存する（新規・更新の両方に用いる）。
	Save(task *Task) error
	// FindAll は全タスクを返す。
	FindAll() ([]*Task, error)
	// FindByID は ID でタスクを取得する。存在しない場合は ErrNotFound を返す。
	FindByID(id string) (*Task, error)
}
