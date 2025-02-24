package main

import (
	"github.com/ddkwork/golibrary/stream"
	"strings"
)

// Code generated by EnumTypesGen - DO NOT EDIT.

type DemoType uint8

const (
	TreeTableType DemoType = iota
	TreeType
	Table2Type
	TableType
	AnimationButtonType
	ContextMenuType
	ColorPickerType
	CardType
	SearchDropDownType
	JsonTreeType
	SvgButtonType
	CodeEditorType
	AsmViewType
	LogViewType
	ComBoxType
	SplitViewType
	MobileType
	ListViewType
	StructViewType
	TerminalType
	StackViewType
	DockViewType
	Gif123Type
	HexEditorType
	ImageEditorType
	MediaPlayerType
	MindType
	PdfViewType
	MapViewType
	ThemeViewType
	SettingsviewType
	SliceviewType
	XyzViewType
	WebViewType
	IconvgViewType
	SvgViewType
	CanvasViewType
	PopMenuType
	TooltipType
	TextfieldType
	MarkdownViewType
	GomitmproxyType
	HyperDbgType
	VstartType
	ExplorerType
	DesignerType
	AiChatType
	EncodingTestType
	GameControlFaceType
	GithubType
	GhipsType
	TaskManagerType
	GitlabType
	SteamType
	BuyTomatoesType
	CcType
	CryptType
	DatabaseType
	DatarecoveryType
	HardInfoHookType
	HardwareIndoType
	DriverToolType
	EnvironmentType
	ErpType
	FleetType
	ImageConvertType
	JetbraType
	JiakaobaodianType
	ManPieceworkType
	MypanType
	NetAdapterType
	NetScanType
	VisualStudiokitType
	C2goType
	VncType
	TodoListType
	DropFileType
	DarkThemeType
)

func (t DemoType) Valid() bool {
	return t >= TreeTableType && t <= DarkThemeType
}

func DemoTypeBy[T stream.Integer](v T) DemoType {
	return DemoType(v)
}

func (t DemoType) AssertBy(name string) DemoType {
	name = strings.TrimSuffix(name, "Type")
	for _, n := range t.EnumTypes() {
		if strings.ToLower(name) == strings.ToLower(n.String()) {
			return n
		}
	}
	panic("InvalidType")
}

func (t DemoType) String() string {
	switch t {
	case TreeTableType:
		return "TreeTable"
	case TreeType:
		return "Tree"
	case Table2Type:
		return "Table2"
	case TableType:
		return "Table"
	case AnimationButtonType:
		return "AnimationButton"
	case ContextMenuType:
		return "ContextMenu"
	case ColorPickerType:
		return "ColorPicker"
	case CardType:
		return "Card"
	case SearchDropDownType:
		return "SearchDropDown"
	case JsonTreeType:
		return "JsonTree"
	case SvgButtonType:
		return "SvgButton"
	case CodeEditorType:
		return "CodeEditor"
	case AsmViewType:
		return "AsmView"
	case LogViewType:
		return "LogView"
	case ComBoxType:
		return "ComBox"
	case SplitViewType:
		return "SplitView"
	case MobileType:
		return "Mobile"
	case ListViewType:
		return "ListView"
	case StructViewType:
		return "StructView"
	case TerminalType:
		return "Terminal"
	case StackViewType:
		return "StackView"
	case DockViewType:
		return "DockView"
	case Gif123Type:
		return "Gif123"
	case HexEditorType:
		return "HexEditor"
	case ImageEditorType:
		return "ImageEditor"
	case MediaPlayerType:
		return "MediaPlayer"
	case MindType:
		return "Mind"
	case PdfViewType:
		return "PdfView"
	case MapViewType:
		return "MapView"
	case ThemeViewType:
		return "ThemeView"
	case SettingsviewType:
		return "Settingsview"
	case SliceviewType:
		return "Sliceview"
	case XyzViewType:
		return "XyzView"
	case WebViewType:
		return "WebView"
	case IconvgViewType:
		return "IconvgView"
	case SvgViewType:
		return "SvgView"
	case CanvasViewType:
		return "CanvasView"
	case PopMenuType:
		return "PopMenu"
	case TooltipType:
		return "Tooltip"
	case TextfieldType:
		return "Textfield"
	case MarkdownViewType:
		return "MarkdownView"
	case GomitmproxyType:
		return "Gomitmproxy"
	case HyperDbgType:
		return "HyperDbg"
	case VstartType:
		return "Vstart"
	case ExplorerType:
		return "Explorer"
	case DesignerType:
		return "Designer"
	case AiChatType:
		return "AiChat"
	case EncodingTestType:
		return "EncodingTest"
	case GameControlFaceType:
		return "GameControlFace"
	case GithubType:
		return "Github"
	case GhipsType:
		return "Ghips"
	case TaskManagerType:
		return "TaskManager"
	case GitlabType:
		return "Gitlab"
	case SteamType:
		return "Steam"
	case BuyTomatoesType:
		return "BuyTomatoes"
	case CcType:
		return "Cc"
	case CryptType:
		return "Crypt"
	case DatabaseType:
		return "Database"
	case DatarecoveryType:
		return "Datarecovery"
	case HardInfoHookType:
		return "HardInfoHook"
	case HardwareIndoType:
		return "HardwareIndo"
	case DriverToolType:
		return "DriverTool"
	case EnvironmentType:
		return "Environment"
	case ErpType:
		return "Erp"
	case FleetType:
		return "Fleet"
	case ImageConvertType:
		return "ImageConvert"
	case JetbraType:
		return "Jetbra"
	case JiakaobaodianType:
		return "Jiakaobaodian"
	case ManPieceworkType:
		return "ManPiecework"
	case MypanType:
		return "Mypan"
	case NetAdapterType:
		return "NetAdapter"
	case NetScanType:
		return "NetScan"
	case VisualStudiokitType:
		return "VisualStudiokit"
	case C2goType:
		return "C2go"
	case VncType:
		return "Vnc"
	case TodoListType:
		return "TodoList"
	case DropFileType:
		return "DropFile"
	case DarkThemeType:
		return "DarkTheme"
	default:
		panic("InvalidType")
	}
}

func (t DemoType) Tooltip() string {
	switch t {
	case TreeTableType:
		return "treeTable"
	case TreeType:
		return "tree"
	case Table2Type:
		return "table2"
	case TableType:
		return "table"
	case AnimationButtonType:
		return "AnimationButton"
	case ContextMenuType:
		return "ContextMenu"
	case ColorPickerType:
		return "colorPicker"
	case CardType:
		return "card"
	case SearchDropDownType:
		return "SearchDropDown"
	case JsonTreeType:
		return "jsonTree"
	case SvgButtonType:
		return "svgButton"
	case CodeEditorType:
		return "codeEditor"
	case AsmViewType:
		return "asmView"
	case LogViewType:
		return "logView"
	case ComBoxType:
		return "comBox"
	case SplitViewType:
		return "splitView"
	case MobileType:
		return "mobile"
	case ListViewType:
		return "listView"
	case StructViewType:
		return "structView"
	case TerminalType:
		return "terminal"
	case StackViewType:
		return "stackView"
	case DockViewType:
		return "dockView"
	case Gif123Type:
		return "gif123"
	case HexEditorType:
		return "hexEditor"
	case ImageEditorType:
		return "imageEditor"
	case MediaPlayerType:
		return "mediaPlayer"
	case MindType:
		return "mind"
	case PdfViewType:
		return "pdfView"
	case MapViewType:
		return "mapView"
	case ThemeViewType:
		return "themeView"
	case SettingsviewType:
		return "settingsview"
	case SliceviewType:
		return "sliceview"
	case XyzViewType:
		return "xyzView"
	case WebViewType:
		return "webView"
	case IconvgViewType:
		return "iconvgView"
	case SvgViewType:
		return "svgView"
	case CanvasViewType:
		return "canvasView"
	case PopMenuType:
		return "popMenu"
	case TooltipType:
		return "tooltip"
	case TextfieldType:
		return "textfield"
	case MarkdownViewType:
		return "markdownView"
	case GomitmproxyType:
		return "gomitmproxy"
	case HyperDbgType:
		return "hyperDbg"
	case VstartType:
		return "vstart"
	case ExplorerType:
		return "explorer"
	case DesignerType:
		return "designer"
	case AiChatType:
		return "aiChat"
	case EncodingTestType:
		return "encodingTest"
	case GameControlFaceType:
		return "Game Control Face"
	case GithubType:
		return "github"
	case GhipsType:
		return "ghips"
	case TaskManagerType:
		return "taskManager"
	case GitlabType:
		return "gitlab"
	case SteamType:
		return "steam"
	case BuyTomatoesType:
		return "BuyTomatoes"
	case CcType:
		return "cc"
	case CryptType:
		return "crypt"
	case DatabaseType:
		return "Database"
	case DatarecoveryType:
		return "datarecovery"
	case HardInfoHookType:
		return "hardInfoHook"
	case HardwareIndoType:
		return "hardwareIndo"
	case DriverToolType:
		return "driverTool"
	case EnvironmentType:
		return "environment"
	case ErpType:
		return "erp"
	case FleetType:
		return "fleet"
	case ImageConvertType:
		return "imageConvert"
	case JetbraType:
		return "jetbra"
	case JiakaobaodianType:
		return "jiakaobaodian"
	case ManPieceworkType:
		return "ManPiecework"
	case MypanType:
		return "mypan"
	case NetAdapterType:
		return "NetAdapter"
	case NetScanType:
		return "netScan"
	case VisualStudiokitType:
		return "VisualStudiokit"
	case C2goType:
		return "c2go"
	case VncType:
		return "vnc"
	case TodoListType:
		return "todoList"
	case DropFileType:
		return "dropFile"
	case DarkThemeType:
		return "darkTheme"
	default:
		panic("InvalidType")
	}
}

func (t DemoType) Names() []string {
	return []string{
		"TreeTable",
		"Tree",
		"Table2",
		"Table",
		"AnimationButton",
		"ContextMenu",
		"ColorPicker",
		"Card",
		"SearchDropDown",
		"JsonTree",
		"SvgButton",
		"CodeEditor",
		"AsmView",
		"LogView",
		"ComBox",
		"SplitView",
		"Mobile",
		"ListView",
		"StructView",
		"Terminal",
		"StackView",
		"DockView",
		"Gif123",
		"HexEditor",
		"ImageEditor",
		"MediaPlayer",
		"Mind",
		"PdfView",
		"MapView",
		"ThemeView",
		"Settingsview",
		"Sliceview",
		"XyzView",
		"WebView",
		"IconvgView",
		"SvgView",
		"CanvasView",
		"PopMenu",
		"Tooltip",
		"Textfield",
		"MarkdownView",
		"Gomitmproxy",
		"HyperDbg",
		"Vstart",
		"Explorer",
		"Designer",
		"AiChat",
		"EncodingTest",
		"GameControlFace",
		"Github",
		"Ghips",
		"TaskManager",
		"Gitlab",
		"Steam",
		"BuyTomatoes",
		"Cc",
		"Crypt",
		"Database",
		"Datarecovery",
		"HardInfoHook",
		"HardwareIndo",
		"DriverTool",
		"Environment",
		"Erp",
		"Fleet",
		"ImageConvert",
		"Jetbra",
		"Jiakaobaodian",
		"ManPiecework",
		"Mypan",
		"NetAdapter",
		"NetScan",
		"VisualStudiokit",
		"C2go",
		"Vnc",
		"TodoList",
		"DropFile",
		"DarkTheme",
	}
}

func (t DemoType) EnumTypes() []DemoType {
	return []DemoType{
		TreeTableType,
		TreeType,
		Table2Type,
		TableType,
		AnimationButtonType,
		ContextMenuType,
		ColorPickerType,
		CardType,
		SearchDropDownType,
		JsonTreeType,
		SvgButtonType,
		CodeEditorType,
		AsmViewType,
		LogViewType,
		ComBoxType,
		SplitViewType,
		MobileType,
		ListViewType,
		StructViewType,
		TerminalType,
		StackViewType,
		DockViewType,
		Gif123Type,
		HexEditorType,
		ImageEditorType,
		MediaPlayerType,
		MindType,
		PdfViewType,
		MapViewType,
		ThemeViewType,
		SettingsviewType,
		SliceviewType,
		XyzViewType,
		WebViewType,
		IconvgViewType,
		SvgViewType,
		CanvasViewType,
		PopMenuType,
		TooltipType,
		TextfieldType,
		MarkdownViewType,
		GomitmproxyType,
		HyperDbgType,
		VstartType,
		ExplorerType,
		DesignerType,
		AiChatType,
		EncodingTestType,
		GameControlFaceType,
		GithubType,
		GhipsType,
		TaskManagerType,
		GitlabType,
		SteamType,
		BuyTomatoesType,
		CcType,
		CryptType,
		DatabaseType,
		DatarecoveryType,
		HardInfoHookType,
		HardwareIndoType,
		DriverToolType,
		EnvironmentType,
		ErpType,
		FleetType,
		ImageConvertType,
		JetbraType,
		JiakaobaodianType,
		ManPieceworkType,
		MypanType,
		NetAdapterType,
		NetScanType,
		VisualStudiokitType,
		C2goType,
		VncType,
		TodoListType,
		DropFileType,
		DarkThemeType,
	}
}

func (t DemoType) SvgFileName() string {
	return t.String() + ".svg"
}
