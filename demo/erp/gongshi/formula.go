package main

import (
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/ux/demo/erp/gongshi/sdk"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ddkwork/ux/demo/erp/gongshi/sdk
var Symbols = interp.Exports{}

type YaegiEngine struct {
	interp *interp.Interpreter
	table  *sdk.TreeTable
}

func NewYaegiEngine(table *sdk.TreeTable) *YaegiEngine {
	i := interp.New(interp.Options{
		GoPath:       "./",
		Unrestricted: true,
	})
	i.Use(stdlib.Symbols)

	engine := &YaegiEngine{interp: i, table: table}
	return engine
}

func (e *YaegiEngine) UpdateRowCell(rowIndex int) {
	row := e.table.GetRow(rowIndex)
	if row == nil {
		panic("行不存在")
	}

	i := interp.New(interp.Options{
		GoPath:       "./",
		Unrestricted: true,
	})
	mylog.Check(i.Use(stdlib.Symbols))
	mylog.Check(i.Use(Symbols))

	for _, cell := range row.RowCells {
		if cell.IsFormula() {
			for _, column := range e.table.Columns {
				if cell.ColumnName == column.Name {
					mylog.Check2(i.Eval(column.Formula))
					runScript := mylog.Check2(i.Eval("RunScript")).Interface().(func(*sdk.TreeTable, int))
					runScript(e.table, rowIndex)
				}
			}
		}
	}
}
