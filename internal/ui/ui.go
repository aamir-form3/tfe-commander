package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/form3tech/f3-tfe/internal/tfe"
	"github.com/jroimartin/gocui"
)

type UI struct {
	gui                       *gocui.Gui
	tfeClient                 *tfe.TFE
	organisations, workspaces *ListView
}

const (
	organisationsViewName = "organisations"
	workspacesViewName    = "workspaces"
	plansViewName         = "plans"
)

func BuildUI(tfeClient *tfe.TFE) (*UI, error) {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	return &UI{
		gui:       gui,
		tfeClient: tfeClient,
	}, nil
}

func (ui *UI) Cleanup() {
	ui.gui.Close()
	*ui = UI{} // wipe pointers to prevent further use
}

func (ui *UI) Launch() error {
	ui.gui.SetManagerFunc(ui.layout)

	if err := ui.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
		return err
	}

	if err := ui.gui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, ui.nextView); err != nil {
		return err
	}

	_, _ = ui.gui.SetCurrentView(organisationsViewName)

	if err := ui.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (ui *UI) layout(g *gocui.Gui) error {
	var err error

	maxX, maxY := g.Size()
	g.FgColor = gocui.ColorGreen
	g.BgColor = gocui.ColorBlack

	v, err := g.SetView(organisationsViewName, 1, 1, 24, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ui.organisations, err = NewListView(ui.gui, v, false)
		if err != nil {
			return err
		}
		orgs, err := ui.tfeClient.Organisations(context.Background())
		if err != nil {
			return err
		}
		ui.organisations.SetItems(orgs)
	}

	v, err = g.SetView(workspacesViewName, 25, 1, 48, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ui.workspaces, err = NewListView(ui.gui, v, false)
		if err != nil {
			return err
		}

		items, err := ui.tfeClient.Workspaces(context.Background())
		if err != nil {
			return err
		}
		ui.workspaces.SetItems(items)
	}

	v, err = g.SetView(plansViewName, 49, 1, maxX-1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
	}

	for _, v := range g.Views() {
		v.Title = strings.Title(v.Name())
	}

	// Set the color to use for window borders & titles
	g.FgColor = gocui.ColorWhite
	return nil
}

func (ui *UI) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (ui *UI) nextView(g *gocui.Gui, v *gocui.View) error {
	plansView, _ := ui.gui.View(plansViewName)
	currentView := g.CurrentView()
	currentViewIdx := -1
	allViews := g.Views()
	for i, view := range allViews {
		if view == currentView {
			fmt.Fprintf(plansView, "Current view %d %s\n", i, view.Name())
			currentViewIdx = i
			break
		}
	}
	currentViewIdx++
	if currentViewIdx >= len(allViews) {
		currentViewIdx = 0
	}
	_, err := g.SetCurrentView(allViews[currentViewIdx].Name())
	fmt.Fprintf(plansView, "Setting view %d %s ~> %v\n", currentViewIdx, allViews[currentViewIdx].Name(), err)
	return err
}
