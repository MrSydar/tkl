package main

import (
	"log"
	"mrsydar/tkl/k360/client"
	"mrsydar/tkl/process"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	window := app.New().NewWindow("tkl")

	var csvPath string

	csvFilePathLabel := widget.NewLabel(textSelectedCsvFile)
	csvFileDialog := dialog.NewFileOpen(
		func(uri fyne.URIReadCloser, err error) {
			if err != nil {
				log.Println("Error: ", err)
			} else {
				csvPath = uri.URI().Path()
				csvFilePathLabel.SetText(textSelectedCsvFile + csvPath)
			}
		},
		window,
	)

	apiIdInput := widget.NewEntry()
	apiIdInput.SetPlaceHolder("API ID")

	apiKeyInput := widget.NewEntry()
	apiKeyInput.SetPlaceHolder("API Key")

	csvFileChooseButton := widget.NewButton(textChooseCsvFile, func() {
		csvFileDialog.Show()
	})

	progressBar := NewProgressBarWithMessage()

	runButton := widget.NewButton(textRun, nil)
	runButton.OnTapped = func() {
		k360Client := client.New(apiIdInput.Text, apiKeyInput.Text)
		go func() {
			disableAll(csvFileChooseButton, runButton, apiIdInput, apiKeyInput)

			process.ProcessInvoices(
				*k360Client,
				csvPath,
				func(message string, recordsNumber, currentRecord int) {
					progressBar.Update(message, float64(currentRecord)/float64(recordsNumber))
				},
			)

			enableAll(csvFileChooseButton, runButton, apiIdInput, apiKeyInput)
		}()
	}

	logFile, err := os.Create("output.log")
	if err != nil {
		log.Fatalf("can't create/truncate errors.log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	content := container.New(layout.NewVBoxLayout(),
		apiIdInput,
		apiKeyInput,
		csvFilePathLabel,
		csvFileChooseButton,
		progressBar,
		runButton,
	)

	window.SetContent(content)

	window.ShowAndRun()
}
