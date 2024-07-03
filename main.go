package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Multi-functional Tool")

	// Create the menu items
	menu := container.NewVBox(
		widget.NewButton("Json", func() {
			showJsonContent(w)
		}),
		widget.NewButton("Base 64", func() {
			showBase64Content(w)
		}),
		widget.NewButton("Password Generator", func() {
			showRandomPasswordContent(w)
		}),
	)
	// Wrap menu in padding
	paddedMenu := container.NewPadded(menu)

	// Initial content
	initialContent := makeBase64UI(w)

	// Create the main layout
	mainContent := container.NewHSplit(paddedMenu, initialContent)
	mainContent.Offset = 0.15 // Adjust the split ratio

	w.SetContent(mainContent)
	w.Resize(fyne.NewSize(1200, 650))
	w.ShowAndRun()
}

func updateContent(w fyne.Window, newContent fyne.CanvasObject) {
	// Retrieve the current split container
	split := w.Content().(*container.Split)
	split.Trailing = newContent
	w.Content().Refresh()
}
