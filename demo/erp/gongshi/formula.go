package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// =============================
// 表格数据结构
// =============================

type FormulaTable struct {
	Rows    []map[string]interface{}
	Columns map[string]*ColumnConfig
}

type ColumnConfig struct {
	Name    string
	Type    string // "text", "number", "formula"
	Formula string // 用户输入的Go代码
}

// =============================
// SDK函数库（预编译的）
// =============================

type SDKFunctions struct {
	rows         []map[string]interface{}
	currentRow   map[string]interface{}
	currentIndex int
}

// Sum 计算指定列的总和
func (sdk *SDKFunctions) Sum(columnName string) float64 {
	total := 0.0
	for _, row := range sdk.rows {
		if val, ok := row[columnName].(float64); ok {
			total += val
		} else if val, ok := row[columnName].(int); ok {
			total += float64(val)
		} else if val, ok := row[columnName].(string); ok {
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				total += f
			}
		}
	}
	return total
}

// SumIf 条件求和
func (sdk *SDKFunctions) SumIf(targetColumn, conditionColumn, conditionValue string) float64 {
	total := 0.0
	for _, row := range sdk.rows {
		if fmt.Sprint(row[conditionColumn]) == conditionValue {
			if val, ok := row[targetColumn].(float64); ok {
				total += val
			} else if val, ok := row[targetColumn].(int); ok {
				total += float64(val)
			} else if val, ok := row[targetColumn].(string); ok {
				if f, err := strconv.ParseFloat(val, 64); err == nil {
					total += f
				}
			}
		}
	}
	return total
}

// CurrentRow 获取当前行数据
func (sdk *SDKFunctions) CurrentRow() map[string]interface{} {
	return sdk.currentRow
}

// GetCell 获取指定单元格值
func (sdk *SDKFunctions) GetCell(rowIndex int, columnName string) interface{} {
	if rowIndex >= 0 && rowIndex < len(sdk.rows) {
		return sdk.rows[rowIndex][columnName]
	}
	return nil
}

// GetCellFloat 安全获取浮点数
func (sdk *SDKFunctions) GetCellFloat(rowIndex int, columnName string) float64 {
	val := sdk.GetCell(rowIndex, columnName)
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	default:
		return 0.0
	}
}

// =============================
// Go代码解释器
// =============================

type FormulaInterpreter struct {
	interpreter *interp.Interpreter
	sdk         *SDKFunctions
}

func NewFormulaInterpreter(rows []map[string]interface{}, currentRowIndex int) *FormulaInterpreter {
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)

	sdk := &SDKFunctions{
		rows:         rows,
		currentRow:   rows[currentRowIndex],
		currentIndex: currentRowIndex,
	}

	// 注入SDK函数到解释器
	i.Use(interp.Exports{
		"sdk/sdk": map[string]reflect.Value{
			"Sum":          reflect.ValueOf(sdk.Sum),
			"SumIf":        reflect.ValueOf(sdk.SumIf),
			"CurrentRow":   reflect.ValueOf(sdk.CurrentRow),
			"GetCell":      reflect.ValueOf(sdk.GetCell),
			"GetCellFloat": reflect.ValueOf(sdk.GetCellFloat),
		},
	})

	return &FormulaInterpreter{
		interpreter: i,
		sdk:         sdk,
	}
}

// Execute 执行Go代码公式
func (fi *FormulaInterpreter) Execute(formulaCode string) (interface{}, error) {
	// 包装用户代码
	wrappedCode := `
package main

import "sdk/sdk"

func calculate() interface{} {
	` + formulaCode + `
}

var result = calculate()
`

	// 执行代码
	_, err := fi.interpreter.Eval(wrappedCode)
	if err != nil {
		return nil, fmt.Errorf("代码执行错误: %v", err)
	}

	// 获取结果
	v, err := fi.interpreter.Eval("result")
	if err != nil {
		return nil, fmt.Errorf("获取结果错误: %v", err)
	}

	return v.Interface(), nil
}

// =============================
// 语法糖转换器（可选）
// =============================

type SugarSyntax struct{}

func (ss *SugarSyntax) Transform(compactCode string) string {
	transformations := map[string]string{
		`@姓名`:          `currentRow()["姓名"]`,
		`@女工日结`:        `getCellFloat(currentIndex, "女工日结")`,
		`SUMIF.*`:      `sumIf`, // 需要更复杂的正则替换
		`currentIndex`: `currentRowIndex`,
	}

	result := compactCode
	for short, long := range transformations {
		result = strings.ReplaceAll(result, short, long)
	}
	return result
}

// =============================
// 表格集成
// =============================

func NewFormulaTable() *FormulaTable {
	return &FormulaTable{
		Rows:    make([]map[string]interface{}, 0),
		Columns: make(map[string]*ColumnConfig),
	}
}

// AddColumn 添加列定义
func (ft *FormulaTable) AddColumn(name, colType, formula string) {
	ft.Columns[name] = &ColumnConfig{
		Name:    name,
		Type:    colType,
		Formula: formula,
	}
}

// AddRow 添加行并自动计算公式列
func (ft *FormulaTable) AddRow(rowData map[string]interface{}) error {
	// 添加行数据
	ft.Rows = append(ft.Rows, rowData)
	rowIndex := len(ft.Rows) - 1

	// 计算所有公式列
	for colName, colConfig := range ft.Columns {
		if colConfig.Type == "formula" && colConfig.Formula != "" {
			err := ft.calculateFormulaColumn(colName, colConfig.Formula, rowIndex)
			if err != nil {
				return fmt.Errorf("列%s计算错误: %v", colName, err)
			}
		}
	}

	return nil
}

func (ft *FormulaTable) calculateFormulaColumn(colName, formula string, rowIndex int) error {
	interpreter := NewFormulaInterpreter(ft.Rows, rowIndex)

	result, err := interpreter.Execute(formula)
	if err != nil {
		// 存储错误信息
		ft.Rows[rowIndex][colName] = "错误: " + err.Error()
		return err
	}

	ft.Rows[rowIndex][colName] = result
	return nil
}

// =============================
// 使用示例
// =============================

func main() {
	// 创建表格
	table := NewFormulaTable()

	// 定义列
	table.AddColumn("姓名", "text", "")
	table.AddColumn("女工日结", "number", "")
	table.AddColumn("应发工资", "formula", `
name := currentRow()["姓名"].(string)
dailyWage := getCellFloat(currentIndex, "女工日结")

switch name {
case "拼车", "三人组":
	return 0.0
case "房东":
	return dailyWage
case "杨萍":
	threeGroupSum := sumIf("女工日结", "姓名", "三人组")
	return threeGroupSum/3 + dailyWage
case "二人组":
	threeGroupSum := sumIf("女工日结", "姓名", "三人组")
	return threeGroupSum/3 + dailyWage/2
default:
	return 0.0
}
`)

	// 添加测试数据
	table.AddRow(map[string]interface{}{"姓名": "拼车", "女工日结": 100})
	table.AddRow(map[string]interface{}{"姓名": "三人组", "女工日结": 300})
	table.AddRow(map[string]interface{}{"姓名": "房东", "女工日结": 150})
	table.AddRow(map[string]interface{}{"姓名": "杨萍", "女工日结": 120})
	table.AddRow(map[string]interface{}{"姓名": "二人组", "女工日结": 200})

	// 打印结果
	fmt.Println("计算结果:")
	for i, row := range table.Rows {
		fmt.Printf("行%d: 姓名=%s, 女工日结=%.1f, 应发工资=%.1f\n",
			i, row["姓名"], row["女工日结"], row["应发工资"])
	}

	// 测试单个公式计算
	fmt.Println("\n测试单个公式:")
	interpreter := NewFormulaInterpreter(table.Rows, 3) // 测试杨萍的行

	result, err := interpreter.Execute(`
name := currentRow()["姓名"].(string)
if name == "杨萍" {
	return "这是杨萍的测试"
}
return "不是杨萍"
`)

	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("结果:", result)
	}
}

// =============================
// 安卓端优化版本
// =============================

// AndroidFormulaInterpreter 安卓专用，带缓存优化
type AndroidFormulaInterpreter struct {
	baseInterpreter *FormulaInterpreter
	cache           map[string]interface{}
}

func NewAndroidFormulaInterpreter(rows []map[string]interface{}, currentRowIndex int) *AndroidFormulaInterpreter {
	return &AndroidFormulaInterpreter{
		baseInterpreter: NewFormulaInterpreter(rows, currentRowIndex),
		cache:           make(map[string]interface{}),
	}
}

func (ai *AndroidFormulaInterpreter) ExecuteWithCache(formula string) (interface{}, error) {
	// 检查缓存
	if cached, exists := ai.cache[formula]; exists {
		return cached, nil
	}

	// 执行并缓存结果
	result, err := ai.baseInterpreter.Execute(formula)
	if err == nil {
		ai.cache[formula] = result
	}

	return result, err
}
