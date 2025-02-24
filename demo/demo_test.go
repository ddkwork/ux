package main

import (
	"github.com/ddkwork/golibrary/mylog"
	"os"
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/golibrary/stream"
)

func TestUpdateAppModule(t *testing.T) {
	if !stream.IsDir("../") {
		return
	}
	mylog.Check(os.Chdir("../"))
	session := stream.RunCommand("git log -1 --format=\"%H\"")
	mylog.Check(os.Chdir("demo"))
	id := mylog.Check2(strconv.Unquote(session.Output.String()))
	mylog.Info("id", id)
	stream.RunCommand("go get github.com/ddkwork/ux@" + id)
}

func TestName(t *testing.T) {
	m := new(safemap.M[string, string])
	m.Set("treeTable", "treeTable")
	m.Set("tree", "tree")
	m.Set("table2", "table2")
	m.Set("table", "table")
	m.Set("colorPicker", "colorPicker")
	m.Set("card", "card")
	m.Set("SearchDropDown", "SearchDropDown")

	m.Set("jsonTree", "jsonTree")
	m.Set("svgButton", "svgButton")
	m.Set("codeEditor", "codeEditor")
	m.Set("asmView", "asmView")
	m.Set("logView", "logView")

	m.Set("comBox", "comBox")
	m.Set("splitView", "splitView")
	m.Set("mobile", "mobile")
	m.Set("listView", "listView")
	m.Set("structView", "structView")
	m.Set("terminal", "terminal")
	m.Set("stackView", "stackView")
	m.Set("dockView", "dockView")
	m.Set("gif123", "gif123")
	m.Set("hexEditor", "hexEditor")
	m.Set("imageEditor", "imageEditor")
	m.Set("mediaPlayer", "mediaPlayer")
	m.Set("mind", "mind")
	m.Set("pdfView", "pdfView")
	m.Set("mapView", "mapView")
	m.Set("themeView", "themeView")
	m.Set("settingsview", "settingsview")
	m.Set("sliceview", "sliceview")
	m.Set("xyzView", "xyzView")
	m.Set("webView", "webView")
	m.Set("iconvgView", "iconvgView")
	m.Set("svgView", "svgView")
	m.Set("canvasView", "canvasView")
	m.Set("popMenu", "popMenu")
	m.Set("tooltip", "tooltip")
	m.Set("textfield", "textfield")
	m.Set("markdownView", "markdownView")
	m.Set("gomitmproxy", "gomitmproxy")
	m.Set("hyperDbg", "hyperDbg")
	m.Set("vstart", "vstart")
	m.Set("explorer", "explorer")
	m.Set("designer", "designer")
	m.Set("aiChat", "aiChat")
	m.Set("encodingTest", "encodingTest")
	m.Set("Game Control Face", "Game Control Face")
	m.Set("github", "github")
	m.Set("ghips", "ghips")
	m.Set("taskManager", "taskManager")
	m.Set("gitlab", "gitlab")
	m.Set("steam", "steam")
	m.Set("BuyTomatoes", "BuyTomatoes")
	m.Set("cc", "cc")
	m.Set("crypt", "crypt")
	m.Set("Database", "Database")
	m.Set("datarecovery", "datarecovery")
	m.Set("hardInfoHook", "hardInfoHook")
	m.Set("hardwareIndo", "hardwareIndo")
	m.Set("driverTool", "driverTool")
	m.Set("environment", "environment")
	m.Set("erp", "erp")
	m.Set("fleet", "fleet")
	m.Set("imageConvert", "imageConvert")
	m.Set("jetbra", "jetbra")
	m.Set("jiakaobaodian", "jiakaobaodian")
	m.Set("ManPiecework", "ManPiecework")
	m.Set("mypan", "mypan")
	m.Set("NetAdapter", "NetAdapter")
	m.Set("netScan", "netScan")
	m.Set("VisualStudiokit", "VisualStudiokit")
	m.Set("c2go", "c2go")
	m.Set("vnc", "vnc")
	m.Set("todoList", "todoList")
	m.Set("dropFile", "dropFile")
	m.Set("darkTheme", "darkTheme")

	stream.NewGeneratedFile().SetPackageName("main").EnumTypes("demo", m)
}
