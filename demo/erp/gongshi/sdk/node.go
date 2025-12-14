package sdk

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/ddkwork/golibrary/std/stream/deepcopy"
	"github.com/ddkwork/golibrary/std/stream/uuid"
)

// Node represents a node in the tree table.
type Node struct {
	ID        uuid.ID    // Unique identifier for the node (UUID)
	Type      string     // Node type (container nodes end with "_container")
	RowCells  []CellData // Row data (including formula columns)
	Children  []*Node    // Child nodes
	parent    *Node      // Parent node
	isOpen    bool       // Whether expanded (only for container nodes)
	GroupKey  string     // Grouping key
	RowNumber int        // Row number (for sorting)
	walkIndex int
}

func newRoot() *Node         { return NewContainerNode("root", nil) }
func (n *Node) IsRoot() bool { return n.parent == nil }

// NewNode creates a new node with the given row cells.
func NewNode(rowCells []CellData) *Node {
	return &Node{
		ID:        newID(),
		Type:      "node",
		RowCells:  rowCells,
		Children:  nil,
		parent:    nil,
		isOpen:    false,
		GroupKey:  "",
		RowNumber: 0,
	}
}

// NewContainerNode creates a new container node.
func NewContainerNode(typeKey string, rowCells []CellData) *Node {
	n := NewNode(rowCells)
	n.Type = typeKey + ContainerKeyPostfix
	n.isOpen = true
	return n
}

// Clone creates a deep copy of the node.
//func (n *Node) Clone() *Node {
//	clone := &Node{
//		ID:        newID(),
//		Type:      n.Type,
//		RowCells:  make([]CellData, len(n.RowCells)),
//		Children:  make([]*Node, len(n.Children)),
//		isOpen:    n.isOpen,
//		GroupKey:  n.GroupKey,
//		RowNumber: n.RowNumber,
//	}
//
//	// Copy row data
//	for i, cell := range n.RowCells {
//		clone.RowCells[i] = cell
//	}
//
//	// Copy child nodes
//	for i, child := range n.Children {
//		cloneChild := child.Clone()
//		cloneChild.parent = clone
//		clone.Children[i] = cloneChild
//	}
//
//	return clone
//}

func (n *Node) Clone() (to *Node) {
	to = deepcopy.Clone(n)
	to.parent = n.parent
	to.ID = newID()
	if n.CanHaveChildren() {
		n.setParents(to.Children, to, true)
	}
	to.OpenAll()
	return
}

// AddChildren adds multiple child nodes.
func (n *Node) AddChildren(children []*Node) {
	for _, child := range children {
		n.AddChild(child)
	}
}

func (n *Node) SetParents(children []*Node, parent *Node) {
	n.setParents(children, parent, false)
}

func (n *Node) setParents(children []*Node, parent *Node, isNewID bool) {
	for _, child := range children {
		child.parent = parent
		if isNewID {
			child.ID = newID()
		}
		if child.CanHaveChildren() {
			n.setParents(child.Children, child, isNewID)
		}
	}
}

func (n *Node) SetChildren(children []*Node) {
	n.Children = children
}

// DataNodes returns a sequence of all data nodes (children of the root).
func (t *TreeTable) DataNodes() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for _, child := range t.Root.Children {
			stack := []*Node{child}
			for len(stack) > 0 {
				n := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				if !yield(n) {
					return
				}

				for i := len(n.Children) - 1; i >= 0; i-- {
					stack = append(stack, n.Children[i])
				}
			}
		}
	}
}

// DataNodesSlice returns a slice of all data nodes.
func (t *TreeTable) DataNodesSlice() []*Node {
	var nodes []*Node
	for node := range t.DataNodes() {
		nodes = append(nodes, node)
	}
	return nodes
}

func (n *Node) Walk() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		n.walkIndex = 0
		if !n.IsRoot() { //todo skip root?应该调用dateNode数据节点方法
			if !yield(n) {
				return
			}
		}
		for i, child := range n.Children {
			child.walkIndex = i
			if !yield(child) { // 迭代索引是为了insert和remove时定位
				break
			}
			if child.CanHaveChildren() {
				// 函数式编程,Walk 方法返回的是一个函数。这个返回的函数接受一个参数（也是一个函数），这个参数就是 yield
				child.Walk()(yield) // 迭代子节点的子节点
			}
		}
	}
}

func (n *Node) WalkContainer() iter.Seq2[int, *Node] {
	return func(yield func(int, *Node) bool) {
		if n.Container() {
			if !yield(0, n) {
				return
			}
		}
		for i, container := range n.Containers() {
			if !yield(i, container) {
				break
			}
			for _, child := range container.Children {
				if child.CanHaveChildren() {
					child.WalkContainer()(yield)
				}
			}
		}
	}
}

func (n *Node) Containers() iter.Seq2[int, *Node] { //适用于分组后的容器节点
	return func(yield func(int, *Node) bool) {
		for i, child := range n.Children {
			if child.Container() { // 迭代当前节点下的所有容器节点
				if !yield(i, child) {
					break
				}
			}
		}
	}
}

// func (n *Node) WalkQueue() iter.Seq[*Node] { // 性能应该不行
//	return func(yield func(*Node) bool) {
//		queue := []*Node{n}
//		for len(queue) > 0 {
//			node := queue[0]
//			queue = queue[1:]
//			if !yield(node) {
//				break
//			}
//			for _, child := range node.Children {
//				queue = append(queue, child) // 这里将子节点添加到队列
//				if child.CanHaveChildren() { // 如果子节点是一个容器，递归地添加它的子节点
//					for _, subChild := range child.Children {
//						queue = append(queue, subChild)
//					}
//				}
//			}
//		}
//	}
// }

func (n *Node) Depth() int {
	count := 0 // root node is 0
	p := n.parent
	for p != nil {
		count++
		p = p.parent
	}
	return count
}
func (t *TreeTable) MaxDepth() int {
	maxDepth := 0
	for node := range t.Root.Walk() {
		childDepth := node.Depth()
		if childDepth > maxDepth {
			maxDepth = childDepth
		}
	}
	return maxDepth
}
func (n *Node) Find() (found *Node) { //todo test
	for child := range n.parent.Walk() {
		if child.ID == n.ID {
			found = child
			break
		}
	}
	return
}

func (n *Node) LenChildren() int { return len(n.Children) }
func (n *Node) LastChild() (lastChild *Node) {
	if n.IsRoot() && n.CanHaveChildren() {
		return n.Children[n.LenChildren()-1]
	}
	return n.parent.Children[n.parent.LenChildren()-1]
}
func (n *Node) IsLastChild() bool { return n.LastChild() == n }
func (n *Node) CopyFrom(from *Node) *Node {
	*n = *from
	return n
}

func (n *Node) ApplyTo(to *Node) *Node {
	*to = *n
	return n
}
func (n *Node) SumChildren() string {
	// container column 0 key is empty string
	k := n.Type
	k = strings.TrimSuffix(k, ContainerKeyPostfix)
	if n.LenChildren() == 0 {
		return k
	}
	k += " (" + fmt.Sprint(n.LenChildren()) + ")"
	return k
}

const ContainerKeyPostfix = "_container"

func (n *Node) UUID() uuid.ID { return n.ID }
func (n *Node) Container() bool {
	return strings.HasSuffix(n.Type, ContainerKeyPostfix)
}

func (n *Node) GetType() string        { return n.Type }
func (n *Node) SetType(typeKey string) { n.Type = typeKey }
func (n *Node) IsOpen() bool           { return n.isOpen && n.Container() }
func (n *Node) SetOpen(open bool)      { n.isOpen = open && n.Container() }
func (n *Node) Parent() *Node          { return n.parent }
func (n *Node) SetParent(parent *Node) { n.parent = parent }

func (n *Node) ResetChildren()        { n.Children = nil }
func (n *Node) CanHaveChildren() bool { return n.HasChildren() }
func (n *Node) HasChildren() bool     { return n.Container() && len(n.Children) > 0 }
func (n *Node) AddChild(child *Node) {
	child.parent = n
	n.Children = append(n.Children, child)
}

// InsertChild inserts a child node at the specified position.
func (n *Node) InsertChild(index int, child *Node) {
	if index < 0 || index > len(n.Children) {
		index = len(n.Children)
	}
	child.parent = n
	n.Children = append(n.Children[:index], append([]*Node{child}, n.Children[index:]...)...)
}

// RemoveChild removes a child node.
func (n *Node) RemoveChild(child *Node) {
	for i, c := range n.Children {
		if c.ID == child.ID {
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			return
		}
	}
}
func (n *Node) InsertAfter(after *Node) {
	after.parent = n.parent
	n.parent.Children = slices.Insert(n.parent.Children, n.Index()+1, after)
}

func (n *Node) Index() int {
	return slices.Index(n.parent.Children, n)
	// for i, child := range n.parent.Children {
	//	if n.ID == child.ID {
	//		return i
	//	}
	// }
	// panic("not found index") // 永远不可能选中root，所以可以放心panic，root不显示，只显示它的children作为rootRows
}

func (n *Node) Remove() {
	for child := range n.parent.Walk() {
		if child.ID == n.ID {
			n.parent.Children = slices.Delete(n.parent.Children, child.walkIndex, child.walkIndex+1)
			break
		}
	}
}
func (t *TreeTable) InsertAfter(after *Node) {
	t.SelectedNode.InsertAfter(after)
}

func (t *TreeTable) Remove() {
	t.SelectedNode.Remove()
}

// IsContainer checks if the node is a container.
func (n *Node) IsContainer() bool {
	return strings.HasSuffix(n.Type, ContainerKeyPostfix)
}
