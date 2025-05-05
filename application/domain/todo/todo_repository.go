package todo

import (
	"context"

	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type TodoRepository interface {
	CreateTodo(todo *Todo) (Todo, error)
	GetTodos() ([]Todo, error)
	GetTodoByID(id string) (Todo, error)
	UpdateTodo(todo Todo) (Todo, error)
	DeleteTodoByID(id string) error
}

// 小文字始まりの構造体は非公開
type todoRepository struct {
	DB *bun.DB
}

// `fx.In` で `bun.DB` を自動 DI
type TodoRepositoryParams struct {
	fx.In
	DB *bun.DB
}

func NewTodoRepository(params TodoRepositoryParams) TodoRepository {
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

// Todo取得
func (r *todoRepository) GetTodoByID(id string) (Todo, error) {
	var todo Todo
	ctx := context.Background()

	err := r.DB.NewSelect().Model(&todo).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
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

func (r *todoRepository) UpdateTodo(todo Todo) (Todo, error) {
	ctx := context.Background()

	_, err := r.DB.NewUpdate().Model(&todo).Where("id = ?", todo.ID).Exec(ctx)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (r *todoRepository) DeleteTodoByID(id string) error {
	ctx := context.Background()

	_, err := r.DB.NewDelete().Model(&Todo{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
