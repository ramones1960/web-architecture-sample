// Package domain はビジネスの中核となるエンティティ・値・ルールを定義する。
// このパッケージは他のどの層にも依存しない（依存方向の最も内側）。
package domain

import (
	"errors"
	"time"
	"unicode/utf8"
)

// タイトルの最大文字数。
const MaxTitleLength = 200

// ドメインエラー。上位層はこれらを判別してHTTPステータス等に変換する。
var (
	// ErrEmptyTitle はタイトルが空のときに返る。
	ErrEmptyTitle = errors.New("task: title must not be empty")
	// ErrTitleTooLong はタイトルが上限を超えたときに返る。
	ErrTitleTooLong = errors.New("task: title is too long")
	// ErrNotFound は対象のタスクが存在しないときに返る。
	ErrNotFound = errors.New("task: not found")
)

// Task はタスク管理（ToDo）の中核エンティティ。
type Task struct {
	ID        string
	Title     string
	Done      bool
	CreatedAt time.Time
}

// NewTask は不変条件を検証したうえで新しい Task を生成する。
// title は前後の空白を除いて空でなく、MaxTitleLength 文字以下でなければならない。
func NewTask(id, title string, createdAt time.Time) (*Task, error) {
	if err := validateTitle(title); err != nil {
		return nil, err
	}
	return &Task{
		ID:        id,
		Title:     title,
		Done:      false,
		CreatedAt: createdAt,
	}, nil
}

// Complete はタスクを完了状態にする。
func (t *Task) Complete() {
	t.Done = true
}

func validateTitle(title string) error {
	if title == "" {
		return ErrEmptyTitle
	}
	if utf8.RuneCountInString(title) > MaxTitleLength {
		return ErrTitleTooLong
	}
	return nil
}
