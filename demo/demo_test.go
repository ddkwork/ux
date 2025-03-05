package main

import (
	"os"
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/mylog"

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
	// go get -u -x all
}

func TestName(t *testing.T) {
	m := safemap.NewOrdered[string, string](func(yield func(string, string) bool) {
		yield("treeTable", "treeTable")
		yield("tree", "tree")
		yield("table2", "table2")
		yield("table", "table")
		yield("AnimationButton", "AnimationButton")
		yield("ContextMenu", "ContextMenu")
		yield("colorPicker", "colorPicker")
		yield("screenshot", "screenshot")
		yield("card", "card")
		yield("dataPicker", "dataPicker")
		yield("resizer", "resizer")
		yield("SearchDropDown", "SearchDropDown")

		yield("jsonTree", "jsonTree")
		yield("svgButton", "svgButton")
		yield("codeEditor", "codeEditor")
		yield("asmView", "asmView")
		yield("logView", "logView")

		yield("comBox", "comBox")
		yield("splitView", "splitView")
		yield("mobile", "mobile")
		yield("listView", "listView")
		yield("structView", "structView")
		yield("terminal", "terminal")
		yield("stackView", "stackView")
		yield("dockView", "dockView")
		yield("gif123", "gif123")
		yield("hexEditor", "hexEditor")
		yield("imageEditor", "imageEditor")
		yield("mediaPlayer", "mediaPlayer")
		yield("mind", "mind")
		yield("pdfView", "pdfView")
		yield("mapView", "mapView")
		yield("themeView", "themeView")
		yield("settingsview", "settingsview")
		yield("sliceview", "sliceview")
		yield("xyzView", "xyzView")
		yield("webView", "webView")
		yield("iconvgView", "iconvgView")
		yield("svgView", "svgView")
		yield("canvasView", "canvasView")
		yield("popMenu", "popMenu")
		yield("tooltip", "tooltip")
		yield("textfield", "textfield")
		yield("markdownView", "markdownView")
		yield("gomitmproxy", "gomitmproxy")
		yield("hyperDbg", "hyperDbg")
		yield("vstart", "vstart")
		yield("explorer", "explorer")
		yield("designer", "designer")
		yield("aiChat", "aiChat")
		yield("encodingTest", "encodingTest")
		yield("Game Control Face", "Game Control Face")
		yield("github", "github")
		yield("ghips", "ghips")
		yield("taskManager", "taskManager")
		yield("gitlab", "gitlab")
		yield("steam", "steam")
		yield("BuyTomatoes", "BuyTomatoes")
		yield("cc", "cc")
		yield("crypt", "crypt")
		yield("Database", "Database")
		yield("datarecovery", "datarecovery")
		yield("hardInfoHook", "hardInfoHook")
		yield("hardwareIndo", "hardwareIndo")
		yield("driverTool", "driverTool")
		yield("environment", "environment")
		yield("erp", "erp")
		yield("fleet", "fleet")
		yield("imageConvert", "imageConvert")
		yield("jetbra", "jetbra")
		yield("jiakaobaodian", "jiakaobaodian")
		yield("ManPiecework", "ManPiecework")
		yield("mypan", "mypan")
		yield("NetAdapter", "NetAdapter")
		yield("netScan", "netScan")
		yield("VisualStudiokit", "VisualStudiokit")
		yield("c2go", "c2go")
		yield("vnc", "vnc")
		yield("todoList", "todoList")
		yield("dropFile", "dropFile")
		yield("darkTheme", "darkTheme")
	})
	stream.NewGeneratedFile().SetPackageName("main").EnumTypes("demo", m)
}
