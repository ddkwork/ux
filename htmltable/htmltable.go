package htmltable

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ddkwork/golibrary/std/mylog"
)

// GeneratedFile 封装文件写入
type GeneratedFile struct {
	*os.File
}

// PKeepSpace 保留空格的写入
func (g *GeneratedFile) PKeepSpace(v ...any) {
	g.p(true, v...)
}

// P 不保留空格的写入（自动修剪）
func (g *GeneratedFile) P(v ...any) {
	g.p(false, v...)
}

// p 底层写入方法
func (g *GeneratedFile) p(keepSpace bool, v ...any) {
	for _, x := range v {
		s, ok := x.(string)
		if ok {
			if strings.Contains(s, "\n") {
				if !keepSpace {
					s = strings.TrimSpace(s)
				}
				s = strings.TrimPrefix(s, "\n")
				x = s
			}
		}
		fmt.Fprint(g, x)
	}
	fmt.Fprintln(g)
}

// 列定义
type Column struct {
	Key      string
	Label    string
	Width    string
	Sortable bool
}

// 树节点
type TreeNode struct {
	ID       string
	Label    string
	Data     map[string]string
	Children []TreeNode
	Expanded bool
}

// 表格配置
type TableConfig struct {
	Title      string
	FileName   string
	HeaderText string
	Columns    []Column
	TreeData   []TreeNode
}

// 生成完整的HTML文件
func GenerateHTMLFile(config TableConfig, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	gf := &GeneratedFile{file}

	// 生成HTML内容
	generateHTMLContent(gf, config)

	return nil
}

// 生成HTML内容
func generateHTMLContent(gf *GeneratedFile, config TableConfig) {
	// DOCTYPE和HTML开始
	gf.P(`<!DOCTYPE html>`)
	gf.P(`<html lang="zh-CN">`)
	gf.P(`<head>`)

	// Meta标签和标题
	gf.PKeepSpace(`    <meta charset="UTF-8">`)
	gf.PKeepSpace(`    <meta name="viewport" content="width=device-width, initial-scale=1.0">`)
	gf.PKeepSpace(`    <title>`, config.Title, `</title>`)

	// 样式
	generateStyles(gf)

	gf.P(`</head>`)
	gf.P(`<body>`)

	// 容器和内容
	generateBodyContent(gf, config)

	// 脚本
	generateScripts(gf, config)

	gf.P(`</body>`)
	gf.P(`</html>`)
}

// 生成CSS样式
func generateStyles(gf *GeneratedFile) {
	gf.PKeepSpace(`    <style>`)
	gf.P(`        * {`)
	gf.P(`            box-sizing: border-box;`)
	gf.P(`            margin: 0;`)
	gf.P(`            padding: 0;`)
	gf.P(`            font-family: 'Segoe UI', 'Microsoft YaHei', sans-serif;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        body {`)
	gf.P(`            background: #1a1a1a;`)
	gf.P(`            min-height: 100vh;`)
	gf.P(`            display: flex;`)
	gf.P(`            flex-direction: column;`)
	gf.P(`            align-items: center;`)
	gf.P(`            padding: 1rem;`)
	gf.P(`            color: #e0e0e0;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .container {`)
	gf.P(`            width: 100%;`)
	gf.P(`            max-width: 1000px;`)
	gf.P(`            background: #2d2d2d;`)
	gf.P(`            border-radius: 8px;`)
	gf.P(`            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);`)
	gf.P(`            overflow: hidden;`)
	gf.P(`            margin-bottom: 1rem;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .search-bar {`)
	gf.P(`            background: #333;`)
	gf.P(`            padding: 0.8rem 1rem;`)
	gf.P(`            display: flex;`)
	gf.P(`            align-items: center;`)
	gf.P(`            border-bottom: 1px solid #444;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .search-input {`)
	gf.P(`            flex-grow: 1;`)
	gf.P(`            background: #2d2d2d;`)
	gf.P(`            border: 1px solid #444;`)
	gf.P(`            border-radius: 20px;`)
	gf.P(`            padding: 0.5rem 1rem;`)
	gf.P(`            color: #e0e0e0;`)
	gf.P(`            font-size: 0.9rem;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .tools-btn {`)
	gf.P(`            margin-left: 0.8rem;`)
	gf.P(`            background: #4CAF50;`)
	gf.P(`            color: white;`)
	gf.P(`            border: none;`)
	gf.P(`            border-radius: 4px;`)
	gf.P(`            padding: 0.5rem 1rem;`)
	gf.P(`            font-size: 0.9rem;`)
	gf.P(`            cursor: pointer;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .header {`)
	gf.P(`            background: #333;`)
	gf.P(`            padding: 0.8rem 1rem;`)
	gf.P(`            border-bottom: 1px solid #444;`)
	gf.P(`            font-size: 0.85rem;`)
	gf.P(`            color: #aaa;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .table-container {`)
	gf.P(`            overflow-x: auto;`)
	gf.P(`            max-height: 500px;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .tree-table {`)
	gf.P(`            width: 100%;`)
	gf.P(`            border-collapse: collapse;`)
	gf.P(`            font-size: 0.8rem;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .tree-table th, .tree-table td {`)
	gf.P(`            border-right: 1px solid #444;`)
	gf.P(`            border-bottom: 1px solid #444;`)
	gf.P(`            padding: 0.2rem 0.5rem;`)
	gf.P(`            height: 24px;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .tree-table th:last-child, .tree-table td:last-child {`)
	gf.P(`            border-right: none;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .tree-table th {`)
	gf.P(`            background-color: #333;`)
	gf.P(`            font-weight: 600;`)
	gf.P(`            color: #e0e0e0;`)
	gf.P(`            text-align: left;`)
	gf.P(`            cursor: pointer;`)
	gf.P(`            user-select: none;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .tree-table th.sortable:hover {`)
	gf.P(`            background-color: #3a3a3a;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .tree-table td {`)
	gf.P(`            vertical-align: middle;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .tree-table tbody tr:nth-child(even) {`)
	gf.P(`            background-color: #2d2d2d;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .tree-table tbody tr:nth-child(odd) {`)
	gf.P(`            background-color: #262626;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .tree-row {`)
	gf.P(`            transition: background-color 0.2s;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .tree-row:hover {`)
	gf.P(`            background-color: #3a3a3a !important;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .tree-cell {`)
	gf.P(`            display: flex;`)
	gf.P(`            align-items: center;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .tree-indent {`)
	gf.P(`            display: inline-block;`)
	gf.P(`            width: 16px;`)
	gf.P(`            height: 100%;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .toggle-btn {`)
	gf.P(`            width: 18px;`)
	gf.P(`            height: 18px;`)
	gf.P(`            border: none;`)
	gf.P(`            background: transparent;`)
	gf.P(`            display: inline-flex;`)
	gf.P(`            align-items: center;`)
	gf.P(`            justify-content: center;`)
	gf.P(`            margin-right: 4px;`)
	gf.P(`            cursor: pointer;`)
	gf.P(`            transition: all 0.2s;`)
	gf.P(`            padding: 0;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .toggle-btn svg {`)
	gf.P(`            width: 16px;`)
	gf.P(`            height: 16px;`)
	gf.P(`            fill: #ffa726;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .toggle-btn:hover svg {`)
	gf.P(`            fill: #ffb74d;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .no-toggle {`)
	gf.P(`            width: 18px;`)
	gf.P(`            display: inline-block;`)
	gf.P(`            margin-right: 2px;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .no-toggle svg {`)
	gf.P(`            width: 16px;`)
	gf.P(`            height: 16px;`)
	gf.P(`            fill: #90a4ae;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .controls {`)
	gf.P(`            display: flex;`)
	gf.P(`            justify-content: space-between;`)
	gf.P(`            padding: 0.8rem 1rem;`)
	gf.P(`            background: #333;`)
	gf.P(`            border-top: 1px solid #444;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        button {`)
	gf.P(`            padding: 0.3rem 0.7rem;`)
	gf.P(`            background: #555;`)
	gf.P(`            color: #e0e0e0;`)
	gf.P(`            border: none;`)
	gf.P(`            border-radius: 3px;`)
	gf.P(`            cursor: pointer;`)
	gf.P(`            transition: background 0.2s;`)
	gf.P(`            font-weight: 500;`)
	gf.P(`            font-size: 0.8rem;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        button:hover {`)
	gf.P(`            background: #666;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .expand-all {`)
	gf.P(`            background: #388e3c;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .expand-all:hover {`)
	gf.P(`            background: #43a047;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .collapse-all {`)
	gf.P(`            background: #d32f2f;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .collapse-all:hover {`)
	gf.P(`            background: #f44336;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .filename {`)
	gf.P(`            background: #333;`)
	gf.P(`            padding: 0.5rem 1rem;`)
	gf.P(`            color: #aaa;`)
	gf.P(`            font-size: 0.8rem;`)
	gf.P(`            border-top: 1px solid #444;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .status-active { color: #4CAF50; }`)
	gf.P(`        .status-inactive { color: #f44336; }`)
	gf.P(`        .status-pending { color: #ff9800; }`)
	gf.P(``)
	gf.P(`        .sort-indicator {`)
	gf.P(`            margin-left: 5px;`)
	gf.P(`            font-size: 0.7rem;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .sort-asc::after {`)
	gf.P(`            content: " ↑";`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        .sort-desc::after {`)
	gf.P(`            content: " ↓";`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        @media (max-width: 768px) {`)
	gf.P(`            body { padding: 0.5rem; }`)
	gf.P(`            .tree-table { font-size: 0.75rem; }`)
	gf.P(`            .tree-table th, .tree-table td { padding: 0.15rem 0.3rem; }`)
	gf.P(`            .tree-indent { width: 12px; }`)
	gf.P(`        }`)
	gf.PKeepSpace(`    </style>`)
}

// 生成页面主体内容
func generateBodyContent(gf *GeneratedFile, config TableConfig) {
	gf.P(`    <div class="container">`)
	gf.P(`        <div class="search-bar">`)
	gf.P(`            <input type="text" class="search-input" placeholder="搜索..." id="searchInput">`)
	gf.P(`            <button class="tools-btn" id="searchBtn">全局搜索</button>`)
	gf.P(`        </div>`)
	gf.P(``)
	gf.P(`        <div class="header">`)
	gf.P(`            `, config.FileName)
	gf.P(`            <span style="float: right; font-size: 0.7rem;">`, config.HeaderText, `。 生成时间: `+time.Now().Format("2006-01-02 15:04:05"), `</span>`)
	gf.P(`        </div>`)
	gf.P(``)
	gf.P(`        <div class="table-container">`)
	gf.P(`            <table class="tree-table">`)
	gf.P(`                <thead id="tableHeader">`)
	gf.P(`                    <tr>`)

	// 生成表头
	for _, col := range config.Columns {
		widthAttr := ""
		if col.Width != "" {
			widthAttr = ` style="width: ` + col.Width + `"`
		}
		sortClass := ""
		if col.Sortable {
			sortClass = ` class="sortable"`
		}
		gf.P(`                        <th`, sortClass, widthAttr, ` data-key="`, col.Key, `" data-sort="none">`, col.Label, `</th>`)
	}

	gf.P(`                    </tr>`)
	gf.P(`                </thead>`)
	gf.P(`                <tbody id="treeTableBody">`)
	gf.P(`                    <!-- 表格内容将通过JavaScript动态生成 -->`)
	gf.P(`                </tbody>`)
	gf.P(`            </table>`)
	gf.P(`        </div>`)
	gf.P(``)
	//gf.P(`        <div class="filename">`, config.FileName, `</div>`)
	gf.P(``)
	gf.P(`        <div class="controls">`)
	gf.P(`            <button class="expand-all" id="expandAll">全部展开</button>`)
	gf.P(`            <button class="collapse-all" id="collapseAll">全部折叠</button>`)
	gf.P(`        </div>`)
	gf.P(`    </div>`)
}

// 生成JavaScript代码
func generateScripts(gf *GeneratedFile, config TableConfig) {
	gf.P(`    <script>`)
	gf.P(`        // 配置数据`)
	gf.P(`        const config = {`)
	gf.P(`            columns: [`)

	// 生成列配置JSON
	for i, col := range config.Columns {
		comma := ","
		if i == len(config.Columns)-1 {
			comma = ""
		}
		gf.P(`                {key: "`, col.Key, `", label: "`, col.Label, `", width: "`, col.Width, `", sortable: `, fmt.Sprintf("%t", col.Sortable), `}`, comma)
	}

	gf.P(`            ],`)
	gf.P(`            treeData: [`)

	// 生成树形数据JSON - 修复了JSON生成问题
	generateTreeDataJSON(gf, config.TreeData, 4)

	gf.P(`            ]`)
	gf.P(`        };`)
	gf.P(``)

	// 生成JavaScript函数
	generateJavaScriptFunctions(gf)

	gf.P(`    </script>`)
}

// 生成树形数据的JSON - 修复版本
func generateTreeDataJSON(gf *GeneratedFile, data []TreeNode, indentLevel int) {
	indent := strings.Repeat(" ", indentLevel)

	for i, node := range data {
		comma := ","
		if i == len(data)-1 {
			comma = ""
		}

		gf.P(indent, `{`)
		gf.P(indent, `    id: "`, node.ID, `",`)
		gf.P(indent, `    label: "`, escapeJSONString(node.Label), `",`)
		gf.P(indent, `    expanded: `, fmt.Sprintf("%t", node.Expanded), `,`)
		gf.P(indent, `    data: {`)

		// 生成数据字段
		dataKeys := make([]string, 0, len(node.Data))
		for k := range node.Data {
			dataKeys = append(dataKeys, k)
		}
		for j, key := range dataKeys {
			dataComma := ","
			if j == len(dataKeys)-1 {
				dataComma = ""
			}
			gf.P(indent, `        "`, key, `": "`, escapeJSONString(node.Data[key]), `"`, dataComma)
		}

		gf.P(indent, `    },`)

		// 递归生成子节点
		if len(node.Children) > 0 {
			gf.P(indent, `    children: [`)
			generateTreeDataJSON(gf, node.Children, indentLevel+8)
			gf.P(indent, `    ]`)
		} else {
			gf.P(indent, `    children: []`)
		}

		gf.P(indent, `}`, comma)
	}
}

// 转义JSON字符串中的特殊字符
func escapeJSONString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}

// 生成JavaScript函数
func generateJavaScriptFunctions(gf *GeneratedFile) {
	gf.P(`        // SVG图标定义`)
	gf.P(`        const svgIcons = {`)
	gf.P(`            folderOpen: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" transform="rotate(90)"><path fill="currentColor" d="M256 8c137 0 248 111 248 248S393 504 256 504 8 393 8 256 119 8 256 8zm113.9 231L234.4 103.5c-9.4-9.4-24.6-9.4-33.9 0l-17 17c-9.4 9.4-9.4 24.6 0 33.9L285.1 256 183.5 357.6c-9.4 9.4-9.4 24.6 0 33.9l17 17c9.4 9.4 24.6 9.4 33.9 0L369.9 273c9.4-9.4 9.4-24.6 0-34z"/></svg>',`)
	gf.P(`            folderClosed: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path fill="currentColor" d="M256 8c137 0 248 111 248 248S393 504 256 504 8 393 8 256 119 8 256 8zm113.9 231L234.4 103.5c-9.4-9.4-24.6-9.4-33.9 0l-17 17c-9.4 9.4-9.4 24.6 0 33.9L285.1 256 183.5 357.6c-9.4 9.4-9.4 24.6 0 33.9l17 17c9.4 9.4 24.6 9.4 33.9 0L369.9 273c9.4-9.4 9.4-24.6 0-34z"/></svg>',`)
	gf.P(`            document: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512" width="12" height="12" fill="#90a4ae"><path d="M0 64C0 28.7 28.7 0 64 0H224V128c0 17.7 14.3 32 32 32H384V448c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V64zm384 64H256V0L384 128z"/></svg>'`)
	gf.P(`        };`)
	gf.P(``)
	gf.P(`        // 排序状态`)
	gf.P(`        let currentSort = {`)
	gf.P(`            key: null,`)
	gf.P(`            direction: 'none' // 'none', 'asc', 'desc'`)
	gf.P(`        };`)
	gf.P(``)
	gf.P(`        // 生成树形表格行的函数`)
	gf.P(`        function generateTreeRows(data, level = 0, parentId = null) {`)
	gf.P(`            let rows = '';`)
	gf.P(`            data.forEach(item => {`)
	gf.P(`                const hasChildren = item.children && item.children.length > 0;`)
	gf.P(`                const indent = level * 16;`)
	gf.P(`                const isExpanded = item.expanded !== false;`)
	gf.P(`                rows += '<tr class="tree-row level-' + level + '" data-id="' + item.id + '" data-parent="' + parentId + '" data-level="' + level + '"' + (!isExpanded && level > 0 ? ' style="display: none;"' : '') + '>';`)
	gf.P(`                config.columns.forEach((col, colIndex) => {`)
	gf.P(`                    if (colIndex === 0) {`)
	gf.P(`                        rows += '<td><div class="tree-cell"><div class="tree-indent" style="width: ' + indent + 'px;"></div>' +`)
	gf.P(`                            (hasChildren ? '<button class="toggle-btn" data-id="' + item.id + '" data-expanded="' + isExpanded + '" title="展开/折叠">' + (isExpanded ? svgIcons.folderOpen : svgIcons.folderClosed) + '</button>' : '<span class="no-toggle" title="无子项">' + svgIcons.document + '</span>') +`)
	gf.P(`                            (item.data[col.key] || item.label) + '</div></td>';`)
	gf.P(`                    } else {`)
	gf.P(`                        const cellValue = item.data[col.key] || '';`)
	gf.P(`                        const statusClass = col.key === 'status' ? getStatusClass(cellValue) : '';`)
	gf.P(`                        rows += '<td><span class="' + statusClass + '">' + cellValue + '</span></td>';`)
	gf.P(`                    }`)
	gf.P(`                });`)
	gf.P(`                rows += '</tr>';`)
	gf.P(`                if (hasChildren && isExpanded) {`)
	gf.P(`                    rows += generateTreeRows(item.children, level + 1, item.id);`)
	gf.P(`                }`)
	gf.P(`            });`)
	gf.P(`            return rows;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        // 获取状态样式类`)
	gf.P(`        function getStatusClass(status) {`)
	gf.P(`            const statusMap = {`)
	gf.P(`                '已完成': 'status-active',`)
	gf.P(`                '进行中': 'status-active',`)
	gf.P(`                '未开始': 'status-inactive',`)
	gf.P(`                '待处理': 'status-pending'`)
	gf.P(`            };`)
	gf.P(`            return statusMap[status] || '';`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        // 初始化表格`)
	gf.P(`        function initTreeTable() {`)
	gf.P(`            // 添加表头排序事件`)
	gf.P(`            document.querySelectorAll('th.sortable').forEach(th => {`)
	gf.P(`                th.addEventListener('click', function() {`)
	gf.P(`                    const key = this.getAttribute('data-key');`)
	gf.P(`                    sortTable(key);`)
	gf.P(`                });`)
	gf.P(`            });`)
	gf.P(`            `)
	gf.P(`            renderTable();`)
	gf.P(`            expandAll();`)
	gf.P(`            `)
	gf.P(`            // 添加展开/折叠事件`)
	gf.P(`            document.querySelectorAll('.toggle-btn').forEach(btn => {`)
	gf.P(`                btn.addEventListener('click', function() {`)
	gf.P(`                    const itemId = this.getAttribute('data-id');`)
	gf.P(`                    const isExpanded = this.getAttribute('data-expanded') === 'true';`)
	gf.P(`                    toggleChildren(itemId, this, isExpanded);`)
	gf.P(`                });`)
	gf.P(`            });`)
	gf.P(`            `)
	gf.P(`            // 按钮事件`)
	gf.P(`            document.getElementById('expandAll').addEventListener('click', expandAll);`)
	gf.P(`            document.getElementById('collapseAll').addEventListener('click', collapseAll);`)
	gf.P(`            document.getElementById('searchBtn').addEventListener('click', performSearch);`)
	gf.P(`            `)
	gf.P(`            // 回车搜索`)
	gf.P(`            document.getElementById('searchInput').addEventListener('keypress', function(e) {`)
	gf.P(`                if (e.key === 'Enter') performSearch();`)
	gf.P(`            });`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        // 排序表格`)
	gf.P(`        function sortTable(key) {`)
	gf.P(`            // 确定新的排序方向`)
	gf.P(`            let newDirection = 'asc';`)
	gf.P(`            if (currentSort.key === key) {`)
	gf.P(`                if (currentSort.direction === 'asc') {`)
	gf.P(`                    newDirection = 'desc';`)
	gf.P(`                } else if (currentSort.direction === 'desc') {`)
	gf.P(`                    newDirection = 'none';`)
	gf.P(`                }`)
	gf.P(`            }`)
	gf.P(`            `)
	gf.P(`            // 更新排序状态`)
	gf.P(`            currentSort.key = key;`)
	gf.P(`            currentSort.direction = newDirection;`)
	gf.P(`            `)
	gf.P(`            // 更新表头样式`)
	gf.P(`            document.querySelectorAll('th.sortable').forEach(th => {`)
	gf.P(`                th.classList.remove('sort-asc', 'sort-desc');`)
	gf.P(`                th.setAttribute('data-sort', 'none');`)
	gf.P(`            });`)
	gf.P(`            `)
	gf.P(`            if (newDirection !== 'none') {`)
	gf.P(`                const currentTh = document.querySelector('th[data-key="' + key + '"]');`)
	gf.P(`                currentTh.classList.add('sort-' + newDirection);`)
	gf.P(`                currentTh.setAttribute('data-sort', newDirection);`)
	gf.P(`            }`)
	gf.P(`            `)
	gf.P(`            // 对数据进行排序`)
	gf.P(`            if (newDirection !== 'none') {`)
	gf.P(`                sortTreeData(config.treeData, key, newDirection);`)
	gf.P(`            }`)
	gf.P(`            `)
	gf.P(`            // 重新渲染表格`)
	gf.P(`            renderTable();`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        // 递归排序树形数据`)
	gf.P(`        function sortTreeData(data, key, direction) {`)
	gf.P(`            data.sort((a, b) => {`)
	gf.P(`                let aValue = a.data[key] || a[key] || '';`)
	gf.P(`                let bValue = b.data[key] || b[key] || '';`)
	gf.P(`                `)
	gf.P(`                // 尝试转换为数字比较`)
	gf.P(`                const aNum = parseFloat(aValue);`)
	gf.P(`                const bNum = parseFloat(bValue);`)
	gf.P(`                if (!isNaN(aNum) && !isNaN(bNum)) {`)
	gf.P(`                    aValue = aNum;`)
	gf.P(`                    bValue = bNum;`)
	gf.P(`                }`)
	gf.P(`                `)
	gf.P(`                if (aValue < bValue) return direction === 'asc' ? -1 : 1;`)
	gf.P(`                if (aValue > bValue) return direction === 'asc' ? 1 : -1;`)
	gf.P(`                return 0;`)
	gf.P(`            });`)
	gf.P(`            `)
	gf.P(`            // 递归排序子节点`)
	gf.P(`            data.forEach(item => {`)
	gf.P(`                if (item.children && item.children.length > 0) {`)
	gf.P(`                    sortTreeData(item.children, key, direction);`)
	gf.P(`                }`)
	gf.P(`            });`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        // 渲染表格`)
	gf.P(`        function renderTable() {`)
	gf.P(`            const tableBody = document.getElementById('treeTableBody');`)
	gf.P(`            tableBody.innerHTML = generateTreeRows(config.treeData);`)
	gf.P(`            `)
	gf.P(`            // 重新绑定展开/折叠事件`)
	gf.P(`            document.querySelectorAll('.toggle-btn').forEach(btn => {`)
	gf.P(`                btn.addEventListener('click', function() {`)
	gf.P(`                    const itemId = this.getAttribute('data-id');`)
	gf.P(`                    const isExpanded = this.getAttribute('data-expanded') === 'true';`)
	gf.P(`                    toggleChildren(itemId, this, isExpanded);`)
	gf.P(`                });`)
	gf.P(`            });`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        // 搜索功能`)
	gf.P(`        function performSearch() {`)
	gf.P(`            const searchTerm = document.getElementById('searchInput').value.toLowerCase();`)
	gf.P(`            if (!searchTerm) {`)
	gf.P(`                // 重置显示所有行`)
	gf.P(`                document.querySelectorAll('.tree-row').forEach(row => {`)
	gf.P(`                    row.style.display = '';`)
	gf.P(`                });`)
	gf.P(`                return;`)
	gf.P(`            }`)
	gf.P(`            `)
	gf.P(`            document.querySelectorAll('.tree-row').forEach(row => {`)
	gf.P(`                const rowText = row.textContent.toLowerCase();`)
	gf.P(`                row.style.display = rowText.includes(searchTerm) ? '' : 'none';`)
	gf.P(`            });`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        // 展开/折叠功能`)
	gf.P(`        function toggleChildren(itemId, button, isExpanded) {`)
	gf.P(`            const rows = document.querySelectorAll('tr[data-parent="' + itemId + '"]');`)
	gf.P(`            const newExpandedState = !isExpanded;`)
	gf.P(`            `)
	gf.P(`            // 更新数据中的展开状态`)
	gf.P(`            updateExpandedState(config.treeData, itemId, newExpandedState);`)
	gf.P(`            `)
	gf.P(`            rows.forEach(row => {`)
	gf.P(`                row.style.display = newExpandedState ? '' : 'none';`)
	gf.P(`                `)
	gf.P(`                // 递归处理子项`)
	gf.P(`                if (newExpandedState) {`)
	gf.P(`                    const childId = row.getAttribute('data-id');`)
	gf.P(`                    const childButton = row.querySelector('.toggle-btn');`)
	gf.P(`                    if (childButton && childButton.getAttribute('data-expanded') === 'true') {`)
	gf.P(`                        toggleChildren(childId, childButton, true);`)
	gf.P(`                    }`)
	gf.P(`                }`)
	gf.P(`            });`)
	gf.P(`            `)
	gf.P(`            button.setAttribute('data-expanded', newExpandedState);`)
	gf.P(`            button.innerHTML = newExpandedState ? svgIcons.folderOpen : svgIcons.folderClosed;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        // 更新数据中的展开状态`)
	gf.P(`        function updateExpandedState(data, itemId, expanded) {`)
	gf.P(`            for (let item of data) {`)
	gf.P(`                if (item.id === itemId) {`)
	gf.P(`                    item.expanded = expanded;`)
	gf.P(`                    return true;`)
	gf.P(`                }`)
	gf.P(`                if (item.children && item.children.length > 0) {`)
	gf.P(`                    if (updateExpandedState(item.children, itemId, expanded)) {`)
	gf.P(`                        return true;`)
	gf.P(`                    }`)
	gf.P(`                }`)
	gf.P(`            }`)
	gf.P(`            return false;`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        function expandAll() {`)
	gf.P(`            setAllExpanded(config.treeData, true);`)
	gf.P(`            renderTable();`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        function collapseAll() {`)
	gf.P(`            setAllExpanded(config.treeData, false);`)
	gf.P(`            renderTable();`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        function setAllExpanded(data, expanded) {`)
	gf.P(`            data.forEach(item => {`)
	gf.P(`                item.expanded = expanded;`)
	gf.P(`                if (item.children && item.children.length > 0) {`)
	gf.P(`                    setAllExpanded(item.children, expanded);`)
	gf.P(`                }`)
	gf.P(`            });`)
	gf.P(`        }`)
	gf.P(``)
	gf.P(`        // 初始化`)
	gf.P(`        document.addEventListener('DOMContentLoaded', initTreeTable);`)
}

// 示例使用
func main() {
	// 创建配置
	config := TableConfig{
		Title:      "树形表格优化版",
		FileName:   "玉溪2025年10月数据表",
		HeaderText: "假设统价2.5元/件，翘靶子4元/件计算，具体以实际为准",
		Columns: []Column{
			{Key: "name", Label: "姓名", Width: "20%", Sortable: true},
			{Key: "count", Label: "件数", Width: "15%", Sortable: true},
			{Key: "qiaobazi", Label: "翘靶子", Width: "15%", Sortable: true},
			{Key: "wage", Label: "工钱", Width: "15%", Sortable: true},
			{Key: "status", Label: "状态", Width: "15%", Sortable: true},
		},
		TreeData: []TreeNode{
			{
				ID:    "1",
				Label: "小风",
				Data: map[string]string{
					"name":     "小风",
					"count":    "85.00",
					"qiaobazi": "11.5",
					"wage":     "258.50",
					"status":   "已完成",
				},
				Children: []TreeNode{
					{
						ID:    "2",
						Label: "10月23日工作",
						Data: map[string]string{
							"name":     "10月23日工作",
							"count":    "45.00",
							"qiaobazi": "6.0",
							"wage":     "142.50",
							"status":   "已完成",
						},
						Children: []TreeNode{},
					},
					{
						ID:    "3",
						Label: "10月25日工作",
						Data: map[string]string{
							"name":     "10月25日工作",
							"count":    "40.00",
							"qiaobazi": "5.5",
							"wage":     "116.00",
							"status":   "已完成",
						},
						Children: []TreeNode{},
					},
				},
			},
			{
				ID:    "4",
				Label: "杨浩",
				Data: map[string]string{
					"name":     "杨浩",
					"count":    "75.00",
					"qiaobazi": "11.5",
					"wage":     "233.50",
					"status":   "进行中",
				},
				Children: []TreeNode{
					{
						ID:    "5",
						Label: "10月23日工作",
						Data: map[string]string{
							"name":     "10月23日工作",
							"count":    "35.00",
							"qiaobazi": "5.5",
							"wage":     "116.00",
							"status":   "已完成",
						},
						Children: []TreeNode{},
					},
				},
			},
		},
	}

	// 生成HTML文件
	err := GenerateHTMLFile(config, "tree_table.html")
	if err != nil {
		fmt.Printf("生成HTML文件失败: %v\n", err)
		return
	}
	fmt.Println("HTML文件生成成功!")
}

func Run(config TableConfig) {
	mylog.Check(GenerateHTMLFile(config, "tree_table.html"))
}
