package main

import (
	"fmt"
	"iter"
	"strconv"
	"testing"
	"time"

	"gioui.org/layout"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
	"github.com/ddkwork/golibrary/std/stream/deepcopy"
	"github.com/ddkwork/ux"
	"github.com/xuri/excelize/v2"
)

func TestGroupBy(*testing.T) {
	t := ux.NewTreeTable(Info{})
	t.TableContext = ux.TableContext[Info]{
		CustomContextMenuItems: func(gtx layout.Context, n *ux.Node[Info]) iter.Seq[ux.ContextMenuItem] {
			return func(yield func(ux.ContextMenuItem) bool) {

			}
		},
		MarshalRowCells: func(n *ux.Node[Info]) (cells []ux.CellData) {
			return ux.MarshalRow(n.Data, func(key string, field any) (value string) {
				switch key {
				case "日期":
				case "Data":
					if n.Container() {
						value = n.SumChildren()
					} else {
						value = FormatTime(field.(time.Time))
					}
				case "姓名":
					value = field.(string)
				case "第几车":
					value = fmt.Sprintf("%d", field.(int))
				case "发车":
					value = fmt.Sprintf("%d", field.(int))
				case "余货":
					value = fmt.Sprintf("%d", field.(int))
				case "女工日结":
					value = fmt.Sprintf("%d", field.(int))
				case "男工车结":
					value = fmt.Sprintf("%d", field.(int))
				case "注释":
					value = field.(string)
				default:
					value = fmt.Sprintf("%v", field)
				}
				return
			})
		},
		UnmarshalRowCells: func(n *ux.Node[Info], rows []ux.CellData) Info {
			return ux.UnmarshalRow[Info](rows, func(key, value string) (field any) {
				return nil
			})
		},
		RowSelectedCallback: func() {

		},
		RowDoubleClickCallback: func() {

		},
		SetRootRowsCallBack: func() {
			//index1 := ux.NewContainerNode(FormatTime(time.Now()), Info{
			//	IndexCard:    "",
			//	Data:         time.Now(),
			//	Lazhi:        0,
			//	Zhongxiaoguo: 0,
			//	Name:         "",
			//	SendCard:     0,
			//	Women:        0,
			//	Note:         "",
			//	YuHuo:        0,
			//})
			//index1.AddChildByData(Info{
			//	IndexCard:    "",
			//	Data:         time.Now(),
			//	Lazhi:        0,
			//	Zhongxiaoguo: 156,
			//	Name:         "芬",
			//	SendCard:     0,
			//	Women:        0,
			//	Note:         "",
			//	YuHuo:        0,
			//})
			//index1.AddChildByData(Info{
			//	IndexCard:    "",
			//	Data:         time.Now(),
			//	Lazhi:        0,
			//	Zhongxiaoguo: 404,
			//	Name:         "其余三人",
			//	SendCard:     540,
			//	Women:        0,
			//	Note:         "",
			//})
			//
			//t.Root.AddChild(index1)
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/01"),
				Name:      "美人椒",
				JianShu:   122,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
				Note:      "4070-伙食费25=4045元",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/01"),
				Name:      "杨友翠",
				JianShu:   111,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
				Note:      "4097.5-伙食费25=4072.5元",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/01"),
				Name:      "小夏",
				JianShu:   67,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
				Note:      "3190-伙食费75=3115元",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/01"),
				Name:      "张太美",
				JianShu:   103,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
				Note:      "3467.5-伙食费70=3397.5元",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/01"),
				Name:      "刘春燕",
				JianShu:   46,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
				Note:      "伙食费",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/01"),
				Name:      "接手余货",
				JianShu:   199,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
				Note:      "已垫付",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/01"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第1车",
				SendCard:  207,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/01"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第2车",
				SendCard:  232,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/01"),
				Name:      "",
				JianShu:   0,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     209,
				Note:      "47+162=209",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/02"),
				Name:      "美人椒",
				JianShu:   157,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/02"),
				Name:      "杨友翠",
				JianShu:   162,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/02"),
				Name:      "小夏",
				JianShu:   102,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/02"),
				Name:      "张太美",
				JianShu:   162,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/02"),
				Name:      "刘春燕",
				JianShu:   80,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
				Note:      "70+25+25+55+75=250",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/02"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第3车",
				SendCard:  546,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/03"),
				Name:      "美人椒",
				JianShu:   170,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/03"),
				Name:      "杨友翠",
				JianShu:   175,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/03"),
				Name:      "小夏",
				JianShu:   62,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/03"),
				Name:      "张太美",
				JianShu:   81,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
				Note:      "件数少的有可能是合并到4号早上计数，建议对账3-4号合并看",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/03"),
				Name:      "刘春燕",
				JianShu:   44,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
				Note:      "2382.5+250=2632.5",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/03"),
				Name:      "游击队",
				JianShu:   404,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
				Note:      "已垫付",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/03"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第4车",
				SendCard:  540,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/03"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第5车",
				SendCard:  533,
				YuHuo:     0,
				Note:      "走掉的男工结账平分到今天，439+546+540=1525/2=762.5元",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/04"),
				Name:      "美人椒",
				JianShu:   131,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/04"),
				Name:      "杨友翠",
				JianShu:   127,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/04"),
				Name:      "小夏",
				JianShu:   152,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/04"),
				Name:      "张太美",
				JianShu:   199,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/04"),
				Name:      "刘春燕",
				JianShu:   103,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/04"),
				Name:      "游击队",
				JianShu:   176,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
				Note:      "已垫付",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/04"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第6车",
				SendCard:  427,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/04"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第7车",
				SendCard:  540,
				YuHuo:     0,
				Note:      "帮忙的男工压盖131件，装车198件，已支付198元给人家",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/05"),
				Name:      "美人椒",
				JianShu:   94,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/05"),
				Name:      "杨友翠",
				JianShu:   95,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/05"),
				Name:      "小夏",
				JianShu:   104,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/05"),
				Name:      "张太美",
				JianShu:   103,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/05"),
				Name:      "刘春燕",
				JianShu:   87,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/05"),
				Name:      "游击队",
				JianShu:   157,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
				Note:      "已垫付",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/05"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第8车",
				SendCard:  241,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/05"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第9车",
				SendCard:  509,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/07"),
				Name:      "美人椒",
				JianShu:   153,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/07"),
				Name:      "杨友翠",
				JianShu:   163,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/07"),
				Name:      "小夏",
				JianShu:   143,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/07"),
				Name:      "张太美",
				JianShu:   0,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/07"),
				Name:      "刘春燕",
				JianShu:   88,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/07"),
				Name:      "",
				JianShu:   0,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/07"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第10车",
				SendCard:  540,
				YuHuo:     6,
				Note:      "记多了1件",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/08"),
				Name:      "美人椒",
				JianShu:   158,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/08"),
				Name:      "杨友翠",
				JianShu:   166,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/08"),
				Name:      "小夏",
				JianShu:   138,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/08"),
				Name:      "张太美",
				JianShu:   104,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/08"),
				Name:      "刘春燕",
				JianShu:   93,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/08"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第11车",
				SendCard:  138,
				YuHuo:     0,
				Note:      "？",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/08"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第12车",
				SendCard:  530,
				YuHuo:     4,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/09"),
				Name:      "美人椒",
				JianShu:   121,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/09"),
				Name:      "杨友翠",
				JianShu:   132,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/09"),
				Name:      "小夏",
				JianShu:   100,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/09"),
				Name:      "张太美",
				JianShu:   115,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/09"),
				Name:      "刘春燕",
				JianShu:   92,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
				Note:      "金宝72元装车费已付",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/09"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第13车",
				SendCard:  530,
				YuHuo:     31,
				Note:      "记少了1件",
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/10"),
				Name:      "美人椒",
				JianShu:   114,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/10"),
				Name:      "杨友翠",
				JianShu:   113,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/10"),
				Name:      "小夏",
				JianShu:   85,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/10"),
				Name:      "张太美",
				JianShu:   111,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/10"),
				Name:      "刘春燕",
				JianShu:   54,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/10"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第14车",
				SendCard:  383,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/11"),
				Name:      "美人椒",
				JianShu:   181,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/11"),
				Name:      "杨友翠",
				JianShu:   181,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/11"),
				Name:      "小夏",
				JianShu:   142,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/11"),
				Name:      "张太美",
				JianShu:   182,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/11"),
				Name:      "刘春燕",
				JianShu:   110,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/11"),
				Name:      "",
				JianShu:   0,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/11"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第15车",
				SendCard:  540,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/11"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第16车",
				SendCard:  504,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/12"),
				Name:      "美人椒",
				JianShu:   136,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/12"),
				Name:      "杨友翠",
				JianShu:   129,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/12"),
				Name:      "小夏",
				JianShu:   117,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/12"),
				Name:      "张太美",
				JianShu:   141,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/12"),
				Name:      "刘春燕",
				JianShu:   103,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/12"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第17车",
				SendCard:  516,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/14"),
				Name:      "美人椒",
				JianShu:   91,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/14"),
				Name:      "杨友翠",
				JianShu:   85,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/14"),
				Name:      "小夏",
				JianShu:   64,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/14"),
				Name:      "张太美",
				JianShu:   86,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/14"),
				Name:      "刘春燕",
				JianShu:   53,
				IndexCard: "",
				SendCard:  0,
				YuHuo:     0,
			})
			t.Root.AddChildByData(Info{
				Data:      ParseTime("2025/10/14"),
				Name:      "",
				JianShu:   0,
				IndexCard: "第18车",
				SendCard:  379,
				YuHuo:     0,
				Note:      "金宝67件装车费已付",
			})

		},
		JsonName:   "treegrid",
		IsDocument: false,
	}
	t.SetRootRowsCallBack()
	t.OriginalRoot = deepcopy.Clone(t.Root) //todo
	t.SizeColumnsToFit(layout.Context{})
	mylog.Call(func() {
		t.GroupBy("日期")
		t.GroupBy("姓名")
		t.MarkDown()
	})
}

func TestFromXls(t *testing.T) {
	f, err := excelize.OpenFile("(数据表)西昌男女工做账统计表.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("西昌男女工做账统计表")
	if err != nil {
		fmt.Println(err)
		return
	}

	//t.Root.AddChildByData(Info{
	//	Data:         time.Time{},
	//	Name:         "",
	//	Lazhi:        0,
	//	Zhongxiaoguo: 0,
	//	JianShu:      0,
	//	IndexCard:    "",
	//	SendCard:     0,
	//	YuHuo:        0,
	//	Women:        0,
	//	Man:          0,
	//	Note:         "",
	//})
	g := stream.NewGeneratedFile()
	for i, row := range rows {
		if i == 0 {
			continue
		}
		g.P("t.Root.AddChildByData(Info{")
		g.P("Data:         ", "ParseTime("+
			strconv.Quote(row[0])+
			")", ",")
		g.P("Name:         ", strconv.Quote(row[1]), ",")
		//g.P("Lazhi:        ", row[2], ",")
		//g.P("Zhongxiaoguo: ", row[3], ",")
		if row[2] == "" {
			g.P("JianShu:    ", 0, ",")
		} else {
			g.P("JianShu:      ", row[2], ",")
		}

		g.P("IndexCard:    ", strconv.Quote(row[3]), ",")

		if row[4] == "" {
			g.P("SendCard:    ", 0, ",")
		} else {
			g.P("SendCard:     ", row[4], ",")
		}
		if row[5] == "" {
			g.P("YuHuo:    ", 0, ",")
		} else {
			g.P("YuHuo:        ", row[5], ",")
		}
		//g.P("Women:        ", row[6], ",")
		//g.P("Man:          ", row[7], ",")
		if len(row) > 8 {
			g.P("Note:         ", strconv.Quote(row[8]), ",")
		}
		g.P("})")
	}
	stream.WriteTruncate("tmp/erp_test.go", g.String())
	for _, row := range rows {
		for _, colCell := range row {
			//print(stream.AlignString(colCell+"|", 15))
			fmt.Printf("%-15s\t|", colCell)
		}
		fmt.Println()
	}
}
