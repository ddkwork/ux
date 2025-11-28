package main

import (
	"strings"
	"testing"
)

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

	for b.Loop() {
		fileParser.ParseContent(jsonContent, "json")
	}
}

func BenchmarkCodeGeneration(b *testing.B) {
	root := NewNode("root", "")
	for range 100 {
		child := NewNode("child", "value")
		root.AddChild(child)
	}

	codeGen := CodeGenerator{}

	for b.Loop() {
		codeGen.GenerateInstanceCode(root, "benchmarkTree")
	}
}
