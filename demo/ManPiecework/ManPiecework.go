package main

import "time"

type (
	BossInfo struct { // 收货地点的提示内容
		Name  string
		Phone string
		Place string // 地点
	}
	DriverInfo struct { // 驾驶员信息,日期的提示内容
		Name        string
		Phone       string
		destination string // 目的地
	}
	// 日结：root容器节点=FoamBox容器节点+PlasticFrame容器节点
	ManPieceworkEditData struct { // 男工计件,用于统计装车
		Time        time.Time `table:"日期"` // 装车日期
		Destination string    `table:"目的地"`

		NumberOfCompleteVehicles int `table:"总件数"`
		FoamBox                  int `table:"泡沫箱"`
		PlasticFrame             int `table:"胶框"`

		Mark1              int `table:"1标"`
		Mark1Bar           int `table:"1标带靶"`
		Mark2              int `table:"2标"`
		Mark3              int `table:"3标"`
		Level3             int `table:"三层"`
		Supermarket3Labels int `table:"超市3标"`
		SmallFruit         int `table:"小果"`
		BigFruit           int `table:"大果"`
		FlowersFruits      int `table:"花果"`

		Mark1_J            int `table:"1标(胶)"`
		Mark2_J            int `table:"2标(胶)"`
		Mark3_J            int `table:"3标(胶)"`
		SmallFruitMark1    int `table:"小果一标(胶)"`
		SmallBunchFruits   int `table:"小串"`
		BunchFruits        int `table:"串果"`
		FineBunchFruits    int `table:"精串"`
		SweetBunchFruits   int `table:"甜串"`
		BunchFruitsAA      int `table:"串AA"`
		AAA                int `table:"AAA"`
		TurnoverBasket     int `table:"周转筐"`
		ElectronicCommerce int `table:"电商"`

		Carpooling int    `table:"拼车"`
		Note       string `table:"备注"`

		// 公用日结结构体
		// 男工报表显示装车费用日结
		// 女工报表显示打包费用日结
		// todo 打包费用没想好怎么显示好一点
		// 打包和装车日结
		// DailySettlement       int       // 日结
		// SettlementDate        time.Date // 结算日期
		// NumberOfCheckoutItems int       // 结账件数,加减人会改变
		// Price                 float32   // 单价
		// TotalPeople           int       // 人数,加减人会改变,根据离职日期计算天数？
		// PerCapitaWage         float64   // 人均工资
		// Settled               float64   // 已结算
		// Unsettled             float64   // 未结算
		// Tags                  []string  // 过滤条件列表
		// Enable                bool      // 是否离职
		// Note                  string    // 备注
	}
	// Salary 底薪
	Cycle struct {
		StartDate   time.Time // 开始日期
		EndDate     time.Time // 结束日期
		Days        int       // 天数
		TotalPeople int       // 人数
	}
)

func (m ManPieceworkEditData) FoamBoxSum() int {
	return m.Mark1 + m.Mark1Bar + m.Mark2 + m.Mark3 + m.Level3 + m.Supermarket3Labels + m.SmallFruit + m.BigFruit + m.FlowersFruits
}

func (m ManPieceworkEditData) PlasticFrameSum() int {
	return m.Mark1_J + m.Mark2_J + m.Mark3_J + m.SmallFruitMark1 +
		m.SmallBunchFruits + m.BunchFruits + m.FineBunchFruits + m.SweetBunchFruits + m.BunchFruitsAA + m.AAA + m.TurnoverBasket + m.ElectronicCommerce
}

func (m ManPieceworkEditData) NumberOfCompleteVehiclesSum() int {
	return m.FoamBoxSum() + m.PlasticFrameSum()
}
