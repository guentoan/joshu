package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
	"time"
)

func makeRSAUI(w fyne.Window) fyne.CanvasObject {
	keySizeText := canvas.NewText("Key Size", theme.ForegroundColor())
	keySize := widget.NewSelect([]string{"512 bit", "1024 bit", "2048 bit", "4096 bit"}, func(s string) {
	})
	keySize.SetSelectedIndex(1)

	generateText := canvas.NewText("", theme.ForegroundColor())
	generateText.TextStyle = fyne.TextStyle{Italic: true}

	privateKeyTitle := canvas.NewText("Private Key", theme.ForegroundColor())
	publicKeyTitle := canvas.NewText("Public Key", theme.ForegroundColor())

	privateKeyTitle.TextSize = 18
	privateKeyTitle.TextStyle = fyne.TextStyle{Bold: true}
	publicKeyTitle.TextSize = 18
	publicKeyTitle.TextStyle = fyne.TextStyle{Bold: true}

	privateKeyTextBox := widget.NewMultiLineEntry()
	publicKeyTextBox := widget.NewMultiLineEntry()

	privateKeyTextBox.SetPlaceHolder("Private key")
	publicKeyTextBox.SetPlaceHolder("Public key")

	privateKeyTextBox.Wrapping = fyne.TextWrapBreak
	publicKeyTextBox.Wrapping = fyne.TextWrapBreak

	privateKey, publickey, err := generateRSAKeys(1024)
	if err != nil {
		generateText.Text = err.Error()
		generateText.Color = colornames.Red
	} else {
		privateKeyTextBox.SetText(privateKey)
		publicKeyTextBox.SetText(publickey)
	}

	copyPrivateKeyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard := w.Clipboard()
		clipboard.SetContent(privateKeyTextBox.Text)
	})

	copyPublicKeyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard := w.Clipboard()
		clipboard.SetContent(publicKeyTextBox.Text)
	})

	generateButton := widget.NewButton("Generate New Keys", func() {
		startTime := time.Now()
		var privateKey, publickey string
		var err error
		defer func() {
			if err != nil {
				return
			}
			endTime := time.Now()
			totalTime := endTime.Sub(startTime)
			generateText.Text = fmt.Sprintf("Generated in %f seconds", totalTime.Seconds())
			generateText.Color = colornames.Green
			generateText.Refresh()
		}()
		generateText.Text = "Generating new keys..."
		generateText.Color = colornames.Red
		generateText.Refresh()

		selectedIndex := keySize.SelectedIndex()
		bitSize := 1024
		switch selectedIndex {
		case 0:
			bitSize = 512
		case 1:
			bitSize = 1024
		case 2:
			bitSize = 2048
		case 3:
			bitSize = 4096
		}

		privateKey, publickey, err = generateRSAKeys(bitSize)
		if err != nil {
			generateText.Text = err.Error()
			generateText.Color = colornames.Red
		} else {
			privateKeyTextBox.SetText(privateKey)
			publicKeyTextBox.SetText(publickey)
			publicKeyTextBox.Refresh()
			privateKeyTextBox.Refresh()
		}
	})
	generateButton.Importance = widget.HighImportance

	RSAEncryptionText := canvas.NewText("RSA Encryption Test", theme.ForegroundColor())
	RSAEncryptionText.TextSize = 24
	RSAEncryptionText.TextStyle = fyne.TextStyle{Bold: true}

	RSAEncryptionInput := widget.NewMultiLineEntry()
	RSAEncryptionInput.SetPlaceHolder("Enter text to encrypt")
	RSAEncryptionInput.Wrapping = fyne.TextWrapBreak
	RSAEncryptionInput.Text = "This is a test!"
	RSAEncryptionTitle := canvas.NewText("Text to encrypt", theme.ForegroundColor())
	RSAEncryptionTitle.TextSize = 14
	RSAEncryptionTitle.TextStyle = fyne.TextStyle{Bold: true}

	RSAEncryptionOutput := widget.NewMultiLineEntry()
	RSAEncryptionOutput.SetPlaceHolder("Enter text to decrypt")
	RSAEncryptionOutput.Wrapping = fyne.TextWrapBreak
	RSAEncryptionOutputTitle := canvas.NewText("Encrypted", theme.ForegroundColor())
	RSAEncryptionOutputTitle.TextSize = 14
	RSAEncryptionOutputTitle.TextStyle = fyne.TextStyle{Bold: true}

	EncryptButton := widget.NewButton("Encrypt", func() {
		message := RSAEncryptionInput.Text
		if len(message) == 0 {
			return
		}

		publickey := publicKeyTextBox.Text
		if len(publickey) == 0 {
			return
		}

		encrypted, err := encryptRSA(message, publickey)
		if err != nil {
			RSAEncryptionOutput.Text = err.Error()
			RSAEncryptionOutput.Refresh()
			return
		}
		RSAEncryptionOutput.SetText(encrypted)
		RSAEncryptionOutput.Refresh()
		RSAEncryptionInput.Text = ""
		RSAEncryptionInput.Refresh()
	})
	EncryptButton.Importance = widget.HighImportance

	DecryptButton := widget.NewButton("Decrypt", func() {
		encryptedMessage := RSAEncryptionOutput.Text
		if len(encryptedMessage) == 0 {
			return
		}

		privateKey := privateKeyTextBox.Text
		if len(privateKey) == 0 {
			return
		}

		decrypted, err := decryptRSA(encryptedMessage, privateKey)
		if err != nil {
			RSAEncryptionInput.Text = err.Error()
			RSAEncryptionInput.Refresh()
			return
		}
		RSAEncryptionInput.SetText(decrypted)
		RSAEncryptionInput.Refresh()
	})
	DecryptButton.Importance = widget.WarningImportance

	RSATestContainer := container.NewVBox(
		RSAEncryptionText,
		container.NewGridWithColumns(2,
			container.NewBorder(RSAEncryptionTitle, EncryptButton,
				container.NewGridWrap(fyne.NewSize(1, 200), layout.NewSpacer()),
				container.NewGridWrap(fyne.NewSize(1, 200), layout.NewSpacer()),
				RSAEncryptionInput),
			container.NewBorder(RSAEncryptionOutputTitle, DecryptButton,
				container.NewGridWrap(fyne.NewSize(1, 200), layout.NewSpacer()),
				container.NewGridWrap(fyne.NewSize(1, 200), layout.NewSpacer()),
				RSAEncryptionOutput),
		),
	)

	genKeyContainer := container.NewVBox(
		container.NewHBox(keySizeText, keySize, generateButton),
		generateText,
		container.NewGridWithColumns(2,
			container.NewBorder(
				container.NewHBox(privateKeyTitle, copyPrivateKeyButton), nil,
				container.NewGridWrap(fyne.NewSize(1, 300), layout.NewSpacer()),
				container.NewGridWrap(fyne.NewSize(1, 300), layout.NewSpacer()),
				privateKeyTextBox,
			),
			container.NewBorder(
				container.NewHBox(publicKeyTitle, copyPublicKeyButton), nil,
				container.NewGridWrap(fyne.NewSize(1, 300), layout.NewSpacer()),
				container.NewGridWrap(fyne.NewSize(1, 300), layout.NewSpacer()),
				publicKeyTextBox,
			),
		),
	)

	header := makeHeader("RSA Key Generator")
	footer := makeFooter()

	return container.NewPadded(container.NewBorder(header, footer, layout.NewSpacer(), layout.NewSpacer(),
		container.NewGridWithRows(2,
			genKeyContainer,
			RSATestContainer,
		)))
}

func generateRSAKeys(bitSize int) (string, string, error) {
	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return "", "", err
	}

	// Export private key as PEM format
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Export public key as PEM format
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}

// encryptRSA encrypts the given message using the RSA public key
func encryptRSA(message string, publicKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return "", errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(message))
	if err != nil {
		return "", err
	}

	// Encode the encrypted bytes to base64 to return as a string
	encryptedBase64 := base64.StdEncoding.EncodeToString(encryptedBytes)
	return encryptedBase64, nil
}

// decryptRSA decrypts the given encrypted message using the RSA private key
func decryptRSA(encryptedMessage string, privateKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// Decode the base64 encoded encrypted message
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedBytes)
	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}
