package ux

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/widget/material"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
)

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

func GetHeader(obj any) []string {
	fields := reflect.VisibleFields(reflect.TypeOf(obj))
	headers := make([]string, 0, len(fields))
	for _, field := range fields {
		headers = append(headers, field.Name)
	}
	return headers
}

//type headerObj struct {
//
//}
//var TestRootRows = []*Node[headerObj]{
//	{
//		RowCells: []ColumnInfo{
//			{ID: 0, Cell: "System Management"},
//			{ID: 1, Cell: "System Overview"},
//			{ID: 2, Cell: "Active"},
//			{ID: 3, Cell: "Manage"},
//			{ID: 4, Cell: "Column 5"},
//			{ID: 5, Cell: "Column 6"},
//		},
//		Children: []*Node{
//			{
//				RowCells: []ColumnInfo{
//					{ID: 0, Cell: "User Settings (1.1)"},
//					{ID: 1, Cell: "User Preferences"},
//					{ID: 2, Cell: "Pending"},
//					{ID: 3, Cell: "Edit"},
//					{ID: 4, Cell: "Value 5"},
//					{ID: 5, Cell: "Value 6"},
//				},
//				Children: []*Node{
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "Role Management (1.1.1)"},
//							{ID: 1, Cell: "Manage User Roles"},
//							{ID: 2, Cell: "Active"},
//							{ID: 3, Cell: "Details"},
//							{ID: 4, Cell: "Value 5"},
//							{ID: 5, Cell: "Value 6"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Add Role (1.1.1.1)"},
//									{ID: 1, Cell: "Create New Role"},
//									{ID: 2, Cell: "Inactive"},
//									{ID: 3, Cell: "Add"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Edit Role (1.1.1.2)"},
//									{ID: 1, Cell: "Modify Existing Role"},
//									{ID: 2, Cell: "Active"},
//									{ID: 3, Cell: "Update"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//						},
//					},
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "Permission Management (1.1.2)"},
//							{ID: 1, Cell: "Manage Permissions"},
//							{ID: 2, Cell: "Active"},
//							{ID: 3, Cell: "Review"},
//							{ID: 4, Cell: "Value 5"},
//							{ID: 5, Cell: "Value 6"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "View Permission (1.1.2.1)"},
//									{ID: 1, Cell: "Check Permissions"},
//									{ID: 2, Cell: "Active"},
//									{ID: 3, Cell: "View"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Modify Permission (1.1.2.2)"},
//									{ID: 1, Cell: "Change Permission Settings"},
//									{ID: 2, Cell: "Inactive"},
//									{ID: 3, Cell: "Alter"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//						},
//					},
//				},
//			},
//			{
//				RowCells: []ColumnInfo{
//					{ID: 0, Cell: "Configuration Management (1.2)"},
//					{ID: 1, Cell: "System Configuration"},
//					{ID: 2, Cell: "Active"},
//					{ID: 3, Cell: "Configure"},
//					{ID: 4, Cell: "Value 5"},
//					{ID: 5, Cell: "Value 6"},
//				},
//				Children: []*Node{
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "System Parameters (1.2.1)"},
//							{ID: 1, Cell: "Adjust Settings"},
//							{ID: 2, Cell: "Active"},
//							{ID: 3, Cell: "Apply"},
//							{ID: 4, Cell: "Value 5"},
//							{ID: 5, Cell: "Value 6"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Update Parameters (1.2.1.1)"},
//									{ID: 1, Cell: "Change Parameters"},
//									{ID: 2, Cell: "Pending"},
//									{ID: 3, Cell: "Execute"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "View Parameters (1.2.1.2)"},
//									{ID: 1, Cell: "Check Configurations"},
//									{ID: 2, Cell: "Active"},
//									{ID: 3, Cell: "Inspect"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//						},
//					},
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "Log Management (1.2.2)"},
//							{ID: 1, Cell: "View System Logs"},
//							{ID: 2, Cell: "Critical"},
//							{ID: 3, Cell: "Monitor"},
//							{ID: 4, Cell: "Value 5"},
//							{ID: 5, Cell: "Value 6"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "View Logs (1.2.2.1)"},
//									{ID: 1, Cell: "Access Log Files"},
//									{ID: 2, Cell: "Active"},
//									{ID: 3, Cell: "Display"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Clear Logs (1.2.2.2)"},
//									{ID: 1, Cell: "Delete Old Entries"},
//									{ID: 2, Cell: "Inactive"},
//									{ID: 3, Cell: "Remove"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//	},
//	{
//		RowCells: []ColumnInfo{
//			{ID: 0, Cell: "System Management"},
//			{ID: 1, Cell: "System Overview"},
//			{ID: 2, Cell: "Active"},
//			{ID: 3, Cell: "Manage"},
//			{ID: 4, Cell: "Column 5"},
//			{ID: 5, Cell: "Column 6"},
//		},
//		Children: []*Node{
//			{
//				RowCells: []ColumnInfo{
//					{ID: 0, Cell: "User Settings (1.1)"},
//					{ID: 1, Cell: "User Preferences"},
//					{ID: 2, Cell: "Pending"},
//					{ID: 3, Cell: "Edit"},
//					{ID: 4, Cell: "Value 5"},
//					{ID: 5, Cell: "Value 6"},
//				},
//				Children: []*Node{
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "Role Management (1.1.1)"},
//							{ID: 1, Cell: "Manage User Roles"},
//							{ID: 2, Cell: "Active"},
//							{ID: 3, Cell: "Details"},
//							{ID: 4, Cell: "Value 5"},
//							{ID: 5, Cell: "Value 6"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Add Role (1.1.1.1)"},
//									{ID: 1, Cell: "Create New Role"},
//									{ID: 2, Cell: "Inactive"},
//									{ID: 3, Cell: "Add"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Edit Role (1.1.1.2)"},
//									{ID: 1, Cell: "Modify Existing Role"},
//									{ID: 2, Cell: "Active"},
//									{ID: 3, Cell: "Update"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//						},
//					},
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "Permission Management (1.1.2)"},
//							{ID: 1, Cell: "Manage Permissions"},
//							{ID: 2, Cell: "Active"},
//							{ID: 3, Cell: "Review"},
//							{ID: 4, Cell: "Value 5"},
//							{ID: 5, Cell: "Value 6"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "View Permission (1.1.2.1)"},
//									{ID: 1, Cell: "Check Permissions"},
//									{ID: 2, Cell: "Active"},
//									{ID: 3, Cell: "View"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Modify Permission (1.1.2.2)"},
//									{ID: 1, Cell: "Change Permission Settings"},
//									{ID: 2, Cell: "Inactive"},
//									{ID: 3, Cell: "Alter"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//						},
//					},
//				},
//			},
//			{
//				RowCells: []ColumnInfo{
//					{ID: 0, Cell: "Configuration Management (1.2)"},
//					{ID: 1, Cell: "System Configuration"},
//					{ID: 2, Cell: "Active"},
//					{ID: 3, Cell: "Configure"},
//					{ID: 4, Cell: "Value 5"},
//					{ID: 5, Cell: "Value 6"},
//				},
//				Children: []*Node{
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "System Parameters (1.2.1)"},
//							{ID: 1, Cell: "Adjust Settings"},
//							{ID: 2, Cell: "Active"},
//							{ID: 3, Cell: "Apply"},
//							{ID: 4, Cell: "Value 5"},
//							{ID: 5, Cell: "Value 6"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Update Parameters (1.2.1.1)"},
//									{ID: 1, Cell: "Change Parameters"},
//									{ID: 2, Cell: "Pending"},
//									{ID: 3, Cell: "Execute"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "View Parameters (1.2.1.2)"},
//									{ID: 1, Cell: "Check Configurations"},
//									{ID: 2, Cell: "Active"},
//									{ID: 3, Cell: "Inspect"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//						},
//					},
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "Log Management (1.2.2)"},
//							{ID: 1, Cell: "View System Logs"},
//							{ID: 2, Cell: "Critical"},
//							{ID: 3, Cell: "Monitor"},
//							{ID: 4, Cell: "Value 5"},
//							{ID: 5, Cell: "Value 6"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "View Logs (1.2.2.1)"},
//									{ID: 1, Cell: "Access Log Files"},
//									{ID: 2, Cell: "Active"},
//									{ID: 3, Cell: "Display"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Clear Logs (1.2.2.2)"},
//									{ID: 1, Cell: "Delete Old Entries"},
//									{ID: 2, Cell: "Inactive"},
//									{ID: 3, Cell: "Remove"},
//									{ID: 4, Cell: "Value 5"},
//									{ID: 5, Cell: "Value 6"},
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//	},
//}

//var TestRootRows1 = []*Node{
//	{
//		RowCells: []ColumnInfo{
//			{ID: 0, Cell: "System Management"},
//			{ID: 1, Cell: "System Overview"},
//			{ID: 2, Cell: "Active"},
//			{ID: 3, Cell: "Manage"},
//		},
//		Children: []*Node{
//			{
//				RowCells: []ColumnInfo{
//					{ID: 0, Cell: "User Settings (1.1)"},
//					{ID: 1, Cell: "User Preferences"},
//					{ID: 2, Cell: "Pending"},
//					{ID: 3, Cell: "Edit"},
//				},
//				Children: []*Node{
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "Role Management (1.1.1)"},
//							{ID: 1, Cell: "Manage User Roles"},
//							{ID: 2, Cell: "Active"},
//							{ID: 3, Cell: "Details"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Add Role (1.1.1.1)"},
//									{ID: 1, Cell: "Create New Role"},
//									{ID: 2, Cell: "Inactive"},
//									{ID: 3, Cell: "Add"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Edit Role (1.1.1.2)"},
//									{ID: 1, Cell: "Modify Existing Role"},
//									{ID: 2, Cell: "Active"},
//									{ID: 3, Cell: "Update"},
//								},
//							},
//						},
//					},
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "Permission Management (1.1.2)"},
//							{ID: 1, Cell: "Manage Permissions"},
//							{ID: 2, Cell: "Active"},
//							{ID: 3, Cell: "Review"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "View Permission (1.1.2.1)"},
//									{ID: 1, Cell: "Check Permissions"},
//									{ID: 2, Cell: "Active"},
//									{ID: 3, Cell: "View"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Modify Permission (1.1.2.2)"},
//									{ID: 1, Cell: "Change Permission Settings"},
//									{ID: 2, Cell: "Inactive"},
//									{ID: 3, Cell: "Alter"},
//								},
//							},
//						},
//					},
//				},
//			},
//			{
//				RowCells: []ColumnInfo{
//					{ID: 0, Cell: "Configuration Management (1.2)"},
//					{ID: 1, Cell: "System Configuration"},
//					{ID: 2, Cell: "Active"},
//					{ID: 3, Cell: "Configure"},
//				},
//				Children: []*Node{
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "System Parameters (1.2.1)"},
//							{ID: 1, Cell: "Adjust Settings"},
//							{ID: 2, Cell: "Active"},
//							{ID: 3, Cell: "Apply"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Update Parameters (1.2.1.1)"},
//									{ID: 1, Cell: "Change Parameters"},
//									{ID: 2, Cell: "Pending"},
//									{ID: 3, Cell: "Execute"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "View Parameters (1.2.1.2)"},
//									{ID: 1, Cell: "Check Configurations"},
//									{ID: 2, Cell: "Active"},
//									{ID: 3, Cell: "Inspect"},
//								},
//							},
//						},
//					},
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "Log Management (1.2.2)"},
//							{ID: 1, Cell: "View System Logs"},
//							{ID: 2, Cell: "Critical"},
//							{ID: 3, Cell: "Monitor"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "View Logs (1.2.2.1)"},
//									{ID: 1, Cell: "Access Log Files"},
//									{ID: 2, Cell: "Active"},
//									{ID: 3, Cell: "Display"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Clear Logs (1.2.2.2)"},
//									{ID: 1, Cell: "Delete Old Entries"},
//									{ID: 2, Cell: "Inactive"},
//									{ID: 3, Cell: "Remove"},
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//	},
//}
//
//var TestRootRows0 = []*Node{
//	{
//		RowCells: []ColumnInfo{
//			{ID: 0, Cell: "系统管理"},
//			{ID: 1, Cell: "描述 A"},
//			{ID: 2, Cell: "状态 A"},
//			{ID: 3, Cell: "操作 A"},
//		},
//		Children: []*Node{
//			{
//				RowCells: []ColumnInfo{
//					{ID: 0, Cell: "用户设置 (1.1)"},
//					{ID: 1, Cell: "描述 B"},
//					{ID: 2, Cell: "状态 B"},
//					{ID: 3, Cell: "操作 B"},
//				},
//				Children: []*Node{
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "角色管理 (1.1.1)"},
//							{ID: 1, Cell: "描述 C"},
//							{ID: 2, Cell: "状态 C"},
//							{ID: 3, Cell: "操作 C"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "添加角色 (1.1.1.1)"},
//									{ID: 1, Cell: "描述 D"},
//									{ID: 2, Cell: "状态 D"},
//									{ID: 3, Cell: "操作 D"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "编辑角色 (1.1.1.2)"},
//									{ID: 1, Cell: "描述 E"},
//									{ID: 2, Cell: "状态 E"},
//									{ID: 3, Cell: "操作 E"},
//								},
//							},
//						},
//					},
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "权限管理 (1.1.2)"},
//							{ID: 1, Cell: "描述 F"},
//							{ID: 2, Cell: "状态 F"},
//							{ID: 3, Cell: "操作 F"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "查看权限 (1.1.2.1)"},
//									{ID: 1, Cell: "描述 G"},
//									{ID: 2, Cell: "状态 G"},
//									{ID: 3, Cell: "操作 G"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "修改权限 (1.1.2.2)"},
//									{ID: 1, Cell: "描述 H"},
//									{ID: 2, Cell: "状态 H"},
//									{ID: 3, Cell: "操作 H"},
//								},
//							},
//						},
//					},
//				},
//			},
//			{
//				RowCells: []ColumnInfo{
//					{ID: 0, Cell: "配置管理 (1.2)"},
//					{ID: 1, Cell: "描述 I"},
//					{ID: 2, Cell: "状态 I"},
//					{ID: 3, Cell: "操作 I"},
//				},
//				Children: []*Node{
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "系统参数 (1.2.1)"},
//							{ID: 1, Cell: "描述 J"},
//							{ID: 2, Cell: "状态 J"},
//							{ID: 3, Cell: "操作 J"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "更新参数 (1.2.1.1)"},
//									{ID: 1, Cell: "描述 K"},
//									{ID: 2, Cell: "状态 K"},
//									{ID: 3, Cell: "操作 K"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "查看参数 (1.2.1.2)"},
//									{ID: 1, Cell: "描述 L"},
//									{ID: 2, Cell: "状态 L"},
//									{ID: 3, Cell: "操作 L"},
//								},
//							},
//						},
//					},
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "日志管理 (1.2.2)"},
//							{ID: 1, Cell: "描述 M"},
//							{ID: 2, Cell: "状态 M"},
//							{ID: 3, Cell: "操作 M"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "查看日志 (1.2.2.1)"},
//									{ID: 1, Cell: "描述 N"},
//									{ID: 2, Cell: "状态 N"},
//									{ID: 3, Cell: "操作 N"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "清除日志 (1.2.2.2)"},
//									{ID: 1, Cell: "描述 O"},
//									{ID: 2, Cell: "状态 O"},
//									{ID: 3, Cell: "操作 O"},
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//	},
//}

//var TestRootRows = []*Node{
//	{
//		RowCells: []ColumnInfo{
//			{ID: 0, Cell: "System Management"},
//			{ID: 1, Cell: "Description 1"},
//			{ID: 2, Cell: "Status 1"},
//			{ID: 3, Cell: "Action 1"},
//		},
//		Children: []*Node{
//			{
//				RowCells: []ColumnInfo{
//					{ID: 0, Cell: "User Settings (1.1)"},
//					{ID: 1, Cell: "Description 2"},
//					{ID: 2, Cell: "Status 2"},
//					{ID: 3, Cell: "Action 2"},
//				},
//				Children: []*Node{
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "Role Management (1.1.1)"},
//							{ID: 1, Cell: "Description 3"},
//							{ID: 2, Cell: "Status 3"},
//							{ID: 3, Cell: "Action 3"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Add Role (1.1.1.1)"},
//									{ID: 1, Cell: "Description 4"},
//									{ID: 2, Cell: "Status 4"},
//									{ID: 3, Cell: "Action 4"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Edit Role (1.1.1.2)"},
//									{ID: 1, Cell: "Description 5"},
//									{ID: 2, Cell: "Status 5"},
//									{ID: 3, Cell: "Action 5"},
//								},
//							},
//						},
//					},
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "Permission Management (1.1.2)"},
//							{ID: 1, Cell: "Description 6"},
//							{ID: 2, Cell: "Status 6"},
//							{ID: 3, Cell: "Action 6"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "View Permission (1.1.2.1)"},
//									{ID: 1, Cell: "Description 7"},
//									{ID: 2, Cell: "Status 7"},
//									{ID: 3, Cell: "Action 7"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Modify Permission (1.1.2.2)"},
//									{ID: 1, Cell: "Description 8"},
//									{ID: 2, Cell: "Status 8"},
//									{ID: 3, Cell: "Action 8"},
//								},
//							},
//						},
//					},
//				},
//			},
//			{
//				RowCells: []ColumnInfo{
//					{ID: 0, Cell: "Configuration Management (1.2)"},
//					{ID: 1, Cell: "Description 9"},
//					{ID: 2, Cell: "Status 9"},
//					{ID: 3, Cell: "Action 9"},
//				},
//				Children: []*Node{
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "System Parameters (1.2.1)"},
//							{ID: 1, Cell: "Description 10"},
//							{ID: 2, Cell: "Status 10"},
//							{ID: 3, Cell: "Action 10"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Update Parameters (1.2.1.1)"},
//									{ID: 1, Cell: "Description 11"},
//									{ID: 2, Cell: "Status 11"},
//									{ID: 3, Cell: "Action 11"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "View Parameters (1.2.1.2)"},
//									{ID: 1, Cell: "Description 12"},
//									{ID: 2, Cell: "Status 12"},
//									{ID: 3, Cell: "Action 12"},
//								},
//							},
//						},
//					},
//					{
//						RowCells: []ColumnInfo{
//							{ID: 0, Cell: "Log Management (1.2.2)"},
//							{ID: 1, Cell: "Description 13"},
//							{ID: 2, Cell: "Status 13"},
//							{ID: 3, Cell: "Action 13"},
//						},
//						Children: []*Node{
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "View Logs (1.2.2.1)"},
//									{ID: 1, Cell: "Description 14"},
//									{ID: 2, Cell: "Status 14"},
//									{ID: 3, Cell: "Action 14"},
//								},
//							},
//							{
//								RowCells: []ColumnInfo{
//									{ID: 0, Cell: "Clear Logs (1.2.2.2)"},
//									{ID: 1, Cell: "Description 15"},
//									{ID: 2, Cell: "Status 15"},
//									{ID: 3, Cell: "Action 15"},
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//	},
//}

type GoroutineList struct {
	cols []Column2
	rows []*packet
	packet
	table      *Table2
	sortColumn int  // 当前排序的列索引
	sortOrder  bool // true 为升序，false 为降序
}

func NewGoroutineList(rows []*packet) *GoroutineList {
	cols := []Column2{
		{
			Name:      "Scheme",
			Width:     0,
			MinWidth:  0,
			Alignment: text.Start,
		},
		{
			Name:      "Method",
			Width:     0,
			MinWidth:  0,
			Alignment: text.Start,
		},
		{
			Name:      "Host",
			Width:     0,
			MinWidth:  0,
			Alignment: text.Start,
		},
		{
			Name:      "Path",
			Width:     0,
			MinWidth:  0,
			Alignment: text.Start,
		},
		{
			Name:      "ContentType",
			Width:     0,
			MinWidth:  0,
			Alignment: text.Start,
		},
		{
			Name:      "ContentLength",
			Width:     0,
			MinWidth:  0,
			Alignment: text.Start,
		},
		{
			Name:      "Status",
			Width:     0,
			MinWidth:  0,
			Alignment: text.Start,
		},
		{
			Name:      "Note",
			Width:     0,
			MinWidth:  0,
			Alignment: text.Start,
		},
		{
			Name:      "Process",
			Width:     0,
			MinWidth:  0,
			Alignment: text.Start,
		},
		{
			Name:      "PadTime",
			Width:     0,
			MinWidth:  0,
			Alignment: text.Start,
		},
	}
	table := NewTable2(cols)
	table.SortedBy = 0
	table.SortOrder = SortAscending
	return &GoroutineList{
		cols:       cols,
		rows:       rows,
		packet:     packet{},
		table:      table,
		sortColumn: 0,
		sortOrder:  false,
	}
}

func (g *GoroutineList) SetGoroutines(gtx layout.Context) {
	cellData := g.GetCellData()               // 获取当前的单元格数据
	g.table.SetColumns(gtx, g.cols, cellData) //调整列宽？没意义，需要根据单元格最大宽度调整当前列的列宽
	g.Sort(g.rows)
}

func (g *GoroutineList) Sort(gs []*packet) {
	// 根据 sortColumn 排序
	slices.SortFunc(gs, func(a, b *packet) int {
		var result int
		switch g.table.Columns[g.table.SortedBy].Name {
		case "Scheme":
			result = strings.Compare(a.Note, b.Note)
		case "Method":
			result = strings.Compare(a.Process, b.Process)
		case "Host":
			result = strings.Compare(a.Host, b.Host)
		case "Path":
			result = strings.Compare(a.Path, b.Path)
		case "ContentType":
			result = strings.Compare(a.ContentType, b.ContentType)
		case "ContentLength":
			result = a.ContentLength
		case "Status":
			result = strings.Compare(a.Status, b.Status)
		case "Note":
			result = strings.Compare(a.Note, b.Note)
		case "PadTime":
			result = int(a.PadTime - b.PadTime)
		default:
			result = 0
		}
		// 如果是降序，则反转结果
		if !g.sortOrder {
			result = -result
		}
		return result
	})
}

func (g *GoroutineList) Update(gtx layout.Context) {
	g.SetGoroutines(gtx)
	g.table.Update(gtx) //得到被点击的列索引
	if clickedColumnIndex, ok := g.table.SortByClickedColumn(); ok {
		if clickedColumnIndex == g.sortColumn {
			g.sortOrder = !g.sortOrder // 切换排序方向
		} else {
			g.sortColumn = clickedColumnIndex // 设置新的排序列
			g.sortOrder = true                // 默认升序
		}
		g.Sort(g.rows)                  // 重新设置排序后的 goroutines
		g.table.clickedColumnIndex = -1 // 重置点击列索引
	}
}

func (g *GoroutineList) Layout(gtx layout.Context) layout.Dimensions {
	g.Update(gtx)
	cellFn := func(gtx layout.Context, row, col int) layout.Dimensions {
		defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
		switch colName := g.table.Columns[col].Name; colName {
		case "Scheme":
			return LabelCell(g.rows[row].Scheme).Layout(gtx)
		case "Method":
			return LabelCell(g.rows[row].Method).Layout(gtx)
		case "Host":
			return LabelCell(g.rows[row].Host).Layout(gtx)
		case "Path":
			return LabelCell(g.rows[row].Path).Layout(gtx)
		case "ContentType":
			return LabelCell(g.rows[row].ContentType).Layout(gtx)
		case "ContentLength":
			return LabelCell(strconv.Itoa(g.rows[row].ContentLength)).Layout(gtx)
		case "Status":
			return LabelCell(g.rows[row].Status).Layout(gtx)
		case "Note":
			return LabelCell(g.rows[row].Note).Layout(gtx)
		case "Process":
			return LabelCell(g.rows[row].Process).Layout(gtx)
		case "PadTime":
			return LabelCell(g.rows[row].PadTime.String()).Layout(gtx)
		}
		return layout.Dimensions{}
	}
	return SimpleTable(
		gtx,
		g.table,
		len(g.rows),
		cellFn,
	)
}

// 在 GoroutineList 结构体中添加一个方法来获取 cellData
func (g *GoroutineList) GetCellData() [][]string {
	cellData := make([][]string, len(g.rows))
	for i, row := range g.rows {
		cellData[i] = []string{
			row.Scheme,
			row.Method,
			row.Host,
			row.Path,
			row.ContentType,
			strconv.Itoa(row.ContentLength),
			row.Status,
			row.Note,
			row.Process,
			row.PadTime.String(),
		}
	}
	return cellData
}

func LabelCell(label string) material.LabelStyle {
	l := material.Label(th.Theme, 12, label)
	l.MaxLines = 1
	return l
}

var Packets = []*packet{
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource1", ContentType: "application/json", ContentLength: 100, Status: http.StatusText(http.StatusOK), Note: "获取资源1", Process: "process1.exe", PadTime: 30},
	{Scheme: "https", Method: http.MethodPost, Host: "example.com", Path: "/api/v1/resource2", ContentType: "application/xml", ContentLength: 150, Status: http.StatusText(http.StatusCreated), Note: "创建资源2", Process: "process2.exe", PadTime: 25},
	{Scheme: "http", Method: http.MethodDelete, Host: "other.com", Path: "/api/v1/resource3", ContentType: "application/json", ContentLength: 200, Status: http.StatusText(http.StatusNoContent), Note: "删除资源3", Process: "process3.exe", PadTime: 20},
	{Scheme: "https", Method: http.MethodPut, Host: "another.com", Path: "/api/v1/resource4", ContentType: "text/plain", ContentLength: 250, Status: http.StatusText(http.StatusOK), Note: "更新资源4", Process: "process4.exe", PadTime: 15},
	{Scheme: "http", Method: http.MethodPatch, Host: "example.com", Path: "/api/v1/resource5", ContentType: "application/json", ContentLength: 300, Status: http.StatusText(http.StatusOK), Note: "修改资源5", Process: "process5.exe", PadTime: 10},
	{Scheme: "http", Method: http.MethodHead, Host: "test.com", Path: "/api/v1/resource6", ContentType: "application/json", ContentLength: 50, Status: http.StatusText(http.StatusOK), Note: "头信息请求", Process: "process6.exe", PadTime: 35},
	{Scheme: "https", Method: http.MethodOptions, Host: "example.org", Path: "/api/v1/resource7", ContentType: "application/json", ContentLength: 120, Status: http.StatusText(http.StatusOK), Note: "选项请求", Process: "process7.exe", PadTime: 40},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource8", ContentType: "application/json", ContentLength: 90, Status: http.StatusText(http.StatusOK), Note: "获取资源8", Process: "process8.exe", PadTime: 30},
	{Scheme: "http", Method: http.MethodGet, Host: "example.net", Path: "/api/v1/resource9", ContentType: "application/json", ContentLength: 110, Status: http.StatusText(http.StatusOK), Note: "获取资源9", Process: "process9.exe", PadTime: 28},
	{Scheme: "http", Method: http.MethodGet, Host: "example.org", Path: "/api/v1/resource10", ContentType: "text/html", ContentLength: 95, Status: http.StatusText(http.StatusOK), Note: "获取资源10", Process: "process10.exe", PadTime: 20},
	{Scheme: "http", Method: http.MethodPost, Host: "sample.com", Path: "/api/v1/resource11", ContentType: "application/json", ContentLength: 75, Status: http.StatusText(http.StatusCreated), Note: "创建资源11", Process: "process11.exe", PadTime: 22},
	{Scheme: "https", Method: http.MethodPut, Host: "sample.com", Path: "/api/v1/resource12", ContentType: "application/json", ContentLength: 160, Status: http.StatusText(http.StatusOK), Note: "更新资源12", Process: "process12.exe", PadTime: 18},
	{Scheme: "http", Method: http.MethodDelete, Host: "example.com", Path: "/api/v1/resource13", ContentType: "application/json", ContentLength: 130, Status: http.StatusText(http.StatusNoContent), Note: "删除资源13", Process: "process13.exe", PadTime: 24},
	{Scheme: "http", Method: http.MethodPatch, Host: "another.com", Path: "/api/v1/resource14", ContentType: "application/json", ContentLength: 140, Status: http.StatusText(http.StatusOK), Note: "修改资源14", Process: "process14.exe", PadTime: 23},
	{Scheme: "http", Method: http.MethodHead, Host: "example.org", Path: "/api/v1/resource15", ContentType: "application/xml", ContentLength: 50, Status: http.StatusText(http.StatusOK), Note: "头信息请求", Process: "process15.exe", PadTime: 19},
	{Scheme: "https", Method: http.MethodOptions, Host: "test.com", Path: "/api/v1/resource16", ContentType: "application/json", ContentLength: 80, Status: http.StatusText(http.StatusOK), Note: "选项请求", Process: "process16.exe", PadTime: 29},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource17", ContentType: "application/json", ContentLength: 110, Status: http.StatusText(http.StatusOK), Note: "获取资源17", Process: "process17.exe", PadTime: 31},
	{Scheme: "https", Method: http.MethodPost, Host: "service.com", Path: "/api/v1/resource18", ContentType: "application/json", ContentLength: 90, Status: http.StatusText(http.StatusCreated), Note: "创建资源18", Process: "process18.exe", PadTime: 27},
	{Scheme: "http", Method: http.MethodDelete, Host: "example.net", Path: "/api/v1/resource19", ContentType: "application/json", ContentLength: 120, Status: http.StatusText(http.StatusNoContent), Note: "删除资源19", Process: "process19.exe", PadTime: 21},
	{Scheme: "https", Method: http.MethodPut, Host: "example.com", Path: "/api/v1/resource20", ContentType: "application/json", ContentLength: 160, Status: http.StatusText(http.StatusOK), Note: "更新资源20", Process: "process20.exe", PadTime: 30},
	{Scheme: "http", Method: http.MethodPatch, Host: "myapp.com", Path: "/api/v1/resource21", ContentType: "application/json", ContentLength: 140, Status: http.StatusText(http.StatusOK), Note: "修改资源21", Process: "process21.exe", PadTime: 12},
	{Scheme: "http", Method: http.MethodHead, Host: "testapp.com", Path: "/api/v1/resource22", ContentType: "application/json", ContentLength: 55, Status: http.StatusText(http.StatusOK), Note: "头信息请求", Process: "process22.exe", PadTime: 11},
	{Scheme: "https", Method: http.MethodOptions, Host: "anotherexample.com", Path: "/api/v1/resource23", ContentType: "application/json", ContentLength: 95, Status: http.StatusText(http.StatusOK), Note: "选项请求", Process: "process23.exe", PadTime: 20},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource24", ContentType: "application/json", ContentLength: 78, Status: http.StatusText(http.StatusOK), Note: "获取资源24", Process: "process24.exe", PadTime: 30},
	{Scheme: "https", Method: http.MethodPost, Host: "example.com", Path: "/api/v1/resource25", ContentType: "application/xml", ContentLength: 118, Status: http.StatusText(http.StatusCreated), Note: "创建资源25", Process: "process25.exe", PadTime: 29},
	{Scheme: "http", Method: http.MethodDelete, Host: "example.com", Path: "/api/v1/resource26", ContentType: "application/json", ContentLength: 55, Status: http.StatusText(http.StatusNoContent), Note: "删除资源26", Process: "process26.exe", PadTime: 28},
	{Scheme: "https", Method: http.MethodPut, Host: "example.net", Path: "/api/v1/resource27", ContentType: "text/plain", ContentLength: 88, Status: http.StatusText(http.StatusOK), Note: "更新资源27", Process: "process27.exe", PadTime: 27},
	{Scheme: "http", Method: http.MethodPatch, Host: "experimental.com", Path: "/api/v1/resource28", ContentType: "application/json", ContentLength: 130, Status: http.StatusText(http.StatusOK), Note: "修改资源28", Process: "process28.exe", PadTime: 20},
	{Scheme: "http", Method: http.MethodHead, Host: "testsite.com", Path: "/api/v1/resource29", ContentType: "application/json", ContentLength: 100, Status: http.StatusText(http.StatusOK), Note: "头信息请求", Process: "process29.exe", PadTime: 30},
	{Scheme: "https", Method: http.MethodOptions, Host: "example.org", Path: "/api/v1/resource30", ContentType: "application/json", ContentLength: 76, Status: http.StatusText(http.StatusOK), Note: "选项请求", Process: "process30.exe", PadTime: 30},

	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource31", ContentType: "application/json", ContentLength: 99, Status: http.StatusText(http.StatusOK), Note: "获取资源31", Process: "process31.exe", PadTime: 32},
	{Scheme: "https", Method: http.MethodPost, Host: "example.com", Path: "/api/v1/resource32", ContentType: "application/json", ContentLength: 152, Status: http.StatusText(http.StatusCreated), Note: "创建资源32", Process: "process32.exe", PadTime: 22},
	{Scheme: "http", Method: http.MethodDelete, Host: "example.com", Path: "/api/v1/resource33", ContentType: "application/json", ContentLength: 202, Status: http.StatusText(http.StatusNoContent), Note: "删除资源33", Process: "process33.exe", PadTime: 20},
	{Scheme: "https", Method: http.MethodPut, Host: "example.com", Path: "/api/v1/resource34", ContentType: "application/json", ContentLength: 250, Status: http.StatusText(http.StatusOK), Note: "更新资源34", Process: "process34.exe", PadTime: 14},
	{Scheme: "http", Method: http.MethodPatch, Host: "example.com", Path: "/api/v1/resource35", ContentType: "application/json", ContentLength: 105, Status: http.StatusText(http.StatusOK), Note: "修改资源35", Process: "process35.exe", PadTime: 31},
	{Scheme: "http", Method: http.MethodHead, Host: "example.com", Path: "/api/v1/resource36", ContentType: "application/json", ContentLength: 98, Status: http.StatusText(http.StatusOK), Note: "头信息请求", Process: "process36.exe", PadTime: 12},
	{Scheme: "https", Method: http.MethodOptions, Host: "example.com", Path: "/api/v1/resource37", ContentType: "application/json", ContentLength: 130, Status: http.StatusText(http.StatusOK), Note: "选项请求", Process: "process37.exe", PadTime: 28},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource38", ContentType: "application/json", ContentLength: 112, Status: http.StatusText(http.StatusOK), Note: "获取资源38", Process: "process38.exe", PadTime: 22},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource39", ContentType: "application/json", ContentLength: 102, Status: http.StatusText(http.StatusOK), Note: "获取资源39", Process: "process39.exe", PadTime: 26},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource40", ContentType: "text/html", ContentLength: 87, Status: http.StatusText(http.StatusOK), Note: "获取资源40", Process: "process40.exe", PadTime: 20},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource41", ContentType: "application/json", ContentLength: 99, Status: http.StatusText(http.StatusOK), Note: "获取资源41", Process: "process41.exe", PadTime: 32},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource42", ContentType: "application/json", ContentLength: 79, Status: http.StatusText(http.StatusOK), Note: "获取资源42", Process: "process42.exe", PadTime: 33},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource43", ContentType: "application/json", ContentLength: 150, Status: http.StatusText(http.StatusOK), Note: "获取资源43", Process: "process43.exe", PadTime: 35},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource44", ContentType: "application/json", ContentLength: 88, Status: http.StatusText(http.StatusOK), Note: "获取资源44", Process: "process44.exe", PadTime: 30},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource45", ContentType: "application/json", ContentLength: 76, Status: http.StatusText(http.StatusOK), Note: "获取资源45", Process: "process45.exe", PadTime: 20},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource46", ContentType: "application/json", ContentLength: 145, Status: http.StatusText(http.StatusOK), Note: "获取资源46", Process: "process46.exe", PadTime: 21},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource47", ContentType: "application/json", ContentLength: 80, Status: http.StatusText(http.StatusOK), Note: "获取资源47", Process: "process47.exe", PadTime: 20},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource48", ContentType: "application/json", ContentLength: 98, Status: http.StatusText(http.StatusOK), Note: "获取资源48", Process: "process48.exe", PadTime: 24},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource49", ContentType: "application/json", ContentLength: 120, Status: http.StatusText(http.StatusOK), Note: "获取资源49", Process: "process49.exe", PadTime: 25},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource50", ContentType: "application/json", ContentLength: 88, Status: http.StatusText(http.StatusOK), Note: "获取资源50", Process: "process50.exe", PadTime: 21},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource51", ContentType: "application/json", ContentLength: 200, Status: http.StatusText(http.StatusOK), Note: "获取资源51", Process: "process51.exe", PadTime: 29},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource52", ContentType: "application/json", ContentLength: 210, Status: http.StatusText(http.StatusOK), Note: "获取资源52", Process: "process52.exe", PadTime: 27},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource53", ContentType: "application/json", ContentLength: 190, Status: http.StatusText(http.StatusOK), Note: "获取资源53", Process: "process53.exe", PadTime: 33},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource54", ContentType: "application/json", ContentLength: 180, Status: http.StatusText(http.StatusOK), Note: "获取资源54", Process: "process54.exe", PadTime: 36},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource55", ContentType: "application/json", ContentLength: 170, Status: http.StatusText(http.StatusOK), Note: "获取资源55", Process: "process55.exe", PadTime: 28},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource56", ContentType: "application/json", ContentLength: 160, Status: http.StatusText(http.StatusOK), Note: "获取资源56", Process: "process56.exe", PadTime: 12},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource57", ContentType: "application/json", ContentLength: 150, Status: http.StatusText(http.StatusOK), Note: "获取资源57", Process: "process57.exe", PadTime: 30},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource58", ContentType: "application/json", ContentLength: 140, Status: http.StatusText(http.StatusOK), Note: "获取资源58", Process: "process58.exe", PadTime: 30},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource59", ContentType: "application/json", ContentLength: 130, Status: http.StatusText(http.StatusOK), Note: "获取资源59", Process: "process59.exe", PadTime: 35},
	{Scheme: "http", Method: http.MethodGet, Host: "example.com", Path: "/api/v1/resource60", ContentType: "application/json", ContentLength: 120, Status: http.StatusText(http.StatusOK), Note: "获取资源60", Process: "process60.exe", PadTime: 38},
}