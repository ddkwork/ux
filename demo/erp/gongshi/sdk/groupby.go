package sdk

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ddkwork/golibrary/std/stream/deepcopy"
	"github.com/ddkwork/ux/demo/erp/gongshi/sdk/field"
)

// SumIf sums values in a column where another column matches a value.
func (t *TreeTable) SumIf(filterColumn, filterValue, sumColumn string) float64 {
	total := 0.0
	for node := range t.DataNodes() {
		filterCell := node.GetCell(filterColumn)
		sumCell := node.GetCell(sumColumn)

		if filterCell != nil && sumCell != nil {
			if fmt.Sprint(filterCell.Value) == filterValue {
				if val, ok := ToFloat(sumCell.Value); ok {
					total += val
				}
			}
		}
	}
	return total
}

// GroupBy 按指定列分组
func (t *TreeTable) _GroupBy_old(columnName string) {
	// 获取所有行
	allRows := t.AllRows()
	if len(allRows) == 0 {
		return
	}

	// 按分组列的值排序，便于分组
	sort.Slice(allRows, func(i, j int) bool {
		cellI := allRows[i].GetCell(columnName)
		cellJ := allRows[j].GetCell(columnName)

		if cellI == nil && cellJ == nil {
			return false
		}
		if cellI == nil {
			return true
		}
		if cellJ == nil {
			return false
		}

		return fmt.Sprintf("%v", cellI.Value) < fmt.Sprintf("%v", cellJ.Value)
	})

	// 创建新的根容器
	root := newRoot()

	// 分组处理
	currentGroup := ""
	var currentGroupContainer *Node

	// 获取所有列配置
	columnConfigs := make(map[string]field.FieldType)
	for _, col := range t.Columns {
		columnConfigs[col.Name] = col.Type
	}

	for _, row := range allRows {
		cell := row.GetCell(columnName)
		groupValue := "未分组"
		if cell != nil {
			groupValue = fmt.Sprintf("%v", cell.Value)
		}

		// 如果分组值改变，创建新的分组容器
		if groupValue != currentGroup {
			currentGroup = groupValue

			// 创建包含所有列的单元格数据
			cells := make([]CellData, len(t.Columns))
			for i, col := range t.Columns {
				cellData := CellData{
					ColumnName: col.Name,
					Type:       col.Type,
				}

				if col.Name == columnName {
					// 分组列显示分组值
					cellData.Value = groupValue
				} else {
					// 其他列根据类型执行聚合
					switch col.Type {
					case field.NumberType:
						// 数字列求和
						sum := 0.0
						for _, r := range allRows {
							if fmt.Sprint(r.GetCell(columnName).Value) == groupValue {
								if v, ok := ToFloat(r.GetCell(col.Name).Value); ok {
									sum += v
								}
							}
						}
						cellData.Value = sum
					case field.TextType:
						// 文本列计数
						count := 0
						for _, r := range allRows {
							if fmt.Sprint(r.GetCell(columnName).Value) == groupValue {
								count++
							}
						}
						cellData.Value = count
					default:
						// 其他类型留空
						cellData.Value = ""
					}
				}

				cells[i] = cellData
			}

			currentGroupContainer = NewContainerNode("group", cells)
			currentGroupContainer.GroupKey = groupValue
			currentGroupContainer.isOpen = true
			root.AddChild(currentGroupContainer)
		}

		// 将行添加到当前分组
		row.parent = nil
		currentGroupContainer.AddChild(row)
	}

	// 更新根节点
	t.Root = root
	t.OriginalRoot = root.Clone() //todo remove
}

func (t *TreeTable) updateOriginalRoot() {
	t.OriginalRoot = deepcopy.Clone(t.Root) //to-do 增删改查调用它
}

// GroupBy 按指定列分组（带聚合信息和稳定顺序）
func (t *TreeTable) GroupBy(columnName string) {
	// 1. 获取所有行
	allRows := t.AllRows()
	if len(allRows) == 0 {
		return
	}

	// 2. 构建分组映射 (分组值 → 行节点列表)
	groupMap := make(map[string][]*Node)
	for _, row := range allRows {
		cell := row.GetCell(columnName)
		groupValue := "未分组"
		if cell != nil {
			groupValue = fmt.Sprintf("%v", cell.Value)
		}
		groupMap[groupValue] = append(groupMap[groupValue], row)
	}

	// 3. 创建新的根容器
	root := newRoot()

	// 4. 收集分组键并排序（确保稳定顺序）
	groupKeys := make([]string, 0, len(groupMap))
	for key := range groupMap {
		groupKeys = append(groupKeys, key)
	}
	sort.Strings(groupKeys) // 按字母顺序排序

	// 5. 为每个分组创建容器节点（带聚合信息）
	for _, groupKey := range groupKeys {
		groupRows := groupMap[groupKey]

		// 创建容器节点的行数据（包含所有列）
		cells := make([]CellData, len(t.Columns))
		for colIdx, col := range t.Columns {
			cellData := CellData{
				ColumnName: col.Name,
				Type:       col.Type,
			}

			if col.Name == columnName {
				// 分组列：显示分组键 + 行数
				cellData.Value = fmt.Sprintf("%s (%d)", groupKey, len(groupRows))
			} else {
				// 其他列：根据类型执行聚合
				switch col.Type {
				case field.NumberType:
					// 数字列求和
					sum := 0.0
					for _, row := range groupRows {
						if cell := row.GetCell(col.Name); cell != nil {
							if v, ok := ToFloat(cell.Value); ok {
								sum += v
							}
						}
					}
					cellData.Value = sum
				case field.TextType:
					// 文本列计数
					cellData.Value = len(groupRows)
				default:
					// 其他类型留空
					cellData.Value = ""
				}
			}
			cells[colIdx] = cellData
		}

		// 创建容器节点（使用聚合信息作为显示名称）
		displayName := fmt.Sprintf("%s (%d)", groupKey, len(groupRows))
		containerNode := NewContainerNode(displayName, cells)
		containerNode.GroupKey = groupKey
		containerNode.isOpen = true

		// 将行节点添加到容器
		for _, row := range groupRows {
			row.parent = nil // 解除原父子关系
			containerNode.AddChild(row)
		}

		// 添加容器到根节点
		root.AddChild(containerNode)
	}

	// 6. 更新表格结构
	//t.Root = root
	//t.OriginalRoot = root.Clone()
	t.groupedRows = root.Children
	t.Root = root
	// 如果有分组结果，展开所有节点以便显示分组内容
	if len(t.rootRows) > 0 {
		t.OpenAll()
	}
}

// GroupBy 按指定列分组（带聚合信息的Map实现）
func (t *TreeTable) _GroupBy__map(columnName string) {
	// 1. 获取所有行
	allRows := t.AllRows()
	if len(allRows) == 0 {
		return
	}

	// 2. 构建分组映射 (分组值 → 行节点列表)
	groupMap := make(map[string][]*Node)
	//groupMap := new(safemap.M[string, []*Node])
	for _, row := range allRows {
		cell := row.GetCell(columnName)
		groupValue := "未分组"
		if cell != nil {
			groupValue = fmt.Sprintf("%v", cell.Value)
		}
		groupMap[groupValue] = append(groupMap[groupValue], row)
	}

	// 3. 创建新的根容器
	root := newRoot()

	// 4. 为每个分组创建容器节点（带聚合信息）
	for groupKey, groupRows := range groupMap {
		// 创建容器节点的行数据（包含所有列）
		cells := make([]CellData, len(t.Columns))
		for colIdx, col := range t.Columns {
			cellData := CellData{
				ColumnName: col.Name,
				Type:       col.Type,
			}

			if col.Name == columnName {
				// 分组列：显示分组键 + 行数
				cellData.Value = fmt.Sprintf("%s (%d)", groupKey, len(groupRows))
			} else {
				// 其他列：根据类型执行聚合
				switch col.Type {
				case field.NumberType:
					// 数字列求和
					sum := 0.0
					for _, row := range groupRows {
						if cell := row.GetCell(col.Name); cell != nil {
							if v, ok := ToFloat(cell.Value); ok {
								sum += v
							}
						}
					}
					cellData.Value = sum
				case field.TextType:
					// 文本列计数
					cellData.Value = len(groupRows)
				default:
					// 其他类型留空
					cellData.Value = ""
				}
			}
			cells[colIdx] = cellData
		}

		// 创建容器节点（使用聚合信息作为显示名称）
		displayName := fmt.Sprintf("%s (%d)", groupKey, len(groupRows))
		containerNode := NewContainerNode(displayName, cells)
		containerNode.GroupKey = groupKey
		containerNode.isOpen = true

		// 将行节点添加到容器
		for _, row := range groupRows {
			row.parent = nil // 解除原父子关系
			containerNode.AddChild(row)
		}

		// 添加容器到根节点
		root.AddChild(containerNode)
	}

	// 5. 更新表格结构
	t.groupedRows = root.Children
	t.Root = root
	// 如果有分组结果，展开所有节点以便显示分组内容
	if len(t.rootRows) > 0 {
		t.OpenAll()
	}
}

// Aggregate 对分组进行聚合计算
func (t *TreeTable) Aggregate(groupColumn, targetColumn, aggType string) map[string]float64 {
	result := make(map[string]float64)

	// 遍历所有分组容器
	for _, node := range t.Root.Children {
		if node.IsContainer() { //&& strings.HasPrefix(node.Type, "group")
			groupKey := node.GroupKey
			if groupKey == "" {
				if cell := node.GetCell(groupColumn); cell != nil {
					groupKey = fmt.Sprintf("%v", cell.Value)
				} else {
					groupKey = "未分组"
				}
			}

			var aggregateValue float64
			count := 0

			// 遍历分组内的所有行
			for row := range node.Walk() {
				if row.IsContainer() {
					continue
				}

				if cell := row.GetCell(targetColumn); cell != nil {
					if val, ok := ToFloat(cell.Value); ok {
						switch aggType {
						case "sum":
							aggregateValue += val
						case "avg":
							aggregateValue += val
							count++
						case "max":
							if count == 0 || val > aggregateValue {
								aggregateValue = val
							}
							count = 1
						case "min":
							if count == 0 || val < aggregateValue {
								aggregateValue = val
							}
							count = 1
						}
					}
				}
			}

			// 处理平均值
			if aggType == "avg" && count > 0 {
				aggregateValue /= float64(count)
			}

			if aggType == "count" {
				rowCount := 0
				for range node.Walk() {
					if !node.IsContainer() {
						rowCount++
					}
				}
				result[groupKey] = float64(rowCount)
			} else {
				result[groupKey] = aggregateValue
			}
		}
	}

	return result
}

// GetGroups 获取所有分组
func (t *TreeTable) GetGroups() []*Node {
	var groups []*Node
	for _, child := range t.Root.Children {
		if child.IsContainer() { //&& strings.HasPrefix(child.Type, "group")
			groups = append(groups, child)
		}
	}
	return groups
}

// ExpandAllGroups 展开所有分组
func (t *TreeTable) ExpandAllGroups() {
	for node := range t.Root.Walk() {
		if node.IsContainer() {
			node.isOpen = true
		}
	}
}

// CollapseAllGroups 折叠所有分组
func (t *TreeTable) CollapseAllGroups() {
	for node := range t.Root.Walk() {
		if node.IsContainer() && !strings.HasPrefix(node.Type, "root") {
			node.isOpen = false
		}
	}
}
func (t *TreeTable) OpenAll() {
	if t.filteredRows != nil {
		for _, row := range t.filteredRows {
			row.OpenAll()
		}
		return
	}
	t.Root.OpenAll()
}

func (t *TreeTable) CloseAll() {
	if t.filteredRows != nil {
		for _, row := range t.filteredRows {
			row.CloseAll()
		}
		return
	}
	t.Root.CloseAll()
}

func (n *Node) OpenAll() {
	for _, node := range n.WalkContainer() {
		node.SetOpen(true)
	}
}

func (n *Node) CloseAll() {
	for _, node := range n.WalkContainer() {
		node.SetOpen(false)
	}
}

// Ungroup 取消分组，回到平面结构
func (t *TreeTable) Ungroup() {
	root := newRoot()

	// 将所有行提取到根节点
	for _, node := range t.Root.Children {
		for row := range node.Walk() {
			if !row.IsContainer() {
				row.parent = nil
				root.AddChild(row)
			}
		}
	}

	t.Root = root
	t.OriginalRoot = root.Clone()
}
