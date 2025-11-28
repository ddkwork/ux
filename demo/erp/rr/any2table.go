package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
)

// ================================
// 泛型N叉树实现 (Go 1.18+)
// ================================

// Node 泛型N叉树节点
type Node[T any] struct {
	Name     string       `json:"name"`
	Value    T           `json:"value,omitempty"`
	Children []*Node[T]  `json:"children,omitempty"`
	Metadata Metadata    `json:"metadata,omitempty"`
}

// Metadata 元数据存储
type Metadata map[string]interface{}

// NewNode 创建新节点
func NewNode[T any](name string, value T) *Node[T] {
	return &Node[T]{
		Name:     name,
		Value:    value,
		Children: make([]*Node[T], 0),
		Metadata: make(Metadata),
	}
}

// AddChild 添加子节点
func (n *Node[T]) AddChild(child *Node[T]) {
	n.Children = append(n.Children, child)
}

// AddMetadata 添加元数据
func (n *Node[T]) AddMetadata(key string, value interface{}) {
	if n.Metadata == nil {
		n.Metadata = make(Metadata)
	}
	n.Metadata[key] = value
}

// FindNode 查找节点（深度优先）
func (n *Node[T]) FindNode(name string) *Node[T] {
	if n.Name == name {
		return n
	}
	for _, child := range n.Children {
		if found := child.FindNode(name); found != nil {
			return found
		}
	}
	return nil
}

// GetDepth 获取节点深度
func (n *Node[T]) GetDepth() int {
	maxDepth := 0
	for _, child := range n.Children {
		if depth := child.GetDepth() + 1; depth > maxDepth {
			maxDepth = depth
		}
	}
	return maxDepth
}

// ================================
// 解析器接口和注册器
// ================================

// Parser 解析器接口
type Parser interface {
	Parse(reader io.Reader) (*Node[string], error)
	SupportsFormat(format string) bool
}

// FileParser 统一文件解析器
type FileParser struct {
	parsers map[string]Parser
}

// NewFileParser 创建新的文件解析器
func NewFileParser() *FileParser {
	fp := &FileParser{
		parsers: make(map[string]Parser),
	}

	// 注册所有解析器
	fp.RegisterParser("json", &JSONParser{})
	fp.RegisterParser("md", &MarkdownParser{})
	fp.RegisterParser("markdown", &MarkdownParser{})
	fp.RegisterParser("txt", &TXTParser{})
	fp.RegisterParser("csv", &CSVParser{})
	fp.RegisterParser("sqlite", &SQLiteParser{})
	fp.RegisterParser("db", &SQLiteParser{})
	fp.RegisterParser("yaml", &YAMLParser{})
	fp.RegisterParser("yml", &YAMLParser{})
	fp.RegisterParser("xml", &XMLParser{})
	fp.RegisterParser("xlsx", &ExcelParser{})
	fp.RegisterParser("xls", &ExcelParser{})

	return fp
}

// RegisterParser 注册解析器
func (fp *FileParser) RegisterParser(format string, parser Parser) {
	fp.parsers[format] = parser
}

// ParseContent 解析内容
func (fp *FileParser) ParseContent(content, format string) (*Node[string], error) {
	parser, exists := fp.parsers[format]
	if !exists {
		return nil, fmt.Errorf("不支持的格式: %s", format)
	}

	reader := strings.NewReader(content)
	return parser.Parse(reader)
}

// GetSupportedFormats 获取支持的格式列表
func (fp *FileParser) GetSupportedFormats() []string {
	formats := make([]string, 0, len(fp.parsers))
	for format := range fp.parsers {
		formats = append(formats, format)
	}
	return formats
}

// ================================
// JSON 解析器
// ================================

type JSONParser struct{}

func (jp *JSONParser) SupportsFormat(format string) bool {
	return format == "json"
}

func (jp *JSONParser) Parse(reader io.Reader) (*Node[string], error) {
	var data interface{}
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	root := NewNode("root", "")
	jp.parseValue(data, root, "root")
	return root, nil
}

func (jp *JSONParser) parseValue(value interface{}, parent *Node[string], key string) {
	switch v := value.(type) {
	case map[string]interface{}:
		node := NewNode(key, "")
		parent.AddChild(node)
		for k, val := range v {
			jp.parseValue(val, node, k)
		}
	case []interface{}:
		node := NewNode(key, "array")
		parent.AddChild(node)
		for i, val := range v {
			jp.parseValue(val, node, "["+strconv.Itoa(i)+"]")
		}
	default:
		valueStr := jp.valueToString(v)
		node := NewNode(key, valueStr)
		parent.AddChild(node)
	}
}

func (jp *JSONParser) valueToString(v interface{}) string {
	if v == nil {
		return "null"
	}
	return fmt.Sprintf("%v", v)
}

// ================================
// Markdown 解析器
// ================================

type MarkdownParser struct{}

func (mp *MarkdownParser) SupportsFormat(format string) bool {
	return format == "md" || format == "markdown"
}

func (mp *MarkdownParser) Parse(reader io.Reader) (*Node[string], error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	root := NewNode("document", "")
	lines := strings.Split(string(content), "\n")

	var currentSection *Node[string]
	var currentSubsection *Node[string]

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 解析标题
		if strings.HasPrefix(line, "#") {
			level := 0
			for strings.HasPrefix(line, "#") {
				level++
				line = strings.TrimPrefix(line, "#")
			}
			line = strings.TrimSpace(line)

			node := NewNode("h"+strconv.Itoa(level), line)
			node.AddMetadata("line", i+1)
			node.AddMetadata("level", level)

			switch level {
			case 1:
				root.AddChild(node)
				currentSection = node
				currentSubsection = nil
			case 2:
				if currentSection != nil {
					currentSection.AddChild(node)
					currentSubsection = node
				} else {
					root.AddChild(node)
					currentSection = node
				}
			default:
				if currentSubsection != nil {
					currentSubsection.AddChild(node)
				} else if currentSection != nil {
					currentSection.AddChild(node)
				} else {
					root.AddChild(node)
				}
			}
		} else if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			// 处理列表项
			content := strings.TrimPrefix(line, "- ")
			content = strings.TrimPrefix(content, "* ")
			listItem := NewNode("list_item", content)
			listItem.AddMetadata("line", i+1)

			if currentSubsection != nil {
				currentSubsection.AddChild(listItem)
			} else if currentSection != nil {
				currentSection.AddChild(listItem)
			} else {
				root.AddChild(listItem)
			}
		} else {
			// 处理段落
			if currentSubsection != nil {
				paragraph := NewNode("paragraph", line)
				paragraph.AddMetadata("line", i+1)
				currentSubsection.AddChild(paragraph)
			} else if currentSection != nil {
				paragraph := NewNode("paragraph", line)
				paragraph.AddMetadata("line", i+1)
				currentSection.AddChild(paragraph)
			} else {
				paragraph := NewNode("paragraph", line)
				paragraph.AddMetadata("line", i+1)
				root.AddChild(paragraph)
			}
		}
	}

	return root, nil
}

// ================================
// 文本文件解析器
// ================================

type TXTParser struct{}

func (tp *TXTParser) SupportsFormat(format string) bool {
	return format == "txt"
}

func (tp *TXTParser) Parse(reader io.Reader) (*Node[string], error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	root := NewNode("text_file", "")
	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		node := NewNode("line_"+strconv.Itoa(i+1), line)
		node.AddMetadata("line_number", i+1)
		node.AddMetadata("length", len(line))
		root.AddChild(node)
	}

	return root, nil
}

// ================================
// CSV 解析器
// ================================

type CSVParser struct{}

func (cp *CSVParser) SupportsFormat(format string) bool {
	return format == "csv"
}

func (cp *CSVParser) Parse(reader io.Reader) (*Node[string], error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return NewNode("csv", "empty"), nil
	}

	root := NewNode("csv", "")
	root.AddMetadata("rows", len(records))
	root.AddMetadata("columns", len(records[0]))

	// 添加表头
	if len(records) > 0 {
		headerNode := NewNode("headers", "")
		for i, header := range records[0] {
			colNode := NewNode("col_"+strconv.Itoa(i+1), header)
			colNode.AddMetadata("index", i)
			headerNode.AddChild(colNode)
		}
		root.AddChild(headerNode)
	}

	// 添加数据行
	for rowIdx, record := range records {
		if rowIdx == 0 {
			continue // 跳过表头
		}

		rowNode := NewNode("row_"+strconv.Itoa(rowIdx), "")
		rowNode.AddMetadata("row_number", rowIdx)

		for colIdx, value := range record {
			cellNode := NewNode("col_"+strconv.Itoa(colIdx+1), value)
			cellNode.AddMetadata("column_index", colIdx)
			rowNode.AddChild(cellNode)
		}

		root.AddChild(rowNode)
	}

	return root, nil
}

// ================================
// SQLite 解析器
// ================================

type SQLiteParser struct{}

func (sp *SQLiteParser) SupportsFormat(format string) bool {
	return format == "sqlite" || format == "db"
}

func (sp *SQLiteParser) Parse(reader io.Reader) (*Node[string], error) {
	// SQLite需要文件路径，这里返回模拟数据
	return sp.parseMockData(), nil
}

func (sp *SQLiteParser) parseMockData() *Node[string] {
	root := NewNode("database", "example.db")

	// 模拟表结构
	tablesNode := NewNode("tables", "")

	usersTable := NewNode("users", "")
	usersTable.AddMetadata("row_count", 3)

	// 添加列信息
	idCol := NewNode("id", "INTEGER PRIMARY KEY")
	nameCol := NewNode("name", "TEXT")
	emailCol := NewNode("email", "TEXT")

	usersTable.AddChild(idCol)
	usersTable.AddChild(nameCol)
	usersTable.AddChild(emailCol)

	// 添加示例数据
	dataNode := NewNode("data", "")

	row1 := NewNode("row_1", "")
	row1.AddChild(NewNode("id", "1"))
	row1.AddChild(NewNode("name", "Alice"))
	row1.AddChild(NewNode("email", "alice@example.com"))

	row2 := NewNode("row_2", "")
	row2.AddChild(NewNode("id", "2"))
	row2.AddChild(NewNode("name", "Bob"))
	row2.AddChild(NewNode("email", "bob@example.com"))

	dataNode.AddChild(row1)
	dataNode.AddChild(row2)

	usersTable.AddChild(dataNode)
	tablesNode.AddChild(usersTable)
	root.AddChild(tablesNode)

	return root
}

// 真实的SQLite解析实现
func (sp *SQLiteParser) ParseFile(dbPath string) (*Node[string], error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	root := NewNode("database", dbPath)

	// 获取所有表
	tables, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return nil, err
	}
	defer tables.Close()

	tablesNode := NewNode("tables", "")

	for tables.Next() {
		var tableName string
		if err := tables.Scan(&tableName); err != nil {
			continue
		}

		tableNode := NewNode(tableName, "")
		tablesNode.AddChild(tableNode)

		// 获取表结构
		columns, err := db.Query("PRAGMA table_info(" + tableName + ")")
		if err == nil {
			for columns.Next() {
				var cid int
				var name, ctype string
				var notnull, pk int
				var dflt_value interface{}

				columns.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk)
				colNode := NewNode(name, ctype)
				tableNode.AddChild(colNode)
			}
			columns.Close()
		}
	}

	root.AddChild(tablesNode)
	return root, nil
}

// ================================
// YAML 解析器
// ================================

type YAMLParser struct{}

func (yp *YAMLParser) SupportsFormat(format string) bool {
	return format == "yaml" || format == "yml"
}

func (yp *YAMLParser) Parse(reader io.Reader) (*Node[string], error) {
	var data interface{}
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	root := NewNode("root", "")
	yp.parseValue(data, root, "root")
	return root, nil
}

func (yp *YAMLParser) parseValue(value interface{}, parent *Node[string], key string) {
	switch v := value.(type) {
	case map[string]interface{}:
		node := NewNode(key, "")
		parent.AddChild(node)
		for k, val := range v {
			yp.parseValue(val, node, k)
		}
	case map[interface{}]interface{}:
		node := NewNode(key, "")
		parent.AddChild(node)
		for k, val := range v {
			yp.parseValue(val, node, fmt.Sprintf("%v", k))
		}
	case []interface{}:
		node := NewNode(key, "sequence")
		parent.AddChild(node)
		for i, val := range v {
			yp.parseValue(val, node, "["+strconv.Itoa(i)+"]")
		}
	default:
		valueStr := yp.valueToString(v)
		node := NewNode(key, valueStr)
		parent.AddChild(node)
	}
}

func (yp *YAMLParser) valueToString(v interface{}) string {
	if v == nil {
		return "null"
	}
	return fmt.Sprintf("%v", v)
}

// ================================
// XML 解析器
// ================================

type XMLParser struct{}

func (xp *XMLParser) SupportsFormat(format string) bool {
	return format == "xml"
}

func (xp *XMLParser) Parse(reader io.Reader) (*Node[string], error) {
	decoder := xml.NewDecoder(reader)
	root := NewNode("xml", "")

	var current *Node[string]
	stack := []*Node[string]{root}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			node := NewNode(se.Name.Local, "")
			for _, attr := range se.Attr {
				node.AddMetadata("attr_"+attr.Name.Local, attr.Value)
			}

			current = node
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.AddChild(node)
			}
			stack = append(stack, node)

		case xml.EndElement:
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}

		case xml.CharData:
			text := strings.TrimSpace(string(se))
			if text != "" && current != nil {
				current.Value = text
			}
		}
	}

	return root, nil
}

// ================================
// Excel 解析器 (模拟实现)
// ================================

type ExcelParser struct{}

func (ep *ExcelParser) SupportsFormat(format string) bool {
	return format == "xlsx" || format == "xls"
}

func (ep *ExcelParser) Parse(reader io.Reader) (*Node[string], error) {
	// 模拟Excel解析 - 实际实现需要使用excelize等库
	return ep.parseMockData(), nil
}

func (ep *ExcelParser) parseMockData() *Node[string] {
	root := NewNode("workbook", "")

	// 模拟工作表
	sheet1 := NewNode("Sheet1", "")
	root.AddChild(sheet1)

	// 模拟表头
	headers := NewNode("headers", "")
	headers.AddChild(NewNode("A1", "Name"))
	headers.AddChild(NewNode("B1", "Age"))
	headers.AddChild(NewNode("C1", "Email"))
	sheet1.AddChild(headers)

	// 模拟数据行
	data := NewNode("data", "")

	row1 := NewNode("row_2", "")
	row1.AddChild(NewNode("A2", "Alice"))
	row1.AddChild(NewNode("B2", "25"))
	row1.AddChild(NewNode("C2", "alice@example.com"))

	row2 := NewNode("row_3", "")
	row2.AddChild(NewNode("A3", "Bob"))
	row2.AddChild(NewNode("B3", "30"))
	row2.AddChild(NewNode("C3", "bob@example.com"))

	data.AddChild(row1)
	data.AddChild(row2)
	sheet1.AddChild(data)

	return root
}

// ================================
// 代码生成器
// ================================

type CodeGenerator struct{}

// GenerateStructCode 生成泛型结构体定义代码
func (cg *CodeGenerator) GenerateStructCode() string {
	return `// Node 泛型N叉树节点
type Node[T any] struct {
    Name     string       \`json:"name"\`
    Value    T            \`json:"value,omitempty"\`
    Children []*Node[T]   \`json:"children,omitempty"\`
    Metadata Metadata     \`json:"metadata,omitempty"\`
}

// Metadata 元数据存储
type Metadata map[string]interface{}

// NewNode 创建新节点
func NewNode[T any](name string, value T) *Node[T] {
    return &Node[T]{
        Name:     name,
        Value:    value,
        Children: make([]*Node[T], 0),
        Metadata: make(Metadata),
    }
}

// AddChild 添加子节点
func (n *Node[T]) AddChild(child *Node[T]) {
    n.Children = append(n.Children, child)
}

// AddMetadata 添加元数据
func (n *Node[T]) AddMetadata(key string, value interface{}) {
    if n.Metadata == nil {
        n.Metadata = make(Metadata)
    }
    n.Metadata[key] = value
}`
}

// GenerateInstanceCode 生成实例化代码
func (cg *CodeGenerator) GenerateInstanceCode(root *Node[string], varName string) string {
	var builder strings.Builder

	builder.WriteString("// 自动生成的N叉树实例代码\n")
	builder.WriteString(fmt.Sprintf("func create%s() *Node[string] {\n", strings.Title(varName)))

	cg.generateNodeCode(root, &builder, 1)

	builder.WriteString(fmt.Sprintf("    return %s\n", varName))
	builder.WriteString("}\n\n")

	return builder.String()
}

func (cg *CodeGenerator) generateNodeCode(node *Node[string], builder *strings.Builder, indentLevel int) {
	indent := strings.Repeat("    ", indentLevel)

	builder.WriteString(fmt.Sprintf("%s%s := NewNode(\"%s\", \"%s\")\n",
		indent, cg.getSafeVarName(node.Name), node.Name, node.Value))

	// 添加元数据
	for key, value := range node.Metadata {
		builder.WriteString(fmt.Sprintf("%s%s.AddMetadata(\"%s\", %#v)\n",
			indent, cg.getSafeVarName(node.Name), key, value))
	}

	// 递归处理子节点
	for _, child := range node.Children {
		cg.generateNodeCode(child, builder, indentLevel+1)
		builder.WriteString(fmt.Sprintf("%s%s.AddChild(%s)\n",
			indent, cg.getSafeVarName(node.Name), cg.getSafeVarName(child.Name)))
	}
}

func (cg *CodeGenerator) getSafeVarName(name string) string {
	// 将名称转换为有效的Go变量名
	safeName := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, name)

	if safeName == "" {
		return "node"
	}

	// 确保不以数字开头
	if safeName[0] >= '0' && safeName[0] <= '9' {
		safeName = "node_" + safeName
	}

	return safeName
}

// ================================
// 工具函数
// ================================

// PrintTree 打印树结构（用于调试）
func PrintTree[T any](node *Node[T], depth int) {
	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s%s: %v\n", indent, node.Name, node.Value)

	for _, child := range node.Children {
		PrintTree(child, depth+1)
	}
}

// TreeToJSON 将树转换为JSON字符串
func TreeToJSON[T any](node *Node[T]) (string, error) {
	data, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ================================
// 主函数和示例
// ================================

func main() {
	// 创建文件解析器
	fileParser := NewFileParser()

	fmt.Printf("支持的格式: %v\n\n", fileParser.GetSupportedFormats())

	// 示例1: 解析JSON
	exampleJSON(fileParser)

	// 示例2: 解析YAML
	exampleYAML(fileParser)

	// 示例3: 解析XML
	exampleXML(fileParser)
}

func exampleJSON(fileParser *FileParser) {
	fmt.Println("=== JSON解析示例 ===")

	jsonContent := `{
		"user": {
			"id": 1,
			"name": "Alice",
			"profile": {
				"age": 25,
				"city": "Beijing"
			},
			"hobbies": ["reading", "swimming"]
		}
	}`

	root, err := fileParser.ParseContent(jsonContent, "json")
	if err != nil {
		log.Fatal("解析JSON失败:", err)
	}

	// 生成代码
	codeGen := CodeGenerator{}

	fmt.Println("=== 实例化代码 ===")
	fmt.Println(codeGen.GenerateInstanceCode(root, "jsonTree"))

	fmt.Println("=== 树结构 ===")
	PrintTree(root, 0)
	fmt.Println()
}

func exampleYAML(fileParser *FileParser) {
	fmt.Println("=== YAML解析示例 ===")

	yamlContent := `
user:
  id: 1
  name: Bob
  profile:
    age: 30
    city: Shanghai
  tags:
    - developer
    - golang
`

	root, err := fileParser.ParseContent(yamlContent, "yaml")
	if err != nil {
		log.Fatal("解析YAML失败:", err)
	}

	fmt.Println("=== 树结构 ===")
	PrintTree(root, 0)
	fmt.Println()
}

func exampleXML(fileParser *FileParser) {
	fmt.Println("=== XML解析示例 ===")

	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<person id="123">
	<name>Charlie</name>
	<age>35</age>
	<address>
		<city>Guangzhou</city>
		<country>China</country>
	</address>
</person>`

	root, err := fileParser.ParseContent(xmlContent, "xml")
	if err != nil {
		log.Fatal("解析XML失败:", err)
	}

	fmt.Println("=== 树结构 ===")
	PrintTree(root, 0)
	fmt.Println()
}

// ================================
// 单元测试
// ================================

func TestJSONParser(t *testing.T) {
	jsonContent := `{"name": "John Doe", "age": 30, "active": true}`

	fileParser := NewFileParser()
	root, err := fileParser.ParseContent(jsonContent, "json")

	if err != nil {
		t.Fatalf("JSON解析失败: %v", err)
	}

	if root == nil {
		t.Fatal("解析结果为空")
	}

	if root.Name != "root" {
		t.Errorf("期望根节点名称为'root'，实际为'%s'", root.Name)
	}

	if len(root.Children) == 0 {
		t.Error("根节点应该包含子节点")
	}
}

func TestYAMLParser(t *testing.T) {
	yamlContent := `
name: John Doe
age: 30
active: true
`

	fileParser := NewFileParser()
	root, err := fileParser.ParseContent(yamlContent, "yaml")

	if err != nil {
		t.Fatalf("YAML解析失败: %v", err)
	}

	if root.Name != "root" {
		t.Errorf("YAML解析结果不符合预期")
	}
}

func TestXMLParser(t *testing.T) {
	xmlContent := `<root><name>John</name><age>30</age></root>`

	fileParser := NewFileParser()
	root, err := fileParser.ParseContent(xmlContent, "xml")

	if err != nil {
		t.Fatalf("XML解析失败: %v", err)
	}

	if root.Name != "xml" {
		t.Errorf("XML解析结果不符合预期")
	}
}

func TestGenericNode(t *testing.T) {
	// 测试字符串类型节点
	strNode := NewNode("test", "value")
	if strNode.Value != "value" {
		t.Errorf("字符串节点值不符合预期")
	}

	// 测试整数类型节点
	intNode := NewNode("test", 42)
	if intNode.Value != 42 {
		t.Errorf("整数节点值不符合预期")
	}

	// 测试添加子节点
	parent := NewNode("parent", "")
	child := NewNode("child", "")
	parent.AddChild(child)

	if len(parent.Children) != 1 {
		t.Errorf("添加子节点失败")
	}

	// 测试查找节点
	if found := parent.FindNode("child"); found == nil {
		t.Errorf("查找节点失败")
	}

	// 测试深度计算
	if depth := parent.GetDepth(); depth != 1 {
		t.Errorf("深度计算错误: 期望1, 实际%d", depth)
	}
}

func TestCodeGenerator(t *testing.T) {
	root := NewNode("root", "")
	child1 := NewNode("child1", "value1")
	child2 := NewNode("child2", "value2")

	root.AddChild(child1)
	root.AddChild(child2)

	codeGen := CodeGenerator{}

	structCode := codeGen.GenerateStructCode()
	if !strings.Contains(structCode, "type Node[T any] struct") {
		t.Error("结构体代码生成失败")
	}

	instanceCode := codeGen.GenerateInstanceCode(root, "testTree")
	if !strings.Contains(instanceCode, "func createTestTree()") {
		t.Error("实例化代码生成失败")
	}
}

func BenchmarkJSONParser(b *testing.B) {
	jsonContent := `{"name": "test", "value": "benchmark"}`
	fileParser := NewFileParser()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileParser.ParseContent(jsonContent, "json")
	}
}

func BenchmarkCodeGeneration(b *testing.B) {
	root := NewNode("root", "")
	for i := 0; i < 100; i++ {
		child := NewNode("child", "value")
		root.AddChild(child)
	}

	codeGen := CodeGenerator{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		codeGen.GenerateInstanceCode(root, "benchmarkTree")
	}
}
