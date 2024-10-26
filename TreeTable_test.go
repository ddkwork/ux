package ux

import (
	"testing"

	"github.com/ddkwork/golibrary/mylog"
)

func TestTreeTable_ContextMenuItem(t1 *testing.T) {
	mylog.Skips = append(mylog.Skips, "patch")
}
func init() {
	mylog.FormatAllFiles()
}
