package ux

import (
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/ddkwork/golibrary/mylog"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type Tabler interface {
	// 获取标题
	GetTitle(i int) (title string)
	// 获取单元格数据
	GetItemText(record any, row, col int) (text string)
	// 获取列宽度
	GetColumnWitdh(i int) (width float32)
	// 获取列属性
	GetColumn(i int) *Column
	GetRow(row int) any
	// 获取列个数
	GetColumnCount() (count int)
	// 获取行数
	Size() (size int)
}

/*
func (t *Table) LayoutHoverTable(gtx layout.Context) layout.Dimensions {
	if len(t.data) == 0 {
		return layout.Dimensions{}
	}
	inset := layout.UniformInset(unit.Dp(2))
	orig := gtx.Constraints
	gtx.Constraints.Min = image.Point{}
	macro := op.Record(gtx.Ops)
	dims := inset.Layout(gtx, layout.Spacer{Height: t.height}.Layout)
	_ = macro.Stop()
	gtx.Constraints = orig
	if t.headerFun == nil {
		t.headerFun = func(gtx layout.Context, index int) layout.Dimensions {
			utils.DrawBackground(gtx, layout.Spacer{}.Layout(gtx).Size, t.theme.Color.TableHeaderBgColor)
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return CellsLabel(t.theme, t.headers[index], true).Layout(gtx)
			})
		}
	}
	if t.dataFun == nil {
		t.dataFun = func(gtx layout.Context, row, col int) layout.Dimensions {
			c := &t.dataContent[row]
			return c.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				if c.Hovered() {
					utils.DrawBackground(gtx, layout.Spacer{}.Layout(gtx).Size, t.theme.Color.DefaultContentBgGrayColor)
				} else {
					utils.DrawBackground(gtx, layout.Spacer{}.Layout(gtx).Size, t.theme.Color.DefaultWindowBgGrayColor)
				}
				NewLine(t.theme).Line(gtx, f32.Pt(0, 0), f32.Pt(float32(gtx.Constraints.Max.X), 0)).Layout(gtx)
				labelDims := layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return CellsLabel(t.theme, fmt.Sprint(t.data[row][col])).Layout(gtx)
				})
				return labelDims
			})
		}
	}
	return component.Table(t.theme.Theme(), &t.grid).Layout(gtx, len(t.data), len(t.data[0]),
		func(axis layout.Axis, index, constraint int) int {
			switch axis {
			case layout.Horizontal:
				return constraint / len(t.headers)
			default:
				return dims.Size.Y
			}
		},
		t.headerFun,
		t.dataFun,
	)
}

*/

type Table struct {
	// todo use this object to draw table
	// height      unit.Dp
	// grid        component.GridState
	// headerFun   layout.ListElement
	// dataFun     outlay.Cell
	// headers     []string
	// data        [][]any
	// dataContent []widget.Bool
	// table       component.TableStyle

	component.GridState

	headerBorder             *widget.Border
	cellBorder               *widget.Border
	headers                  []*widget.Clickable
	cells                    []*widget.Clickable                    // 单元格单击事件
	SelectionChangedCallback func(gtx layout.Context, row, col int) // num 单击了几次
	DoubleClickCallback      func(gtx layout.Context, row, col int) `json:"-"` // 双击编辑节点或行结构体数据
	rowChan                  chan int                               // 存放需要刷新的行
	rowChanSize              int
	lastRowIdx               int // 用于去重
	HeaderCellIndex          int // 哪列被点击了

	creatCellCallback func(gtx layout.Context, row, col int) layout.Dimensions

	rowIdx int // 左键的索引
	colIdx int

	Tabler

	cellsAreas []*component.ContextArea
	menu       *ContextMenu
}

func NewTable(table Tabler) *Table {
	return &Table{
		GridState: component.GridState{},
		headerBorder: &widget.Border{
			Color: color.NRGBA{R: 76, G: 76, B: 76, A: 255},
			//Color: color.NRGBA{
			//	R: 32,
			//	G: 32,
			//	B: 32,
			//	A: 255,
			//},
			Width: unit.Dp(0.2),
		},
		cellBorder: &widget.Border{
			Color: color.NRGBA{
				R: 32,
				G: 32,
				B: 32,
				A: 255,
			},
			Width: unit.Dp(0),
		},
		headers:                  nil,
		cells:                    nil,
		SelectionChangedCallback: nil,
		DoubleClickCallback:      nil,
		rowChan:                  nil,
		rowChanSize:              0,
		lastRowIdx:               -1,
		HeaderCellIndex:          0,
		creatCellCallback:        nil,
		rowIdx:                   -1,
		colIdx:                   -1,
		Tabler:                   table,
		cellsAreas:               nil,
		menu:                     nil,
	}
}

// 创建自定义单元格回调函数
func (m *Table) SetCreateCellCallback(f func(gtx layout.Context, row, col int) layout.Dimensions) *Table {
	m.creatCellCallback = f
	return m
}

func (m *Table) SetRowChannel(ch chan int, size int) *Table {
	m.rowChan = ch
	m.rowChanSize = size
	return m
}

func (m *Table) GetSelectedCell() (int, int) {
	return m.rowIdx, m.colIdx
}

func (m *Table) SetMenu(menu *ContextMenu) *Table {
	m.menu = menu
	return m
}

func (m *Table) Layout(gtx layout.Context) layout.Dimensions {
	if len(m.headers) != m.GetColumnCount() {
		m.headers = make([]*widget.Clickable, m.GetColumnCount())
	}
	if len(m.cells) != m.Size()*m.GetColumnCount() {
		m.cells = make([]*widget.Clickable, m.Size()*m.GetColumnCount())
	}
	mylog.CheckNil(m.menu) // todo，把增删改查，节点转换，复制行列设置为默认菜单，并添加每个菜单子项目的回调字段
	if len(m.cellsAreas) != m.Size()*m.GetColumnCount() {
		m.cellsAreas = make([]*component.ContextArea, m.Size()*m.GetColumnCount())
	}
	return component.Table(th.Theme, &m.GridState).Layout(gtx, m.Size(), m.GetColumnCount(),
		/*单元格动态尺寸计算，适应可见区域，支持不同轴方向 dimensioner outlay.Dimensioner*/ func(axis layout.Axis, index, constraint int) int {
			used := float32(0)
			defCount := 0
			defWidth := 0
			n := m.GetColumnCount()
			for i := 0; i < n; i++ {
				if m.GetColumnWitdh(i) == 0 {
					defCount++
					continue
				}
				used += m.GetColumnWitdh(i)
			}

			if defCount > 0 {
				free := (constraint - int(used)) / defCount
				if free > defWidth {
					defWidth = free
				}
			}
			switch axis {
			case layout.Horizontal:
				width := m.GetColumnWitdh(index)
				if width == 0 {
					return gtx.Dp(unit.Dp(defWidth))
				}
				return gtx.Dp(unit.Dp(width))
			default:
				return gtx.Dp(unit.Dp(27)) // 行高
			}
		},
		/*表头绘制 headingFunc DefaultDraw.ListElement*/ func(gtx C, col int) D {
			DrawColumnDivider(gtx, col)                                                       // 为表头绘制列分隔条
			paint.FillShape(gtx.Ops, ColorHeaderFg, clip.Rect{Max: gtx.Constraints.Max}.Op()) // 表头背景色

			//return component.Resize{}.Layout(gtx, func(gtx DefaultDraw.Context) DefaultDraw.Dimensions {
			//
			//})

			return m.headerBorder.Layout(gtx, func(gtx C) D {
				column := m.GetColumn(col)
				click := m.headers[col]
				if click == nil {
					click = new(widget.Clickable)
					m.headers[col] = click
				}
				// component.Rect{}.Layout(gtx)
				if click.Clicked(gtx) {
					m.HeaderCellIndex = col
					if column.cb != nil {
						column.cb(col) // todo 表头每列双击弹出对应每列位置的搜索框，单机排序，右击复制列到剪切板
					}
				}

				// button := material.Button(th.Theme, RowSelectedCallback, m.GetTitle(col))
				// button.Background = color.NRGBA(colornames.Grey500)
				// return DefaultDraw.Center.Layout(gtx, button.Layout)
				return material.Clickable(gtx, click, func(gtx C) D {
					return layout.UniformInset(0).Layout(gtx, func(gtx C) D {
						DrawColumnDivider(gtx, col) // 为每列绘制列分隔条
						body1 := material.Body1(th.Theme, m.GetTitle(col))
						body1.MaxLines = 1
						body1.Truncator = "..."
						body1.Color = th.Color.DefaultTextWhiteColor
						return layout.Center.Layout(gtx, body1.Layout)
					})
				})
			})
		},
		/*渲染body单元格 cellFunc outlay.Cell*/ func(gtx C, row, col int) D {
			DrawCrosswalk(gtx, row)         // 绘制斑马线
			DrawColumnDivider(gtx, col)     // 为每列绘制列分隔条
			if m.creatCellCallback != nil { // todo remove ? 看看是的呀着色会不会用到
				return m.creatCellCallback(gtx, row, col)
			}
			if m.lastRowIdx != row {
				m.lastRowIdx = row
				if m.rowChan != nil {
					if len(m.rowChan) < m.rowChanSize {
						m.rowChan <- row
					}
				}
			}
			record := m.GetRow(row)
			txt := m.GetItemText(record, row, col)
			idx := row*m.GetColumnCount() + col
			cell := m.cells[idx]
			if cell == nil {
				cell = &widget.Clickable{}
				m.cells[idx] = cell
			}
			for {
				click, ok := cell.Update(gtx)
				if !ok {
					break
				}
				m.rowIdx = row
				m.colIdx = col
				switch click.NumClicks {
				case 1:
					if m.SelectionChangedCallback != nil {
						m.SelectionChangedCallback(gtx, m.rowIdx, m.colIdx)
						gtx.Execute(op.InvalidateCmd{})
					}
				case 2: // todo https://github.com/chapar-rest/chapar/issues/17
					if m.DoubleClickCallback != nil {
						m.DoubleClickCallback(gtx, m.rowIdx, m.colIdx)
						gtx.Execute(op.InvalidateCmd{})
					}
				}
			}
			return material.Clickable(gtx, cell, func(gtx C) D {
				cellMenu := m.cellsAreas[idx]
				if cellMenu == nil {
					cellMenu = &component.ContextArea{
						Activation:       pointer.ButtonSecondary,
						AbsolutePosition: true,
						PositionHint:     0,
					}
					m.cellsAreas[idx] = cellMenu
				}
				m.menu.Clicked(gtx) // callback

				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						return layout.UniformInset(0).Layout(gtx, func(gtx C) D {
							cellText := material.Body1(th.Theme, txt)
							cellText.Color = th.Color.DefaultTextWhiteColor
							if m.rowIdx == row {
								if m.colIdx == col {
									cellText.Color = color.NRGBA{R: 0, G: 0, B: 255, A: 255}
								} else {
									cellText.Color = color.NRGBA{R: 0, G: 0, B: 180, A: 255}
									HighlightRow(gtx) // IsRowSelected
								}
							}
							cellText.MaxLines = 1
							cellText.Truncator = "..."
							return layout.Center.Layout(gtx, cellText.Layout)
						})
					}),
					layout.Expanded(func(gtx C) D {
						return cellMenu.Layout(gtx, func(gtx C) D {
							m.rowIdx = row
							m.colIdx = col
							gtx.Constraints.Max.X = 500
							gtx.Constraints.Max.Y = 1400
							return m.drawContextArea(gtx)
						})
					}),
				)
			})
		},
	)
}

func (m *Table) drawContextArea(gtx C) D {
	return layout.Center.Layout(gtx, func(gtx C) D { // 重置min x y 到0，并根据max x y 计算弹出菜单的合适大小
		// mylog.Struct("todo",gtx.Constraints)
		menuStyle := component.Menu(th.Theme, &m.menu.MenuState)
		menuStyle.SurfaceStyle = component.SurfaceStyle{
			Theme: th.Theme,
			ShadowStyle: component.ShadowStyle{
				CornerRadius: 18, // 弹出菜单的椭圆角度
				Elevation:    0,
				// AmbientColor:  color.NRGBA(colornames.Blue400),
				// PenumbraColor: color.NRGBA(colornames.Blue400),
				// UmbraColor:    color.NRGBA(colornames.Blue400),
			},
			Fill: color.NRGBA{R: 50, G: 50, B: 50, A: 255}, // 弹出菜单的背景色
		}
		return menuStyle.Layout(gtx)
	})
}

//func drawColumnDivider(gtx C, col int, color color.NRGBA) { //绘制列分隔条,todo最后一列没绘制到
//	if col > 0 {
//		dividerWidth := 1
//		tallestHeight := gtx.Constraints.Min.Y
//		stack3 := clip.Rect{Max: image.Pt(dividerWidth, tallestHeight)}.Push(gtx.Ops)
//		paint.Fill(gtx.Ops, color)
//		stack3.Pop()
//	}
//}

func HighlightRow(gtx C) { // 高亮选中行为蓝色
	paint.FillShape(gtx.Ops, color.NRGBA(colornames.Blue400), clip.Rect{Max: gtx.Constraints.Max}.Op())
}

func DrawCrosswalk(gtx C, row int) { // 绘制斑马线
	if row%2 == 0 {
		paint.FillShape(gtx.Ops, color.NRGBA{
			R: 42,
			G: 42,
			B: 42,
			A: 255,
		}, clip.Rect{Max: gtx.Constraints.Max}.Op())
	} else {
		paint.FillShape(gtx.Ops, color.NRGBA{
			R: 32,
			G: 32,
			B: 32,
			A: 255,
		}, clip.Rect{Max: gtx.Constraints.Max}.Op())
	}
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
