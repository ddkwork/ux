package ux

import (
	"testing"
)

func TestNewPrompt(t *testing.T) {
	NewPrompt("Save", "Do you want to save the changes? (Tips: you can always save the changes using CMD/CTRL+s)", ModalTypeWarn,
		[]Option{{Text: "Yes"}, {Text: "No"}, {Text: "Cancel"}}...,
	)
	//func(selectedOption string, remember bool) {
	//			if selectedOption == "Cancel" {
	//				c.view.HidePrompt(id)
	//				return
	//			}
	//
	//			if selectedOption == "Yes" {
	//				c.saveEnvironment(id)
	//			}
	//
	//			c.view.CloseTab(id)
	//			c.state.ReloadEnvironment(id, state.SourceController)
	//		},
}
