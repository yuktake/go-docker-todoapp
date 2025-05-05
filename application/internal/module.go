package internal

import (
	"github.com/yuktake/todo-webapp/internal/echoapp"
	"go.uber.org/fx"
)

var Module = fx.Module("internal",
	fx.Provide(
		echoapp.NewEcho,
	),
)
