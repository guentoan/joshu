package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/image/colornames"
	"strconv"
)

func makeBcryptUI(w fyne.Window) fyne.CanvasObject {
	header := makeHeader("Bcrypt Generator")
	footer := makeFooter()

	EncryptTile := canvas.NewText("Encrypt", theme.ForegroundColor())
	EncryptTile.TextSize = 18
	EncryptTile.TextStyle = fyne.TextStyle{Bold: true}

	EncryptSubTitle := canvas.NewText("Encrypt some text. The result shown will be a Bcrypt encrypted hash.", theme.ForegroundColor())
	EncryptSubTitle.TextSize = 14
	EncryptSubTitle.TextStyle = fyne.TextStyle{Italic: true}

	EncryptInput := widget.NewEntry()
	EncryptInput.SetPlaceHolder("Enter text to encrypt")

	EncryptRound := widget.NewEntry()
	EncryptRound.SetPlaceHolder("Enter number of rounds")
	EncryptRound.Text = "12" // default value

	EncryptOutput := canvas.NewText("", theme.ForegroundColor())
	EncryptOutput.TextSize = 14
	EncryptOutput.TextStyle = fyne.TextStyle{Italic: true}
	EncryptOutput.Color = colornames.Red

	coppyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard := w.Clipboard()
		clipboard.SetContent(EncryptOutput.Text)
	})

	EncryptButton := widget.NewButton("Encrypt", func() {
		EncryptOutput.Text = "Hashing..."
		EncryptOutput.Color = colornames.Red
		EncryptOutput.Refresh()

		password := EncryptInput.Text
		cost, err := parseInt(EncryptRound.Text)
		if err != nil {
			EncryptOutput.Text = err.Error()
			EncryptOutput.Color = colornames.Red
			EncryptOutput.Refresh()
			return
		}

		hash, err := EncryptPassword(password, cost)
		if err != nil {
			EncryptOutput.Text = err.Error()
			EncryptOutput.Color = colornames.Red
			EncryptOutput.Refresh()
			return
		}
		EncryptOutput.Text = hash
		EncryptOutput.Color = colornames.Green
		EncryptOutput.Refresh()
	})
	EncryptButton.Importance = widget.HighImportance

	encryptContent := container.NewVBox(
		EncryptTile,
		EncryptSubTitle,
		container.NewGridWithColumns(2, EncryptInput,
			container.NewGridWithColumns(2,
				EncryptRound,
				EncryptButton,
			),
		),
		container.NewGridWithColumns(2,
			EncryptOutput,
			coppyButton,
		),
	)

	DecryptTile := canvas.NewText("Decrypt", theme.ForegroundColor())
	DecryptTile.TextSize = 18
	DecryptTile.TextStyle = fyne.TextStyle{Bold: true}

	DecryptSubTitle := canvas.NewText("Test your Bcrypt hash against some plaintext, to see if they match.", theme.ForegroundColor())
	DecryptSubTitle.TextSize = 14
	DecryptSubTitle.TextStyle = fyne.TextStyle{Italic: true}

	DecryptHashInput := widget.NewEntry()
	DecryptHashInput.SetPlaceHolder("Hash to check")

	DecryptInput := widget.NewEntry()
	DecryptInput.SetPlaceHolder("String to check against")

	DecryptOutput := canvas.NewText("", theme.ForegroundColor())
	DecryptOutput.TextSize = 14
	DecryptOutput.TextStyle = fyne.TextStyle{Italic: true}

	DecryptButton := widget.NewButton("Decrypt", func() {
		DecryptOutput.Text = "Checking..."
		DecryptOutput.Color = colornames.Red
		DecryptOutput.Refresh()

		hash := DecryptHashInput.Text
		password := DecryptInput.Text
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err != nil {
			DecryptOutput.Text = fmt.Sprintf("Not a match! Error: %v", err)
			DecryptOutput.Color = colornames.Red
			DecryptOutput.Refresh()
			return
		}
		DecryptOutput.Text = "Passwords match"
		DecryptOutput.Color = colornames.Green
		DecryptOutput.Refresh()
	})
	DecryptButton.Importance = widget.HighImportance

	decryptContent := container.NewVBox(
		DecryptTile,
		DecryptSubTitle,
		DecryptHashInput,
		DecryptInput,
		DecryptButton,
		DecryptOutput,
	)

	content := container.NewBorder(header, footer, nil, nil,
		container.NewVBox(encryptContent, decryptContent),
	)
	paddedContent := container.NewPadded(content)
	return paddedContent
}

func showBcryptContent(w fyne.Window) {
	content := makeBcryptUI(w)
	updateContent(w, content)
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func EncryptPassword(password string, cost int) (string, error) {
	// Generate a bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
