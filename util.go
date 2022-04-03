package main

import "fyne.io/fyne/v2"

func enableAll(components ...fyne.Disableable) {
	for _, widget := range components {
		widget.Enable()
	}
}

func disableAll(components ...fyne.Disableable) {
	for _, widget := range components {
		widget.Disable()
	}
}
