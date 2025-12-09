package main

import (
	"fmt"
	"sync"

	"github.com/ddkwork/golibrary/std/stream"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// ==================== 内存表格 ====================
type MemoryTable struct {
	mu   sync.RWMutex
	rows []map[string]interface{}
}

func NewMemoryTable() *MemoryTable {
	return &MemoryTable{rows: make([]map[string]interface{}, 0)}
}

func (t *MemoryTable) AddRow(row map[string]interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.rows = append(t.rows, row)
}

func (t *MemoryTable) RowCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.rows)
}

func (t *MemoryTable) GetRow(i int) map[string]interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if i < 0 || i >= len(t.rows) {
		return nil
	}
	row := make(map[string]interface{})
	for k, v := range t.rows[i] {
		row[k] = v
	}
	return row
}

func (t *MemoryTable) SetValue(i int, col string, val interface{}) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if i < 0 || i >= len(t.rows) {
		return fmt.Errorf("索引越界")
	}
	t.rows[i][col] = val
	return nil
}

func (t *MemoryTable) SumIf(field string, crit interface{}, sumField string) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var sum float64
	for _, r := range t.rows {
		if v, ok := r[field]; ok && fmt.Sprint(v) == fmt.Sprint(crit) {
			if sv, ok := r[sumField].(float64); ok {
				sum += sv
			} else if iv, ok := r[sumField].(int); ok {
				sum += float64(iv)
			}
		}
	}
	return sum
}

// ==================== Yaegi 引擎 ====================
type YaegiEngine struct{}

func NewYaegiEngine() *YaegiEngine {
	return &YaegiEngine{}
}

// 计算单行公式 - 真正的脚本化版本
func (e *YaegiEngine) CalculateRow(script string, row map[string]interface{}, tableData []map[string]interface{}) (float64, error) {
	// 创建新的解释器实例
	i := interp.New(interp.Options{
		GoPath:       "./",
		Unrestricted: true,
	})
	i.Use(stdlib.Symbols)

	// 构建完整的 Go 代码
	code := fmt.Sprintf(`
		package main
		
		import "fmt"
		
		// 当前行数据
		var row = %#v
		
		// 表格数据
		var tableData = %#v
		
		// SUMIF 函数
		func SUMIF(field string, criteria interface{}, sumField string) float64 {
			var sum float64
			for _, r := range tableData {
				if fmt.Sprint(r[field]) == fmt.Sprint(criteria) {
					if val, ok := r[sumField].(float64); ok {
						sum += val
					} else if val, ok := r[sumField].(int); ok {
						sum += float64(val)
					}
				}
			}
			return sum
		}
		
		// 执行用户脚本
		func calc() float64 {
			%s
		}
	`, row, tableData, script) //妙，调用sdk api动态生成代码运算的数据而不是使用反射导出sdk api方法，那样会得到很多错误
	//思考，有没有别的场景不是传递这些参数呢？

	stream.WriteGoFile("tmp/main.go", code) //执行goland语法检查
	// 执行代码
	v, err := i.Eval(code)
	if err != nil {
		return 0, fmt.Errorf("执行失败: %v", err)
	}

	// 获取结果
	result, ok := v.Interface().(float64)
	if !ok {
		return 0, fmt.Errorf("结果类型错误")
	}

	return result, nil
}

// ==================== 主程序 ====================
func main() {
	// 创建表格
	table := NewMemoryTable()

	// 添加数据
	table.AddRow(map[string]interface{}{"姓名": "拼车", "女工日结": 0.0})
	table.AddRow(map[string]interface{}{"姓名": "三人组", "女工日结": 900.0})
	table.AddRow(map[string]interface{}{"姓名": "房东", "女工日结": 350.0})
	table.AddRow(map[string]interface{}{"姓名": "杨萍", "女工日结": 200.0})
	table.AddRow(map[string]interface{}{"姓名": "二人组", "女工日结": 600.0})

	// 创建 Yaegi 引擎
	engine := NewYaegiEngine()

	// 真正的脚本化公式 - 完全写在字符串中
	script := `
		// 获取当前行数据
		name := row["姓名"].(string)
		
		// 安全获取女工日结值
		var nvGong float64
		switch v := row["女工日结"].(type) {
		case float64:
			nvGong = v
		case int:
			nvGong = float64(v)
		case int64:
			nvGong = float64(v)
		default:
			nvGong = 0.0
		}
		
		// 计算三人组总和
		sanRenZuSum := SUMIF("姓名", "三人组", "女工日结")
		
		// 使用 Go 的 switch 语句
		switch name {
		case "拼车", "三人组":
			return 0.0
		case "房东":
			return nvGong
		case "杨萍":
			return (sanRenZuSum / 3.0) + nvGong
		case "二人组":
			return (sanRenZuSum / 3.0) + (nvGong / 2.0)
		default:
			return 0.0
		}
	`

	// 计算
	fmt.Println("=== 计算结果 ===")
	fmt.Printf("%-10s | %-10s | %-10s\n", "姓名", "女工日结", "计算结果")
	fmt.Println("-----------|------------|------------")

	// 获取表格数据

	for i := 0; i < table.RowCount(); i++ {
		row := table.GetRow(i)

		// 使用 Yaegi 计算
		result, err := engine.CalculateRow(script, row, table.rows)
		if err != nil {
			fmt.Printf("行 %d 错误: %v\n", i, err)
			result = 0.0
		}

		table.SetValue(i, "计算结果", result)

		fmt.Printf("%-10s | %-10.0f | %-10.2f\n",
			row["姓名"], row["女工日结"], result)
	}

	// 验证
	sanRenZuSum := table.SumIf("姓名", "三人组", "女工日结")
	fmt.Printf("\n验证: 三人组总和=%.0f\n", sanRenZuSum)
	fmt.Printf("杨萍应得: %.0f/3 + 200 = %.0f\n", sanRenZuSum, sanRenZuSum/3+200)
	fmt.Printf("二人组应得: %.0f/3 + 600/2 = %.0f\n", sanRenZuSum, sanRenZuSum/3+300)
}
