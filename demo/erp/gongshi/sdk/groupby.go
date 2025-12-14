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
func (t *TreeTable) GroupBy(columnName string) error {
	// 获取所有行
	allRows := t.AllRows()
	if len(allRows) == 0 {
		return nil
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

	for _, row := range allRows {
		cell := row.GetCell(columnName)
		groupValue := "未分组"
		if cell != nil {
			groupValue = fmt.Sprintf("%v", cell.Value)
		}

		// 如果分组值改变，创建新的分组容器
		if groupValue != currentGroup {
			currentGroup = groupValue
			currentGroupContainer = NewContainerNode("group", []CellData{
				{ColumnName: columnName, Value: groupValue, Type: field.TextType},
			})
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
	t.OriginalRoot = root.Clone()

	return nil
}

func (t *TreeTable) updateOriginalRoot() {
	t.OriginalRoot = deepcopy.Clone(t.Root) //to-do 增删改查调用它
}

func (t *TreeTable) GroupBy2(field string) {
	// 按照字段名称找到分组的列
	//mylog.CheckNil(t.OriginalRoot)
	//t.Root = t.OriginalRoot
	//t.groupedRows = nil
	//groupColumIndex := -1
	//for i, cell := range t.header.columnCells {
	//	if field == cell.Name {
	//		groupColumIndex = i
	//		mylog.Trace(i, "即将按"+strconv.Quote(cell.Name)+"列进行分组 ")
	//		break
	//	}
	//}
	//if groupColumIndex == -1 {
	//	mylog.Check("未找到分组列,请检查输入的字段名称是否正确")
	//	return
	//}
	//
	//// 构建分组映射
	//groupMap := make(map[string][]*Node)
	//for _, n := range t.RootRows() {
	//	cells := t.MarshalRowCells(n)
	//	// 检查列索引是否有效
	//	if groupColumIndex >= len(cells) {
	//		continue
	//	}
	//	// 获取指定列的值
	//	groupValue := cells[groupColumIndex].Value
	//	// 将节点添加到对应的分组中
	//	groupMap[groupValue] = append(groupMap[groupValue], n)
	//}
	//
	//// 创建新的根节点
	//var zero T
	//root := newRoot(zero)
	//for groupKey, groupNodes := range groupMap {
	//	// 创建容器节点
	//	containerNode := NewContainerNode(groupKey, zero)
	//	for _, node := range groupNodes {
	//		// 将节点添加到容器节点中
	//		containerNode.AddChild(node)
	//	}
	//	// 将容器节点添加到新的根节点中
	//	root.AddChild(containerNode)
	//}

	// 更新树形表格的根节点和根行
	//t.groupedRows = root.Children
	//t.Root = root
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
		if node.IsContainer() && strings.HasPrefix(node.Type, "group") {
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
		if child.IsContainer() && strings.HasPrefix(child.Type, "group") {
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
