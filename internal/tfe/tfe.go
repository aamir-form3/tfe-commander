package tfe

import (
	"fmt"

	"github.com/hashicorp/go-tfe"
)

type TFE struct {
	client                *tfe.Client
	selectedOrganizations map[string]bool
	selectedWorkspaces    map[Workspace]bool
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
	res := make([]Workspace, 0, len(t.selectedWorkspaces))
	for ws := range t.selectedWorkspaces {
		res = append(res, ws)
	}
	return res
}

func (t *TFE) SelectWorkspace(ws Workspace) {
	t.selectedWorkspaces[ws] = true
}

func (t *TFE) DeselectWorkspace(ws Workspace) {
	delete(t.selectedWorkspaces, ws)
}

func (t *TFE) IsWorkspaceSelected(ws Workspace) bool {
	return t.selectedWorkspaces[ws]
}

func (t *TFE) SelectedOrganizations() []string {
	orgs := make([]string, 0, len(t.selectedOrganizations))
	for org := range t.selectedOrganizations {
		orgs = append(orgs, org)
	}
	return orgs
}

func (t *TFE) SelectOrganization(org string) {
	t.selectedOrganizations[org] = true
}

func (t *TFE) DeselectOrganization(org string) {
	delete(t.selectedOrganizations, org)
}

func (t *TFE) IsOrganizationSelected(org string) bool {
	return t.selectedOrganizations[org]
}
