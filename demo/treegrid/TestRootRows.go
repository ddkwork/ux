package main

var TestRootRows = []*Node{
	{
		RowCells: []ColumnInfo{
			{ID: 0, Cell: "System Management"},
			{ID: 1, Cell: "System Overview"},
			{ID: 2, Cell: "Active"},
			{ID: 3, Cell: "Manage"},
			{ID: 4, Cell: "Column 5"},
			{ID: 5, Cell: "Column 6"},
		},
		Children: []*Node{
			{
				RowCells: []ColumnInfo{
					{ID: 0, Cell: "User Settings (1.1)"},
					{ID: 1, Cell: "User Preferences"},
					{ID: 2, Cell: "Pending"},
					{ID: 3, Cell: "Edit"},
					{ID: 4, Cell: "Value 5"},
					{ID: 5, Cell: "Value 6"},
				},
				Children: []*Node{
					{
						RowCells: []ColumnInfo{
							{ID: 0, Cell: "Role Management (1.1.1)"},
							{ID: 1, Cell: "Manage User Roles"},
							{ID: 2, Cell: "999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999"},
							{ID: 3, Cell: "Details"},
							{ID: 4, Cell: "Value 5"},
							{ID: 5, Cell: "Value 6"},
						},
						Children: []*Node{
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "Add Role (1.1.1.1)"},
									{ID: 1, Cell: "Create New Role"},
									{ID: 2, Cell: "Inactive"},
									{ID: 3, Cell: "Add"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "Edit Role (1.1.1.2)"},
									{ID: 1, Cell: "Modify Existing Role"},
									{ID: 2, Cell: "Active"},
									{ID: 3, Cell: "Update"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
						},
					},
					{
						RowCells: []ColumnInfo{
							{ID: 0, Cell: "Permission Management (1.1.2)"},
							{ID: 1, Cell: "Manage Permissions"},
							{ID: 2, Cell: "Active"},
							{ID: 3, Cell: "Review"},
							{ID: 4, Cell: "Value 5"},
							{ID: 5, Cell: "Value 6"},
						},
						Children: []*Node{
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "View Permission (1.1.2.1)"},
									{ID: 1, Cell: "Check Permissions"},
									{ID: 2, Cell: "Active"},
									{ID: 3, Cell: "View"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "Modify Permission (1.1.2.2)"},
									{ID: 1, Cell: "Change Permission Settings"},
									{ID: 2, Cell: "Inactive"},
									{ID: 3, Cell: "Alter"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
						},
					},
				},
			},
			{
				RowCells: []ColumnInfo{
					{ID: 0, Cell: "Configuration Management (1.2)"},
					{ID: 1, Cell: "System Configuration"},
					{ID: 2, Cell: "Active"},
					{ID: 3, Cell: "Configure"},
					{ID: 4, Cell: "Value 5"},
					{ID: 5, Cell: "Value 6"},
				},
				Children: []*Node{
					{
						RowCells: []ColumnInfo{
							{ID: 0, Cell: "System Parameters (1.2.1)"},
							{ID: 1, Cell: "Adjust Settings"},
							{ID: 2, Cell: "Active"},
							{ID: 3, Cell: "Apply"},
							{ID: 4, Cell: "Value 5"},
							{ID: 5, Cell: "Value 6"},
						},
						Children: []*Node{
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "Update Parameters (1.2.1.1)"},
									{ID: 1, Cell: "Change Parameters"},
									{ID: 2, Cell: "Pending"},
									{ID: 3, Cell: "Execute"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "View Parameters (1.2.1.2)"},
									{ID: 1, Cell: "Check Configurations"},
									{ID: 2, Cell: "Active"},
									{ID: 3, Cell: "Inspect"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
						},
					},
					{
						RowCells: []ColumnInfo{
							{ID: 0, Cell: "Log Management (1.2.2)"},
							{ID: 1, Cell: "View System Logs"},
							{ID: 2, Cell: "Critical"},
							{ID: 3, Cell: "Monitor"},
							{ID: 4, Cell: "Value 5"},
							{ID: 5, Cell: "Value 6"},
						},
						Children: []*Node{
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "View Logs (1.2.2.1)"},
									{ID: 1, Cell: "Access Log Files"},
									{ID: 2, Cell: "Active"},
									{ID: 3, Cell: "Display"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "Clear Logs (1.2.2.2)"},
									{ID: 1, Cell: "Delete Old Entries"},
									{ID: 2, Cell: "Inactive"},
									{ID: 3, Cell: "Remove"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		RowCells: []ColumnInfo{
			{ID: 0, Cell: "System Management"},
			{ID: 1, Cell: "System Overview"},
			{ID: 2, Cell: "Active"},
			{ID: 3, Cell: "Manage"},
			{ID: 4, Cell: "Column 5"},
			{ID: 5, Cell: "Column 6"},
		},
		Children: []*Node{
			{
				RowCells: []ColumnInfo{
					{ID: 0, Cell: "User Settings (1.1)"},
					{ID: 1, Cell: "User Preferences"},
					{ID: 2, Cell: "Pending"},
					{ID: 3, Cell: "Edit"},
					{ID: 4, Cell: "Value 5"},
					{ID: 5, Cell: "Value 6"},
				},
				Children: []*Node{
					{
						RowCells: []ColumnInfo{
							{ID: 0, Cell: "Role Management (1.1.1)"},
							{ID: 1, Cell: "Manage User Roles"},
							{ID: 2, Cell: "Active"},
							{ID: 3, Cell: "Details"},
							{ID: 4, Cell: "Value 5"},
							{ID: 5, Cell: "Value 6"},
						},
						Children: []*Node{
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "Add Role (1.1.1.1)"},
									{ID: 1, Cell: "Create New Role"},
									{ID: 2, Cell: "Inactive"},
									{ID: 3, Cell: "Add"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "Edit Role (1.1.1.2)"},
									{ID: 1, Cell: "Modify Existing Role"},
									{ID: 2, Cell: "Active"},
									{ID: 3, Cell: "Update"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
						},
					},
					{
						RowCells: []ColumnInfo{
							{ID: 0, Cell: "Permission Management (1.1.2)"},
							{ID: 1, Cell: "Manage Permissions"},
							{ID: 2, Cell: "Active"},
							{ID: 3, Cell: "Review"},
							{ID: 4, Cell: "Value 5"},
							{ID: 5, Cell: "Value 6"},
						},
						Children: []*Node{
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "View Permission (1.1.2.1)"},
									{ID: 1, Cell: "Check Permissions"},
									{ID: 2, Cell: "Active"},
									{ID: 3, Cell: "View"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "Modify Permission (1.1.2.2)"},
									{ID: 1, Cell: "Change Permission Settings"},
									{ID: 2, Cell: "Inactive"},
									{ID: 3, Cell: "Alter"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
						},
					},
				},
			},
			{
				RowCells: []ColumnInfo{
					{ID: 0, Cell: "Configuration Management (1.2)"},
					{ID: 1, Cell: "System Configuration"},
					{ID: 2, Cell: "Active"},
					{ID: 3, Cell: "Configure"},
					{ID: 4, Cell: "Value 5"},
					{ID: 5, Cell: "Value 6"},
				},
				Children: []*Node{
					{
						RowCells: []ColumnInfo{
							{ID: 0, Cell: "System Parameters (1.2.1)"},
							{ID: 1, Cell: "Adjust Settings"},
							{ID: 2, Cell: "Active"},
							{ID: 3, Cell: "Apply"},
							{ID: 4, Cell: "Value 5"},
							{ID: 5, Cell: "Value 6"},
						},
						Children: []*Node{
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "Update Parameters (1.2.1.1)"},
									{ID: 1, Cell: "Change Parameters"},
									{ID: 2, Cell: "Pending"},
									{ID: 3, Cell: "Execute"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "View Parameters (1.2.1.2)"},
									{ID: 1, Cell: "Check Configurations"},
									{ID: 2, Cell: "Active"},
									{ID: 3, Cell: "Inspect"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
						},
					},
					{
						RowCells: []ColumnInfo{
							{ID: 0, Cell: "Log Management (1.2.2)"},
							{ID: 1, Cell: "View System Logs"},
							{ID: 2, Cell: "Critical"},
							{ID: 3, Cell: "Monitor"},
							{ID: 4, Cell: "Value 5"},
							{ID: 5, Cell: "Value 6"},
						},
						Children: []*Node{
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "View Logs (1.2.2.1)"},
									{ID: 1, Cell: "Access Log Files"},
									{ID: 2, Cell: "Active"},
									{ID: 3, Cell: "Display"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
							{
								RowCells: []ColumnInfo{
									{ID: 0, Cell: "Clear Logs (1.2.2.2)"},
									{ID: 1, Cell: "Delete Old Entries"},
									{ID: 2, Cell: "Inactive"},
									{ID: 3, Cell: "Remove"},
									{ID: 4, Cell: "Value 5"},
									{ID: 5, Cell: "Value 6"},
								},
							},
						},
					},
				},
			},
		},
	},
}
