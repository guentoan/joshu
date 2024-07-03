package main

import (
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func showJsonContent(w fyne.Window) {
	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Enter JSON here...")

	output := widget.NewMultiLineEntry()

	beautifyButton := widget.NewButton("Beautify", func() {
		var prettyJSON string
		err := jsonBeautify(input.Text, &prettyJSON)
		if err != nil {
			output.SetText(fmt.Sprintf("Error: %v", err))
		} else {
			output.SetText(prettyJSON)
		}
	})

	container.NewVBox(input, beautifyButton, output)
}

func jsonBeautify(input string, output *string) error {
	var jsonObj map[string]interface{}
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
