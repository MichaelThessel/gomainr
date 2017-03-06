package app

import (
	"github.com/jroimartin/gocui"
)

type keyConfig struct {
	views   *[]string
	key     interface{}
	mod     gocui.Modifier
	handler func(*gocui.Gui, *gocui.View) error
}

// setKeyBindings sets up the keyboard shortcuts
func (a *App) setKeyBindings() error {
	var kc = []keyConfig{
		{
			&selectableViews,
			gocui.KeyCtrlQ,
			gocui.ModNone,
			a.quit,
		},
		{
			&selectableViews,
			gocui.KeyTab,
			gocui.ModNone,
			a.wrapEditor,
		},
		{
			&selectableViews,
			gocui.KeyArrowDown,
			gocui.ModNone,
			a.nextEditor,
		},
		{
			&selectableViews,
			gocui.KeyArrowUp,
			gocui.ModNone,
			a.prevEditor,
		},
		{
			&selectableViews,
			gocui.KeyCtrlSlash,
			gocui.ModNone,
			a.search,
		},
		{
			&selectableViews,
			gocui.KeyCtrlK,
			gocui.ModNone,
			a.scrollUp,
		},
		{
			&selectableViews,
			gocui.KeyCtrlJ,
			gocui.ModNone,
			a.scrollDown,
		},
		{
			&selectableViews,
			gocui.KeyCtrlS,
			gocui.ModNone,
			a.saveModal,
		},
		{
			&selectableViews,
			gocui.KeyCtrlL,
			gocui.ModNone,
			a.loadModal,
		},
		{
			&[]string{viewSave},
			gocui.KeyEnter,
			gocui.ModNone,
			a.save,
		},
		{
			&[]string{viewLoad},
			gocui.KeyEnter,
			gocui.ModNone,
			a.load,
		},
		{
			&[]string{viewSave, viewLoad},
			gocui.KeyCtrlQ,
			gocui.ModNone,
			a.closeModal,
		},
	}

	for _, shortcut := range kc {
		for _, view := range *shortcut.views {
			if err := a.gui.SetKeybinding(
				view,
				shortcut.key,
				shortcut.mod,
				shortcut.handler,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
