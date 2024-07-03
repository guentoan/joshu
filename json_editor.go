package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	jsonrepair "github.com/RealAlexandreAI/json-repair"
)

const (
	oneLineStatus  = 1
	beautifyStatus = 2
)

var status = oneLineStatus

func makeJsonEditorUI(w fyne.Window) fyne.CanvasObject {
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Enter JSON here...")
	input.Wrapping = fyne.TextWrapBreak

	leftToolbar := widget.NewToolbar(
		// one line toolbar
		widget.NewToolbarAction(theme.MenuIcon(), func() {
			status = oneLineStatus
			var oneLineJSON string
			err := oneLineJson(input.Text, &oneLineJSON)
			if err != nil {
				return
			} else {
				input.SetText(oneLineJSON)
			}
			input.Refresh()
		}),
		// beautify toolbar
		widget.NewToolbarAction(theme.ListIcon(), func() {
			status = beautifyStatus
			var prettyJSON string
			err := jsonBeautify(input.Text, &prettyJSON)
			if err != nil {
				return
			} else {
				input.SetText(prettyJSON)
			}
			input.Refresh()
		}),
		widget.NewToolbarSeparator(),
		// copy toolbar
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			w.Clipboard().SetContent(input.Text)
		}),
		// clear toolbar
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {
			input.SetText(w.Clipboard().Content())
			input.Refresh()
		}),
		widget.NewToolbarSeparator(),
		// repair toolbar
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			if len(input.Text) == 0 {
				return
			}
			repairedJson, err := jsonrepair.RepairJSON(input.Text)
			if err != nil {
				return
			}

			if status == oneLineStatus {
				input.SetText(repairedJson)
			} else {
				var prettyJSON string
				err = jsonBeautify(repairedJson, &prettyJSON)
				if err != nil {
					input.SetText(repairedJson)
					input.Refresh()
					return
				}

				input.SetText(prettyJSON)
			}

			input.Refresh()
		}),
		// reset toolbar
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			input.SetText("")
			input.Refresh()
		}),
	)

	contentContainer := container.NewBorder(leftToolbar, nil, nil, nil, input)
	header := makeHeader("Json Editor")
	footer := makeFooter()

	content := container.NewBorder(header, footer, nil, nil, contentContainer)
	paddedContent := container.NewPadded(content)
	return paddedContent
}

func jsonBeautify(input string, output *string) error {
	var jsonObj interface{}
	err := json.Unmarshal([]byte(input), &jsonObj)
	if err != nil {
		return err
	}

	prettyJSONBytes, err := json.MarshalIndent(jsonObj, "", "  ")
	if err != nil {
		return err
	}

	*output = string(prettyJSONBytes)
	return nil
}

func oneLineJson(input string, output *string) error {
	var jsonObj interface{}
	err := json.Unmarshal([]byte(input), &jsonObj)
	if err != nil {
		return err
	}

	prettyJSONBytes, err := json.Marshal(jsonObj)
	if err != nil {
		return err
	}

	*output = string(prettyJSONBytes)
	return nil
}
