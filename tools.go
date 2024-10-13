//go:build tools
// +build tools

package ux

import (
	gogio "gioui.org/cmd/gogio"
	_ "gioui.org/example"
	_ "gioui.org/example/7gui/counter"
	_ "gioui.org/example/7gui/temperature"
	_ "gioui.org/example/7gui/timer"
	_ "gioui.org/example/bidi"
	_ "gioui.org/example/color-grid"
	_ "gioui.org/example/colorpicker"
	_ "gioui.org/example/component"
	_ "gioui.org/example/component/applayout"
	_ "gioui.org/example/component/icon"
	_ "gioui.org/example/component/pages"
	_ "gioui.org/example/component/pages/about"
	_ "gioui.org/example/component/pages/appbar"
	_ "gioui.org/example/component/pages/discloser"
	_ "gioui.org/example/component/pages/menu"
	_ "gioui.org/example/component/pages/navdrawer"
	_ "gioui.org/example/component/pages/textfield"
	_ "gioui.org/example/customdeco"
	_ "gioui.org/example/explorer"
	_ "gioui.org/example/fps-table"
	_ "gioui.org/example/galaxy"
	_ "gioui.org/example/glfw"
	_ "gioui.org/example/gophers"
	_ "gioui.org/example/haptic"
	_ "gioui.org/example/hello"
	_ "gioui.org/example/kitchen"
	_ "gioui.org/example/life"
	_ "gioui.org/example/markdown"
	_ "gioui.org/example/multiwindow"
	_ "gioui.org/example/notify"
	_ "gioui.org/example/opacity"
	_ "gioui.org/example/opengl"
	_ "gioui.org/example/outlay/fan"
	_ "gioui.org/example/outlay/fan/cribbage"
	_ "gioui.org/example/outlay/fan/cribbage/cmd"
	_ "gioui.org/example/outlay/fan/playing"
	_ "gioui.org/example/outlay/fan/widget"
	_ "gioui.org/example/outlay/fan/widget/boring"
	_ "gioui.org/example/outlay/grid"
	_ "gioui.org/example/tabs"
	_ "gioui.org/example/textfeatures"
)

func init() {
	gogio.AndroidPermissions["networkSecurityConfig"] = []string{
		"network_security_config.xml",
		`<?xml version="1.0" encoding="utf-8"?>
<network-security-config>
   <domain-config cleartextTrafficPermitted="true">
       <domain includeSubdomains="true">localhost</domain>
       <domain includeSubdomains="true">127.0.0.1</domain>
   </domain-config>
</network-security-config>
`,
	}
}

func main() {
	m := map[string]string{}
	filepath.Walk("D:\\111111111111111111111111\\giotest\\up\\gio-example", func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" {
			m[filepath.Dir(path)] = filepath.Base(path)
		}
		return err
	})
	for s := range m {
		_, after, found := strings.Cut(s, "gio-example")
		if !found {
			continue
		}
		after = strings.ReplaceAll(after, "\\", "/")
		after = "gioui.org/example" + after
		after = "_ " + strconv.Quote(after)
		println(after)
	}
}

/*
This file locks gogio as a dependency so that its version will
stay in sync with the version of gio that we use in our go.mod.
*/
