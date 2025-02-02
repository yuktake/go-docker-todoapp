package domain

import (
	"go.uber.org/fx"

	"github.com/yuktake/todo-webapp/domain/todo"
	"github.com/yuktake/todo-webapp/domain/user"
)

// Fx Module
var Module = fx.Module("domain",
	fx.Provide(todo.NewTodoRepository),
	fx.Provide(user.NewUserRepository),
)
