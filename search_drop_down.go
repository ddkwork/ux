package ux

import (
	"github.com/ddkwork/ux/internal/keys"
	"strings"

	"github.com/ddkwork/ux/resources/images"

	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/ddkwork/ux/widget/material"
	"github.com/ddkwork/ux/x/component"
)

type SearchDropDown struct {
	menuContextArea component.ContextArea
	list            *widget.List

	listStyle material.ListStyle

	results []*SearchResult
	input   *TextField

	loaderFunc func() []Item

	OnSelectResult func(result *SearchResult)
}

func NewSearchDropDown() *SearchDropDown {
	c := &SearchDropDown{
		input: NewTextField("", "Search..."),
		menuContextArea: component.ContextArea{
			Activation:       pointer.ButtonPrimary,
			AbsolutePosition: true,
		},
		list: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
	}

	c.input.SetOnTextChange(c.onSearch)
	c.input.OnKeyPress = c.onKeyPress

	return c
}

func (c *SearchDropDown) onKeyPress(gtx layout.Context, k key.Name) {
	switch k {
	case key.NameEscape:
		c.menuContextArea.Dismiss()
		c.results = nil
		c.input.SetText("")
	case key.NameReturn:
		if len(c.results) > 0 {
			if c.OnSelectResult != nil {
				c.OnSelectResult(c.results[0])
				c.menuContextArea.Dismiss()
				c.input.SetText("")
				c.results = nil
			}
		}
	}
}

func (c *SearchDropDown) onSearch(query string) {
	if c.loaderFunc == nil {
		return
	}

	items := c.loaderFunc()
	results := FuzzySearch(items, query, 100)
	if len(results) == 0 {
		return
	}

	for _, item := range results {
		switch item.Item.Kind {
		case KindEnv:
			item.Icon = images.NavigationMenuIcon
		case KindRequest:
			item.Icon = images.ActionSwapHorizIcon
		case KindWorkspace:
			item.Icon = images.NavigationAppsIcon
		case KindProtoFile:
			item.Icon = images.FileFolderIcon
		case KindCollection:
			item.Icon = images.NavigationChevronRightIcon
		}
	}

	c.results = results
	c.menuContextArea.Active()
}

func (c *SearchDropDown) SetLoader(fn func() []Item) {
	c.loaderFunc = fn
}

func (c *SearchDropDown) resultItemLayout(gtx layout.Context, item *SearchResult) layout.Dimensions {
	if item.Clickable.Clicked(gtx) {
		if c.OnSelectResult != nil {
			c.OnSelectResult(item)
		}
	}

	return Clickable(gtx, &item.Clickable, 0, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top:    unit.Dp(5),
			Bottom: unit.Dp(5),
			Left:   unit.Dp(10),
			Right:  unit.Dp(5),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Dp(18)
					return images.Layout(gtx, item.Icon, th.ContrastFg, 0)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Label(th, th.TextSize, item.Item.Title).Layout(gtx)
				}),
			)
		})
	})
}

func (c *SearchDropDown) Update(gtx layout.Context) {
	keys.OnKey(gtx, c, key.Filter{Required: key.ModShortcut, Name: "F"}, func() {
		gtx.Execute(key.FocusCmd{Tag: c.input})
	})
}

// Layout the SearchDropDown.
func (c *SearchDropDown) Layout(gtx layout.Context) layout.Dimensions {
	c.Update(gtx)

	inputDims := c.input.Layout(gtx)
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return inputDims
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return c.menuContextArea.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				offset := layout.Inset{
					Top:  unit.Dp(float32(inputDims.Size.Y)/gtx.Metric.PxPerDp + 1),
					Left: unit.Dp(0),
				}
				return offset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					sf := component.Surface(th)
					// sf.Fill = theme.DropDownMenuBgColor
					sf.Elevation = unit.Dp(2)
					return sf.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						gtx.Constraints.Max.Y = gtx.Dp(200)

						if len(c.results) == 0 {
							gtx.Constraints.Min.Y = gtx.Dp(50)
							return layout.Inset{
								Top:    unit.Dp(5),
								Bottom: unit.Dp(5),
								Left:   unit.Dp(10),
								Right:  unit.Dp(5),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return material.Label(th, th.TextSize, "Search in workspace").Layout(gtx)
							})
						}

						gtx.Constraints.Min.Y = gtx.Dp(200)
						c.listStyle = material.List(th, c.list)

						return c.listStyle.Layout(gtx, len(c.results), func(gtx layout.Context, index int) layout.Dimensions {
							return c.resultItemLayout(gtx, c.results[index])
						})
					})
				})
			})
		}),
	)
}

// ////////////////////////////////
const (
	ApiVersion = "v1"

	KindConfig      = "Config"
	KindWorkspace   = "Workspace"
	KindProtoFile   = "ProtoFile"
	KindEnv         = "Environment"
	KindRequest     = "Request"
	KindPreferences = "Preferences"
	KindCollection  = "Collection"
)

type MetaData struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type KeyValue struct {
	ID     string `yaml:"id"`
	Key    string `yaml:"key"`
	Value  string `yaml:"value"`
	Enable bool   `yaml:"enable"`
}

func KeyValuesToText(values []KeyValue) string {
	var text string
	for _, v := range values {
		text += v.Key + ": " + v.Value + "\n"
	}
	return text
}

func TextToKeyValue(txt string) []KeyValue {
	values := make([]KeyValue, 0)
	lines := strings.SplitSeq(txt, "\n")
	for line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		values = append(values, KeyValue{
			Key:   strings.TrimSpace(parts[0]),
			Value: strings.TrimSpace(parts[1]),
		})
	}

	return values
}

func FindKeyValue(values []KeyValue, key string) string {
	for _, v := range values {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}
