//go:build !android && !ios

package component

import "github.com/ddkwork/ux/widget/material"

// PlatformTooltip creates a tooltip styled to the current platform
// (desktop or mobile) by choosing based on the OS. This choice may
// not always be appropriate as it only uses the OS to decide.
func PlatformTooltip(th *material.Theme, text string) Tooltip {
	return DesktopTooltip(th, text)
}
