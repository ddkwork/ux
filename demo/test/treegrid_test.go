package main

import (
	"testing"
)

func TestCalculateLevelColumnWidth(t *testing.T) {
	nodes := []*Node{
		{RowCells: []string{"Root", "Value1", "Description1"}},
		{RowCells: []string{"Child1", "Value2", "Description2"}},
		{RowCells: []string{"Child2Long", "Value3", "Description3"}},
	}

	maxWidths := CalculateLevelColumnMaxCellWidth(nodes)

	if maxWidths[0] != len("Child2Long") {
		t.Errorf("Expected max width for first column to be %d, got %d", len("Child2Long"), maxWidths[0])
	}
}

func TestCalculateAverageNonLevelWidth(t *testing.T) {
	nodes := []*Node{
		{RowCells: []string{"Root", "Value1", "Description1"}},
		{RowCells: []string{"Child1", "Value2", "Description2"}},
		{RowCells: []string{"Child2Long", "Value3", "Description3"}},
	}

	average := CalculateAverageNonLevelWidth(nodes)

	expected := (len("Value1") + len("Description1") + len("Value2") + len("Description2") + len("Value3") + len("Description3")) / 6
	if average != expected {
		t.Errorf("Expected average width to be %d, got %d", expected, average)
	}
}

func TestAdjustLevelColumn(t *testing.T) {
	// 创建 Node 实例
	nodes := []*Node{
		{RowCells: []string{"Root", "Value1", "Description1"}},
		{RowCells: []string{"Child1", "Value2", "Description2"}},
		{RowCells: []string{"Child2LongName", "Value3", "Description3"}},
	}

	// 获取最大列宽
	maxWidths := CalculateLevelColumnMaxCellWidth(nodes)

	// 使用动态获取的最大宽度进行测试
	tests := []struct {
		currentWidth  int
		maxWidth      int
		depth         int
		indentWidth   int
		expectedWidth int
		expectedFill  int
	}{
		{len(nodes[0].RowCells[0]), maxWidths[0], 0, indentBase, maxWidths[0], maxWidths[0] - (len(nodes[0].RowCells[0]) + (0 * indentBase))}, // Root
		{len(nodes[1].RowCells[0]), maxWidths[0], 1, indentBase, maxWidths[0], maxWidths[0] - (len(nodes[1].RowCells[0]) + (1 * indentBase))}, // Child1
		{len(nodes[2].RowCells[0]), maxWidths[0], 2, indentBase, maxWidths[0], maxWidths[0] - (len(nodes[1].RowCells[0]) + (2 * indentBase))}, // Child2LongName
	}

	for _, test := range tests {
		actualWidth, actualFill := AdjustLevelColumn(test.currentWidth, test.maxWidth, test.depth, test.indentWidth)

		if actualWidth != test.expectedWidth {
			t.Errorf("Expected width %d, got %d for currentWidth %d",
				test.expectedWidth, actualWidth, test.currentWidth)
		}

		if actualFill != test.expectedFill {
			t.Errorf("Expected fill %d, got %d for currentWidth %d",
				test.expectedFill, actualFill, test.currentWidth)
		}
	}
}
