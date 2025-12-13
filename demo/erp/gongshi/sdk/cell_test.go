package sdk

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/ddkwork/golibrary/std/assert"
	"github.com/ddkwork/golibrary/std/stream/uuid"
	"github.com/ddkwork/ux/demo/erp/gongshi/sdk/field"
)

func TestCellData_AsBool(t *testing.T) {
	type fields struct {
		ColumnName string
		Value      any
		Type       field.FieldType
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CellData{
				ColumnName: tt.fields.ColumnName,
				Value:      tt.fields.Value,
				Type:       tt.fields.Type,
			}
			if got := c.AsBool(); got != tt.want {
				t.Errorf("AsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCellData_AsFloat(t *testing.T) {
	type fields struct {
		ColumnName string
		Value      any
		Type       field.FieldType
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CellData{
				ColumnName: tt.fields.ColumnName,
				Value:      tt.fields.Value,
				Type:       tt.fields.Type,
			}
			if got := c.AsFloat(); got != tt.want {
				t.Errorf("AsFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCellData_AsInt(t *testing.T) {
	type fields struct {
		ColumnName string
		Value      any
		Type       field.FieldType
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CellData{
				ColumnName: tt.fields.ColumnName,
				Value:      tt.fields.Value,
				Type:       tt.fields.Type,
			}
			if got := c.AsInt(); got != tt.want {
				t.Errorf("AsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCellData_AsString(t *testing.T) {
	type fields struct {
		ColumnName string
		Value      any
		Type       field.FieldType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CellData{
				ColumnName: tt.fields.ColumnName,
				Value:      tt.fields.Value,
				Type:       tt.fields.Type,
			}
			if got := c.AsString(); got != tt.want {
				t.Errorf("AsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCellData_AsTime(t *testing.T) {
	type fields struct {
		ColumnName string
		Value      any
		Type       field.FieldType
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CellData{
				ColumnName: tt.fields.ColumnName,
				Value:      tt.fields.Value,
				Type:       tt.fields.Type,
			}
			if got := c.AsTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCellData_DetectValueType(t *testing.T) {
	type fields struct {
		ColumnName string
		Value      any
		Type       field.FieldType
	}
	tests := []struct {
		name   string
		fields fields
		want   field.FieldType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CellData{
				ColumnName: tt.fields.ColumnName,
				Value:      tt.fields.Value,
				Type:       tt.fields.Type,
			}
			if got := c.DetectValueType(); got != tt.want {
				t.Errorf("DetectValueType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCellData_IsAttachment(t *testing.T) {
	type fields struct {
		ColumnName string
		Value      any
		Type       field.FieldType
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CellData{
				ColumnName: tt.fields.ColumnName,
				Value:      tt.fields.Value,
				Type:       tt.fields.Type,
			}
			if got := c.IsAttachment(); got != tt.want {
				t.Errorf("IsAttachment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCellData_IsFormula(t *testing.T) {
	type fields struct {
		ColumnName string
		Value      any
		Type       field.FieldType
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CellData{
				ColumnName: tt.fields.ColumnName,
				Value:      tt.fields.Value,
				Type:       tt.fields.Type,
			}
			if got := c.IsFormula(); got != tt.want {
				t.Errorf("IsFormula() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCellData_IsLink(t *testing.T) {
	type fields struct {
		ColumnName string
		Value      any
		Type       field.FieldType
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CellData{
				ColumnName: tt.fields.ColumnName,
				Value:      tt.fields.Value,
				Type:       tt.fields.Type,
			}
			if got := c.IsLink(); got != tt.want {
				t.Errorf("IsLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCellData_IsSelect(t *testing.T) {
	type fields struct {
		ColumnName string
		Value      any
		Type       field.FieldType
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CellData{
				ColumnName: tt.fields.ColumnName,
				Value:      tt.fields.Value,
				Type:       tt.fields.Type,
			}
			if got := c.IsSelect(); got != tt.want {
				t.Errorf("IsSelect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_GetCell(t *testing.T) {
	type fields struct {
		ID        uuid.ID
		Type      string
		RowCells  []CellData
		Children  []*Node
		parent    *Node
		isOpen    bool
		GroupKey  string
		RowNumber int
	}
	type args struct {
		colName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CellData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				ID:        tt.fields.ID,
				Type:      tt.fields.Type,
				RowCells:  tt.fields.RowCells,
				Children:  tt.fields.Children,
				parent:    tt.fields.parent,
				isOpen:    tt.fields.isOpen,
				GroupKey:  tt.fields.GroupKey,
				RowNumber: tt.fields.RowNumber,
			}
			if got := n.GetCell(tt.args.colName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCell() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_SetCellValue(t *testing.T) {
	type fields struct {
		ID        uuid.ID
		Type      string
		RowCells  []CellData
		Children  []*Node
		parent    *Node
		isOpen    bool
		GroupKey  string
		RowNumber int
	}
	type args struct {
		colName string
		value   any
		table   *TreeTable
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{
				ID:        tt.fields.ID,
				Type:      tt.fields.Type,
				RowCells:  tt.fields.RowCells,
				Children:  tt.fields.Children,
				parent:    tt.fields.parent,
				isOpen:    tt.fields.isOpen,
				GroupKey:  tt.fields.GroupKey,
				RowNumber: tt.fields.RowNumber,
			}
			n.SetCellValue(tt.args.colName, tt.args.value, tt.args.table)
		})
	}
}
func TestTreeTable_GetCellByRowIndex(t1 *testing.T) {
	t := tableDemo()
	assert.Equal(t1, "三人组", t.GetCellByRowIndex(0, "姓名").AsString())
	assert.Equal(t1, 2966.30, t.GetCellByRowIndex(0, "女工日结").AsFloat())
}

func TestTreeTable_SetCellValue(t1 *testing.T) {
	type fields struct {
		Root             *Node
		OriginalRoot     *Node
		Columns          []ColumnDefinition
		columnMap        map[string]*ColumnDefinition
		SelectedNode     *Node
		once             sync.Once
		OnRowSelected    func(n *Node)
		OnRowDoubleClick func(n *Node)
		filteredRows     []*Node
		groupedRows      []*Node
		rootRows         []*Node
	}
	type args struct {
		rowIndex int
		colName  string
		value    any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TreeTable{
				Root:             tt.fields.Root,
				OriginalRoot:     tt.fields.OriginalRoot,
				Columns:          tt.fields.Columns,
				columnMap:        tt.fields.columnMap,
				SelectedNode:     tt.fields.SelectedNode,
				once:             tt.fields.once,
				OnRowSelected:    tt.fields.OnRowSelected,
				OnRowDoubleClick: tt.fields.OnRowDoubleClick,
				filteredRows:     tt.fields.filteredRows,
				groupedRows:      tt.fields.groupedRows,
				rootRows:         tt.fields.rootRows,
			}
			if got := t.SetCellValue(tt.args.rowIndex, tt.args.colName, tt.args.value); got != tt.want {
				t1.Errorf("SetCellValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_detectTypeFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want field.FieldType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectTypeFromString(tt.args.s); got != tt.want {
				t.Errorf("detectTypeFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inferType(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want field.FieldType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inferType(tt.args.v); got != tt.want {
				t.Errorf("inferType() = %v, want %v", got, tt.want)
			}
		})
	}
}
