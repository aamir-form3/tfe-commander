package component

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type ListView struct {
	view              *gocui.View
	items             []string
	multiSelect       bool
	selection         map[int]bool
	position          int
	selectionCallback func([]string)
}

func NewListView(gui *gocui.Gui, view *gocui.View, multiSelect bool, selectionCallback func([]string)) (*ListView, error) {
	res := &ListView{
		view:              view,
		multiSelect:       multiSelect,
		selectionCallback: selectionCallback,
	}

	err := gui.SetKeybinding(view.Name(), gocui.KeyArrowDown, gocui.ModNone, res.down)
	if err != nil {
		return nil, err
	}
	err = gui.SetKeybinding(view.Name(), gocui.KeyArrowUp, gocui.ModNone, res.up)
	if err != nil {
		return nil, err
	}
	err = gui.SetKeybinding(view.Name(), gocui.KeyCtrlA, gocui.ModNone, res.selectAll)
	if err != nil {
		return nil, err
	}
	err = gui.SetKeybinding(view.Name(), gocui.KeyEsc, gocui.ModNone, res.clearSelection)
	if err != nil {
		return nil, err
	}
	err = gui.SetKeybinding(view.Name(), gocui.KeySpace, gocui.ModNone, res.space)
	if err != nil {
		return nil, err
	}
	err = gui.SetKeybinding(view.Name(), gocui.KeyEnter, gocui.ModNone, res.enter)
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
		if l.selection[idx] {
			fmt.Fprintf(l.view, "\x1b[4m")
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
	l.selection = map[int]bool{}
	l.selectionUpdated()
}

func (l *ListView) GetSelection() []string {
	res := make([]string, 0, len(l.selection))
	for idx, val := range l.selection {
		if val {
			res = append(res, l.items[idx])
		}
	}
	return res
}

func (l *ListView) selectCurrent() {
	if !l.selection[l.position] {
		l.ToggleCurrent()
	}
}

func (l *ListView) ToggleCurrent() {
	defer l.selectionUpdated()
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

func (l *ListView) selectAll(gui *gocui.Gui, view *gocui.View) error {
	if l.multiSelect {
		for i, _ := range l.items {
			l.selection[i] = true
		}
		l.selectionUpdated()
	}
	return nil
}

func (l *ListView) clearSelection(gui *gocui.Gui, view *gocui.View) error {
	l.selection = map[int]bool{}
	l.selectionUpdated()
	l.UpdateBuffer()
	return nil
}

func (l *ListView) space(gui *gocui.Gui, view *gocui.View) error {
	l.ToggleCurrent()
	return nil
}

func (l *ListView) enter(gui *gocui.Gui, view *gocui.View) error {
	l.selection = map[int]bool{
		l.position: true,
	}
	l.selectionUpdated()
	return nil
}

func (l *ListView) selectionUpdated() {
	if l.selectionCallback != nil {
		selection := make([]string, 0)
		for itemIdx := range l.selection {
			selection = append(selection, l.items[itemIdx])
		}
		l.selectionCallback(selection)
	}
	l.UpdateBuffer()
}
