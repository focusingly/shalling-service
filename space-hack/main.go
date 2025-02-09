package main

import (
	"space-domain/model"

	"space-api/db"

	"gorm.io/gen"
)

func main() {
	genBizCodes()
	genExtraCodes()
}

// 拓展 Sql
type Querier interface {
	/*
		SELECT * FROM @@table
		{{if len(condList) > 0}}
			WHERE
			{{for index, cond := range condList}}
				{{if index < len(condList)-1}}
					@@cond.Column = @cond.Val AND
				{{else}}
					@@cond.Column = @cond.Val
				{{end}}
			{{end}}
		{{end}}
		{{if len(sortList) > 0}}
			ORDER BY
			{{for index, sort := range sortList}}
				{{if index < len(sortList)-1}}
					{{if sort.Desc}}
						@@sort.Column DESC,
					{{else}}
						@@sort.Column ASC,
					{{end}}
				{{else}}
					{{if sort.Desc}}
						@@sort.Column DESC
					{{else}}
						@@sort.Column ASC
					{{end}}
				{{end}}
			{{end}}
		{{end}}
	*/
	SelectWithSorts(condList []model.WhereCond, sortList []model.SortColumn) ([]gen.T, error)
}

func genBizCodes() {
	db := db.GetBizDB()

	bizTbs := model.GetBizMigrateTables()
	db.AutoMigrate()
	g := gen.NewGenerator(gen.Config{
		OutPath:           "../space-domain/dao/biz",
		OutFile:           "gen.go",
		WithUnitTest:      false,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.UseDB(db)

	g.ApplyBasic(
		bizTbs...,
	)
	g.ApplyInterface(func(Querier) {}, bizTbs...)
	g.Execute()
}

func genExtraCodes() {
	db := db.GetExtraDB()
	extraBizTbs := model.GetExtraHelperMigrateTables()
	db.AutoMigrate(extraBizTbs...)
	g := gen.NewGenerator(gen.Config{
		OutPath:           "../space-domain/dao/extra",
		OutFile:           "gen.go",
		WithUnitTest:      false,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.UseDB(db)
	g.ApplyBasic(
		extraBizTbs...,
	)
	g.ApplyInterface(func(Querier) {}, extraBizTbs...)

	g.Execute()
}
