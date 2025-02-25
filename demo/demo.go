package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"net/http"
	"os"
	"time"

	"github.com/ddkwork/golibrary/stream"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/outlay"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/terminal"
)

var appBar *ux.AppBar

func main() {
	// stream.FileServerWindowsDisk()
	// stream.FileServer()
	// 安卓的类似功能是文件闪传，上传速度只有FileServerWindowsDisk的1%左右
	// 这种本地起文件服务的方式对方下载很快

	w := ux.NewWindow("gio demo")
	panel := ux.NewPanel(w)

	hPanel := ux.NewHPanel(w)
	panel.AddChild(hPanel.Layout)

	tipIconButtons := []*ux.TipIconButton{
		ux.NewTooltipButton(ux.IconBack, "action code", nil),
		ux.NewTooltipButton(ux.IconFavorite, "action code", nil),
		ux.NewTooltipButton(ux.IconDone, "action bug report", nil),
		ux.NewTooltipButton(ux.IconDelete, "action build", nil),
		ux.NewTooltipButton(ux.IconClose, "action build", nil),
		ux.NewTooltipButton(ux.IconArrowDropDown, "action build", nil),
		ux.NewTooltipButton(ux.IconNaviLeft, "action build", nil),
		ux.NewTooltipButton(ux.IconNaviRight, "action build", nil),
		ux.NewTooltipButton(ux.IconFileFolder, "action build", nil),
		ux.NewTooltipButton(ux.IconUpload, "action build", nil),
		ux.NewTooltipButton(ux.IconDownload, "action build", nil),
	}

	appBar = ux.InitAppBar(hPanel, tipIconButtons, speechTxt)
	appBar.Search.SetOnChanged(func(text string) {
		println(text)
	})

	m := new(safemap.M[DemoType, ux.Widget])
	for _, Type := range TreeTableType.EnumTypes() {
		switch Type {
		case TreeTableType:
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
			t := ux.NewTreeTable(packet{})
			t.TableContext = ux.TableContext[packet]{
				ContextMenuItems: func(gtx layout.Context, n *ux.Node[packet]) (items []ux.ContextMenuItem) {
					return []ux.ContextMenuItem{
						{
							Title: "delete file",
							Icon:  nil,
							Can:   func() bool { return stream.IsFilePath(n.Data.Path) }, // n是当前渲染的行,它的元数据是路径才显示
							Do: func() {
								mylog.Check(os.Remove(t.SelectedNode.Data.Path))
								t.Remove(gtx)
							},
							AppendDivider: false,
							Clickable:     widget.Clickable{},
						},
						{
							Title: "delete directory",
							Icon:  nil,
							Can:   func() bool { return stream.IsDir(n.Data.Path) }, // n是当前渲染的行,它的元数据是目录才显示
							Do: func() {
								mylog.Check(os.RemoveAll(t.SelectedNode.Data.Path))
								t.Remove(gtx)
							},
							AppendDivider: false,
							Clickable:     widget.Clickable{},
						},
					}
				},
				MarshalRowCells: func(n *ux.Node[packet]) (cells []ux.CellData) {
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
						{Text: n.Data.Scheme, FgColor: ux.Orange100},
						{Text: n.Data.Method, FgColor: ux.ColorPink},
						{Text: n.Data.Host},
						{Text: n.Data.Path},
						{Text: n.Data.ContentType},
						{Text: fmt.Sprintf("%d", n.Data.ContentLength)},
						{Text: n.Data.Status},
						{Text: n.Data.Note},
						{Text: n.Data.Process},
						{Text: fmt.Sprintf("%s", n.Data.PadTime)},
					}
				},
				UnmarshalRowCells: func(n *ux.Node[packet], values []ux.CellData) {
					mylog.Todo("unmarshal row cells for edit node")
				},
				RowSelectedCallback: func() {
					mylog.Struct("selected node", t.SelectedNode.Data)
				},
				RowDoubleClickCallback: func() {
					mylog.Info("node:", t.SelectedNode.Data.Path, " double clicked")
				},
				LongPressCallback: nil,
				SetRootRowsCallBack: func() {
					for i := 0; i < 100; i++ {
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
						var node *ux.Node[packet]
						if i%10 == 3 {
							node = ux.NewContainerNode(fmt.Sprintf("Row %d", i+1), data)
							t.Root.AddChild(node)
							for j := 0; j < 5; j++ {
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
									node.AddChild(subNode)
									for k := 0; k < 2; k++ {
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
									node.AddChild(subNode)
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
			// t.SetRootRowsCallBack()//已经在layout内once一次，避免每个实例都要写一遍
			appBar.Search.SetOnChanged(func(text string) {
				// todo 这里可以设计一个类似aggrid的高级搜索功能：把n叉树的元数据结构体取出来，然后通过反射结构体布局一个所有字段值的过滤综合条件，最后设置过滤结果填充到表格的过滤rows中
				t.Filter(text)
			})
			m.Set(TreeTableType, t.Layout)
		case TreeType:
			rootNodes := []*ux.TreeNode{
				{
					Title: "Root 0",
				},
				{
					Title: "Root 1",
					Children: []*ux.TreeNode{
						{
							Title: "Child 1.1",
							Children: []*ux.TreeNode{
								{Title: "Grandchild 1.1.1"},
								{Title: "Grandchild 1.1.2"},
							},
							ClickCallback: func(gtx layout.Context) {
								println("clicked")
							},
						},
						{
							Title: "Child 1.2",
							Children: []*ux.TreeNode{
								{Title: "Grandchild 1.2.1"},
							},
						},
					},
				},

				{
					Title: "Root 2",
					Children: []*ux.TreeNode{
						{
							Title: "Child 2.1",
							Children: []*ux.TreeNode{
								{Title: "Grandchild 2.1.1"},
							},
						},
					},
				},
				{
					Title: "Root 3",
				},
			}
			tree := ux.NewTree()
			tree.SetNodes(rootNodes)
			tree.OnClick(func(gtx layout.Context, node *ux.TreeNode) {
				fmt.Println("node:", node.Title, " clicked")
			})
			m.Set(TreeType, tree.Layout)
		case Table2Type:
			t := table2(Packets)
			m.Set(Table2Type, t.Layout)
		case TableType:
			type BasicFan struct {
				ID          int     `json:"id"`          // id
				CompanyID   int     `json:"companyId"`   // 公司id
				CompanyName string  `json:"companyName"` // 公司名称
				PlantID     int     `json:"plantId"`     // 场站id
				PlantName   string  `json:"plantName"`   // 场站名称
				StagingID   int     `json:"stagingId"`   // 工期id
				StagingName string  `json:"stagingName"` // 工期名称
				CircuitID   int     `json:"circuitId"`   // 集电线id
				CircuitName string  `json:"circuitName"` // 集电线名称
				FanName     string  `json:"fanName"`     // 风机名称
				PowerField  string  `json:"powerField"`  // 电量计算的原始点
				FanCode     string  `json:"fanCode"`     // 风机编码
				InnerCode   string  `json:"innerCode"`   // 内部编码
				ModelID     int     `json:"modelId"`     // 型号id
				ModelName   string  `json:"modelName"`   // 型号名称
				Status      int     `json:"status"`      // 1 运行 2 调试 3 未接入
				StartSpeed  float64 `json:"startSpeed"`  // 切入风速(m/s)
				StopSpeed   float64 `json:"stopSpeed"`   // 切出风速(m/s)
				FanCap      float64 `json:"fanCap"`      // 装机容量
				// Host        string  `json:"host"`

				IsParadigm   int    `json:"isParadigm"`   // 是否是标杆
				FanLocalType string `json:"fanLocalType"` // fan_local_type 海风陆风
			}
			var fans []*BasicFan
			for i := 0; i < 100*10000; i++ {
				fans = append(fans, &BasicFan{ID: i + 1, CompanyID: 1, CompanyName: "company1", PlantName: "plant1", StagingName: "staging1", CircuitName: "circuit1", FanName: "#1风机"})
			}
			datatable := ux.NewDataTable(fans, nil, nil)
			table := ux.NewTable(datatable)
			table.SelectionChangedCallback = func(gtx layout.Context, row, col int) {
				data := table.GetRow(row)
				mylog.Struct("todo", data) // todo check data
			}

			table.DoubleClickCallback = func(gtx layout.Context, row, col int) {
				// data := table.GetRow(row)
				// mylog.Struct("todo",data)
				mylog.Info("double click for edit row data")
			}

			contextMenu := ux.NewContextMenu()
			contextMenu.AddItem(ux.ContextMenuItem{
				Title: "addRow",
				//Icon:  ux.IconAdd,//todo
				Can: func() bool {
					return true
				},
				Do: func() {
					mylog.Info("add row")
				},
				AppendDivider: false,
				Clickable:     widget.Clickable{},
			})
			contextMenu.AddItem(ux.ContextMenuItem{
				Title: "deleteRow",
				//Icon:  ux.IconDelete,//todo
				Can: func() bool {
					return true
				},
				Do: func() {
					mylog.Info("delete row")
				},
				AppendDivider: true,
				Clickable:     widget.Clickable{},
			})

			table.SetMenu(contextMenu)

			//header := ux.NewContainer(w)
			//button2 := ux.NewButton(w, "滚动到第100行", func() {
			//	//table.GridState.Vertical.Last
			//	log.Println(table.Vertical.First, table.VScrollbar.ScrollDistance(), table.Vertical.OffsetAbs, table.Vertical.Length)
			//	table.Vertical.First = 0
			//	table.Vertical.Offset = (table.Vertical.Length / (table.Size())) * (100 - 1)
			//})
			//header.Add(layout.Rigid(button2.Layout))
			//header.Add(layout.Rigid(button2.Layout))
			//header.Add(layout.Rigid(button2.Layout))
			m.Set(TableType, table.Layout)
		case SearchDropDownType:
			dropDown := ux.NewSearchDropDown()
			dropDown.SetOnChanged(func(value string) {
				println(dropDown.GetSelected())
			})
			dropDown.SetWidth(unit.Dp(300))
			dropDown.SetOptions([]*ux.SearchDropDownOption{
				{
					Text:       "aa",
					Value:      "xx",
					Identifier: "f",
					Icon:       nil,
					IconColor:  color.NRGBA{},
				},
				{
					Text:       "bb",
					Value:      "yy",
					Identifier: "f",
					Icon:       nil,
					IconColor:  color.NRGBA{},
				},
				{
					Text:       "cc",
					Value:      "zz",
					Identifier: "f",
					Icon:       nil,
					IconColor:  color.NRGBA{},
				},
			})
			m.Set(SearchDropDownType, dropDown.Layout)
		case IconvgViewType:
			m.Set(IconvgViewType, ux.NewIconView().Layout)
		case StructViewType:
			type Object struct {
				MachineID string
				RegCode   string
				Version   string
				Website   string
				SimpleMid string

				SimpleMid1  string
				SimpleMid2  string
				SimpleMid3  string
				SimpleMid4  string
				SimpleMid5  string
				SimpleMid6  string
				SimpleMid7  string
				SimpleMid8  string
				SimpleMid9  string
				SimpleMid10 string
				SimpleMid11 string
				SimpleMid12 string
				SimpleMid13 string
				SimpleMid14 string
				SimpleMid15 string
				SimpleMid16 string
				SimpleMid17 string
			}

			object := Object{
				MachineID:   "1111-2222-3333-4444",
				RegCode:     "aaa-bbb-ccc-ddd",
				Version:     "1.1.1",
				Website:     "https://www.baidu.com",
				SimpleMid:   "2222-3333-4444-5555",
				SimpleMid1:  "",
				SimpleMid2:  "",
				SimpleMid3:  "",
				SimpleMid4:  "",
				SimpleMid5:  "",
				SimpleMid6:  "",
				SimpleMid7:  "",
				SimpleMid8:  "",
				SimpleMid9:  "",
				SimpleMid10: "",
				SimpleMid11: "",
				SimpleMid12: "",
				SimpleMid13: "",
				SimpleMid14: "",
				SimpleMid15: "",
				SimpleMid16: "",
				SimpleMid17: "",
			}

			form := ux.NewStructView(Object{}, func() (elems []ux.CellData) {
				return []ux.CellData{
					{Text: object.MachineID},
					{Text: object.RegCode},
					{Text: object.Version},
					{Text: object.Website},
					{Text: object.SimpleMid},
					{Text: object.SimpleMid1},
					{Text: object.SimpleMid2},
					{Text: object.SimpleMid3},
					{Text: object.SimpleMid4},
					{Text: object.SimpleMid5},
					{Text: object.SimpleMid6},
					{Text: object.SimpleMid7},
					{Text: object.SimpleMid8},
					{Text: object.SimpleMid9},
					{Text: object.SimpleMid10},
					{Text: object.SimpleMid11},
					{Text: object.SimpleMid12},
					{Text: object.SimpleMid13},
					{Text: object.SimpleMid14},
					{Text: object.SimpleMid15},
					{Text: object.SimpleMid16},
					{Text: object.SimpleMid17},
				}
			})
			userName := ux.NewInput("please input username")
			password := ux.NewInput("please input password")
			email := ux.NewInput("please input email")

			form.Add("username", userName.Layout)
			form.Add("password", password.Layout)
			form.Add("email", email.Layout)
			dropDown := ux.NewDropDown(SuperRecovery2Type.Names()...)
			form.InsertAt(0, "choose a app", dropDown.Layout)
			// form.Add("", ux.BlueButton(&clickable, "submit", unit.Dp(100)).Layout)
			m.Set(StructViewType, form.Layout)
		case ColorPickerType:
			m.Set(ColorPickerType, ux.NewColorPicker().Layout)
		case CardType:
			f := &ux.FlowWrap{
				Cards: nil,
				List: widget.List{
					Scrollbar: widget.Scrollbar{},
					List: layout.List{
						// Axis:        layout.Vertical,
						ScrollToEnd: false,
						Alignment:   0,
						Position:    layout.Position{},
					},
				},
				Wrap:       outlay.FlowWrap{},
				Contextual: nil,
				Loaded:     false,
			}
			m.Set(CardType, f.Layout)
		case MobileType:
		case SvgButtonType:
			m.Set(SvgButtonType, ux.NewButton("Hex Editor", nil).SetRectIcon(true).SetSVGIcon(ux.CircledChevronRight).Layout)
		case CodeEditorType:
			m.Set(CodeEditorType, ux.NewCodeEditor(tabGo, ux.CodeLanguageGolang).Layout)
		case AsmViewType:
		case LogViewType:
			m.Set(LogViewType, ux.LogView()) // todo 日志没有对齐，控制台是对齐的，增加滚动条
		case ComBoxType:
			// m.Set(ComBoxType, combox(w))//newselect
		case SplitViewType:
			sp := ux.NewSplit(ux.Split{
				Ratio:  0, // 布局比例，0 表示居中，-1 表示完全靠左，1 表示完全靠右
				Bar:    10,
				Axis:   layout.Horizontal,
				First:  ux.NewCodeEditor(tabGo, ux.CodeLanguageGolang).Layout,
				Second: ux.NewCodeEditor(tabGo, ux.CodeLanguageGolang).Layout,
			})
			m.Set(SplitViewType, sp.Layout)
		case ListViewType:
		case JsonTreeType:
		case AnimationButtonType:
			newButtonAnimation := ux.NewButtonAnimation("animation button", ux.IconBack, func() {
				mylog.Info("animation button clicked")
			})
			m.Set(AnimationButtonType, newButtonAnimation.Layout) // todo bug
		case TerminalType: // todo 控制台被接管了
			if mylog.IsWindows() {
				continue // todo bug
			}
			screen, settings := terminal.Demo()
			m.Set(TerminalType, ux.NewTabItem("Tab 5", func(gtx layout.Context) layout.Dimensions {
				return terminal.Console(screen, settings)(gtx)
			}).LayoutContent)
		case StackViewType: // todo stackview
		case DockViewType: // todo dockview
		case Gif123Type: // todo gif123
		case HexEditorType: // todo hex editor
		case ContextMenuType: // todo contextmenu
			m.Set(ContextMenuType, NewMenuObj().Layout)
		case ImageEditorType: // todo 图片编辑器
		case MediaPlayerType: // todo 媒体播放器
		case MindType: // todo 思维导图
		case PdfViewType: // todo pdf
		case MapViewType: // todo 地图
		case ThemeViewType: // todo 主题编辑器
		case SettingsviewType:
		case SliceviewType: // todo 切片器
		case XyzViewType:
		case WebViewType: // todo webview
		case SvgViewType: // todo svg
		case CanvasViewType:
		case PopMenuType:
		case TooltipType: // todo tooltip
		case TextfieldType: // todo textfield
		case MarkdownViewType: // todo markdown
		case GomitmproxyType: // todo gomitmproxy
		case HyperDbgType: // todo hyperdbg
		case VstartType: // todo vstart
		case ExplorerType: // todo 文件管理器
		case DesignerType:
		case AiChatType: // todo 机器人聊天
		case EncodingTestType: // todo 编码测试
		case GameControlFaceType: // todo 游戏控制器面板
		case GithubType: // todo github
		case GhipsType: // todo ghips
		case TaskManagerType: // todo 任务管理器
		case GitlabType:
		case SteamType: // todo steam
		case BuyTomatoesType: // todo 番茄
		case CcType:
		case CryptType: // todo 加密解密
		case DatabaseType: // todo 数据库
		case DatarecoveryType: // todo 数据恢复
		case HardInfoHookType: // todo 硬件信息
		case HardwareIndoType: // todo 硬件信息
		case DriverToolType: // todo 驱动工具
		case EnvironmentType: // todo 环境变量
		case ErpType: // todo 电子商务
		case FleetType: // todo 代码编辑器
		case ImageConvertType:
		case JetbraType: // todo jb crack
		case JiakaobaodianType:
		case ManPieceworkType:
		case MypanType: // todo 网盘上传下载
		case NetAdapterType: // todo 网络适配器
		case NetScanType: // todo 网络扫描
		case VisualStudiokitType: // todo visual studio kit,cmake generator
		case C2goType: // todo c2go
		case VncType:
		case TodoListType: // todo 待办事项
		case DropFileType: // todo 拖拽文件
		case DarkThemeType:
		}
	}

	vtab := ux.NewTabView(layout.Vertical)
	for k, v := range m.Range() {
		tab := ux.NewTabItem(k.String(), v)
		vtab.AddTab(tab)
	}
	// mylog.Success("test append log")
	// mylog.Warning("test append log")
	// mylog.Trace("test append log")

	// htab := ux.NewTabView(layout.Horizontal)

	// buildTabItems(htab)
	// buildTabItems(vtab)

	// panel.AddChild(htab.Layout)
	panel.AddChild(vtab.Layout)

	//app.FileDropCallback(func(files []string) {
	//	for _, file := range files {
	//		println(file)
	//	}
	//})

	ux.Run(panel)
}

//go:embed demo.go
var tabGo string

//go:embed speech.txt
var speechTxt string

// https://github.com/tibold/svg-explorer-extension
// https://github.com/gio-eui/md3-icons
// D:\workspace\workspace\branch\gui\packaging
/*
```go
    import "github.com/gio-eui/md3-icons/icons/toggle/check_box"

    var CheckBox *widget.Icon
    CheckBox, _ = widget.NewIcon(mdiToggleCheckBox.Ivg)
*/
