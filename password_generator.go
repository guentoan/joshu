package main

import (
	"crypto/rand"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"math/big"
)

const (
	lowerCase = "abcdefghijklmnopqrstuvwxyz"
	upperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers   = "1234567890"
	special   = "`~!@#$%^&*()-=_+[]{}|;':\",./<>?"
	hexChars  = "123456789ABCDEF"
)

var passwordTypes = map[string]string{
	"Memorable Passwords":         "memorable_pwd",
	"Strong Passwords":            "strong_pwd",
	"Fort Knox Passwords":         "ft_knox_pwd",
	"CodeIgniter Encryption Keys": "ci_key",
	"160-bit WPA Key":             "160_wpa",
	"504-bit WPA Key":             "504_wpa",
	"64-bit WEP Keys":             "64_wep",
	"128-bit WEP Keys":            "128_wep",
	"152-bit WEP Keys":            "152_wep",
	"256-bit WEP Keys":            "256_wep",
}

var supportPasswordTypes = []string{
	"Memorable Passwords",
	"Strong Passwords",
	"Fort Knox Passwords",
	"CodeIgniter Encryption Keys",
	"160-bit WPA Key",
	"504-bit WPA Key",
	"64-bit WEP Keys",
	"128-bit WEP Keys",
	"152-bit WEP Keys",
	"256-bit WEP Keys",
}

func makeRandomPasswordUI(w fyne.Window) fyne.CanvasObject {
	header := makeHeader("Password Generator")
	footer := makeFooter()

	generateButton := widget.NewButton("Generate Passwords and Keys", func() {
		showRandomPasswordContent(w)
	})

	passwordBlock := makePasswordUI(w)

	content := container.NewBorder(header, footer, nil, nil,
		container.NewVBox(generateButton, passwordBlock))
	paddedContent := container.NewPadded(content)
	scrollable := container.NewVScroll(paddedContent)
	return scrollable
}
func showRandomPasswordContent(w fyne.Window) {
	content := makeRandomPasswordUI(w)
	updateContent(w, content)
}

func generateRandomPasswords(count int, passType string) []string {
	var passwords []string
	for i := 0; i < count; i++ {
		pt := passwordTypes[passType]
		key, err := GetKey(pt)
		if err != nil {
			continue
		}
		passwords = append(passwords, key)
	}
	return passwords
}

func makePasswordUI(w fyne.Window) fyne.CanvasObject {
	content := container.NewVBox()
	// Generate 3 random passwords for each type
	for _, pt := range supportPasswordTypes {
		passwords := generateRandomPasswords(3, pt)

		blockTitle := canvas.NewText(pt, theme.ForegroundColor())
		blockTitle.TextSize = 14
		blockTitle.TextStyle = fyne.TextStyle{Bold: true}

		var passBlocks []fyne.CanvasObject
		for _, pw := range passwords {
			passwordText := widget.NewLabel(pw)
			passwordText.Wrapping = fyne.TextWrapBreak
			passwordText.Alignment = fyne.TextAlignCenter

			// Copy button
			coppyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
				clipboard := w.Clipboard()
				clipboard.SetContent(pw)
			})
			coppyButton.Resize(fyne.NewSize(20, 20))

			passwordTextContainer := container.NewVBox(passwordText, coppyButton)
			passBlocks = append(passBlocks, passwordTextContainer)
		}

		blockContent := container.NewVBox(
			blockTitle,
			container.NewGridWithColumns(3, passBlocks...),
		)

		paddedBlockContent := container.NewPadded(blockContent)

		content.Add(paddedBlockContent)
	}

	paddedContent := container.NewPadded(content)
	return paddedContent
}

// Random returns a random float64 between 0 and 1.
func Random() float64 {
	maxNumber := big.NewInt(1 << 53)
	n, _ := rand.Int(rand.Reader, maxNumber)
	return float64(n.Int64()) / (1 << 53)
}

// KeyGen generates a random key of the specified length.
func KeyGen(length int, useLowerCase, useUpperCase, useNumbers, useSpecial, useHex bool) string {
	var chars string
	var key string

	if useLowerCase {
		chars += lowerCase
	}
	if useUpperCase {
		chars += upperCase
	}
	if useNumbers {
		chars += numbers
	}
	if useSpecial {
		chars += special
	}
	if useHex {
		chars += hexChars
	}

	for i := 0; i < length; i++ {
		index := int(Random() * float64(len(chars)))
		key += string(chars[index])
	}

	return key
}

// GetKey returns a key based on the strength specified.
func GetKey(strength string) (string, error) {
	switch strength {
	case "memorable_pwd":
		return KeyGen(10, true, true, true, false, false), nil
	case "strong_pwd":
		return KeyGen(15, true, true, true, true, false), nil
	case "ft_knox_pwd":
		return KeyGen(30, true, true, true, true, false), nil
	case "ci_key":
		return KeyGen(32, true, true, true, false, false), nil
	case "160_wpa":
		return KeyGen(20, true, true, true, true, false), nil
	case "504_wpa":
		return KeyGen(63, true, true, true, true, false), nil
	case "64_wep":
		return KeyGen(5, false, false, false, false, true), nil
	case "128_wep":
		return KeyGen(13, false, false, false, false, true), nil
	case "152_wep":
		return KeyGen(16, false, false, false, false, true), nil
	case "256_wep":
		return KeyGen(29, false, false, false, false, true), nil
	default:
		return "", fmt.Errorf("no such strength \"%s\"", strength)
	}
}
