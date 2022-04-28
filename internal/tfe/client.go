package tfe

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-tfe"

	"github.com/form3tech/f3-tfe/internal/log"
)

func (t *TFE) Organisations(ctx context.Context) ([]string, error) {
	l, err := t.client.Organizations.List(ctx, tfe.OrganizationListOptions{})
	if err != nil {
		return nil, err
	}

	res := make([]string, len(l.Items))
	for i, org := range l.Items {
		res[i] = org.Name
	}
	return res, nil
}

func (t *TFE) Workspaces(ctx context.Context) ([]Workspace, error) {
	res := []Workspace{}
	fmt.Fprintf(log.Writer, "listing workspaces for %d orgs\n", len(t.selectedOrganizations))
	for _, org := range t.selectedOrganizations {
		fmt.Fprintf(log.Writer, "listing workspaces for \"%s\"\n", org)
		i := 0
		for {
			l, err := t.client.Workspaces.List(ctx, org, tfe.WorkspaceListOptions{
				ListOptions: tfe.ListOptions{
					PageNumber: i,
					PageSize:   100,
				},
			})
			if err != nil {
				return nil, err
			}

			for _, ws := range l.Items {
				res = append(res, Workspace{
					OrganizationName: org,
					Name:             ws.Name,
				})
			}
			if l.NextPage <= 0 {
				break
			}
			i = l.NextPage
		}
	}
	return res, nil
}
