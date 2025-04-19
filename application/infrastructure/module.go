package infrastructure

import (
	"github.com/yuktake/todo-webapp/infrastructure/db"
	"go.uber.org/fx"
)

// Fx Module
var Module = fx.Module("infrastructure",
	fx.Provide(
		db.NewDBConfig,
		db.InitDB,
		db.NewBunDB,
	),
)
