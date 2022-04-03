package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/widget"
)

type LabelWriter struct {
	label *widget.Label

	timeoutMilliseconds int
	lastLogTime         time.Time
}

func (lw *LabelWriter) Write(p []byte) (n int, err error) {
	if lw.label == nil {
		return 0, fmt.Errorf("label cannot be nil")
	}

	time.Sleep(time.Until(lw.lastLogTime.Add(time.Millisecond * time.Duration(lw.timeoutMilliseconds))))
	lw.lastLogTime = time.Now()

	lw.label.SetText(
		lw.label.Text + string(p),
	)

	return len(p), nil
}
