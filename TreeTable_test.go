package ux_test

import (
	"testing"

	"github.com/ddkwork/golibrary/stream"

	"github.com/ddkwork/golibrary"
)

//
//import (
//	"fmt"
//	"image"
//	"net/http"
//	"os"
//	"slices"
//	"strings"
//	"testing"
//	"time"
//
//	"github.com/ddkwork/golibrary"
//
//	"github.com/ddkwork/golibrary/assert"
//
//	"gioui.org/gpu/headless"
//	"gioui.org/layout"
//	"gioui.org/op"
//	"gioui.org/unit"
//	"gioui.org/widget"
//	"github.com/ddkwork/golibrary/mylog"
//	"github.com/ddkwork/golibrary/stream"
//	"github.com/ddkwork/ux"
//)

func TestName(t *testing.T) {
	stream.UpdateModsByWorkSpace(false)
	return
	golibrary.StaticCheck()
}

//func BenchmarkTransposeMatrix(b *testing.B) {
//	columnCells := ux.InitHeader(packet{})
//	rows := make([][]packet, 0, len(columnCells))
//	for range 1000000 {
//		rows = append(rows, []packet{
//			{
//				Scheme: "http",
//				Method: "GET",
//			},
//			{
//				Host:          "www.example.com",
//				Path:          "/path/to/file.html",
//				ContentType:   "text/html",
//				ContentLength: 1024,
//				Status:        "200 OK",
//				Note:          "This is a note",
//				Process:       "chrome.exe",
//				PadTime:       10 * time.Millisecond,
//			},
//		})
//	}
//	for i := 0; i < b.N; i++ {
//		ux.TransposeMatrix[packet](rows) // BenchmarkTransposeMatrix-8   	       7	 161059457 ns/op 单独测试就很夸张
//	}
//}
//
//func TransposeMatrix2[T any](rows [][]T) (columns [][]T) {
//	if len(rows) == 0 {
//		return [][]T{}
//	}
//	numColumns := len(rows[0])
//	columns = make([][]T, numColumns)
//	for i := range columns {
//		columns[i] = make([]T, 0, len(rows))
//		for _, row := range rows {
//			columns[i] = append(columns[i], row[i])
//		}
//	}
//	return columns
//}
//
//func Benchmark_SizeColumnsToFit(b *testing.B) {
//	const scale = 1.5
//	size := image.Point{X: 1200 * scale, Y: 600 * scale}
//	w := mylog.Check2(headless.NewWindow(size.X, size.Y))
//	gtx := layout.Context{
//		Ops: new(op.Ops),
//		Metric: unit.Metric{
//			PxPerDp: scale,
//			PxPerSp: scale,
//		},
//		Constraints: layout.Exact(size),
//	}
//	mylog.Check(w.Frame(gtx.Ops))
//	t := treeTable()
//	for i := 0; i < b.N; i++ {
//		// 矩阵置换只执行一次
//		// Benchmark_SizeColumnsToFit-8   	  846088	      1685 ns/op  看起来也不差啊
//		// Benchmark_SizeColumnsToFit-8   	  299866	      4998 ns/op
//		// Benchmark_SizeColumnsToFit-8   	  226672	      4573 ns/op
//		// Benchmark_SizeColumnsToFit-8   	  217288	      5086 ns/op
//		// Benchmark_SizeColumnsToFit-8   	  278792	      4514 ns/op
//
//		// 恢复每次矩阵置换之后,渲染不那么丝滑了，可在节点获得焦点，鼠标悬停高亮行的时候明显感觉到卡顿
//		// Benchmark_SizeColumnsToFit-8   	   16428	    116786 ns/op
//		// Benchmark_SizeColumnsToFit-8   	   31593	     32610 ns/op
//		// Benchmark_SizeColumnsToFit-8   	   30961	     32604 ns/op
//		// Benchmark_SizeColumnsToFit-8   	   31825	     49652 ns/op
//		// Benchmark_SizeColumnsToFit-8   	   18141	     55213 ns/op
//		// Benchmark_SizeColumnsToFit-8   	   24626	     42476 ns/op
//		// Benchmark_SizeColumnsToFit-8   	   24418	     44423 ns/op
//		t.SizeColumnsToFit(gtx)
//	}
//}
//
//func BenchmarkLabelWidth(b *testing.B) {
//	const scale = 1.5
//	size := image.Point{X: 1200 * scale, Y: 600 * scale}
//	w := mylog.Check2(headless.NewWindow(size.X, size.Y))
//	gtx := layout.Context{
//		Ops: new(op.Ops),
//		Metric: unit.Metric{
//			PxPerDp: scale,
//			PxPerSp: scale,
//		},
//		Constraints: layout.Exact(size),
//	}
//	mylog.Check(w.Frame(gtx.Ops))
//	for i := 0; i < b.N; i++ {
//		// BenchmarkLabelWidth-8   	  628802	      1854 ns/op 这个似乎不怎么耗时,字符串越长需要的时间越多
//		ux.LabelWidth(gtx, "Hello, world! xxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
//	}
//}
//
//func BenchmarkWalk(b *testing.B) {
//	t := treeTable()
//	for i := 0; i < b.N; i++ {
//		t.Root.Walk() // BenchmarkWalk-8   	46173380	        27.27 ns/op
//	}
//}
//
//func BenchmarkMaxDepth(b *testing.B) {
//	t := treeTable()
//	for i := 0; i < b.N; i++ {
//		t.Root.MaxDepth() // BenchmarkMaxDepth-8   	 8747222	       151.8 ns/op
//	}
//}
//
//type packet struct {
//	Scheme        string        // 请求协议
//	Method        string        // 请求方式
//	Host          string        // 请求主机
//	Path          string        // 请求路径
//	ContentType   string        // 收发都有
//	ContentLength int           // 收发都有
//	Status        string        // 返回的状态码文本
//	Note          string        // 注释
//	Process       string        // 进程
//	PadTime       time.Duration // 请求到返回耗时
//}
//
//func treeTable() *ux.TreeTable[packet] {
//	t := ux.NewTreeTable(packet{})
//	t.TableContext = ux.TableContext[packet]{
//		ContextMenuItems: func(gtx layout.Context, n *ux.Node[packet]) (items []ux.ContextMenuItem) {
//			return []ux.ContextMenuItem{
//				{
//					Title: "delete file",
//					Icon:  nil,
//					Can:   func() bool { return stream.IsFilePath(n.Data.Path) }, // n是当前渲染的行,它的元数据是路径才显示
//					Do: func() {
//						mylog.Check(os.Remove(t.SelectedNode.Data.Path))
//						t.Remove(gtx)
//					},
//					AppendDivider: false,
//					Clickable:     widget.Clickable{},
//				},
//				{
//					Title: "delete directory",
//					Icon:  nil,
//					Can:   func() bool { return stream.IsDir(n.Data.Path) }, // n是当前渲染的行,它的元数据是目录才显示
//					Do: func() {
//						mylog.Check(os.RemoveAll(t.SelectedNode.Data.Path))
//						t.Remove(gtx)
//					},
//					AppendDivider: false,
//					Clickable:     widget.Clickable{},
//				},
//			}
//		},
//		MarshalRowCells: func(n *ux.Node[packet]) (columnCells []ux.CellData) {
//			if n.Container() {
//				n.Data.Scheme = n.SumChildren()
//				sumBytes := 0
//				sumTime := time.Duration(0)
//				n.Data.ContentLength = sumBytes
//				n.Data.PadTime = sumTime
//				for _, node := range n.Walk() {
//					sumBytes += node.Data.ContentLength
//					sumTime += node.Data.PadTime
//				}
//				n.Data.ContentLength = sumBytes
//				n.Data.PadTime = sumTime
//			}
//			return []ux.CellData{
//				{Text: n.Data.Scheme, FgColor: ux.Orange100},
//				{Text: n.Data.Method, FgColor: ux.ColorPink},
//				{Text: n.Data.Host},
//				{Text: n.Data.Path},
//				{Text: n.Data.ContentType},
//				{Text: fmt.Sprintf("%d", n.Data.ContentLength)},
//				{Text: n.Data.Status},
//				{Text: n.Data.Note},
//				{Text: n.Data.Process},
//				{Text: fmt.Sprintf("%s", n.Data.PadTime)},
//			}
//		},
//		UnmarshalRowCells: nil,
//		RowSelectedCallback: func() {
//			mylog.Struct("selected node", t.SelectedNode.Data)
//		},
//		RowDoubleClickCallback: func() {
//			mylog.Info("node:", t.SelectedNode.Data.Path, " double clicked")
//		},
//		LongPressCallback: nil,
//		SetRootRowsCallBack: func() {
//			for i := 0; i < 100; i++ {
//				data := packet{
//					Scheme:        "Row" + fmt.Sprint(i+1),
//					Method:        http.MethodGet,
//					Host:          "example.com",
//					Path:          fmt.Sprintf("/api/v%d/resource", i+1),
//					ContentType:   "application/json",
//					ContentLength: i + 100,
//					Status:        http.StatusText(http.StatusOK),
//					Note:          fmt.Sprintf("获取资源%d", i+1),
//					Process:       fmt.Sprintf("process%d.exe", i+1),
//					PadTime:       time.Duration(i+1) * time.Second,
//				}
//				var node *ux.Node[packet]
//				if i%10 == 3 {
//					node = ux.NewContainerNode(fmt.Sprintf("Row %d", i+1), data)
//					t.Root.AddChild(node)
//					for j := 0; j < 5; j++ {
//						subData := packet{
//							Scheme:        "Row" + fmt.Sprint(j+1),
//							Method:        http.MethodGet,
//							Host:          "example.com",
//							Path:          fmt.Sprintf("/api/v%d/resource%d", i+1, j+1),
//							ContentType:   "application/json",
//							ContentLength: i + 100 + j + 1,
//							Status:        http.StatusText(http.StatusOK),
//							Note:          fmt.Sprintf("获取资源%d-%d", i+1, j+1),
//							Process:       fmt.Sprintf("process%d-%d.exe", i+1, j+1),
//							PadTime:       time.Duration(i+1+j+1) * time.Second,
//						}
//						if j < 2 {
//							subNode := ux.NewContainerNode("Sub Row "+fmt.Sprint(j+1), subData)
//							node.AddChild(subNode)
//							for k := 0; k < 2; k++ {
//								subSubData := packet{
//									Scheme:        "Sub Sub Row" + fmt.Sprint(k+1),
//									Method:        http.MethodGet,
//									Host:          "example.com",
//									Path:          fmt.Sprintf("/api/v%d/resource%d-%d", i+1, j+1, k+1),
//									ContentType:   "application/json",
//									ContentLength: i + 100 + j + 1 + k + 1,
//									Status:        http.StatusText(http.StatusOK),
//									Note:          fmt.Sprintf("获取资源%d-%d-%d", i+1, j+1, k+1),
//									Process:       fmt.Sprintf("process%d-%d-%d.exe", i+1, j+1, k+1),
//									PadTime:       time.Duration(i+1+j+1+k+1) * time.Second,
//								}
//								subSubNode := ux.NewNode(subSubData)
//								subNode.AddChild(subSubNode)
//							}
//						} else {
//							subData.Scheme = "Sub Row" + fmt.Sprint(j+1)
//							subNode := ux.NewNode(subData)
//							node.AddChild(subNode)
//						}
//					}
//				} else {
//					t.Root.AddChild(ux.NewNode(data))
//				}
//			}
//			t.Root.OpenAll()
//			t.Format()
//		},
//		JsonName:   "demo",
//		IsDocument: true,
//	}
//	return t
//}
//
////////////////////////////////
//
//// BenchmarkInsert tests the performance of the Insert function.
//func BenchmarkInsert(b *testing.B) {
//	// Example slice and values to insert.
//	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//	v := []int{100, 101, 102, 103, 104}
//	i := 5
//
//	for n := 0; n < b.N; n++ {
//		sCopy := make([]int, len(s))
//		copy(sCopy, s)
//		slices.Insert(sCopy, i, v...) // BenchmarkInsert-8   	12175423	       104.5 ns/op
//	}
//}
//
//// BenchmarkInsertAtEnd tests the performance of the Insert function when inserting at the end.
//func BenchmarkInsertAtEnd(b *testing.B) {
//	// Example slice and values to insert.
//	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//	v := []int{100, 101, 102, 103, 104}
//	i := len(s)
//
//	for n := 0; n < b.N; n++ {
//		sCopy := make([]int, len(s))
//		copy(sCopy, s)
//		slices.Insert(sCopy, i, v...) // BenchmarkInsertAtEnd-8   	13798304	        90.37 ns/op
//	}
//}
//
//// BenchmarkInsertWithResize tests the performance of the Insert function when resizing is needed.
//func BenchmarkInsertWithResize(b *testing.B) {
//	// Example slice and values to insert.
//	s := make([]int, 10, 10) // Create a slice with length 10 and capacity 10
//	v := []int{100, 101, 102, 103, 104}
//	i := 5
//
//	for n := 0; n < b.N; n++ {
//		sCopy := make([]int, len(s), cap(s))
//		copy(sCopy, s)
//		slices.Insert(sCopy, i, v...) // BenchmarkInsertWithResize-8   	11250734	       105.6 ns/op
//	}
//}
//
//func TestTreeTable_Filter(t1 *testing.T) {
//	t := treeTable()
//	t.SetRootRowsCallBack()
//	t.Filter("ok")
//	return
//	assert.True(t1, strings.EqualFold("row96", "Row96"))
//	fmt.Println(strings.EqualFold("GoLang", "golang"))
//	fmt.Println(strings.EqualFold("golang", "GoLang"))
//}
