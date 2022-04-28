package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/form3tech/f3-tfe/internal/log"
	"github.com/form3tech/f3-tfe/internal/tfe"
	"github.com/form3tech/f3-tfe/internal/ui/component"
	"github.com/jroimartin/gocui"
)

type UI struct {
	gui                       *gocui.Gui
	tfeClient                 *tfe.TFE
	organisations, workspaces *component.ListView
	logsView, plansView       *gocui.View
	showLog                   bool
}

const (
	organisationsViewName = "organisations"
	workspacesViewName    = "workspaces"
	plansViewName         = "plans"
	logsViewName          = "logs"
)

func BuildUI(tfeClient *tfe.TFE) (*UI, error) {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	return &UI{
		gui:       gui,
		tfeClient: tfeClient,
		showLog:   true,
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

	if err := ui.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (ui *UI) layout(g *gocui.Gui) error {
	var err error

	if ui.gui.CurrentView() == nil {
		_, _ = ui.gui.SetCurrentView(organisationsViewName)
	}

	maxX, maxY := g.Size()
	g.FgColor = gocui.ColorGreen
	g.BgColor = gocui.ColorBlack

	v, err := g.SetView(organisationsViewName, 0, 0, 24, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ui.organisations, err = component.NewListView(ui.gui, v, false, ui.organisationsSelectedUpdated)
		if err != nil {
			return err
		}
		orgs, err := ui.tfeClient.Organisations(context.Background())
		if err != nil {
			return err
		}
		ui.organisations.SetItems(orgs)
	}

	v, err = g.SetView(workspacesViewName, 25, 0, 60, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ui.workspaces, err = component.NewListView(ui.gui, v, true, ui.workspacesSelectedUpdated)
		if err != nil {
			return err
		}
	}

	logsHeight := 0
	if ui.showLog == true {
		logsHeight = 8
	}
	v, err = g.SetView(plansViewName, 61, 0, maxX-1, maxY-1-logsHeight)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		ui.plansView = v
		v.Autoscroll = true
	}

	if ui.showLog == true {
		v, err = g.SetView(logsViewName, 61, maxY-logsHeight, maxX-1, maxY-1)
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			log.Writer = v
			ui.logsView = v
			v.Autoscroll = true
		}
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
	currentView := g.CurrentView()
	currentViewIdx := -1
	allViews := g.Views()
	for i, view := range allViews {
		if view == currentView {
			ui.log("Current view %d %s", i, view.Name())
			currentViewIdx = i
			break
		}
	}
	currentViewIdx++
	if currentViewIdx >= len(allViews) {
		currentViewIdx = 0
	}
	_, err := g.SetCurrentView(allViews[currentViewIdx].Name())
	ui.log("Setting view %d %s ~> %v", currentViewIdx, allViews[currentViewIdx].Name(), err)
	return err
}

func (ui *UI) log(fmt_ string, args ...interface{}) {
	if ui.logsView != nil {
		ui.showLog = true
		fmt.Fprintf(ui.logsView, fmt_, args...)
		fmt.Fprintln(ui.logsView)
	}
}

func (ui *UI) organisationsSelectedUpdated(orgs []string) {
	ui.tfeClient.SelectOrganizations(orgs)
	wss, err := ui.tfeClient.Workspaces(context.Background())
	if err != nil {
		ui.log("Failed to fetch workspaces: %v", err)
	}
	if ui.workspaces != nil {
		items := make([]string, len(wss))
		for i, ws := range wss {
			items[i] = ws.Name
		}
		ui.workspaces.SetItems(items)
	}
}

func (ui *UI) workspacesSelectedUpdated(wss []string) {
}
