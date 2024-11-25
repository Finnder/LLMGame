package ui

import (
	"fyne.io/fyne/v2/widget"
)

func CreateLabel(text string) *widget.Label {
	return widget.NewLabel(text)
}
