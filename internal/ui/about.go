package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type about struct {
	nameLabel *widget.Label
	app       fyne.App
}

func newAbout(app fyne.App) *about {
	return &about{app: app}
}

func (a *about) buildUI() *fyne.Container {
	a.nameLabel = newBoldLabel("Hermes")

	spacer := &layout.Spacer{}
	return container.NewVBox(
		spacer,
		container.NewHBox(spacer, a.nameLabel, spacer),
	)
}

func (a *about) tabItem() *container.TabItem {
	return &container.TabItem{Text: "About", Icon: theme.InfoIcon(), Content: a.buildUI()}
}
