package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Create(app fyne.App) *container.AppTabs {
	return &container.AppTabs{Items: []*container.TabItem{
		newAbout(app).tabItem(),
	}}
}
