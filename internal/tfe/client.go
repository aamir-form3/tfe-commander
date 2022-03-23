package tfe

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-tfe"
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

func (t *TFE) Workspaces(ctx context.Context) ([]string, error) {
	res := []string{}
	for org := range t.selectedOrganizations {
		l, err := t.client.Workspaces.List(ctx, org, tfe.WorkspaceListOptions{})
		if err != nil {
			return nil, err
		}

		for _, ws := range l.Items {
			res = append(res, fmt.Sprintf("%s:%s", org, ws.Name))
		}

	}
	return res, nil
}
