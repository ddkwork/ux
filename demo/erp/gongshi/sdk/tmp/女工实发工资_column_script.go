package main

import (
	"fmt"

	"github.com/ddkwork/ux/demo/erp/gongshi/sdk"
)

func RunScript(t *TreeTable, rowIndex int) {
	nameVal, ok := t.GetCellByRowIndex(rowIndex, "姓名")
	if !ok {
		return
	}
	name := fmt.Sprintf("%v", nameVal)

	nvGongVal, ok := t.GetCellByRowIndex(rowIndex, "女工日结")
	if !ok {
		return
	}
	nvGong, _ := ToFloat(nvGongVal)

	sanRenZuSum := t.SumIf("姓名", "三人组", "女工日结")
	switch name {
	case "拼车", "三人组":
		t.SetCellValue(rowIndex, "计算结果", 0.0)
	case "房东":
		t.SetCellValue(rowIndex, "计算结果", nvGong)
	case "杨萍":
		t.SetCellValue(rowIndex, "计算结果", (sanRenZuSum/3.0)+nvGong)
	case "二人组":
		t.SetCellValue(rowIndex, "计算结果", (sanRenZuSum/3.0)+(nvGong/2.0))
	default:
		t.SetCellValue(rowIndex, "计算结果", 0.0)
	}
}
