package tfe

import (
	"context"
	"fmt"

	"github.com/form3tech/f3-tfe/internal/convert"

	"github.com/hashicorp/go-tfe"

	"github.com/form3tech/f3-tfe/internal/log"
)

const (
	defaultPage     = 1
	defaultPageSize = 100
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
	var res []Workspace
	fmt.Fprintf(log.Writer, "listing workspaces for %d orgs\n", len(t.selectedOrganizations))
	for _, org := range t.selectedOrganizations {
		fmt.Fprintf(log.Writer, "listing workspaces for \"%s\"\n", org)
		page := 0
		for {
			l, err := t.OrganizationWorkspaces(ctx, org, page)
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
			page = l.NextPage
		}
	}
	return res, nil
}

func (t *TFE) OrganizationWorkspaces(ctx context.Context, org string, pageInfo ...int) (*tfe.WorkspaceList, error) {
	page := defaultPage
	if len(pageInfo) > 0 {
		page = pageInfo[0]
	}
	pageSize := defaultPageSize
	if len(pageInfo) > 1 {
		pageSize = pageInfo[1]
	}
	result, err := t.client.Workspaces.List(ctx, org, tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: page,
			PageSize:   pageSize,
		},
		Search: convert.StringPtr("development-calliope-txb-gateway"),
	})
	return result, err
}

func (t *TFE) WorkspaceVariables(ctx context.Context, workspaceId string) (*tfe.VariableList, error) {
	vl, err := t.client.Variables.List(ctx, workspaceId, tfe.VariableListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   20,
		},
	})
	return vl, err
}

func (t *TFE) AddVariable(ctx context.Context, workspaceId, key, value string) (*tfe.Variable, error) {
	category := tfe.CategoryTerraform
	return t.client.Variables.Create(ctx, workspaceId, tfe.VariableCreateOptions{
		Key:         convert.StringPtr(key),
		Value:       convert.StringPtr(value),
		Description: convert.StringPtr("Added by utility"),
		Category:    &category,
	})
}
