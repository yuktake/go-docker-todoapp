package service

import (
	"go.uber.org/fx"
)

// Fx Module
var Module = fx.Module(
	"service",
	fx.Provide(NewTodoService),
	fx.Provide(NewUserService),
	fx.Provide(NewAuthService),
)
