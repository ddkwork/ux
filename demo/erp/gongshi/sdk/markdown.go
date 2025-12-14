package sdk

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ddkwork/golibrary/std/stream"
)

// ToMarkdown exports the table to Markdown format.
func (t *TreeTable) ToMarkdown(name string) string {
	var sb strings.Builder
	sb.WriteString("# Tree Table Structure\n\n")

	// Write header
	sb.WriteString("| ")
	for _, col := range t.Columns {
		sb.WriteString(fmt.Sprintf("%s | ", col.Name))
	}
	sb.WriteString("\n|")
	for range t.Columns {
		sb.WriteString("--------|")
	}
	sb.WriteString("\n")

	// Write data rows
	for node := range t.DataNodes() {
		// Calculate relative depth (root is depth 0)
		depth := max(node.Depth()-1, 0)

		// Create indentation and icon
		indent := strings.Repeat("&nbsp;&nbsp;&nbsp;", depth)
		icon := "📄"
		if node.IsContainer() {
			if node.isOpen {
				icon = "📂"
			} else {
				icon = "📁"
			}
		}

		sb.WriteString("| ")

		// Merge icon and first column
		if len(t.Columns) > 0 {
			firstCol := t.Columns[0]
			cell := node.GetCell(firstCol.Name)
			value := "-"
			if cell != nil {
				value = fmt.Sprintf("%v", cell.Value)
			}
			// Combine icon, indent and value
			sb.WriteString(fmt.Sprintf("%s%s %s", indent, icon, value))

			// Add remaining columns
			for i := 1; i < len(t.Columns); i++ {
				col := t.Columns[i]
				//println(col.Name)
				cell := node.GetCell(col.Name)
				value := "-"
				if cell != nil {
					value = fmt.Sprintf("%v", cell.Value)
				}
				sb.WriteString(fmt.Sprintf(" | %s", value))
			}
		} else {
			// No columns case
			sb.WriteString(fmt.Sprintf("%s%s", indent, icon))
		}

		sb.WriteString(" |\n")
	}

	if !stream.IsAndroid() {
		stream.WriteTruncate(filepath.Join("tmp", name+".md"), sb.String())
	}
	return sb.String()
}
