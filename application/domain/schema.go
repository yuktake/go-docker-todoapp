package domain

import (
	"github.com/uptrace/bun"
)

type IndexedModel interface {
	Indexes() []func(*bun.DB) *bun.CreateIndexQuery
}
