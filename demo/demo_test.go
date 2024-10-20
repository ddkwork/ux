package main

import (
	"testing"

	"github.com/ddkwork/golibrary/stream"
)

func TestName(t *testing.T) {
	stream.NewGeneratedFile().SetPackageName("main").Types("demo",
		[]string{
			"treeTable",
			"tree",
			"table2",
			"table",
			"colorPicker",
			"card",
			"SearchDropDown",

			"jsonTree",
			"svgButton",
			"codeEditor",
			"asmView",
			"logView",

			"comBox",
			"splitView",
			"mobile",
			"listView",
			"structView",
			"terminal",
			"stackView",
			"dockView",
			"gif123",
			"hexEditor",
			"imageEditor",
			"mediaPlayer",
			"mind", // 思维导图，流程图
			"pdfView",
			"mapView",
			"themeView",
			"settingsview",
			"sliceview",
			"xyzView", // 3d游戏编程
			"webView", // network_security_config D:\workspace\workspace\apk\androidbuild.go
			"iconvgView",
			"svgView",
			"canvasView",
			"popMenu",
			"tooltip",
			"textfield",
			"markdownView", // need sync
			"gomitmproxy",
			"hyperDbg",
			"vstart",
			"explorer",
			"designer",
			"aiChat",
			"encodingTest",
			"Game Control Face",
			"github",
			"ghips",
			"taskManager",
			"gitlab",
			"steam",
			"BuyTomatoes",
			"cc",
			"crypt",
			"Database",
			"datarecovery",
			"hardInfoHook",
			"hardwareIndo",
			"driverTool",
			"environment",
			"erp",
			"fleet",
			"imageConvert",
			"jetbra",
			"jiakaobaodian",
			"ManPiecework",
			"mypan",
			"NetAdapter",
			"netScan",
			"VisualStudiokit",
			"c2go",
			"vnc",
			"todoList",
			"dropFile",
			"darkTheme",
		},
		nil)
}
