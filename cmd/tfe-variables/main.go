package main

import (
	"context"
	"fmt"

	"github.com/form3tech/f3-tfe/internal/tfe"
	"github.com/form3tech/f3-tfe/internal/util"
)

func main() {
	var err error
	util.Must(tfe.Configure())
	client, err := tfe.New()
	util.Must(err)
	ctx := context.Background()
	orgs, err := client.Organisations(ctx)
	util.Must(err)

	fmt.Printf("Organizations : %+v\n", orgs)

	wss, err := client.OrganizationWorkspaces(ctx, orgs[5], tfe.NewSearchInfo())
	util.Must(err)

	fmt.Printf("Workspaces : %+v\n", wss.Items[0].ID)

	wvs, err := client.WorkspaceVariables(ctx, wss.Items[0].ID)
	util.Must(err)

	fmt.Printf("Workspaces variables : %+v\n", wvs.Items[0].Value)

	nwv, err := client.AddVariable(ctx, wss.Items[0].ID, tfe.NewVariableOption("test-key", "test-value"))
	util.Must(err)
	fmt.Printf("Workspaces variables : %+v\n", nwv.ID)

}
