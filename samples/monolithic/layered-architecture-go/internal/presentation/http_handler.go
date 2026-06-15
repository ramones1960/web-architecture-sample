// Package presentation は HTTP の入出力を担う層。
// リクエストの解釈・レスポンスの整形を行い、ユースケースはアプリケーション層へ委譲する。
package presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"example.com/layered-architecture-go/internal/application"
	"example.com/layered-architecture-go/internal/domain"
)

// taskResponse は API が返す JSON 表現。ドメインエンティティを直接公開しないことで、
// 内部モデルと外部契約を分離する。
type taskResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

type createTaskRequest struct {
	Title string `json:"title"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func toResponse(t *domain.Task) taskResponse {
	return taskResponse{
		ID:        t.ID,
		Title:     t.Title,
		Done:      t.Done,
		CreatedAt: t.CreatedAt,
	}
}

// Handler は HTTP ハンドラ群をまとめ、アプリケーション層に依存する。
type Handler struct {
	svc *application.TaskService
}

// NewHandler はハンドラを生成する。
func NewHandler(svc *application.TaskService) *Handler {
	return &Handler{svc: svc}
}

// Routes はルーティング済みの http.Handler を返す。
// Go 1.22+ の net/http パターン（メソッド + パス変数）を利用する。
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", h.listTasks)
	mux.HandleFunc("POST /tasks", h.createTask)
	mux.HandleFunc("POST /tasks/{id}/complete", h.completeTask)
	return mux
}

func (h *Handler) listTasks(w http.ResponseWriter, _ *http.Request) {
	tasks, err := h.svc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	resp := make([]taskResponse, 0, len(tasks))
	for _, t := range tasks {
		resp = append(resp, toResponse(t))
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid JSON body"))
		return
	}
	task, err := h.svc.Create(req.Title)
	if err != nil {
		writeError(w, statusForError(err), err)
		return
	}
	writeJSON(w, http.StatusCreated, toResponse(task))
}

func (h *Handler) completeTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := h.svc.Complete(id)
	if err != nil {
		writeError(w, statusForError(err), err)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(task))
}

// statusForError はドメインエラーを HTTP ステータスへ対応付ける。
func statusForError(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrEmptyTitle), errors.Is(err, domain.ErrTitleTooLong):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, errorResponse{Error: err.Error()})
}
