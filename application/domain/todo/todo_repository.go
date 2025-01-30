package todo

import (
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type TodoRepository interface {
	CreateTodo() error
	GetTodos() ([]Todo, error)
}

// 小文字始まりの構造体は非公開
type todoRepository struct {
	DB *bun.DB
}

// `fx.In` で `bun.DB` を自動 DI
type todoRepositoryParams struct {
	fx.In
	DB *bun.DB
}

func NewTodoRepository(params todoRepositoryParams) TodoRepository {
	return &todoRepository{DB: params.DB}
}

// Todo作成
func (r *todoRepository) CreateTodo() error {
	return nil
}

func (r *todoRepository) GetTodos() ([]Todo, error) {
	return nil, nil
}
