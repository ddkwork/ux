package main

import (
	"fmt"
	"github.com/ddkwork/golibrary/stream"
	"testing"
)

// 取色
// https://products.eptimize.app/zh/color-convert/rgb-to-rgba
func TestName(t *testing.T) {
	//48d4fe 255不透明，测试gio颜色拾取器pick到的表格2选中背景色是否正确
	b := stream.NewHexString("2b2b2b")
	fmt.Println(b.Bytes())
}
