package db

import (
	"github.com/uptrace/bun"
	"github.com/yuktake/todo-webapp/domain"
)

type IndexedModel = domain.IndexedModel

// Schemaに生成するテーブルの構造体を定義
func getIndexedModels() []IndexedModel {
	return []IndexedModel{
		(*User)(nil),
		(*Todo)(nil),
	}
}

func ModelsToByte(db *bun.DB) []byte {
	var data []byte
	models := getIndexedModels()

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

func IndexesToByte(db *bun.DB) []byte {
	var data []byte
	models := getIndexedModels()
	for _, model := range models {
		idxCreators := model.Indexes()
		for _, idxCreator := range idxCreators {
			idx := idxCreator(db)
			rawQuery, err := idx.AppendQuery(db.Formatter(), nil)
			if err != nil {
				panic(err)
			}
			data = append(data, rawQuery...)
			data = append(data, ";\n"...)
		}
	}

	return data
}
