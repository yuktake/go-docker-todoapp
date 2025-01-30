package domain

import (
	"go.uber.org/fx"

	"github.com/yuktake/todo-webapp/domain/todo"
)

// Fx Module
var Module = fx.Module("domain",
	fx.Provide(todo.NewTodoRepository),
)
