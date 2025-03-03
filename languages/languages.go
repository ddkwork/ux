package languages

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

var configs []*chroma.Config

func init() {
	for _, lexer := range lexers.GlobalLexerRegistry.Lexers {
		configs = append(configs, lexer.Config())
	}
}

type Language struct {
	Name       string
	Extensions []string
	MimeTypes  []string
}

type Languages struct {
	Languages []Language
	Names     []string
	Tooltips  []string
}

// todo https://github.com/abadojack/whatlanggo
func NewLanguages() *Languages {
	l := &Languages{
		Languages: make([]Language, 0),
		Names:     make([]string, 0),
		Tooltips:  make([]string, 0),
	}

	for _, config := range configs {
		l.Languages = append(l.Languages, Language{Name: config.Name, Extensions: config.Filenames, MimeTypes: config.MimeTypes})
	}

	for _, language := range l.Languages {
		l.Names = append(l.Names, strings.NewReplacer(
			"'", "",
			"C++", "Cpp",
			"C#", "CSharp",
			"-", "",
		).Replace(language.Name))

		b := stream.NewBuffer("")
		b.WriteString(language.Name)
		b.Indent(2)
		b.WriteString(fmt.Sprintf("%v", language.Extensions))
		b.Indent(2)
		b.WriteString(fmt.Sprintf("%v", language.MimeTypes))
		l.Tooltips = append(l.Tooltips, b.String())
	}

	return l
}

func GetTokens(code *stream.Buffer, language LanguagesKind) ([]chroma.Token, *chroma.Style) {
	lexer := lexers.Get(language.String())
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)
	style := styles.Get("dracula")
	if style == nil {
		style = styles.Fallback
	}
	iterator := mylog.Check2(lexer.Tokenise(nil, code.String()))
	return iterator.Tokens(), style
}

func (l *Languages) CodeFile2Language(path string) LanguagesKind {
	mylog.Check(stream.IsFilePath(path))
	for _, language := range l.Languages {
		for _, extension := range language.Extensions {
			extension = strings.TrimPrefix(extension, "*")
			if filepath.Ext(path) == extension {
				return InvalidLanguagesKind.AssertKind(language.Name)
			}
		}
	}
	return InvalidLanguagesKind
}

func Code2Language(code string) LanguagesKind {
	lexer := lexers.Analyse(code)
	mylog.CheckNil(lexer)
	l := chroma.Coalesce(lexer)
	if l == nil {
		mylog.Check("未知语言类型")
		return InvalidLanguagesKind
	}
	kind := InvalidLanguagesKind.AssertKind(l.Config().Name)
	if kind == InvalidLanguagesKind {
		mylog.Warning(l.Config().Name + " 未注册到 chroma 库中")
	}
	return kind
}

func CodeFile2Language(path string) LanguagesKind {
	lexer := lexers.Get(path)
	if lexer == nil {
		mylog.CheckIgnore("不支持的文件类型")
		return InvalidLanguagesKind
	}
	return InvalidLanguagesKind.AssertKind(lexer.Config().Name)
}
