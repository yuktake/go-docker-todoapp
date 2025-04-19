package db

import (
	"github.com/uptrace/bun"
	"github.com/yuktake/todo-webapp/domain"
)

type IndexedModel = domain.IndexedModel

func ModelsToByte(db *bun.DB, models []IndexedModel) []byte {
	var data []byte
	for _, model := range models {
		query := db.NewCreateTable().Model(model).WithForeignKeys()
		rawQuery, err := query.AppendQuery(db.Formatter(), nil)
		if err != nil {
			panic(err)
		}
		data = append(data, rawQuery...)
		data = append(data, ";\n"...)
	}
	return data
}

func IndexesToByte(db *bun.DB, idxCreators []func(*bun.DB) *bun.CreateIndexQuery) []byte {
	var data []byte
	for _, idxCreator := range idxCreators {
		idx := idxCreator(db)
		rawQuery, err := idx.AppendQuery(db.Formatter(), nil)
		if err != nil {
			panic(err)
		}
		data = append(data, rawQuery...)
		data = append(data, ";\n"...)
	}
	return data
}
