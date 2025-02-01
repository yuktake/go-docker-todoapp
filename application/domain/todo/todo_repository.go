package todo

import (
	"context"

	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type TodoRepository interface {
	CreateTodo(todo *Todo) (Todo, error)
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
func (r *todoRepository) CreateTodo(todo *Todo) (Todo, error) {
	ctx := context.Background()

	_, err := r.DB.NewInsert().Model(todo).Returning("*").Exec(ctx)
	if err != nil {
		return Todo{}, err
	}

	return *todo, nil
}

func (r *todoRepository) GetTodos() ([]Todo, error) {
	// リポジトリからデータを取得
	var todos []Todo
	ctx := context.Background()

	err := r.DB.NewSelect().Model(&todos).Order("created_at").Scan(ctx)

	if err != nil {
		return nil, err
	}

	return todos, nil
}
