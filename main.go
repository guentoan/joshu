package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	a := app.New()
	w := a.NewWindow("助手 - Developer's Assistant")

	tabs := container.NewAppTabs(
		container.NewTabItem("Json Editor", makeJsonEditorUI(w)),
		container.NewTabItem("Base 64", makeBase64UI(w)),
		container.NewTabItem("Password Generator", makeRandomPasswordUI(w)),
		container.NewTabItem("Bcrypt", makeBcryptUI(w)),
		container.NewTabItem("RSA Generator", makeRSAUI(w)),
	)

	tabs.SetTabLocation(container.TabLocationLeading)
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1500, 850))
	w.ShowAndRun()
}
