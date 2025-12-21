package main

import (
	"fmt"
	"iter"
	"time"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/resources/images"
)

type (
	Info struct { // todo加入公告，写入男女工计件价格，工具栏加入 日结和结账按钮，按下面的两个东西分组rows
		Data time.Time `table:"日期"` // 结合每日余货核查记账准确性，日产量是否达标等等
		Name string    `table:"姓名"` // 结账按姓名分组，减去伙食费就是每个女工的工资条，吃饭打卡设计？

		// todi 添加更多的分类，如果价格不同的话，需要解决这里占据太宽的问题，子节点？
		Lazhi        int `table:"蜡纸"`  // 加入大果做成celltype，然后form录入明细？这样看不到分类汇总和小计
		Zhongxiaoguo int `table:"中小果"` // 为了日结公式的使用，价格不同的必须有独立的列？
		JianShu      int `table:"件数"`

		// 胶框小标
		// 胶框无标
		// 胶框小框
		// 串果
		// 泡沫箱
		// 周转框
		// 这么多列需要三折叠屏才合适，或者横屏编辑，需要支持动态删除列，近年来的场景不需要这么多分类，都是统价

		IndexCard string `table:"第几车"` // 实现celltype？自动自增
		SendCard  int    `table:"发车"`  // todo 把分类件数明细，驾驶员姓名，电话，车牌，接货人和电话，接货点等布局一个from作为celltype
		YuHuo     int    `table:"余货"`  // todo 核查计数准确性，自动校验后写入图标或者相差的件数

		Women int `table:"女工日结"` // todo 公式引入
		Man   int `table:"男工车结"` // todo 结账form加入来去路费，平常用车耗油，小费等用于汇总结账需要

		Note string `table:"备注"`
		// todo 底部布局一个每列汇总
		// 垂直布局一个结账form
	}
)

// todo加入截图功能，截屏整个画布，截屏所有女工的日结，截屏结账分类汇总
func main() {
	// stream.GitProxy(true)
	// return
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
		JsonName:      "treegrid",
		CreatMarkdown: false,
	}
	panel := ux.NewPanel()
	hPanel := ux.NewHPanel()
	panel.AddChild(hPanel)

	type toolbarButtons struct {
		redo      *widget.Clickable // 撤销
		next      *widget.Clickable // 前进
		add       *widget.Clickable // 添加记录
		field     *widget.Clickable // 字段管理
		child     *widget.Clickable // 子记录
		filter    *widget.Clickable // 筛选
		sort      *widget.Clickable // 排序
		group     *widget.Clickable // 分组
		check     *widget.Clickable // 按日期分组，核对每日所有发车总件数+余货=女工计数总件数
		bill      *widget.Clickable // 结账,按姓名分组
		rowHeight *widget.Clickable // 行高
		notice    *widget.Clickable // 公告
		find      *widget.Clickable // 查找
		ai        *widget.Clickable // ai分析数据
		remind    *widget.Clickable // 自动提醒
		export    *widget.Clickable // 导出数据
		form      *widget.Clickable // 生成表单
		share     *widget.Clickable // 分享视图
		screen    *widget.Clickable // 截屏

	}
	bar := toolbarButtons{
		redo:      &widget.Clickable{},
		next:      &widget.Clickable{},
		add:       &widget.Clickable{},
		field:     &widget.Clickable{},
		child:     &widget.Clickable{},
		filter:    &widget.Clickable{},
		sort:      &widget.Clickable{},
		group:     &widget.Clickable{},
		check:     &widget.Clickable{},
		bill:      &widget.Clickable{},
		rowHeight: &widget.Clickable{},
		notice:    &widget.Clickable{},
		find:      &widget.Clickable{},
		ai:        &widget.Clickable{},
		remind:    &widget.Clickable{},
		export:    &widget.Clickable{},
		form:      &widget.Clickable{},
		share:     &widget.Clickable{},
		screen:    &widget.Clickable{},
	}
	t.TableContext.GroupCallback = func(gtx layout.Context) {
		if bar.group.Clicked(gtx) { // todo pop window
			t.GroupBy("姓名")
		}
	}

	InitAppBar(hPanel, func(yield func(style ux.ButtonStyle) bool) {
		yield(ux.Button(bar.redo, images.ContentUndoIcon, "撤销"))
		yield(ux.Button(bar.next, images.ContentForwardIcon, "前进"))
		yield(ux.Button(bar.add, images.ContentAddIcon, "添加记录"))
		yield(ux.Button(bar.field, images.ActionSettingsIcon, "字段管理"))
		yield(ux.Button(bar.child, images.ContentAddIcon, "子记录"))
		yield(ux.Button(bar.filter, images.ImageFilterTiltShiftIcon, "筛选"))
		yield(ux.Button(bar.sort, images.ContentSortIcon, "排序"))
		yield(ux.Button(bar.group, images.ActionGroupWorkIcon, "分组"))
		yield(ux.Button(bar.check, images.NavigationCheckIcon, "日结"))
		yield(ux.Button(bar.bill, images.ActionRecordVoiceOverIcon, "结账"))
		yield(ux.Button(bar.rowHeight, images.ActionHighlightOffIcon, "行高"))
		yield(ux.Button(bar.notice, images.MapsSatelliteIcon, "公告"))
		yield(ux.Button(bar.find, images.ActionSearchIcon, "查找"))
		yield(ux.Button(bar.ai, images.SvgIconDatabase, "ai分析数据"))
		yield(ux.Button(bar.remind, images.ToggleCheckBoxIcon, "自动提醒"))
		yield(ux.Button(bar.export, images.CommunicationImportExportIcon, "导出数据"))
		yield(ux.Button(bar.form, images.ContentTextFormatIcon, "生成表单"))
		yield(ux.Button(bar.share, images.FileFolderSharedIcon, "分享视图"))
		yield(ux.Button(bar.screen, images.DeviceScreenRotationIcon, "截屏"))
	})
	panel.AddChild(t)
	ux.Run("智能表格编辑器", panel)
}

type AppBar struct {
	Search *ux.Input
	About  ux.ButtonStyle
}

func InitAppBar(hPanel *ux.Panel, toolBars iter.Seq[ux.ButtonStyle]) *AppBar {
	search := ux.NewInput("请输入搜索关键字...").SetIcon(images.ActionSearchIcon).SetRadius(16)
	hPanel.AddChildFlexed(1, search) // todo 太多之后apk需要管理溢出

	for toolbar := range toolBars {
		toolbar.Border.Width = 0.3
		hPanel.AddChild(toolbar)
	}
	return &AppBar{
		Search: search,
		// About:  about,
	}
}

const timeLayout = `2006/01/02`

func FormatTime(t time.Time) string {
	return t.Format(timeLayout)
}

func ParseTime(t string) time.Time {
	return mylog.Check2(time.Parse(timeLayout, t))
}
