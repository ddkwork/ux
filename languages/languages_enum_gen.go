package languages

import "github.com/ddkwork/golibrary/stream"

type LanguagesKind byte

const (
	AbapKind LanguagesKind = iota
	AbnfKind
	ActionScriptKind
	ActionScript3Kind
	AdaKind
	AgdaKind
	AlKind
	AlloyKind
	Angular2Kind
	AntlrKind
	ApacheConfKind
	AplKind
	AppleScriptKind
	ArangoDbAqlKind
	ArduinoKind
	ArmAsmKind
	AutoHotkeyKind
	AutoItKind
	AwkKind
	BallerinaKind
	BashKind
	BashSessionKind
	BatchfileKind
	BibTeXKind
	BicepKind
	BlitzBasicKind
	BnfKind
	BqnKind
	BrainfuckKind
	CSharpKind
	CppKind
	CKind
	CapnProtoKind
	CassandraCqlKind
	CeylonKind
	CfEngine3Kind
	CfstatementKind
	ChaiScriptKind
	ChapelKind
	CheetahKind
	ClojureKind
	CMakeKind
	CobolKind
	CoffeeScriptKind
	CommonLispKind
	CoqKind
	CrystalKind
	CssKind
	CueKind
	CythonKind
	DKind
	DartKind
	DaxKind
	DesktopFileKind
	DiffKind
	DjangoJinjaKind
	DnsKind
	DockerKind
	DtdKind
	DylanKind
	EbnfKind
	ElixirKind
	ElmKind
	EmacsLispKind
	ErlangKind
	FactorKind
	FennelKind
	FishKind
	ForthKind
	FortranKind
	FortranFixedKind
	FSharpKind
	GasKind
	GdScriptKind
	GdScript3Kind
	GherkinKind
	GlslKind
	GnuplotKind
	GoTemplateKind
	GraphQlKind
	GroffKind
	GroovyKind
	HandlebarsKind
	HareKind
	HaskellKind
	HclKind
	HexdumpKind
	HlbKind
	HlslKind
	HolyCKind
	HtmlKind
	HyKind
	IdrisKind
	IgorKind
	IniKind
	IoKind
	IsCdhcpdKind
	JKind
	JavaKind
	JavaScriptKind
	JsonKind
	JuliaKind
	JungleKind
	KotlinKind
	LighttpdConfigurationFileKind
	LlvmKind
	LuaKind
	MakefileKind
	MakoKind
	MasonKind
	MaterializeSqlDialectKind
	MathematicaKind
	MatlabKind
	McfunctionKind
	MesonKind
	MetalKind
	MiniZincKind
	MlirKind
	Modula2Kind
	MonkeyCKind
	MorrowindScriptKind
	MyghtyKind
	MySqlKind
	NasmKind
	NaturalKind
	NdisasmKind
	NewspeakKind
	NginxConfigurationFileKind
	NimKind
	NixKind
	ObjectiveCKind
	ObjectPascalKind
	OCamlKind
	OctaveKind
	OdinKind
	OnesEnterpriseKind
	OpenEdgeAblKind
	OpenScadKind
	OrgModeKind
	PacmanConfKind
	PerlKind
	PhpKind
	PigKind
	PkgConfigKind
	PlPgSqlKind
	PlaintextKind
	PlutusCoreKind
	PonyKind
	PostgreSqlSqlDialectKind
	PostScriptKind
	PovRayKind
	PowerQueryKind
	PowerShellKind
	PrologKind
	PromelaKind
	PromQlKind
	PropertiesKind
	ProtocolBufferKind
	PrqlKind
	PslKind
	PuppetKind
	PythonKind
	Python2Kind
	QBasicKind
	QmlKind
	RKind
	RacketKind
	RagelKind
	ReactKind
	ReasonMlKind
	RegKind
	RegoKind
	RexxKind
	RpmSpecKind
	RubyKind
	RustKind
	SasKind
	SassKind
	ScalaKind
	SchemeKind
	ScilabKind
	ScssKind
	SedKind
	SieveKind
	SmaliKind
	SmalltalkKind
	SmartyKind
	SnobolKind
	SolidityKind
	SourcePawnKind
	SparqlKind
	SqlKind
	SquidConfKind
	StandardMlKind
	StasKind
	StylusKind
	SwiftKind
	SystemdKind
	SystemverilogKind
	TableGenKind
	TalKind
	TasmKind
	TclKind
	TcshKind
	TermcapKind
	TerminfoKind
	TerraformKind
	TeXKind
	ThriftKind
	TomlKind
	TradingViewKind
	TransactSqlKind
	TuringKind
	TurtleKind
	TwigKind
	TypeScriptKind
	TypoScriptKind
	TypoScriptCssDataKind
	TypoScriptHtmlDataKind
	UcodeKind
	VKind
	VShellKind
	ValaKind
	VbNetKind
	VerilogKind
	VhdlKind
	VhsKind
	VimLKind
	VueKind
	WdteKind
	WebGpuShadingLanguageKind
	WhileyKind
	XmlKind
	XorgKind
	YamlKind
	YangKind
	Z80AssemblyKind
	ZedKind
	ZigKind
	CaddyfileKind
	CaddyfileDirectivesKind
	GenshiTextKind
	GenshiHtmlKind
	GenshiKind
	GoHtmlTemplateKind
	GoTextTemplateKind
	GoKind
	HaxeKind
	HttpKind
	MarkdownKind
	PhtmlKind
	RakuKind
	ReStructuredTextKind
	SvelteKind
	InvalidLanguagesKind
)

func ConvertInteger2LanguagesKind[T stream.Integer](v T) LanguagesKind {
	return LanguagesKind(v)
}

func (k LanguagesKind) AssertKind(kinds string) LanguagesKind {
	for _, kind := range k.Kinds() {
		if kinds == kind.String() {
			return kind
		}
	}
	return InvalidLanguagesKind
}

func (k LanguagesKind) String() string {
	switch k {
	case AbapKind:
		return "Abap"
	case AbnfKind:
		return "Abnf"
	case ActionScriptKind:
		return "ActionScript"
	case ActionScript3Kind:
		return "ActionScript3"
	case AdaKind:
		return "Ada"
	case AgdaKind:
		return "Agda"
	case AlKind:
		return "Al"
	case AlloyKind:
		return "Alloy"
	case Angular2Kind:
		return "Angular2"
	case AntlrKind:
		return "Antlr"
	case ApacheConfKind:
		return "ApacheConf"
	case AplKind:
		return "Apl"
	case AppleScriptKind:
		return "AppleScript"
	case ArangoDbAqlKind:
		return "ArangoDbAql"
	case ArduinoKind:
		return "Arduino"
	case ArmAsmKind:
		return "ArmAsm"
	case AutoHotkeyKind:
		return "AutoHotkey"
	case AutoItKind:
		return "AutoIt"
	case AwkKind:
		return "Awk"
	case BallerinaKind:
		return "Ballerina"
	case BashKind:
		return "Bash"
	case BashSessionKind:
		return "BashSession"
	case BatchfileKind:
		return "Batchfile"
	case BibTeXKind:
		return "BibTeX"
	case BicepKind:
		return "Bicep"
	case BlitzBasicKind:
		return "BlitzBasic"
	case BnfKind:
		return "Bnf"
	case BqnKind:
		return "Bqn"
	case BrainfuckKind:
		return "Brainfuck"
	case CSharpKind:
		return "CSharp"
	case CppKind:
		return "Cpp"
	case CKind:
		return "C"
	case CapnProtoKind:
		return "CapnProto"
	case CassandraCqlKind:
		return "CassandraCql"
	case CeylonKind:
		return "Ceylon"
	case CfEngine3Kind:
		return "CfEngine3"
	case CfstatementKind:
		return "Cfstatement"
	case ChaiScriptKind:
		return "ChaiScript"
	case ChapelKind:
		return "Chapel"
	case CheetahKind:
		return "Cheetah"
	case ClojureKind:
		return "Clojure"
	case CMakeKind:
		return "CMake"
	case CobolKind:
		return "Cobol"
	case CoffeeScriptKind:
		return "CoffeeScript"
	case CommonLispKind:
		return "CommonLisp"
	case CoqKind:
		return "Coq"
	case CrystalKind:
		return "Crystal"
	case CssKind:
		return "Css"
	case CueKind:
		return "Cue"
	case CythonKind:
		return "Cython"
	case DKind:
		return "D"
	case DartKind:
		return "Dart"
	case DaxKind:
		return "Dax"
	case DesktopFileKind:
		return "DesktopFile"
	case DiffKind:
		return "Diff"
	case DjangoJinjaKind:
		return "DjangoJinja"
	case DnsKind:
		return "Dns"
	case DockerKind:
		return "Docker"
	case DtdKind:
		return "Dtd"
	case DylanKind:
		return "Dylan"
	case EbnfKind:
		return "Ebnf"
	case ElixirKind:
		return "Elixir"
	case ElmKind:
		return "Elm"
	case EmacsLispKind:
		return "EmacsLisp"
	case ErlangKind:
		return "Erlang"
	case FactorKind:
		return "Factor"
	case FennelKind:
		return "Fennel"
	case FishKind:
		return "Fish"
	case ForthKind:
		return "Forth"
	case FortranKind:
		return "Fortran"
	case FortranFixedKind:
		return "FortranFixed"
	case FSharpKind:
		return "FSharp"
	case GasKind:
		return "Gas"
	case GdScriptKind:
		return "GdScript"
	case GdScript3Kind:
		return "GdScript3"
	case GherkinKind:
		return "Gherkin"
	case GlslKind:
		return "Glsl"
	case GnuplotKind:
		return "Gnuplot"
	case GoTemplateKind:
		return "GoTemplate"
	case GraphQlKind:
		return "GraphQl"
	case GroffKind:
		return "Groff"
	case GroovyKind:
		return "Groovy"
	case HandlebarsKind:
		return "Handlebars"
	case HareKind:
		return "Hare"
	case HaskellKind:
		return "Haskell"
	case HclKind:
		return "Hcl"
	case HexdumpKind:
		return "Hexdump"
	case HlbKind:
		return "Hlb"
	case HlslKind:
		return "Hlsl"
	case HolyCKind:
		return "HolyC"
	case HtmlKind:
		return "Html"
	case HyKind:
		return "Hy"
	case IdrisKind:
		return "Idris"
	case IgorKind:
		return "Igor"
	case IniKind:
		return "Ini"
	case IoKind:
		return "Io"
	case IsCdhcpdKind:
		return "IsCdhcpd"
	case JKind:
		return "J"
	case JavaKind:
		return "Java"
	case JavaScriptKind:
		return "JavaScript"
	case JsonKind:
		return "Json"
	case JuliaKind:
		return "Julia"
	case JungleKind:
		return "Jungle"
	case KotlinKind:
		return "Kotlin"
	case LighttpdConfigurationFileKind:
		return "LighttpdConfigurationFile"
	case LlvmKind:
		return "Llvm"
	case LuaKind:
		return "Lua"
	case MakefileKind:
		return "Makefile"
	case MakoKind:
		return "Mako"
	case MasonKind:
		return "Mason"
	case MaterializeSqlDialectKind:
		return "MaterializeSqlDialect"
	case MathematicaKind:
		return "Mathematica"
	case MatlabKind:
		return "Matlab"
	case McfunctionKind:
		return "Mcfunction"
	case MesonKind:
		return "Meson"
	case MetalKind:
		return "Metal"
	case MiniZincKind:
		return "MiniZinc"
	case MlirKind:
		return "Mlir"
	case Modula2Kind:
		return "Modula2"
	case MonkeyCKind:
		return "MonkeyC"
	case MorrowindScriptKind:
		return "MorrowindScript"
	case MyghtyKind:
		return "Myghty"
	case MySqlKind:
		return "MySql"
	case NasmKind:
		return "Nasm"
	case NaturalKind:
		return "Natural"
	case NdisasmKind:
		return "Ndisasm"
	case NewspeakKind:
		return "Newspeak"
	case NginxConfigurationFileKind:
		return "NginxConfigurationFile"
	case NimKind:
		return "Nim"
	case NixKind:
		return "Nix"
	case ObjectiveCKind:
		return "ObjectiveC"
	case ObjectPascalKind:
		return "ObjectPascal"
	case OCamlKind:
		return "OCaml"
	case OctaveKind:
		return "Octave"
	case OdinKind:
		return "Odin"
	case OnesEnterpriseKind:
		return "OnesEnterprise"
	case OpenEdgeAblKind:
		return "OpenEdgeAbl"
	case OpenScadKind:
		return "OpenScad"
	case OrgModeKind:
		return "OrgMode"
	case PacmanConfKind:
		return "PacmanConf"
	case PerlKind:
		return "Perl"
	case PhpKind:
		return "Php"
	case PigKind:
		return "Pig"
	case PkgConfigKind:
		return "PkgConfig"
	case PlPgSqlKind:
		return "PlPgSql"
	case PlaintextKind:
		return "Plaintext"
	case PlutusCoreKind:
		return "PlutusCore"
	case PonyKind:
		return "Pony"
	case PostgreSqlSqlDialectKind:
		return "PostgreSqlSqlDialect"
	case PostScriptKind:
		return "PostScript"
	case PovRayKind:
		return "PovRay"
	case PowerQueryKind:
		return "PowerQuery"
	case PowerShellKind:
		return "PowerShell"
	case PrologKind:
		return "Prolog"
	case PromelaKind:
		return "Promela"
	case PromQlKind:
		return "PromQl"
	case PropertiesKind:
		return "Properties"
	case ProtocolBufferKind:
		return "ProtocolBuffer"
	case PrqlKind:
		return "Prql"
	case PslKind:
		return "Psl"
	case PuppetKind:
		return "Puppet"
	case PythonKind:
		return "Python"
	case Python2Kind:
		return "Python2"
	case QBasicKind:
		return "QBasic"
	case QmlKind:
		return "Qml"
	case RKind:
		return "R"
	case RacketKind:
		return "Racket"
	case RagelKind:
		return "Ragel"
	case ReactKind:
		return "React"
	case ReasonMlKind:
		return "ReasonMl"
	case RegKind:
		return "Reg"
	case RegoKind:
		return "Rego"
	case RexxKind:
		return "Rexx"
	case RpmSpecKind:
		return "RpmSpec"
	case RubyKind:
		return "Ruby"
	case RustKind:
		return "Rust"
	case SasKind:
		return "Sas"
	case SassKind:
		return "Sass"
	case ScalaKind:
		return "Scala"
	case SchemeKind:
		return "Scheme"
	case ScilabKind:
		return "Scilab"
	case ScssKind:
		return "Scss"
	case SedKind:
		return "Sed"
	case SieveKind:
		return "Sieve"
	case SmaliKind:
		return "Smali"
	case SmalltalkKind:
		return "Smalltalk"
	case SmartyKind:
		return "Smarty"
	case SnobolKind:
		return "Snobol"
	case SolidityKind:
		return "Solidity"
	case SourcePawnKind:
		return "SourcePawn"
	case SparqlKind:
		return "Sparql"
	case SqlKind:
		return "Sql"
	case SquidConfKind:
		return "SquidConf"
	case StandardMlKind:
		return "StandardMl"
	case StasKind:
		return "Stas"
	case StylusKind:
		return "Stylus"
	case SwiftKind:
		return "Swift"
	case SystemdKind:
		return "Systemd"
	case SystemverilogKind:
		return "Systemverilog"
	case TableGenKind:
		return "TableGen"
	case TalKind:
		return "Tal"
	case TasmKind:
		return "Tasm"
	case TclKind:
		return "Tcl"
	case TcshKind:
		return "Tcsh"
	case TermcapKind:
		return "Termcap"
	case TerminfoKind:
		return "Terminfo"
	case TerraformKind:
		return "Terraform"
	case TeXKind:
		return "TeX"
	case ThriftKind:
		return "Thrift"
	case TomlKind:
		return "Toml"
	case TradingViewKind:
		return "TradingView"
	case TransactSqlKind:
		return "TransactSql"
	case TuringKind:
		return "Turing"
	case TurtleKind:
		return "Turtle"
	case TwigKind:
		return "Twig"
	case TypeScriptKind:
		return "TypeScript"
	case TypoScriptKind:
		return "TypoScript"
	case TypoScriptCssDataKind:
		return "TypoScriptCssData"
	case TypoScriptHtmlDataKind:
		return "TypoScriptHtmlData"
	case UcodeKind:
		return "Ucode"
	case VKind:
		return "V"
	case VShellKind:
		return "VShell"
	case ValaKind:
		return "Vala"
	case VbNetKind:
		return "VbNet"
	case VerilogKind:
		return "Verilog"
	case VhdlKind:
		return "Vhdl"
	case VhsKind:
		return "Vhs"
	case VimLKind:
		return "VimL"
	case VueKind:
		return "Vue"
	case WdteKind:
		return "Wdte"
	case WebGpuShadingLanguageKind:
		return "WebGpuShadingLanguage"
	case WhileyKind:
		return "Whiley"
	case XmlKind:
		return "Xml"
	case XorgKind:
		return "Xorg"
	case YamlKind:
		return "Yaml"
	case YangKind:
		return "Yang"
	case Z80AssemblyKind:
		return "Z80Assembly"
	case ZedKind:
		return "Zed"
	case ZigKind:
		return "Zig"
	case CaddyfileKind:
		return "Caddyfile"
	case CaddyfileDirectivesKind:
		return "CaddyfileDirectives"
	case GenshiTextKind:
		return "GenshiText"
	case GenshiHtmlKind:
		return "GenshiHtml"
	case GenshiKind:
		return "Genshi"
	case GoHtmlTemplateKind:
		return "GoHtmlTemplate"
	case GoTextTemplateKind:
		return "GoTextTemplate"
	case GoKind:
		return "Go"
	case HaxeKind:
		return "Haxe"
	case HttpKind:
		return "Http"
	case MarkdownKind:
		return "Markdown"
	case PhtmlKind:
		return "Phtml"
	case RakuKind:
		return "Raku"
	case ReStructuredTextKind:
		return "ReStructuredText"
	case SvelteKind:
		return "Svelte"
	default:
		return "InvalidLanguagesKind"
	}
}

func (k LanguagesKind) Tooltip() string {
	switch k {
	case AbapKind:
		return "ABAP  [*.abap *.ABAP]  [text/x-abap]"
	case AbnfKind:
		return "ABNF  [*.abnf]  [text/x-abnf]"
	case ActionScriptKind:
		return "ActionScript  [*.as]  [application/x-actionscript text/x-actionscript text/actionscript]"
	case ActionScript3Kind:
		return "ActionScript 3  [*.as]  [application/x-actionscript3 text/x-actionscript3 text/actionscript3]"
	case AdaKind:
		return "Ada  [*.adb *.ads *.ada]  [text/x-ada]"
	case AgdaKind:
		return "Agda  [*.agda]  [text/x-agda]"
	case AlKind:
		return "AL  [*.al *.dal]  [text/x-al]"
	case AlloyKind:
		return "Alloy  [*.als]  [text/x-alloy]"
	case Angular2Kind:
		return "Angular2  []  []"
	case AntlrKind:
		return "ANTLR  []  []"
	case ApacheConfKind:
		return "ApacheConf  [.htaccess apache.conf apache2.conf]  [text/x-apacheconf]"
	case AplKind:
		return "APL  [*.apl]  []"
	case AppleScriptKind:
		return "AppleScript  [*.applescript]  []"
	case ArangoDbAqlKind:
		return "ArangoDB AQL  [*.aql]  [text/x-aql]"
	case ArduinoKind:
		return "Arduino  [*.ino]  [text/x-arduino]"
	case ArmAsmKind:
		return "ArmAsm  [*.s *.S]  [text/x-armasm text/x-asm]"
	case AutoHotkeyKind:
		return "AutoHotkey  [*.ahk *.ahkl]  [text/x-autohotkey]"
	case AutoItKind:
		return "AutoIt  [*.au3]  [text/x-autoit]"
	case AwkKind:
		return "Awk  [*.awk]  [application/x-awk]"
	case BallerinaKind:
		return "Ballerina  [*.bal]  [text/x-ballerina]"
	case BashKind:
		return "Bash  [*.sh *.ksh *.bash *.ebuild *.eclass .env *.env *.exheres-0 *.exlib *.zsh *.zshrc .bashrc bashrc .bash_* bash_* zshrc .zshrc PKGBUILD]  [application/x-sh application/x-shellscript]"
	case BashSessionKind:
		return "Bash Session  [*.sh-session]  [text/x-sh]"
	case BatchfileKind:
		return "Batchfile  [*.bat *.cmd]  [application/x-dos-batch]"
	case BibTeXKind:
		return "BibTeX  [*.bib]  [text/x-bibtex]"
	case BicepKind:
		return "Bicep  [*.bicep]  []"
	case BlitzBasicKind:
		return "BlitzBasic  [*.bb *.decls]  [text/x-bb]"
	case BnfKind:
		return "BNF  [*.bnf]  [text/x-bnf]"
	case BqnKind:
		return "BQN  [*.bqn]  []"
	case BrainfuckKind:
		return "Brainfuck  [*.bf *.b]  [application/x-brainfuck]"
	case CSharpKind:
		return "C#  [*.cs]  [text/x-csharp]"
	case CppKind:
		return "C++  [*.cpp *.hpp *.c++ *.h++ *.cc *.hh *.cxx *.hxx *.C *.H *.cp *.CPP *.tpp]  [text/x-c++hdr text/x-c++src]"
	case CKind:
		return "C  [*.c *.h *.idc *.x[bp]m]  [text/x-chdr text/x-csrc image/x-xbitmap image/x-xpixmap]"
	case CapnProtoKind:
		return "Cap'n Proto  [*.capnp]  []"
	case CassandraCqlKind:
		return "Cassandra CQL  [*.cql]  [text/x-cql]"
	case CeylonKind:
		return "Ceylon  [*.ceylon]  [text/x-ceylon]"
	case CfEngine3Kind:
		return "CFEngine3  [*.cf]  []"
	case CfstatementKind:
		return "cfstatement  []  []"
	case ChaiScriptKind:
		return "ChaiScript  [*.chai]  [text/x-chaiscript application/x-chaiscript]"
	case ChapelKind:
		return "Chapel  [*.chpl]  []"
	case CheetahKind:
		return "Cheetah  [*.tmpl *.spt]  [application/x-cheetah application/x-spitfire]"
	case ClojureKind:
		return "Clojure  [*.clj *.edn]  [text/x-clojure application/x-clojure application/edn]"
	case CMakeKind:
		return "CMake  [*.cmake CMakeLists.txt]  [text/x-cmake]"
	case CobolKind:
		return "COBOL  [*.cob *.COB *.cpy *.CPY]  [text/x-cobol]"
	case CoffeeScriptKind:
		return "CoffeeScript  [*.coffee]  [text/coffeescript]"
	case CommonLispKind:
		return "Common Lisp  [*.cl *.lisp]  [text/x-common-lisp]"
	case CoqKind:
		return "Coq  [*.v]  [text/x-coq]"
	case CrystalKind:
		return "Crystal  [*.cr]  [text/x-crystal]"
	case CssKind:
		return "CSS  [*.css]  [text/css]"
	case CueKind:
		return "CUE  [*.cue]  [text/x-cue]"
	case CythonKind:
		return "Cython  [*.pyx *.pxd *.pxi]  [text/x-cython application/x-cython]"
	case DKind:
		return "D  [*.d *.di]  [text/x-d]"
	case DartKind:
		return "Dart  [*.dart]  [text/x-dart]"
	case DaxKind:
		return "Dax  [*.dax]  []"
	case DesktopFileKind:
		return "Desktop file  [*.desktop]  [application/x-desktop]"
	case DiffKind:
		return "Diff  [*.diff *.patch]  [text/x-diff text/x-patch]"
	case DjangoJinjaKind:
		return "Django/Jinja  []  [application/x-django-templating application/x-jinja]"
	case DnsKind:
		return "dns  [*.zone]  [text/dns]"
	case DockerKind:
		return "Docker  [Dockerfile Dockerfile.* *.Dockerfile *.docker]  [text/x-dockerfile-config]"
	case DtdKind:
		return "DTD  [*.dtd]  [application/xml-dtd]"
	case DylanKind:
		return "Dylan  [*.dylan *.dyl *.intr]  [text/x-dylan]"
	case EbnfKind:
		return "EBNF  [*.ebnf]  [text/x-ebnf]"
	case ElixirKind:
		return "Elixir  [*.ex *.eex *.exs]  [text/x-elixir]"
	case ElmKind:
		return "Elm  [*.elm]  [text/x-elm]"
	case EmacsLispKind:
		return "EmacsLisp  [*.el]  [text/x-elisp application/x-elisp]"
	case ErlangKind:
		return "Erlang  [*.erl *.hrl *.es *.escript]  [text/x-erlang]"
	case FactorKind:
		return "Factor  [*.factor]  [text/x-factor]"
	case FennelKind:
		return "Fennel  [*.fennel]  [text/x-fennel application/x-fennel]"
	case FishKind:
		return "Fish  [*.fish *.load]  [application/x-fish]"
	case ForthKind:
		return "Forth  [*.frt *.fth *.fs]  [application/x-forth]"
	case FortranKind:
		return "Fortran  [*.f03 *.f90 *.f95 *.F03 *.F90 *.F95]  [text/x-fortran]"
	case FortranFixedKind:
		return "FortranFixed  [*.f *.F]  [text/x-fortran]"
	case FSharpKind:
		return "FSharp  [*.fs *.fsi]  [text/x-fsharp]"
	case GasKind:
		return "GAS  [*.s *.S]  [text/x-gas]"
	case GdScriptKind:
		return "GDScript  [*.gd]  [text/x-gdscript application/x-gdscript]"
	case GdScript3Kind:
		return "GDScript3  [*.gd]  [text/x-gdscript application/x-gdscript]"
	case GherkinKind:
		return "Gherkin  [*.feature *.FEATURE]  [text/x-gherkin]"
	case GlslKind:
		return "GLSL  [*.vert *.frag *.geo]  [text/x-glslsrc]"
	case GnuplotKind:
		return "Gnuplot  [*.plot *.plt]  [text/x-gnuplot]"
	case GoTemplateKind:
		return "Go Template  [*.gotmpl *.go.tmpl]  []"
	case GraphQlKind:
		return "GraphQL  [*.graphql *.graphqls]  []"
	case GroffKind:
		return "Groff  [*.[1-9] *.1p *.3pm *.man]  [application/x-troff text/troff]"
	case GroovyKind:
		return "Groovy  [*.groovy *.gradle]  [text/x-groovy]"
	case HandlebarsKind:
		return "Handlebars  [*.handlebars *.hbs]  []"
	case HareKind:
		return "Hare  [*.ha]  [text/x-hare]"
	case HaskellKind:
		return "Haskell  [*.hs]  [text/x-haskell]"
	case HclKind:
		return "HCL  [*.hcl]  [application/x-hcl]"
	case HexdumpKind:
		return "Hexdump  []  []"
	case HlbKind:
		return "HLB  [*.hlb]  []"
	case HlslKind:
		return "HLSL  [*.hlsl *.hlsli *.cginc *.fx *.fxh]  [text/x-hlsl]"
	case HolyCKind:
		return "HolyC  [*.HC *.hc *.HH *.hh *.hc.z *.HC.Z]  [text/x-chdr text/x-csrc image/x-xbitmap image/x-xpixmap]"
	case HtmlKind:
		return "HTML  [*.html *.htm *.xhtml *.xslt]  [text/html application/xhtml+xml]"
	case HyKind:
		return "Hy  [*.hy]  [text/x-hy application/x-hy]"
	case IdrisKind:
		return "Idris  [*.idr]  [text/x-idris]"
	case IgorKind:
		return "Igor  [*.ipf]  [text/ipf]"
	case IniKind:
		return "INI  [*.ini *.cfg *.inf *.service *.socket .gitconfig .editorconfig pylintrc .pylintrc]  [text/x-ini text/inf]"
	case IoKind:
		return "Io  [*.io]  [text/x-iosrc]"
	case IsCdhcpdKind:
		return "ISCdhcpd  [dhcpd.conf]  []"
	case JKind:
		return "J  [*.ijs]  [text/x-j]"
	case JavaKind:
		return "Java  [*.java]  [text/x-java]"
	case JavaScriptKind:
		return "JavaScript  [*.js *.jsm *.mjs *.cjs]  [application/javascript application/x-javascript text/x-javascript text/javascript]"
	case JsonKind:
		return "JSON  [*.json *.avsc]  [application/json]"
	case JuliaKind:
		return "Julia  [*.jl]  [text/x-julia application/x-julia]"
	case JungleKind:
		return "Jungle  [*.jungle]  [text/x-jungle]"
	case KotlinKind:
		return "Kotlin  [*.kt]  [text/x-kotlin]"
	case LighttpdConfigurationFileKind:
		return "Lighttpd configuration file  []  [text/x-lighttpd-conf]"
	case LlvmKind:
		return "LLVM  [*.ll]  [text/x-llvm]"
	case LuaKind:
		return "Lua  [*.lua *.wlua]  [text/x-lua application/x-lua]"
	case MakefileKind:
		return "Makefile  [*.mak *.mk Makefile makefile Makefile.* GNUmakefile BSDmakefile Justfile justfile .justfile]  [text/x-makefile]"
	case MakoKind:
		return "Mako  [*.mao]  [application/x-mako]"
	case MasonKind:
		return "Mason  [*.m *.mhtml *.mc *.mi autohandler dhandler]  [application/x-mason]"
	case MaterializeSqlDialectKind:
		return "Materialize SQL dialect  []  [text/x-materializesql]"
	case MathematicaKind:
		return "Mathematica  [*.cdf *.m *.ma *.mt *.mx *.nb *.nbp *.wl]  [application/mathematica application/vnd.wolfram.mathematica application/vnd.wolfram.mathematica.package application/vnd.wolfram.cdf]"
	case MatlabKind:
		return "Matlab  [*.m]  [text/matlab]"
	case McfunctionKind:
		return "mcfunction  [*.mcfunction]  []"
	case MesonKind:
		return "Meson  [meson.build meson_options.txt]  [text/x-meson]"
	case MetalKind:
		return "Metal  [*.metal]  [text/x-metal]"
	case MiniZincKind:
		return "MiniZinc  [*.mzn *.dzn *.fzn]  [text/minizinc]"
	case MlirKind:
		return "MLIR  [*.mlir]  [text/x-mlir]"
	case Modula2Kind:
		return "Modula-2  [*.def *.mod]  [text/x-modula2]"
	case MonkeyCKind:
		return "MonkeyC  [*.mc]  [text/x-monkeyc]"
	case MorrowindScriptKind:
		return "MorrowindScript  []  []"
	case MyghtyKind:
		return "Myghty  [*.myt autodelegate]  [application/x-myghty]"
	case MySqlKind:
		return "MySQL  [*.sql]  [text/x-mysql text/x-mariadb]"
	case NasmKind:
		return "NASM  [*.asm *.ASM *.nasm]  [text/x-nasm]"
	case NaturalKind:
		return "Natural  [*.NSN *.NSP *.NSS *.NSH *.NSG *.NSL *.NSA *.NSM *.NSC *.NS7]  [text/x-natural]"
	case NdisasmKind:
		return "NDISASM  []  [text/x-disasm]"
	case NewspeakKind:
		return "Newspeak  [*.ns2]  [text/x-newspeak]"
	case NginxConfigurationFileKind:
		return "Nginx configuration file  [nginx.conf]  [text/x-nginx-conf]"
	case NimKind:
		return "Nim  [*.nim *.nimrod]  [text/x-nim]"
	case NixKind:
		return "Nix  [*.nix]  [text/x-nix]"
	case ObjectiveCKind:
		return "Objective-C  [*.m *.h]  [text/x-objective-c]"
	case ObjectPascalKind:
		return "ObjectPascal  [*.pas *.pp *.inc *.dpr *.dpk *.lpr *.lpk]  [text/x-pascal]"
	case OCamlKind:
		return "OCaml  [*.ml *.mli *.mll *.mly]  [text/x-ocaml]"
	case OctaveKind:
		return "Octave  [*.m]  [text/octave]"
	case OdinKind:
		return "Odin  [*.odin]  [text/odin]"
	case OnesEnterpriseKind:
		return "OnesEnterprise  [*.EPF *.epf *.ERF *.erf]  [application/octet-stream]"
	case OpenEdgeAblKind:
		return "OpenEdge ABL  [*.p *.cls *.w *.i]  [text/x-openedge application/x-openedge]"
	case OpenScadKind:
		return "OpenSCAD  [*.scad]  [text/x-scad]"
	case OrgModeKind:
		return "Org Mode  [*.org]  [text/org]"
	case PacmanConfKind:
		return "PacmanConf  [pacman.conf]  []"
	case PerlKind:
		return "Perl  [*.pl *.pm *.t]  [text/x-perl application/x-perl]"
	case PhpKind:
		return "PHP  [*.php *.php[345] *.inc]  [text/x-php]"
	case PigKind:
		return "Pig  [*.pig]  [text/x-pig]"
	case PkgConfigKind:
		return "PkgConfig  [*.pc]  []"
	case PlPgSqlKind:
		return "PL/pgSQL  []  [text/x-plpgsql]"
	case PlaintextKind:
		return "plaintext  [*.txt]  [text/plain]"
	case PlutusCoreKind:
		return "Plutus Core  [*.plc]  [text/x-plutus-core application/x-plutus-core]"
	case PonyKind:
		return "Pony  [*.pony]  []"
	case PostgreSqlSqlDialectKind:
		return "PostgreSQL SQL dialect  []  [text/x-postgresql]"
	case PostScriptKind:
		return "PostScript  [*.ps *.eps]  [application/postscript]"
	case PovRayKind:
		return "POVRay  [*.pov *.inc]  [text/x-povray]"
	case PowerQueryKind:
		return "PowerQuery  [*.pq]  [text/x-powerquery]"
	case PowerShellKind:
		return "PowerShell  [*.ps1 *.psm1 *.psd1]  [text/x-powershell]"
	case PrologKind:
		return "Prolog  [*.ecl *.prolog *.pro *.pl]  [text/x-prolog]"
	case PromelaKind:
		return "Promela  [*.pml *.prom *.prm *.promela *.pr *.pm]  [text/x-promela]"
	case PromQlKind:
		return "PromQL  [*.promql]  []"
	case PropertiesKind:
		return "properties  [*.properties]  [text/x-java-properties]"
	case ProtocolBufferKind:
		return "Protocol Buffer  [*.proto]  []"
	case PrqlKind:
		return "PRQL  [*.prql]  [application/prql]"
	case PslKind:
		return "PSL  [*.psl *.BATCH *.TRIG *.PROC]  [text/x-psl]"
	case PuppetKind:
		return "Puppet  [*.pp]  []"
	case PythonKind:
		return "Python  [*.py *.pyi *.pyw *.jy *.sage *.sc SConstruct SConscript *.bzl BUCK BUILD BUILD.bazel WORKSPACE *.tac]  [text/x-python application/x-python text/x-python3 application/x-python3]"
	case Python2Kind:
		return "Python 2  []  [text/x-python2 application/x-python2]"
	case QBasicKind:
		return "QBasic  [*.BAS *.bas]  [text/basic]"
	case QmlKind:
		return "QML  [*.qml *.qbs]  [application/x-qml application/x-qt.qbs+qml]"
	case RKind:
		return "R  [*.S *.R *.r .Rhistory .Rprofile .Renviron]  [text/S-plus text/S text/x-r-source text/x-r text/x-R text/x-r-history text/x-r-profile]"
	case RacketKind:
		return "Racket  [*.rkt *.rktd *.rktl]  [text/x-racket application/x-racket]"
	case RagelKind:
		return "Ragel  []  []"
	case ReactKind:
		return "react  [*.jsx *.react]  [text/jsx text/typescript-jsx]"
	case ReasonMlKind:
		return "ReasonML  [*.re *.rei]  [text/x-reasonml]"
	case RegKind:
		return "reg  [*.reg]  [text/x-windows-registry]"
	case RegoKind:
		return "Rego  [*.rego]  []"
	case RexxKind:
		return "Rexx  [*.rexx *.rex *.rx *.arexx]  [text/x-rexx]"
	case RpmSpecKind:
		return "RPMSpec  [*.spec]  [text/x-rpm-spec]"
	case RubyKind:
		return "Ruby  [*.rb *.rbw Rakefile *.rake *.gemspec *.rbx *.duby Gemfile Vagrantfile]  [text/x-ruby application/x-ruby]"
	case RustKind:
		return "Rust  [*.rs *.rs.in]  [text/rust text/x-rust]"
	case SasKind:
		return "SAS  [*.SAS *.sas]  [text/x-sas text/sas application/x-sas]"
	case SassKind:
		return "Sass  [*.sass]  [text/x-sass]"
	case ScalaKind:
		return "Scala  [*.scala]  [text/x-scala]"
	case SchemeKind:
		return "Scheme  [*.scm *.ss]  [text/x-scheme application/x-scheme]"
	case ScilabKind:
		return "Scilab  [*.sci *.sce *.tst]  [text/scilab]"
	case ScssKind:
		return "SCSS  [*.scss]  [text/x-scss]"
	case SedKind:
		return "Sed  [*.sed *.[gs]sed]  [text/x-sed]"
	case SieveKind:
		return "Sieve  [*.siv *.sieve]  []"
	case SmaliKind:
		return "Smali  [*.smali]  [text/smali]"
	case SmalltalkKind:
		return "Smalltalk  [*.st]  [text/x-smalltalk]"
	case SmartyKind:
		return "Smarty  [*.tpl]  [application/x-smarty]"
	case SnobolKind:
		return "Snobol  [*.snobol]  [text/x-snobol]"
	case SolidityKind:
		return "Solidity  [*.sol]  []"
	case SourcePawnKind:
		return "SourcePawn  [*.sp *.inc]  [text/x-sourcepawn]"
	case SparqlKind:
		return "SPARQL  [*.rq *.sparql]  [application/sparql-query]"
	case SqlKind:
		return "SQL  [*.sql]  [text/x-sql]"
	case SquidConfKind:
		return "SquidConf  [squid.conf]  [text/x-squidconf]"
	case StandardMlKind:
		return "Standard ML  [*.sml *.sig *.fun]  [text/x-standardml application/x-standardml]"
	case StasKind:
		return "stas  [*.stas]  []"
	case StylusKind:
		return "Stylus  [*.styl]  [text/x-styl]"
	case SwiftKind:
		return "Swift  [*.swift]  [text/x-swift]"
	case SystemdKind:
		return "SYSTEMD  [*.automount *.device *.dnssd *.link *.mount *.netdev *.network *.path *.scope *.service *.slice *.socket *.swap *.target *.timer]  [text/plain]"
	case SystemverilogKind:
		return "systemverilog  [*.sv *.svh]  [text/x-systemverilog]"
	case TableGenKind:
		return "TableGen  [*.td]  [text/x-tablegen]"
	case TalKind:
		return "Tal  [*.tal]  [text/x-uxntal]"
	case TasmKind:
		return "TASM  [*.asm *.ASM *.tasm]  [text/x-tasm]"
	case TclKind:
		return "Tcl  [*.tcl *.rvt]  [text/x-tcl text/x-script.tcl application/x-tcl]"
	case TcshKind:
		return "Tcsh  [*.tcsh *.csh]  [application/x-csh]"
	case TermcapKind:
		return "Termcap  [termcap termcap.src]  []"
	case TerminfoKind:
		return "Terminfo  [terminfo terminfo.src]  []"
	case TerraformKind:
		return "Terraform  [*.tf]  [application/x-tf application/x-terraform]"
	case TeXKind:
		return "TeX  [*.tex *.aux *.toc]  [text/x-tex text/x-latex]"
	case ThriftKind:
		return "Thrift  [*.thrift]  [application/x-thrift]"
	case TomlKind:
		return "TOML  [*.toml Pipfile poetry.lock]  [text/x-toml]"
	case TradingViewKind:
		return "TradingView  [*.tv]  [text/x-tradingview]"
	case TransactSqlKind:
		return "Transact-SQL  []  [text/x-tsql]"
	case TuringKind:
		return "Turing  [*.turing *.tu]  [text/x-turing]"
	case TurtleKind:
		return "Turtle  [*.ttl]  [text/turtle application/x-turtle]"
	case TwigKind:
		return "Twig  [*.twig]  [application/x-twig]"
	case TypeScriptKind:
		return "TypeScript  [*.ts *.tsx *.mts *.cts]  [text/x-typescript]"
	case TypoScriptKind:
		return "TypoScript  [*.ts]  [text/x-typoscript]"
	case TypoScriptCssDataKind:
		return "TypoScriptCssData  []  []"
	case TypoScriptHtmlDataKind:
		return "TypoScriptHtmlData  []  []"
	case UcodeKind:
		return "ucode  [*.uc]  [application/x.ucode text/x.ucode]"
	case VKind:
		return "V  [*.v *.vv v.mod]  [text/x-v]"
	case VShellKind:
		return "V shell  [*.vsh]  [text/x-vsh]"
	case ValaKind:
		return "Vala  [*.vala *.vapi]  [text/x-vala]"
	case VbNetKind:
		return "VB.net  [*.vb *.bas]  [text/x-vbnet text/x-vba]"
	case VerilogKind:
		return "verilog  [*.v]  [text/x-verilog]"
	case VhdlKind:
		return "VHDL  [*.vhdl *.vhd]  [text/x-vhdl]"
	case VhsKind:
		return "VHS  [*.tape]  []"
	case VimLKind:
		return "VimL  [*.vim .vimrc .exrc .gvimrc _vimrc _exrc _gvimrc vimrc gvimrc]  [text/x-vim]"
	case VueKind:
		return "vue  [*.vue]  [text/x-vue application/x-vue]"
	case WdteKind:
		return "WDTE  [*.wdte]  []"
	case WebGpuShadingLanguageKind:
		return "WebGPU Shading Language  [*.wgsl]  [text/wgsl]"
	case WhileyKind:
		return "Whiley  [*.whiley]  [text/x-whiley]"
	case XmlKind:
		return "XML  [*.xml *.xsl *.rss *.xslt *.xsd *.wsdl *.wsf *.svg *.csproj *.vcxproj *.fsproj]  [text/xml application/xml image/svg+xml application/rss+xml application/atom+xml]"
	case XorgKind:
		return "Xorg  [xorg.conf]  []"
	case YamlKind:
		return "YAML  [*.yaml *.yml]  [text/x-yaml]"
	case YangKind:
		return "YANG  [*.yang]  [application/yang]"
	case Z80AssemblyKind:
		return "Z80 Assembly  [*.z80 *.asm]  []"
	case ZedKind:
		return "Zed  [*.zed]  [text/zed]"
	case ZigKind:
		return "Zig  [*.zig]  [text/zig]"
	case CaddyfileKind:
		return "Caddyfile  [Caddyfile*]  []"
	case CaddyfileDirectivesKind:
		return "Caddyfile Directives  []  []"
	case GenshiTextKind:
		return "Genshi Text  []  [application/x-genshi-text text/x-genshi]"
	case GenshiHtmlKind:
		return "Genshi HTML  []  [text/html+genshi]"
	case GenshiKind:
		return "Genshi  [*.kid]  [application/x-genshi application/x-kid]"
	case GoHtmlTemplateKind:
		return "Go HTML Template  []  []"
	case GoTextTemplateKind:
		return "Go Text Template  []  []"
	case GoKind:
		return "Go  [*.go]  [text/x-gosrc]"
	case HaxeKind:
		return "Haxe  [*.hx *.hxsl]  [text/haxe text/x-haxe text/x-hx]"
	case HttpKind:
		return "HTTP  []  []"
	case MarkdownKind:
		return "markdown  [*.md *.mkd *.markdown]  [text/x-markdown]"
	case PhtmlKind:
		return "PHTML  [*.phtml *.php *.php[345] *.inc]  [application/x-php application/x-httpd-php application/x-httpd-php3 application/x-httpd-php4 application/x-httpd-php5 text/x-php]"
	case RakuKind:
		return "Raku  [*.pl *.pm *.nqp *.p6 *.6pl *.p6l *.pl6 *.6pm *.p6m *.pm6 *.t *.raku *.rakumod *.rakutest *.rakudoc]  [text/x-perl6 application/x-perl6 text/x-raku application/x-raku]"
	case ReStructuredTextKind:
		return "reStructuredText  [*.rst *.rest]  [text/x-rst text/prs.fallenstein.rst]"
	case SvelteKind:
		return "Svelte  [*.svelte]  [application/x-svelte]"
	default:
		return "InvalidLanguagesKind"
	}
}

func (k LanguagesKind) Keys() []string {
	return []string{
		"Abap",
		"Abnf",
		"ActionScript",
		"ActionScript3",
		"Ada",
		"Agda",
		"Al",
		"Alloy",
		"Angular2",
		"Antlr",
		"ApacheConf",
		"Apl",
		"AppleScript",
		"ArangoDbAql",
		"Arduino",
		"ArmAsm",
		"AutoHotkey",
		"AutoIt",
		"Awk",
		"Ballerina",
		"Bash",
		"BashSession",
		"Batchfile",
		"BibTeX",
		"Bicep",
		"BlitzBasic",
		"Bnf",
		"Bqn",
		"Brainfuck",
		"CSharp",
		"Cpp",
		"C",
		"CapnProto",
		"CassandraCql",
		"Ceylon",
		"CfEngine3",
		"Cfstatement",
		"ChaiScript",
		"Chapel",
		"Cheetah",
		"Clojure",
		"CMake",
		"Cobol",
		"CoffeeScript",
		"CommonLisp",
		"Coq",
		"Crystal",
		"Css",
		"Cue",
		"Cython",
		"D",
		"Dart",
		"Dax",
		"DesktopFile",
		"Diff",
		"DjangoJinja",
		"Dns",
		"Docker",
		"Dtd",
		"Dylan",
		"Ebnf",
		"Elixir",
		"Elm",
		"EmacsLisp",
		"Erlang",
		"Factor",
		"Fennel",
		"Fish",
		"Forth",
		"Fortran",
		"FortranFixed",
		"FSharp",
		"Gas",
		"GdScript",
		"GdScript3",
		"Gherkin",
		"Glsl",
		"Gnuplot",
		"GoTemplate",
		"GraphQl",
		"Groff",
		"Groovy",
		"Handlebars",
		"Hare",
		"Haskell",
		"Hcl",
		"Hexdump",
		"Hlb",
		"Hlsl",
		"HolyC",
		"Html",
		"Hy",
		"Idris",
		"Igor",
		"Ini",
		"Io",
		"IsCdhcpd",
		"J",
		"Java",
		"JavaScript",
		"Json",
		"Julia",
		"Jungle",
		"Kotlin",
		"LighttpdConfigurationFile",
		"Llvm",
		"Lua",
		"Makefile",
		"Mako",
		"Mason",
		"MaterializeSqlDialect",
		"Mathematica",
		"Matlab",
		"Mcfunction",
		"Meson",
		"Metal",
		"MiniZinc",
		"Mlir",
		"Modula2",
		"MonkeyC",
		"MorrowindScript",
		"Myghty",
		"MySql",
		"Nasm",
		"Natural",
		"Ndisasm",
		"Newspeak",
		"NginxConfigurationFile",
		"Nim",
		"Nix",
		"ObjectiveC",
		"ObjectPascal",
		"OCaml",
		"Octave",
		"Odin",
		"OnesEnterprise",
		"OpenEdgeAbl",
		"OpenScad",
		"OrgMode",
		"PacmanConf",
		"Perl",
		"Php",
		"Pig",
		"PkgConfig",
		"PlPgSql",
		"Plaintext",
		"PlutusCore",
		"Pony",
		"PostgreSqlSqlDialect",
		"PostScript",
		"PovRay",
		"PowerQuery",
		"PowerShell",
		"Prolog",
		"Promela",
		"PromQl",
		"Properties",
		"ProtocolBuffer",
		"Prql",
		"Psl",
		"Puppet",
		"Python",
		"Python2",
		"QBasic",
		"Qml",
		"R",
		"Racket",
		"Ragel",
		"React",
		"ReasonMl",
		"Reg",
		"Rego",
		"Rexx",
		"RpmSpec",
		"Ruby",
		"Rust",
		"Sas",
		"Sass",
		"Scala",
		"Scheme",
		"Scilab",
		"Scss",
		"Sed",
		"Sieve",
		"Smali",
		"Smalltalk",
		"Smarty",
		"Snobol",
		"Solidity",
		"SourcePawn",
		"Sparql",
		"Sql",
		"SquidConf",
		"StandardMl",
		"Stas",
		"Stylus",
		"Swift",
		"Systemd",
		"Systemverilog",
		"TableGen",
		"Tal",
		"Tasm",
		"Tcl",
		"Tcsh",
		"Termcap",
		"Terminfo",
		"Terraform",
		"TeX",
		"Thrift",
		"Toml",
		"TradingView",
		"TransactSql",
		"Turing",
		"Turtle",
		"Twig",
		"TypeScript",
		"TypoScript",
		"TypoScriptCssData",
		"TypoScriptHtmlData",
		"Ucode",
		"V",
		"VShell",
		"Vala",
		"VbNet",
		"Verilog",
		"Vhdl",
		"Vhs",
		"VimL",
		"Vue",
		"Wdte",
		"WebGpuShadingLanguage",
		"Whiley",
		"Xml",
		"Xorg",
		"Yaml",
		"Yang",
		"Z80Assembly",
		"Zed",
		"Zig",
		"Caddyfile",
		"CaddyfileDirectives",
		"GenshiText",
		"GenshiHtml",
		"Genshi",
		"GoHtmlTemplate",
		"GoTextTemplate",
		"Go",
		"Haxe",
		"Http",
		"Markdown",
		"Phtml",
		"Raku",
		"ReStructuredText",
		"Svelte",
		"InvalidLanguagesKind",
	}
}

func (k LanguagesKind) Kinds() []LanguagesKind {
	return []LanguagesKind{
		AbapKind,
		AbnfKind,
		ActionScriptKind,
		ActionScript3Kind,
		AdaKind,
		AgdaKind,
		AlKind,
		AlloyKind,
		Angular2Kind,
		AntlrKind,
		ApacheConfKind,
		AplKind,
		AppleScriptKind,
		ArangoDbAqlKind,
		ArduinoKind,
		ArmAsmKind,
		AutoHotkeyKind,
		AutoItKind,
		AwkKind,
		BallerinaKind,
		BashKind,
		BashSessionKind,
		BatchfileKind,
		BibTeXKind,
		BicepKind,
		BlitzBasicKind,
		BnfKind,
		BqnKind,
		BrainfuckKind,
		CSharpKind,
		CppKind,
		CKind,
		CapnProtoKind,
		CassandraCqlKind,
		CeylonKind,
		CfEngine3Kind,
		CfstatementKind,
		ChaiScriptKind,
		ChapelKind,
		CheetahKind,
		ClojureKind,
		CMakeKind,
		CobolKind,
		CoffeeScriptKind,
		CommonLispKind,
		CoqKind,
		CrystalKind,
		CssKind,
		CueKind,
		CythonKind,
		DKind,
		DartKind,
		DaxKind,
		DesktopFileKind,
		DiffKind,
		DjangoJinjaKind,
		DnsKind,
		DockerKind,
		DtdKind,
		DylanKind,
		EbnfKind,
		ElixirKind,
		ElmKind,
		EmacsLispKind,
		ErlangKind,
		FactorKind,
		FennelKind,
		FishKind,
		ForthKind,
		FortranKind,
		FortranFixedKind,
		FSharpKind,
		GasKind,
		GdScriptKind,
		GdScript3Kind,
		GherkinKind,
		GlslKind,
		GnuplotKind,
		GoTemplateKind,
		GraphQlKind,
		GroffKind,
		GroovyKind,
		HandlebarsKind,
		HareKind,
		HaskellKind,
		HclKind,
		HexdumpKind,
		HlbKind,
		HlslKind,
		HolyCKind,
		HtmlKind,
		HyKind,
		IdrisKind,
		IgorKind,
		IniKind,
		IoKind,
		IsCdhcpdKind,
		JKind,
		JavaKind,
		JavaScriptKind,
		JsonKind,
		JuliaKind,
		JungleKind,
		KotlinKind,
		LighttpdConfigurationFileKind,
		LlvmKind,
		LuaKind,
		MakefileKind,
		MakoKind,
		MasonKind,
		MaterializeSqlDialectKind,
		MathematicaKind,
		MatlabKind,
		McfunctionKind,
		MesonKind,
		MetalKind,
		MiniZincKind,
		MlirKind,
		Modula2Kind,
		MonkeyCKind,
		MorrowindScriptKind,
		MyghtyKind,
		MySqlKind,
		NasmKind,
		NaturalKind,
		NdisasmKind,
		NewspeakKind,
		NginxConfigurationFileKind,
		NimKind,
		NixKind,
		ObjectiveCKind,
		ObjectPascalKind,
		OCamlKind,
		OctaveKind,
		OdinKind,
		OnesEnterpriseKind,
		OpenEdgeAblKind,
		OpenScadKind,
		OrgModeKind,
		PacmanConfKind,
		PerlKind,
		PhpKind,
		PigKind,
		PkgConfigKind,
		PlPgSqlKind,
		PlaintextKind,
		PlutusCoreKind,
		PonyKind,
		PostgreSqlSqlDialectKind,
		PostScriptKind,
		PovRayKind,
		PowerQueryKind,
		PowerShellKind,
		PrologKind,
		PromelaKind,
		PromQlKind,
		PropertiesKind,
		ProtocolBufferKind,
		PrqlKind,
		PslKind,
		PuppetKind,
		PythonKind,
		Python2Kind,
		QBasicKind,
		QmlKind,
		RKind,
		RacketKind,
		RagelKind,
		ReactKind,
		ReasonMlKind,
		RegKind,
		RegoKind,
		RexxKind,
		RpmSpecKind,
		RubyKind,
		RustKind,
		SasKind,
		SassKind,
		ScalaKind,
		SchemeKind,
		ScilabKind,
		ScssKind,
		SedKind,
		SieveKind,
		SmaliKind,
		SmalltalkKind,
		SmartyKind,
		SnobolKind,
		SolidityKind,
		SourcePawnKind,
		SparqlKind,
		SqlKind,
		SquidConfKind,
		StandardMlKind,
		StasKind,
		StylusKind,
		SwiftKind,
		SystemdKind,
		SystemverilogKind,
		TableGenKind,
		TalKind,
		TasmKind,
		TclKind,
		TcshKind,
		TermcapKind,
		TerminfoKind,
		TerraformKind,
		TeXKind,
		ThriftKind,
		TomlKind,
		TradingViewKind,
		TransactSqlKind,
		TuringKind,
		TurtleKind,
		TwigKind,
		TypeScriptKind,
		TypoScriptKind,
		TypoScriptCssDataKind,
		TypoScriptHtmlDataKind,
		UcodeKind,
		VKind,
		VShellKind,
		ValaKind,
		VbNetKind,
		VerilogKind,
		VhdlKind,
		VhsKind,
		VimLKind,
		VueKind,
		WdteKind,
		WebGpuShadingLanguageKind,
		WhileyKind,
		XmlKind,
		XorgKind,
		YamlKind,
		YangKind,
		Z80AssemblyKind,
		ZedKind,
		ZigKind,
		CaddyfileKind,
		CaddyfileDirectivesKind,
		GenshiTextKind,
		GenshiHtmlKind,
		GenshiKind,
		GoHtmlTemplateKind,
		GoTextTemplateKind,
		GoKind,
		HaxeKind,
		HttpKind,
		MarkdownKind,
		PhtmlKind,
		RakuKind,
		ReStructuredTextKind,
		SvelteKind,
		InvalidLanguagesKind,
	}
}

func (k LanguagesKind) SvgFileName() string {
	switch k {
	case AbapKind:
		return "Abap"
	case AbnfKind:
		return "Abnf"
	case ActionScriptKind:
		return "ActionScript"
	case ActionScript3Kind:
		return "ActionScript3"
	case AdaKind:
		return "Ada"
	case AgdaKind:
		return "Agda"
	case AlKind:
		return "Al"
	case AlloyKind:
		return "Alloy"
	case Angular2Kind:
		return "Angular2"
	case AntlrKind:
		return "Antlr"
	case ApacheConfKind:
		return "ApacheConf"
	case AplKind:
		return "Apl"
	case AppleScriptKind:
		return "AppleScript"
	case ArangoDbAqlKind:
		return "ArangoDbAql"
	case ArduinoKind:
		return "Arduino"
	case ArmAsmKind:
		return "ArmAsm"
	case AutoHotkeyKind:
		return "AutoHotkey"
	case AutoItKind:
		return "AutoIt"
	case AwkKind:
		return "Awk"
	case BallerinaKind:
		return "Ballerina"
	case BashKind:
		return "Bash"
	case BashSessionKind:
		return "BashSession"
	case BatchfileKind:
		return "Batchfile"
	case BibTeXKind:
		return "BibTeX"
	case BicepKind:
		return "Bicep"
	case BlitzBasicKind:
		return "BlitzBasic"
	case BnfKind:
		return "Bnf"
	case BqnKind:
		return "Bqn"
	case BrainfuckKind:
		return "Brainfuck"
	case CSharpKind:
		return "CSharp"
	case CppKind:
		return "Cpp"
	case CKind:
		return "C"
	case CapnProtoKind:
		return "CapnProto"
	case CassandraCqlKind:
		return "CassandraCql"
	case CeylonKind:
		return "Ceylon"
	case CfEngine3Kind:
		return "CfEngine3"
	case CfstatementKind:
		return "Cfstatement"
	case ChaiScriptKind:
		return "ChaiScript"
	case ChapelKind:
		return "Chapel"
	case CheetahKind:
		return "Cheetah"
	case ClojureKind:
		return "Clojure"
	case CMakeKind:
		return "CMake"
	case CobolKind:
		return "Cobol"
	case CoffeeScriptKind:
		return "CoffeeScript"
	case CommonLispKind:
		return "CommonLisp"
	case CoqKind:
		return "Coq"
	case CrystalKind:
		return "Crystal"
	case CssKind:
		return "Css"
	case CueKind:
		return "Cue"
	case CythonKind:
		return "Cython"
	case DKind:
		return "D"
	case DartKind:
		return "Dart"
	case DaxKind:
		return "Dax"
	case DesktopFileKind:
		return "DesktopFile"
	case DiffKind:
		return "Diff"
	case DjangoJinjaKind:
		return "DjangoJinja"
	case DnsKind:
		return "Dns"
	case DockerKind:
		return "Docker"
	case DtdKind:
		return "Dtd"
	case DylanKind:
		return "Dylan"
	case EbnfKind:
		return "Ebnf"
	case ElixirKind:
		return "Elixir"
	case ElmKind:
		return "Elm"
	case EmacsLispKind:
		return "EmacsLisp"
	case ErlangKind:
		return "Erlang"
	case FactorKind:
		return "Factor"
	case FennelKind:
		return "Fennel"
	case FishKind:
		return "Fish"
	case ForthKind:
		return "Forth"
	case FortranKind:
		return "Fortran"
	case FortranFixedKind:
		return "FortranFixed"
	case FSharpKind:
		return "FSharp"
	case GasKind:
		return "Gas"
	case GdScriptKind:
		return "GdScript"
	case GdScript3Kind:
		return "GdScript3"
	case GherkinKind:
		return "Gherkin"
	case GlslKind:
		return "Glsl"
	case GnuplotKind:
		return "Gnuplot"
	case GoTemplateKind:
		return "GoTemplate"
	case GraphQlKind:
		return "GraphQl"
	case GroffKind:
		return "Groff"
	case GroovyKind:
		return "Groovy"
	case HandlebarsKind:
		return "Handlebars"
	case HareKind:
		return "Hare"
	case HaskellKind:
		return "Haskell"
	case HclKind:
		return "Hcl"
	case HexdumpKind:
		return "Hexdump"
	case HlbKind:
		return "Hlb"
	case HlslKind:
		return "Hlsl"
	case HolyCKind:
		return "HolyC"
	case HtmlKind:
		return "Html"
	case HyKind:
		return "Hy"
	case IdrisKind:
		return "Idris"
	case IgorKind:
		return "Igor"
	case IniKind:
		return "Ini"
	case IoKind:
		return "Io"
	case IsCdhcpdKind:
		return "IsCdhcpd"
	case JKind:
		return "J"
	case JavaKind:
		return "Java"
	case JavaScriptKind:
		return "JavaScript"
	case JsonKind:
		return "Json"
	case JuliaKind:
		return "Julia"
	case JungleKind:
		return "Jungle"
	case KotlinKind:
		return "Kotlin"
	case LighttpdConfigurationFileKind:
		return "LighttpdConfigurationFile"
	case LlvmKind:
		return "Llvm"
	case LuaKind:
		return "Lua"
	case MakefileKind:
		return "Makefile"
	case MakoKind:
		return "Mako"
	case MasonKind:
		return "Mason"
	case MaterializeSqlDialectKind:
		return "MaterializeSqlDialect"
	case MathematicaKind:
		return "Mathematica"
	case MatlabKind:
		return "Matlab"
	case McfunctionKind:
		return "Mcfunction"
	case MesonKind:
		return "Meson"
	case MetalKind:
		return "Metal"
	case MiniZincKind:
		return "MiniZinc"
	case MlirKind:
		return "Mlir"
	case Modula2Kind:
		return "Modula2"
	case MonkeyCKind:
		return "MonkeyC"
	case MorrowindScriptKind:
		return "MorrowindScript"
	case MyghtyKind:
		return "Myghty"
	case MySqlKind:
		return "MySql"
	case NasmKind:
		return "Nasm"
	case NaturalKind:
		return "Natural"
	case NdisasmKind:
		return "Ndisasm"
	case NewspeakKind:
		return "Newspeak"
	case NginxConfigurationFileKind:
		return "NginxConfigurationFile"
	case NimKind:
		return "Nim"
	case NixKind:
		return "Nix"
	case ObjectiveCKind:
		return "ObjectiveC"
	case ObjectPascalKind:
		return "ObjectPascal"
	case OCamlKind:
		return "OCaml"
	case OctaveKind:
		return "Octave"
	case OdinKind:
		return "Odin"
	case OnesEnterpriseKind:
		return "OnesEnterprise"
	case OpenEdgeAblKind:
		return "OpenEdgeAbl"
	case OpenScadKind:
		return "OpenScad"
	case OrgModeKind:
		return "OrgMode"
	case PacmanConfKind:
		return "PacmanConf"
	case PerlKind:
		return "Perl"
	case PhpKind:
		return "Php"
	case PigKind:
		return "Pig"
	case PkgConfigKind:
		return "PkgConfig"
	case PlPgSqlKind:
		return "PlPgSql"
	case PlaintextKind:
		return "Plaintext"
	case PlutusCoreKind:
		return "PlutusCore"
	case PonyKind:
		return "Pony"
	case PostgreSqlSqlDialectKind:
		return "PostgreSqlSqlDialect"
	case PostScriptKind:
		return "PostScript"
	case PovRayKind:
		return "PovRay"
	case PowerQueryKind:
		return "PowerQuery"
	case PowerShellKind:
		return "PowerShell"
	case PrologKind:
		return "Prolog"
	case PromelaKind:
		return "Promela"
	case PromQlKind:
		return "PromQl"
	case PropertiesKind:
		return "Properties"
	case ProtocolBufferKind:
		return "ProtocolBuffer"
	case PrqlKind:
		return "Prql"
	case PslKind:
		return "Psl"
	case PuppetKind:
		return "Puppet"
	case PythonKind:
		return "Python"
	case Python2Kind:
		return "Python2"
	case QBasicKind:
		return "QBasic"
	case QmlKind:
		return "Qml"
	case RKind:
		return "R"
	case RacketKind:
		return "Racket"
	case RagelKind:
		return "Ragel"
	case ReactKind:
		return "React"
	case ReasonMlKind:
		return "ReasonMl"
	case RegKind:
		return "Reg"
	case RegoKind:
		return "Rego"
	case RexxKind:
		return "Rexx"
	case RpmSpecKind:
		return "RpmSpec"
	case RubyKind:
		return "Ruby"
	case RustKind:
		return "Rust"
	case SasKind:
		return "Sas"
	case SassKind:
		return "Sass"
	case ScalaKind:
		return "Scala"
	case SchemeKind:
		return "Scheme"
	case ScilabKind:
		return "Scilab"
	case ScssKind:
		return "Scss"
	case SedKind:
		return "Sed"
	case SieveKind:
		return "Sieve"
	case SmaliKind:
		return "Smali"
	case SmalltalkKind:
		return "Smalltalk"
	case SmartyKind:
		return "Smarty"
	case SnobolKind:
		return "Snobol"
	case SolidityKind:
		return "Solidity"
	case SourcePawnKind:
		return "SourcePawn"
	case SparqlKind:
		return "Sparql"
	case SqlKind:
		return "Sql"
	case SquidConfKind:
		return "SquidConf"
	case StandardMlKind:
		return "StandardMl"
	case StasKind:
		return "Stas"
	case StylusKind:
		return "Stylus"
	case SwiftKind:
		return "Swift"
	case SystemdKind:
		return "Systemd"
	case SystemverilogKind:
		return "Systemverilog"
	case TableGenKind:
		return "TableGen"
	case TalKind:
		return "Tal"
	case TasmKind:
		return "Tasm"
	case TclKind:
		return "Tcl"
	case TcshKind:
		return "Tcsh"
	case TermcapKind:
		return "Termcap"
	case TerminfoKind:
		return "Terminfo"
	case TerraformKind:
		return "Terraform"
	case TeXKind:
		return "TeX"
	case ThriftKind:
		return "Thrift"
	case TomlKind:
		return "Toml"
	case TradingViewKind:
		return "TradingView"
	case TransactSqlKind:
		return "TransactSql"
	case TuringKind:
		return "Turing"
	case TurtleKind:
		return "Turtle"
	case TwigKind:
		return "Twig"
	case TypeScriptKind:
		return "TypeScript"
	case TypoScriptKind:
		return "TypoScript"
	case TypoScriptCssDataKind:
		return "TypoScriptCssData"
	case TypoScriptHtmlDataKind:
		return "TypoScriptHtmlData"
	case UcodeKind:
		return "Ucode"
	case VKind:
		return "V"
	case VShellKind:
		return "VShell"
	case ValaKind:
		return "Vala"
	case VbNetKind:
		return "VbNet"
	case VerilogKind:
		return "Verilog"
	case VhdlKind:
		return "Vhdl"
	case VhsKind:
		return "Vhs"
	case VimLKind:
		return "VimL"
	case VueKind:
		return "Vue"
	case WdteKind:
		return "Wdte"
	case WebGpuShadingLanguageKind:
		return "WebGpuShadingLanguage"
	case WhileyKind:
		return "Whiley"
	case XmlKind:
		return "Xml"
	case XorgKind:
		return "Xorg"
	case YamlKind:
		return "Yaml"
	case YangKind:
		return "Yang"
	case Z80AssemblyKind:
		return "Z80Assembly"
	case ZedKind:
		return "Zed"
	case ZigKind:
		return "Zig"
	case CaddyfileKind:
		return "Caddyfile"
	case CaddyfileDirectivesKind:
		return "CaddyfileDirectives"
	case GenshiTextKind:
		return "GenshiText"
	case GenshiHtmlKind:
		return "GenshiHtml"
	case GenshiKind:
		return "Genshi"
	case GoHtmlTemplateKind:
		return "GoHtmlTemplate"
	case GoTextTemplateKind:
		return "GoTextTemplate"
	case GoKind:
		return "Go"
	case HaxeKind:
		return "Haxe"
	case HttpKind:
		return "Http"
	case MarkdownKind:
		return "Markdown"
	case PhtmlKind:
		return "Phtml"
	case RakuKind:
		return "Raku"
	case ReStructuredTextKind:
		return "ReStructuredText"
	case SvelteKind:
		return "Svelte"
	default:
		return "InvalidLanguagesKind"
	}
}
