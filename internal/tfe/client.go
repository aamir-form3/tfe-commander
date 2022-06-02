package tfe

import (
	"context"
	"errors"
	"fmt"

	"github.com/form3tech/f3-tfe/internal/convert"

	"github.com/hashicorp/go-tfe"

	"github.com/form3tech/f3-tfe/internal/log"
)

const (
	defaultPage     = 1
	defaultPageSize = 100
)

type SearchInfo struct {
	Page     int
	PageSize int
	Tags     *string
	Search   *string
	Include  *string
}

func NewSearchInfo() SearchInfo {
	return SearchInfo{
		Page:     defaultPage,
		PageSize: defaultPageSize,
	}
}

func NewVariableOption(key, value string) tfe.VariableCreateOptions {
	defaultCategory := tfe.CategoryTerraform
	return tfe.VariableCreateOptions{
		Key:         &key,
		Value:       &value,
		Category:    &defaultCategory,
		Description: convert.StringPtr("Added by TFE utility"),
	}
}

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
		search := NewSearchInfo()
		for {
			search.Page = page
			l, err := t.OrganizationWorkspaces(ctx, org, search)
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

func (t *TFE) OrganizationWorkspaces(ctx context.Context, org string, search SearchInfo) (*tfe.WorkspaceList, error) {
	page := defaultPage
	if search.Page > 0 {
		page = search.Page
	}
	pageSize := defaultPageSize
	if search.PageSize > 0 {
		pageSize = search.PageSize
	}
	result, err := t.client.Workspaces.List(ctx, org, tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: page,
			PageSize:   pageSize,
		},
		Tags:    search.Tags,
		Search:  search.Search,
		Include: search.Include,
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

func (t *TFE) AddVariable(ctx context.Context, workspaceId string, options tfe.VariableCreateOptions) (*tfe.Variable, error) {
	return t.client.Variables.Create(ctx, workspaceId, options)
}

func (t *TFE) UpdateVariable(ctx context.Context, workspaceId, variableId string, options tfe.VariableUpdateOptions) (*tfe.Variable, error) {
	return t.client.Variables.Update(ctx, workspaceId, variableId, options)
}
func (t *TFE) AddVariableToWorkspaces(ctx context.Context, workspaceId string, options tfe.VariableCreateOptions) (*tfe.Variable, error) {
	return t.client.Variables.Create(ctx, workspaceId, options)
}

func (t *TFE) GetVariableWithKey(ctx context.Context, workspaceId, key string) (*tfe.Variable, error) {
	options := tfe.VariableListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: defaultPage,
			PageSize:   defaultPageSize,
		},
	}
	page := 0
	for {
		options.PageNumber = page
		vl, err := t.client.Variables.List(ctx, workspaceId, options)
		if err != nil {
			return nil, err
		}
		for _, item := range vl.Items {
			if item.Key == key {
				return item, nil
			}
		}
		if vl.NextPage <= 0 {
			return nil, errors.New("not found with key")
		}
		page = vl.NextPage
	}
}
