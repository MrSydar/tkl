package main

import (
	"fmt"

	"fyne.io/fyne/v2/widget"
)

type ProgressBarWithText struct {
	widget.ProgressBar

	message string
}

func (p *ProgressBarWithText) Update(message string, progress float64) {
	p.message = message
	p.SetValue(progress)
}

func NewProgressBarWithMessage() *ProgressBarWithText {
	bar := &ProgressBarWithText{
		message: "Progress",
	}
	bar.ExtendBaseWidget(bar)

	bar.TextFormatter = func() string {
		return fmt.Sprintf("%v : %.0f%%", bar.message, bar.Value*100)
	}

	return bar
}
