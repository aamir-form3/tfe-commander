package main

import (
	"github.com/form3tech/f3-tfe/internal/tfe"
	"github.com/form3tech/f3-tfe/internal/ui"
	"github.com/form3tech/f3-tfe/internal/util"
	"github.com/wvdschel-f3/gocui"
)

var client *tfe.TFE

func main() {
	var err error

	util.Must(tfe.Configure())
	client, err = tfe.New()
	util.Must(err)

	g, err := ui.BuildUI(client)
	util.Must(err)
	defer g.Cleanup()

	err = g.Launch()
	if err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}
