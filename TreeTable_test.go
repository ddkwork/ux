package ux_test

import (
	"fmt"
	"image"
	"iter"
	"net/http"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"gioui.org/gpu/headless"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/std/assert"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/resources/colors"
)

func TestName(t *testing.T) {
	if stream.IsRunningOnGitHubActions() {
		return
	}
	mylog.Call(func() {
		stream.UpdateAllLocalRep()
		// fakeError.Walk("")
	})
}

func BenchmarkTransposeMatrix(b *testing.B) { // todo
	columnCells := ux.InitHeader(packet{})
	rows := make([][]packet, 0, len(columnCells))
	for range 1000000 {
		rows = append(rows, []packet{
			{
				Scheme: "http",
				Method: "GET",
			},
			{
				Host:          "www.example.com",
				Path:          "/path/to/file.html",
				ContentType:   "text/html",
				ContentLength: 1024,
				Status:        "200 OK",
				Note:          "This is a note",
				Process:       "chrome.exe",
				PadTime:       10 * time.Millisecond,
			},
		})
	}
	for b.Loop() {
		for range ux.TransposeMatrix[packet](rows) {

		}
	}
}

func BenchmarkXX(b *testing.B) {
	matrix := generateMatrix(1000, 1000)
	for b.Loop() {
		ux.TransposeMatrix[int](matrix)
	}
}

func generateMatrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}
	for i := range rows {
		for j := range cols {
			matrix[i][j] = i*cols + j
		}
	}
	return matrix
}

func Benchmark_SizeColumnsToFit(b *testing.B) { // todo
	const scale = 1.5
	size := image.Point{X: 1200 * scale, Y: 600 * scale}
	w := mylog.Check2(headless.NewWindow(size.X, size.Y))
	gtx := layout.Context{
		Ops: new(op.Ops),
		Metric: unit.Metric{
			PxPerDp: scale,
			PxPerSp: scale,
		},
		Constraints: layout.Exact(size),
	}
	mylog.Check(w.Frame(gtx.Ops))
	t := tableDemo()
	// println(t.String())
	// return
	for b.Loop() {
		t.SizeColumnsToFit(gtx)
	}
}

func BenchmarkLabelWidth(b *testing.B) { // todo
	const scale = 1.5
	size := image.Point{X: 1200 * scale, Y: 600 * scale}
	w := mylog.Check2(headless.NewWindow(size.X, size.Y))
	gtx := layout.Context{
		Ops: new(op.Ops),
		Metric: unit.Metric{
			PxPerDp: scale,
			PxPerSp: scale,
		},
		Constraints: layout.Exact(size),
	}
	mylog.Check(w.Frame(gtx.Ops))
	for b.Loop() {
		// BenchmarkLabelWidth-8   	  628802	      1854 ns/op 这个似乎不怎么耗时,字符串越长需要的时间越多
		ux.LabelWidth(gtx, "Hello, world! xxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	}
}

func BenchmarkWalk(b *testing.B) {
	t := tableDemo()
	for b.Loop() {
		t.Root.Walk() // BenchmarkWalk-8   	46173380	        27.27 ns/op
	}
}

func BenchmarkMaxDepth(b *testing.B) { // todo
	t := tableDemo()
	for b.Loop() {
		t.MaxDepth() // BenchmarkMaxDepth-8   	 8747222	       151.8 ns/op
	}
}

type packet struct {
	Scheme        string        // 请求协议
	Method        string        // 请求方式
	Host          string        // 请求主机
	Path          string        // 请求路径
	ContentType   string        // 收发都有
	ContentLength int           // 收发都有
	Status        string        // 返回的状态码文本
	Note          string        // 注释
	Process       string        // 进程
	PadTime       time.Duration // 请求到返回耗时
}

func tableDemo() *ux.TreeTable[packet] {
	t := ux.NewTreeTable(packet{})
	t.TableContext = ux.TableContext[packet]{
		CustomContextMenuItems: func(gtx layout.Context, n *ux.Node[packet]) iter.Seq[ux.ContextMenuItem] {
			return func(yield func(ux.ContextMenuItem) bool) {
				yield(ux.ContextMenuItem{
					Title: "delete file",
					Icon:  nil,
					Can:   func() bool { return stream.IsFilePath(n.Data.Path) }, // n是当前渲染的行,它的元数据是路径才显示
					Do: func() {
						mylog.Check(os.Remove(t.SelectedNode.Data.Path))
						t.Remove(gtx)
					},
					AppendDivider: false,
					Clickable:     widget.Clickable{},
				})
				yield(ux.ContextMenuItem{
					Title: "delete directory",
					Icon:  nil,
					Can:   func() bool { return stream.IsDir(n.Data.Path) }, // n是当前渲染的行,它的元数据是目录才显示
					Do: func() {
						mylog.Check(os.RemoveAll(t.SelectedNode.Data.Path))
						t.Remove(gtx)
					},
					AppendDivider: false,
					Clickable:     widget.Clickable{},
				})
			}
		},
		MarshalRowCells: func(n *ux.Node[packet]) (columnCells []ux.CellData) {
			if n.Container() {
				n.Data.Scheme = n.SumChildren()
				sumBytes := 0
				sumTime := time.Duration(0)
				n.Data.ContentLength = sumBytes
				n.Data.PadTime = sumTime
				for _, node := range n.Walk() {
					sumBytes += node.Data.ContentLength
					sumTime += node.Data.PadTime
				}
				n.Data.ContentLength = sumBytes
				n.Data.PadTime = sumTime
			}
			return []ux.CellData{
				{Value: n.Data.Scheme, FgColor: colors.Orange100},
				{Value: n.Data.Method, FgColor: colors.ColorPink},
				{Value: n.Data.Host},
				{Value: n.Data.Path},
				{Value: n.Data.ContentType},
				{Value: fmt.Sprintf("%d", n.Data.ContentLength)},
				{Value: n.Data.Status},
				{Value: n.Data.Note},
				{Value: n.Data.Process},
				{Value: fmt.Sprintf("%s", n.Data.PadTime)},
			}
		},
		UnmarshalRowCells: nil,
		RowSelectedCallback: func() {
			mylog.Struct(t.SelectedNode.Data)
		},
		RowDoubleClickCallback: func() {
			mylog.Info("node:", t.SelectedNode.Data.Path, " double clicked")
		},
		SetRootRowsCallBack: func() {
			for i := range 100 {
				data := packet{
					Scheme:        "Row" + fmt.Sprint(i+1),
					Method:        http.MethodGet,
					Host:          "example.com",
					Path:          fmt.Sprintf("/api/v%d/resource", i+1),
					ContentType:   "application/json",
					ContentLength: i + 100,
					Status:        http.StatusText(http.StatusOK),
					Note:          fmt.Sprintf("获取资源%d", i+1),
					Process:       fmt.Sprintf("process%d.exe", i+1),
					PadTime:       time.Duration(i+1) * time.Second,
				}
				var parent *ux.Node[packet]
				if i%10 == 3 {
					parent = ux.NewContainerNode(fmt.Sprintf("Row %d", i+1), data)
					t.Root.AddChild(parent)
					for j := range 5 {
						subData := packet{
							Scheme:        "Row" + fmt.Sprint(j+1),
							Method:        http.MethodGet,
							Host:          "example.com",
							Path:          fmt.Sprintf("/api/v%d/resource%d", i+1, j+1),
							ContentType:   "application/json",
							ContentLength: i + 100 + j + 1,
							Status:        http.StatusText(http.StatusOK),
							Note:          fmt.Sprintf("获取资源%d-%d", i+1, j+1),
							Process:       fmt.Sprintf("process%d-%d.exe", i+1, j+1),
							PadTime:       time.Duration(i+1+j+1) * time.Second,
						}
						if j < 2 {
							subNode := ux.NewContainerNode("Sub Row "+fmt.Sprint(j+1), subData)
							parent.AddChild(subNode)
							for k := range 2 {
								subSubData := packet{
									Scheme:        "Sub Sub Row" + fmt.Sprint(k+1),
									Method:        http.MethodGet,
									Host:          "example.com",
									Path:          fmt.Sprintf("/api/v%d/resource%d-%d", i+1, j+1, k+1),
									ContentType:   "application/json",
									ContentLength: i + 100 + j + 1 + k + 1,
									Status:        http.StatusText(http.StatusOK),
									Note:          fmt.Sprintf("获取资源%d-%d-%d", i+1, j+1, k+1),
									Process:       fmt.Sprintf("process%d-%d-%d.exe", i+1, j+1, k+1),
									PadTime:       time.Duration(i+1+j+1+k+1) * time.Second,
								}
								subSubNode := ux.NewNode(subSubData)
								subNode.AddChild(subSubNode)
							}
						} else {
							subData.Scheme = "Sub Row" + fmt.Sprint(j+1)
							subNode := ux.NewNode(subData)
							parent.AddChild(subNode)
						}
					}
				} else {
					t.Root.AddChild(ux.NewNode(data))
				}
			}
			t.Root.OpenAll()
		},
		JsonName:   "demo",
		IsDocument: true,
	}

	const scale = 1.5
	size := image.Point{X: 1200 * scale, Y: 600 * scale}
	gtx := layout.Context{
		Ops: new(op.Ops),
		Metric: unit.Metric{
			PxPerDp: scale,
			PxPerSp: scale,
		},
		Constraints: layout.Exact(size),
	}
	t.Layout(gtx)
	return t
}

// ////////////////////////////

// BenchmarkInsert tests the performance of the Insert function.
func BenchmarkInsert(b *testing.B) {
	// Example slice and values to insert.
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	v := []int{100, 101, 102, 103, 104}
	i := 5

	for b.Loop() {
		sCopy := make([]int, len(s))
		copy(sCopy, s)
		slices.Insert(sCopy, i, v...) // BenchmarkInsert-8   	12175423	       104.5 ns/op
	}
}

// BenchmarkInsertAtEnd tests the performance of the Insert function when inserting at the end.
func BenchmarkInsertAtEnd(b *testing.B) {
	// Example slice and values to insert.
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	v := []int{100, 101, 102, 103, 104}
	i := len(s)

	for b.Loop() {
		sCopy := make([]int, len(s))
		copy(sCopy, s)
		slices.Insert(sCopy, i, v...) // BenchmarkInsertAtEnd-8   	13798304	        90.37 ns/op
	}
}

// BenchmarkInsertWithResize tests the performance of the Insert function when resizing is needed.
func BenchmarkInsertWithResize(b *testing.B) {
	// Example slice and values to insert.
	s := make([]int, 10, 10) // Create a slice with length 10 and capacity 10
	v := []int{100, 101, 102, 103, 104}
	i := 5

	for b.Loop() {
		sCopy := make([]int, len(s), cap(s))
		copy(sCopy, s)
		slices.Insert(sCopy, i, v...) // BenchmarkInsertWithResize-8   	11250734	       105.6 ns/op
	}
}

func TestTreeTable_Filter(t1 *testing.T) { // todo
	t := tableDemo()
	t.SetRootRowsCallBack()
	t.Filter("ok")
	return
	assert.True(t1, strings.EqualFold("row96", "Row96"))
	fmt.Println(strings.EqualFold("GoLang", "golang"))
	fmt.Println(strings.EqualFold("golang", "GoLang"))
}

func TestTreeTable_MaxDepth(t1 *testing.T) {
	demo := tableDemo()
	assert.Equal(t1, demo.MaxDepth(), 3)
}

func TestNode_Depth(t *testing.T) {
	demo := tableDemo()
	for _, n := range demo.Root.Walk() {
		if n.Data.Process == "process4-1-1.exe" {
			assert.Equal(t, n.Depth(), 3)
		}
	}
}

func TestTreeTable_updateRowNumber(t1 *testing.T) {
	demo := tableDemo()
	// mylog.Struct(demo.Root.LastChild())
	assert.Equal(t1, demo.Root.LastChild().RowNumber, ux.CountTableRows(demo.RootRows()))
}
