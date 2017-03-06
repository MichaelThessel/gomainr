package app

import (
	"strings"

	"github.com/jroimartin/gocui"
)

// quit handles quit keyboard shortcut
func (a *App) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// nextEditor focusses next editor (no-wrap)
func (a *App) nextEditor(g *gocui.Gui, v *gocui.View) error {
	return a.switchEditor(true, false)
}

// nextEditor focusses previous editor (no-wrap)
func (a *App) prevEditor(g *gocui.Gui, v *gocui.View) error {
	return a.switchEditor(false, false)
}

// nextEditor focusses next editor (wrap)
func (a *App) wrapEditor(g *gocui.Gui, v *gocui.View) error {
	return a.switchEditor(true, true)
}

// switchEditor changes focus to next/previous editor
func (a *App) switchEditor(forward bool, wrap bool) error {
	var index int

	if forward {
		index = a.currentView + 1
		if index > len(selectableViews)-1 {
			if wrap {
				index = 0
			} else {
				return nil
			}
		}
	} else {
		index = a.currentView - 1
		if index < 0 {
			if wrap {
				index = len(selectableViews) - 1
			} else {
				return nil
			}
		}
	}

	return a.setSelectableView(index)
}

// scrollUp scrolls the result list up
func (a *App) scrollUp(g *gocui.Gui, v *gocui.View) error {
	list, _ := a.gui.View(viewDomain)
	x, y := list.Origin()
	y--
	if y >= 0 {
		list.SetOrigin(x, y)
	}

	return nil
}

// scrollUp scrolls the result list down
func (a *App) scrollDown(g *gocui.Gui, v *gocui.View) error {
	list, _ := a.gui.View(viewDomain)
	x, y := list.Origin()
	y++
	_, bMax := list.Size()
	max := len(strings.Split(list.Buffer(), "\n")) - bMax - 1
	if y <= max {
		list.SetOrigin(x, y)
	}

	return nil
}
