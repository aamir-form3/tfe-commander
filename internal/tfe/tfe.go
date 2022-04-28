package tfe

import (
	"fmt"

	"github.com/hashicorp/go-tfe"
)

type TFE struct {
	client                *tfe.Client
	selectedOrganizations []string
	selectedWorkspaces    []Workspace
}

type Workspace struct {
	OrganizationName, Name string
}

func New() (*TFE, error) {
	tfeClient, err := tfe.NewClient(&tfe.Config{
		Token:   token,
		Address: fmt.Sprintf("https://%s", domain),
	})

	return &TFE{
		client: tfeClient,
	}, err
}

func (t *TFE) SelectedWorkspaces() []Workspace {
	res := make([]Workspace, len(t.selectedWorkspaces))
	copy(res, t.selectedWorkspaces)
	return res
}

func (t *TFE) SelectWorkspaces(ws []Workspace) {
	t.selectedWorkspaces = make([]Workspace, len(ws))
	copy(t.selectedWorkspaces, ws)
}

func (t *TFE) SelectedOrganizations() []string {
	orgs := make([]string, len(t.selectedOrganizations))
	copy(orgs, t.selectedOrganizations)
	return orgs
}

func (t *TFE) SelectOrganizations(orgs []string) {
	t.selectedOrganizations = make([]string, len(orgs))
	copy(t.selectedOrganizations, orgs)
}
