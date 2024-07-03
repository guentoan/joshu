package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	a := app.New()
	w := a.NewWindow("Multi-functional Tool")

	tabs := container.NewAppTabs(
		container.NewTabItem("Base 64", makeBase64UI(w)),
		container.NewTabItem("Password Generator", makeRandomPasswordUI(w)),
		container.NewTabItem("Bcrypt", makeBcryptUI(w)),
	)

	tabs.SetTabLocation(container.TabLocationLeading)
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1300, 750))
	w.ShowAndRun()
}
