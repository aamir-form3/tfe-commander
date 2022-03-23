package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type ListView struct {
	view        *gocui.View
	items       []string
	multiSelect bool
	selection   map[int]bool
	position    int
}

func NewListView(gui *gocui.Gui, view *gocui.View, multiSelect bool) (*ListView, error) {
	res := &ListView{
		view:        view,
		multiSelect: multiSelect,
	}

	err := gui.SetKeybinding(view.Name(), gocui.KeyArrowDown, gocui.ModNone, res.down)
	if err != nil {
		return nil, err
	}
	err = gui.SetKeybinding(view.Name(), gocui.KeyArrowUp, gocui.ModNone, res.up)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (l *ListView) UpdateBuffer() {
	l.view.Clear()
	for idx, item := range l.items {
		if idx == l.position {
			fmt.Fprintf(l.view, "\033[3%d;4%dm", l.view.BgColor-1, l.view.FgColor-1)
		}
		fmt.Fprintln(l.view, item)
		fmt.Fprintf(l.view, "\033[3%d;4%dm", l.view.FgColor-1, l.view.BgColor-1)
	}
}

func (l *ListView) SetItems(items []string) {
	l.items = items
	if l.position > len(l.items) {
		l.position = len(l.items) - 1
	}
	l.UpdateBuffer()
}

func (l *ListView) SelectCurrent() {
	if !l.selection[l.position] {
		l.ToggleCurrent()
	}
}

func (l *ListView) ToggleCurrent() {
	if l.position < len(l.items) {
		for itemIdx := range l.selection {
			if l.position == itemIdx {
				delete(l.selection, itemIdx)
				return
			}
		}
		if !l.multiSelect {
			l.selection = map[int]bool{l.position: true}
		} else {
			l.selection[l.position] = true
		}
	}
}

func (l *ListView) down(gui *gocui.Gui, view *gocui.View) error {
	if l.position < len(l.items)-1 {
		l.position += 1
	}
	l.UpdateBuffer()
	return nil
}

func (l *ListView) up(gui *gocui.Gui, view *gocui.View) error {
	if l.position > 0 {
		l.position -= 1
	}
	l.UpdateBuffer()
	return nil
}
