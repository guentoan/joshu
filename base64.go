package main

import (
	"encoding/base64"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/url"
	"time"
)

func makeFooter() fyne.CanvasObject {
	now := time.Now()
	footerText := fmt.Sprintf("Copyright © %d ToanNV. All rights reserved.", now.Year())
	u, _ := url.Parse("https://toannv.me")
	footer := widget.NewHyperlinkWithStyle(footerText, u, fyne.TextAlignCenter, fyne.TextStyle{})

	return footer
}

func makeHeader(title string) fyne.CanvasObject {
	header := canvas.NewText(title, theme.PrimaryColor())
	header.TextSize = 24
	header.Alignment = fyne.TextAlignCenter

	return header
}

func makeBase64UI(w fyne.Window) fyne.CanvasObject {
	header := makeHeader("Base64 Encoder/Decoder")
	footer := makeFooter()

	input := widget.NewEntry()
	input.MultiLine = true
	input.Wrapping = fyne.TextWrapBreak
	input.SetPlaceHolder("Input Text Or Read from Clipboard")

	output := widget.NewEntry()
	output.MultiLine = true
	output.Wrapping = fyne.TextWrapBreak
	output.SetPlaceHolder("Output Result")

	encodeButton := widget.NewButtonWithIcon("Encode", theme.MediaFastForwardIcon(), func() {
		if input.Text == "" {
			input.Text = w.Clipboard().Content()
			input.Refresh()
		}
		out := base64.StdEncoding.EncodeToString([]byte(input.Text))
		output.Text = out
		output.Refresh()
	})
	encodeButton.Importance = widget.HighImportance

	clearButton := widget.NewButtonWithIcon("Clear", theme.ContentClearIcon(), func() {
		output.Text = ""
		output.Refresh()
		input.Text = ""
		input.Refresh()
	})
	clearButton.Importance = widget.MediumImportance

	decodeButton := widget.NewButtonWithIcon("Decode", theme.MediaFastRewindIcon(), func() {
		if input.Text == "" {
			input.Text = w.Clipboard().Content()
			input.Refresh()
		}
		out, err := base64.StdEncoding.DecodeString(input.Text)
		if err == nil {
			output.Text = string(out)
		} else {
			output.Text = err.Error()
		}
		output.Text = string(out)
		output.Refresh()
	})
	decodeButton.Importance = widget.HighImportance

	copyButton := widget.NewButtonWithIcon("Copy to Clipboard", theme.ContentCopyIcon(), func() {
		clipboard := w.Clipboard()
		clipboard.SetContent(output.Text)
	})
	copyButton.Importance = widget.WarningImportance

	content := container.NewBorder(header, footer, nil, nil,
		container.NewGridWithRows(2,
			container.NewBorder(
				nil,
				container.NewGridWithColumns(4, encodeButton, decodeButton, copyButton, clearButton),
				nil,
				nil,
				input),
			output,
		),
	)
	paddedContent := container.NewPadded(content)

	return paddedContent
}
