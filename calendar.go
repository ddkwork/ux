package ux

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"time"

	"github.com/ddkwork/ux/widget/material"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type OnCalendarDateClick func(t time.Time)

type monthButton struct {
	Month time.Month
	widget.Clickable
}

type yearButton struct {
	Year int
	widget.Clickable
}

type cellItem struct {
	widget.Clickable
	time.Time
}

// space between months and years dropdown in the header
var spaceBetweenHeaderDropdowns = unit.Dp(32)

const dropdownWidth = unit.Dp(120)

var allMonthsButtonsArr = [12]monthButton{
	{Month: 1},
	{Month: 2},
	{Month: 3},
	{Month: 4},
	{Month: 5},
	{Month: 6},
	{Month: 7},
	{Month: 8},
	{Month: 9},
	{Month: 10},
	{Month: 11},
	{Month: 12},
}

var (
	allYearsButtonsSlice  []yearButton
	monthsHeaderRowHeight = unit.Dp(64)
	viewHeaderHeight      = unit.Dp(32)
)

func SetAllYearsButtonsSlice(startTime, endTime time.Time) {
	allYearsButtonsSlice = GetYearsRangeButtons(startTime.Year(), endTime.Year())
}

func init() {
	startTime := time.Now().AddDate(-100, 0, 0)
	endTime := time.Now().AddDate(101, 0, 0)
	SetAllYearsButtonsSlice(startTime, endTime)
}

type Calendar struct {
	time               time.Time
	btnDropdownMonth   widget.Clickable
	monthsList         layout.List
	yearsList          layout.List
	btnDropdownYear    widget.Clickable
	initialized        bool
	ShowMonthsDropdown bool
	showYearsDropdown  bool
	fullView           widget.Clickable
	viewList           layout.List
	OnCalendarDateClick
	weekdays       [7]time.Weekday
	FirstDayOfWeek time.Weekday
	cellItemsArr   []*cellItem
	maxWidth       int
	layout.Inset
}

func (c *Calendar) SetTime(t time.Time) {
	c.time = t
}

func (c *Calendar) Time() time.Time {
	return c.time
}

func (c *Calendar) Layout(gtx layout.Context) layout.Dimensions {
	if !c.initialized {
		now := time.Now()
		if c.Time().IsZero() {
			c.SetTime(now)
		}
		c.cellItemsArr = make([]*cellItem, 0)
		for range 42 {
			c.cellItemsArr = append(c.cellItemsArr, &cellItem{})
		}
		c.weekdays = [7]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}
		c.initialized = true
	}
	if gtx.Constraints.Max.X > gtx.Constraints.Max.Y {
		gtx.Constraints.Max.X = gtx.Constraints.Max.Y
	}

	c.maxWidth = gtx.Constraints.Max.X - gtx.Dp(c.Inset.Left+c.Inset.Right)

	firstDay := int(c.FirstDayOfWeek)
	c.weekdays[0] = c.FirstDayOfWeek
	for i := 1; i < 7; i++ {
		firstDay++
		firstDay %= 7
		c.weekdays[i] = time.Weekday(firstDay)
	}

	if c.fullView.Clicked(gtx) {
		if !c.btnDropdownMonth.Pressed() {
			c.ShowMonthsDropdown = false
		}
		if !c.btnDropdownYear.Pressed() {
			c.showYearsDropdown = false
		}
	}
	d := c.fullView.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return c.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			flex := layout.Flex{Axis: layout.Vertical}
			return flex.Layout(gtx,
				layout.Rigid(c.drawViewHeader),
				layout.Rigid(c.drawHeaderRow),
				layout.Rigid(c.drawBodyRows),
			)
		})
	})

	if c.ShowMonthsDropdown {
		c.drawMonthsDropdownItems(gtx)
	}
	if c.showYearsDropdown {
		c.drawYearsDropdownItems(gtx)
	}
	return d
}

func (c *Calendar) drawHeaderRow(gtx layout.Context) layout.Dimensions {
	flexChildren := make([]layout.FlexChild, len(c.weekdays))
	columnWidth := c.maxWidth / 7
	for i, day := range c.weekdays {
		dayStr := strings.ToUpper(day.String()[0:3])
		flexChildren[i] = c.drawHeaderColumn(gtx, dayStr, columnWidth)
	}
	flex := layout.Flex{}
	mac := op.Record(gtx.Ops)
	d := flex.Layout(gtx, flexChildren...)
	call := mac.Stop()
	rect := clip.Rect{Max: d.Size}
	paint.FillShape(gtx.Ops, th.ContrastBg, rect.Op())
	call.Add(gtx.Ops)
	return d
}

func (c *Calendar) drawHeaderColumn(gtx layout.Context, day string, columnWidth int) layout.FlexChild {
	gtx.Constraints.Min.Y, gtx.Constraints.Max.Y = gtx.Dp(monthsHeaderRowHeight), gtx.Dp(monthsHeaderRowHeight)
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X, gtx.Constraints.Max.X = columnWidth, columnWidth
		inset := layout.UniformInset(16)
		if c.maxWidth < gtx.Dp(500) {
			inset.Top, inset.Bottom, inset.Right, inset.Left = 8, 8, 8, 8
		}
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, th.TextSize, day)
				label.Color = th.ContrastFg
				label.MaxLines = 1
				if c.maxWidth < gtx.Dp(500) {
					label.Text = day[:1]
					label.TextSize = unit.Sp(14)
				}
				return label.Layout(gtx)
			})
		})
	})
}

func (c *Calendar) drawColumn(gtx layout.Context, columnWidth int, btn *cellItem) layout.FlexChild {
	dayStr := fmt.Sprintf("%d", btn.Day())
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		bgColor := th.Bg
		txtColor := th.Fg
		txtColor.A = 210
		if btn.Month() != c.Time().Month() { // 非本月份日期
			bgColor = color.NRGBA(colornames.Grey800)
			txtColor.A = 100
		}
		if c.Time().Month() == btn.Month() {
			if btn.Clicked(gtx) {
				if c.OnCalendarDateClick != nil {
					c.OnCalendarDateClick(btn.Time)
				}
			}
			now := time.Now()
			isEqual := btn.Day() == now.Day() && now.Year() == btn.Year() && now.Month() == btn.Month()
			if btn.Hovered() || isEqual {
				bgColor = th.ContrastBg
				bgColor.A = 240
				txtColor = th.ContrastFg
				txtColor.A = 240
			}
		}
		return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			mac := op.Record(gtx.Ops)
			gtx.Constraints.Min.X, gtx.Constraints.Max.X = columnWidth, columnWidth
			gtx.Constraints.Min.Y, gtx.Constraints.Max.Y = columnWidth, columnWidth
			d := layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				center := layout.N
				txtSize := th.TextSize
				txtSize *= 1.5
				d := center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Label(th, txtSize, dayStr)
					label.MaxLines = 1
					label.Color = txtColor
					label.Alignment = text.Middle
					if c.maxWidth < gtx.Dp(500) {
						label.TextSize = unit.Sp(14)
					}
					return label.Layout(gtx)
				})
				return d
			})
			call := mac.Stop()
			rect := clip.Rect{Max: d.Size}
			paint.FillShape(gtx.Ops, bgColor, rect.Op())
			call.Add(gtx.Ops)
			return d
		})
	})
}

func (c *Calendar) drawBodyRows(gtx layout.Context) layout.Dimensions {
	flex := layout.Flex{Axis: layout.Vertical}
	t := c.Time()
	columnWidth := c.maxWidth / 7
	startDate := firstDayOfWeek(t, time.Monday)
	lastDate := lastDayOfWeek(t, time.Monday)
	day := startDate
	endDate := lastDate

	allRows := make([]layout.FlexChild, 0)
	cellItemsArr := make([]*cellItem, 0)
	for rowIndex, cellIndex := 0, 0; day.Before(endDate) || rowIndex < 6; rowIndex++ {
		var shouldBreak bool
		for i := range 7 {
			// Minimum row cellIndex is 3 (ie at least 4 fields) and max row cellIndex is 5 (ie at least 6 fields)
			if (rowIndex == 4 || rowIndex == 5) && i == 0 {
				if day.Month() != c.Time().Month() {
					allRows = allRows[:rowIndex]
					shouldBreak = true
					break
				}
			}
			cellItemsArr = append(cellItemsArr, c.cellItemsArr[cellIndex])
			cellItemsArr[cellIndex].Time = day
			cellIndex++
			day = day.AddDate(0, 0, 1)
			if i == 0 {
				allRows = append(allRows, layout.FlexChild{})
			}
		}
		if shouldBreak {
			break
		}
	}
	cellIndex := 0
	for rowIndex := range allRows {
		var flexChildren []layout.FlexChild
		for range 7 {
			flexChildren = append(flexChildren, c.drawColumn(gtx, columnWidth, cellItemsArr[cellIndex]))
			cellIndex++
		}
		flexChild := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			flex := layout.Flex{}
			gtx.Constraints.Min.Y, gtx.Constraints.Max.Y = columnWidth, columnWidth
			return flex.Layout(gtx, flexChildren...)
		})
		allRows[rowIndex] = flexChild
	}
	c.viewList.Axis = layout.Vertical
	d := c.viewList.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
		return flex.Layout(gtx, allRows...)
	})
	return d
}

func (c *Calendar) OnMonthButtonClick(gtx layout.Context, month *monthButton) {
	t := c.Time()
	t = time.Date(t.Year(), month.Month, t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	c.SetTime(t)
	gtx.Execute(op.InvalidateCmd{})
}

func (c *Calendar) OnYearButtonClick(gtx layout.Context, year *yearButton) {
	t := c.Time()
	t = time.Date(year.Year, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	c.SetTime(t)
	gtx.Execute(op.InvalidateCmd{})
}

func (c *Calendar) drawMonthsDropdownItems(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.Y = (c.maxWidth / 7) * 4
	op.Offset(image.Point{
		X: gtx.Dp(16),
		Y: gtx.Dp(viewHeaderHeight) + gtx.Dp(8),
	}).Add(gtx.Ops)
	layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			mac := op.Record(gtx.Ops)
			gtx.Constraints.Min.X = gtx.Dp(dropdownWidth)
			c.monthsList.Axis = layout.Vertical
			border := widget.Border{
				Color:        th.ContrastBg,
				CornerRadius: 0,
				Width:        unit.Dp(1),
			}
			d := border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				d := c.monthsList.Layout(gtx, len(allMonthsButtonsArr), func(gtx layout.Context, index int) layout.Dimensions {
					bgColor := RowColor(index)
					txtColor := th.Fg
					isSelected := c.Time().Month().String() == allMonthsButtonsArr[index].Month.String()
					if allMonthsButtonsArr[index].Hovered() || isSelected {
						bgColor = th.DropDownItemHoveredGrayColor
						txtColor = th.Bg
					}
					if allMonthsButtonsArr[index].Clicked(gtx) {
						c.ShowMonthsDropdown = false
						c.OnMonthButtonClick(gtx, &allMonthsButtonsArr[index])
					}
					mac := op.Record(gtx.Ops)
					d := allMonthsButtonsArr[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						inset := layout.Inset{Top: 8, Bottom: 8, Left: 16, Right: 16}
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							txt := allMonthsButtonsArr[index].Month.String()
							label := material.Label(th, th.TextSize, txt)
							label.Alignment = text.Start
							label.Color = txtColor
							return label.Layout(gtx)
						})
					})
					call := mac.Stop()
					rect := clip.Rect{Max: d.Size}
					paint.FillShape(gtx.Ops, bgColor, rect.Op())
					call.Add(gtx.Ops)
					return d
				})
				return d
			})
			call := mac.Stop()
			rect := clip.Rect{Max: d.Size}
			paint.FillShape(gtx.Ops, th.Bg, rect.Op())
			call.Add(gtx.Ops)
			return d
		}),
	)
	return layout.Dimensions{}
}

func (c *Calendar) drawYearsDropdownItems(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.Y = (c.maxWidth / 7) * 4
	op.Offset(image.Point{
		X: gtx.Dp(16) + gtx.Dp(dropdownWidth) + gtx.Dp(spaceBetweenHeaderDropdowns),
		Y: gtx.Dp(viewHeaderHeight) + gtx.Dp(8),
	}).Add(gtx.Ops)
	layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			mac := op.Record(gtx.Ops)
			gtx.Constraints.Min.X = gtx.Dp(dropdownWidth)
			c.yearsList.Axis = layout.Vertical
			border := widget.Border{Color: th.ContrastBg, CornerRadius: 0, Width: unit.Dp(1)}
			d := border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				d := c.yearsList.Layout(gtx, len(allYearsButtonsSlice), func(gtx layout.Context, index int) layout.Dimensions {
					bgColor := RowColor(index)
					txtColor := th.Fg
					isSelected := c.Time().Year() == allYearsButtonsSlice[index].Year
					if allYearsButtonsSlice[index].Hovered() || isSelected {
						bgColor = th.DropDownItemHoveredGrayColor
						txtColor = th.Bg
					}
					if allYearsButtonsSlice[index].Clicked(gtx) {
						c.showYearsDropdown = false
						c.OnYearButtonClick(gtx, &allYearsButtonsSlice[index])
					}
					mac := op.Record(gtx.Ops)
					d := allYearsButtonsSlice[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						inset := layout.Inset{Top: 8, Bottom: 8, Left: 16, Right: 16}
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							txt := fmt.Sprintf("%d", allYearsButtonsSlice[index].Year)
							label := material.Label(th, th.TextSize, txt)
							label.Alignment = text.Start
							label.Color = txtColor
							return label.Layout(gtx)
						})
					})
					call := mac.Stop()
					rect := clip.Rect{Max: d.Size}
					paint.FillShape(gtx.Ops, bgColor, rect.Op())
					call.Add(gtx.Ops)
					return d
				})
				return d
			})
			call := mac.Stop()
			rect := clip.Rect{Max: d.Size}
			paint.FillShape(gtx.Ops, th.Bg, rect.Op())
			call.Add(gtx.Ops)
			return d
		}),
	)
	return layout.Dimensions{}
}

func (c *Calendar) drawViewHeader(gtx layout.Context) layout.Dimensions {
	month := c.Time().Month().String()
	year := fmt.Sprintf("%d", c.Time().Year())
	gtx.Constraints.Max.Y, gtx.Constraints.Min.Y = gtx.Dp(viewHeaderHeight), gtx.Dp(viewHeaderHeight)
	flex := layout.Flex{Spacing: layout.SpaceEnd, Alignment: layout.Middle}
	d := flex.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if c.btnDropdownMonth.Clicked(gtx) {
				c.ShowMonthsDropdown = !c.ShowMonthsDropdown
				c.showYearsDropdown = false
				if c.ShowMonthsDropdown {
					for i, eachButton := range allMonthsButtonsArr {
						if eachButton.Month.String() == c.Time().Month().String() {
							c.monthsList.Position.First = i
							c.monthsList.Position.Offset = -32
							break
						}
					}
				}
			}
			gtx.Constraints.Min.X = gtx.Dp(dropdownWidth)
			d := c.btnDropdownMonth.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				flex := layout.Flex{Spacing: layout.SpaceBetween}
				return flex.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Label(th, th.TextSize, month).Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						downIcon, _ := widget.NewIcon(icons.NavigationArrowDropDown)
						return downIcon.Layout(gtx, th.ContrastBg)
					}),
				)
			})
			return d
		}),
		layout.Rigid(layout.Spacer{Width: spaceBetweenHeaderDropdowns}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if c.btnDropdownYear.Clicked(gtx) {
				c.ShowMonthsDropdown = false
				c.showYearsDropdown = !c.showYearsDropdown
				if c.showYearsDropdown {
					c.scrollToSelectedYear()
				}
			}
			d := c.btnDropdownYear.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(dropdownWidth)
				flex := layout.Flex{Spacing: layout.SpaceBetween}
				return flex.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Label(th, th.TextSize, year)
						return label.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						downIcon, _ := widget.NewIcon(icons.NavigationArrowDropDown)
						return downIcon.Layout(gtx, th.ContrastBg)
					}),
				)
			})
			return d
		}),
	)
	return d
}

func (c *Calendar) scrollToSelectedYear() {
	for i, eachYear := range allYearsButtonsSlice {
		if eachYear.Year == c.Time().Year() {
			c.yearsList.Position.First = i
			c.yearsList.Position.Offset = -32
			break
		}
	}
}

//////////////////////////

// Ref https://stackoverflow.com/questions/36830212/get-the-first-and-last-day-of-current-month-in-go-golang
func beginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func endOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}

func firstDayOfWeek(tm time.Time, weekStartDay time.Weekday) time.Time {
	// tm = time.Date(tm.Year(), tm.Month(), 1, 0, 0, 0, 0, tm.Location())
	tm = beginningOfMonth(tm)
	for tm.Weekday() != weekStartDay {
		tm = tm.AddDate(0, 0, -1)
	}
	return tm
}

func lastDayOfWeek(tm time.Time, weekStartDay time.Weekday) time.Time {
	tm = endOfMonth(tm)
	weekEndDay := (weekStartDay + 6) % 7
	for tm.Weekday() != weekEndDay {
		tm = tm.AddDate(0, 0, 1)
	}
	return tm
}

// GetYearsRangeButtons returns slice of yearButton with year range between startYear and upto but not including lastYear
func GetYearsRangeButtons(startYear, endYear int) []yearButton {
	yearsRange := make([]yearButton, 0)
	for currentYear := startYear; currentYear < endYear; currentYear++ {
		yearsRange = append(yearsRange, yearButton{Year: currentYear})
	}
	return yearsRange
}
