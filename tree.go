package ux

import (
	"fmt"
	"image"

	"github.com/ddkwork/ux/widget/material"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
)

type (
	ClickAction1 func(gtx layout.Context, node *TreeNode)
	Tree         struct {
		nodes       []*TreeNode
		width       unit.Dp
		clickedNode *TreeNode
		click       ClickAction1
	}
)

func NewTree() *Tree {
	return &Tree{
		width: unit.Dp(200),
	}
}

func (t *Tree) OnClick(fun ClickAction1) *Tree {
	t.click = fun
	return t
}

func (t *Tree) SetWidth(width unit.Dp) *Tree {
	t.width = width
	return t
}

func (t *Tree) SetNodes(nodes []*TreeNode) *Tree {
	for _, node := range nodes {
		t.setClick(node)
	}
	t.setPath(nodes, []int{})
	t.nodes = nodes
	return t
}

func (t *Tree) AddTopNode(newNode *TreeNode) error {
	t.setClick(newNode)
	t.nodes = append(t.nodes, newNode)
	t.setPath(t.nodes, []int{})
	return nil
}

func (t *Tree) AddSonNode(newNode *TreeNode) error {
	if t.clickedNode == nil {
		return fmt.Errorf("no node rowSelected")
	}
	t.setClick(newNode)
	path := t.clickedNode.Path
	parentNode := mylog.Check2(t.getNode(t.nodes, path))

	parentNode.Children = append(parentNode.Children, newNode)
	t.setPath(t.nodes, []int{})
	return nil
}

func (t *Tree) DeleteNode(newNode *TreeNode) error {
	if t.clickedNode == nil {
		return fmt.Errorf("no node rowSelected")
	}
	t.setClick(newNode)
	path := t.clickedNode.Path
	parentNode := mylog.Check2(t.getNode(t.nodes, path))

	parentNode.IsDeleted = true
	return nil
}

func (t *Tree) getNode(nodes []*TreeNode, paths []int) (*TreeNode, error) {
	if nodes == nil {
		nodes = t.nodes
	}
	for i, path := range paths {
		if len(nodes) <= path {
			return nil, fmt.Errorf("路径错误: 节点索引超出范围")
		}
		if i < len(paths)-1 { // 检查是否是最后一个路径值
			if nodes[path].Children != nil {
				return t.getNode(nodes[path].Children, paths[i+1:])
			}
			return nodes[path], nil // 返回最后一个路径值对应的节点
		}
	}
	return nil, fmt.Errorf("路径错误: 路径为空")
}

func (t *Tree) setPath(nodes []*TreeNode, path []int) {
	if nodes == nil {
		nodes = t.nodes
	}
	for key, node := range nodes {
		var sign []int
		if len(path) == 0 {
			sign = []int{key}
		} else {
			sign = append(path, key)
		}
		node.Path = sign
		if len(node.Children) > 0 {
			t.setPath(node.Children, sign)
		}
	}
}

func (t *Tree) setClick(nodes *TreeNode) {
	nodes.clickable = &widget.Clickable{}
	if len(nodes.Children) > 0 {
		for _, child := range nodes.Children {
			t.setClick(child)
		}
	}
}

type CallbackFun1 func(gtx layout.Context)

type TreeNode struct {
	Id       int
	Title    string
	ParentId int
	Icon     *widget.Icon
	Children []*TreeNode
	Expanded bool
	// selected      bool
	clickable     *widget.Clickable
	ClickCallback CallbackFun1
	Path          []int
	IsDeleted     bool
}

func (t *Tree) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// 这里可以添加头部或者其他固定的内容
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, t.renderTree(gtx, t.nodes)...)
		}),
	)
}

func (t *Tree) renderTree(gtx layout.Context, nodes []*TreeNode) []layout.FlexChild {
	if len(nodes) == 0 {
		return []layout.FlexChild{}
	}
	var dims []layout.FlexChild
	for _, node := range nodes {
		if node.IsDeleted {
			continue
		}
		dims = append(dims, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.renderNode(gtx, node, 0, true)
		}))
	}
	return dims
}

func (t *Tree) renderNode(gtx layout.Context, node *TreeNode, depth int, isParent bool) layout.Dimensions {
	// 渲节点标题
	bgColor := th.Bg
	bgColor = RowColor(depth + 1)

	if node.clickable.Clicked(gtx) {
		node.Expanded = !node.Expanded
		t.clickedNode = node
		if node.ClickCallback != nil {
			node.ClickCallback(gtx)
		}
		if t.click != nil {
			t.click(gtx, node)
		}
	}
	if node.clickable.Hovered() {
		bgColor = th.Color.TreeHoveredBgColor
	}
	var sonItems []layout.FlexChild
	// 绘制展开/折叠图标
	if len(node.Children) > 0 {
		sonItems = append(sonItems, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return node.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(1)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Dp(th.Size.DefaultIconSize)
					svg := SvgIconCircledChevronRight
					if node.Expanded {
						svg = SvgIconCircledChevronDown
					}
					return iconButtonSmall(new(widget.Clickable), svg, "").Layout(gtx)
				})
			})
		}))
	}
	sonItems = append(sonItems, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Dp(t.width)
		return node.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Label(th.Theme, 12, node.Title).Layout(gtx)
			})
		})
	}))

	items := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Background{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// defer clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}, gtx.Dp(th.Size.DefaultElementRadiusSize)).Push(gtx.Ops).Pop()
				if t.clickedNode == node {
					bgColor = th.Color.TreeClickedBgColor
				}
				defer clip.Rect{
					Max: image.Point{
						X: gtx.Constraints.Max.X,
						Y: gtx.Constraints.Min.Y,
					},
				}.Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, bgColor)
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(25))
				return layout.Inset{Left: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Left: unit.Dp(depth * 16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, sonItems...)
					})
				})
			})
		}),
	}
	// 递归渲染子节点
	if node.Expanded && len(node.Children) > 0 {
		var dims []layout.FlexChild
		for _, child := range node.Children {
			if child.IsDeleted {
				continue
			}
			dims = append(dims, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return t.renderNode(gtx, child, depth+1, false)
			}))
		}
		items = append(items, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// return layout.Inset{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, dims...)
			// })
		}))
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
}

func (t *Tree) SetCurrentNode(node *TreeNode) {
	t.clickedNode = node
}

func (t *Tree) GetCurrentNode() *TreeNode {
	return t.clickedNode
}

func (t *Tree) MinTree(gtx layout.Context, nodes []*TreeNode) {
	if nodes == nil {
		nodes = t.nodes
	}
	for _, node := range nodes {
		if node.Expanded {
			node.Expanded = false
		}
		if len(node.Children) > 0 {
			t.MinTree(gtx, node.Children)
		}
	}
	gtx.Execute(op.InvalidateCmd{})
}
