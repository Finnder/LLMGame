package ui

import (
	"fyne.io/fyne/v2/widget"
)

// CreateButton creates a customizable button with a label and action.
func CreateButton(label string, onClick func()) *widget.Button {
	return widget.NewButton(label, onClick)
}
