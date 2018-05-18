// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package cloud

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/apitype"
	"github.com/pulumi/pulumi/pkg/backend"
	"github.com/pulumi/pulumi/pkg/engine"
	"github.com/pulumi/pulumi/pkg/operations"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/config"
	"github.com/pulumi/pulumi/pkg/resource/deploy"
	"github.com/pulumi/pulumi/pkg/tokens"
	"github.com/pulumi/pulumi/pkg/workspace"
)

// Stack is a cloud stack.  This simply adds some cloud-specific properties atop the standard backend stack interface.
type Stack interface {
	backend.Stack
	CloudURL() string            // the URL to the cloud containing this stack.
	OrgName() string             // the organization that owns this stack.
	CloudName() string           // the PPC in which this stack is running.
	RunLocally() bool            // true if previews/updates/destroys targeting this stack run locally.
	ConsoleURL() (string, error) // the URL to view the stack's information on Pulumi.com
}

// cloudStack is a cloud stack descriptor.
type cloudStack struct {
	name      backend.StackReference // the stack's name.
	cloudURL  string                 // the URL to the cloud containing this stack.
	orgName   string                 // the organization that owns this stack.
	cloudName string                 // the PPC in which this stack is running.
	config    config.Map             // the stack's config bag.
	snapshot  *deploy.Snapshot       // a snapshot representing the latest deployment state.
	b         *cloudBackend          // a pointer to the backend this stack belongs to.
}

type cloudBackendReference struct {
	name  tokens.QName
	owner string
	b     *cloudBackend
}

func (c cloudBackendReference) String() string {
	curUser, err := c.b.client.GetPulumiAccountName(context.Background())
	if err != nil {
		curUser = ""
	}

	if c.owner == curUser {
		return string(c.name)
	}

	return fmt.Sprintf("%s/%s", c.owner, c.name)
}

func (c cloudBackendReference) StackName() tokens.QName {
	return c.name
}

func newStack(apistack apitype.Stack, b *cloudBackend) Stack {
	// Create a fake snapshot out of this stack.
	// TODO[pulumi/pulumi-service#249]: add time, version, etc. info to the manifest.
	stackName := apistack.StackName

	var resources []*resource.State
	for _, res := range apistack.Resources {
		resources = append(resources, resource.NewState(
			res.Type,
			res.URN,
			res.Custom,
			false,
			res.ID,
			resource.NewPropertyMapFromMap(res.Inputs),
			resource.NewPropertyMapFromMap(res.Outputs),
			res.Parent,
			res.Protect,
			// TODO(swgillespie) provide an actual list of dependencies
			[]resource.URN{},
		))
	}
	snapshot := deploy.NewSnapshot(deploy.Manifest{}, resources)

	// Now assemble all the pieces into a stack structure.
	return &cloudStack{
		name: cloudBackendReference{
			owner: apistack.OrgName,
			name:  stackName,
			b:     b,
		},
		cloudURL:  b.CloudURL(),
		orgName:   apistack.OrgName,
		cloudName: apistack.CloudName,
		config:    nil, // TODO[pulumi/pulumi-service#249]: add the config variables.
		snapshot:  snapshot,
		b:         b,
	}
}

// managedCloudName is the name used to refer to the cloud in the Pulumi Service that owns all of an organization's
// managed stacks. All engine operations for a managed stack--previews, updates, destroys, etc.--run locally.
const managedCloudName = "pulumi"

func (s *cloudStack) Name() backend.StackReference { return s.name }
func (s *cloudStack) Config() config.Map           { return s.config }
func (s *cloudStack) Snapshot() *deploy.Snapshot   { return s.snapshot }
func (s *cloudStack) Backend() backend.Backend     { return s.b }
func (s *cloudStack) CloudURL() string             { return s.cloudURL }
func (s *cloudStack) OrgName() string              { return s.orgName }
func (s *cloudStack) CloudName() string            { return s.cloudName }
func (s *cloudStack) RunLocally() bool             { return s.cloudName == managedCloudName }

func (s *cloudStack) Remove(ctx context.Context, force bool) (bool, error) {
	return backend.RemoveStack(ctx, s, force)
}

func (s *cloudStack) Preview(ctx context.Context, proj *workspace.Project, root string, m backend.UpdateMetadata,
	opts backend.UpdateOptions, scopes backend.CancellationScopeSource) (engine.ResourceChanges, error) {
	return backend.PreviewStack(ctx, s, proj, root, m, opts, scopes)
}

func (s *cloudStack) Update(ctx context.Context, proj *workspace.Project, root string, m backend.UpdateMetadata,
	opts backend.UpdateOptions, scopes backend.CancellationScopeSource) (engine.ResourceChanges, error) {
	return backend.UpdateStack(ctx, s, proj, root, m, opts, scopes)
}

func (s *cloudStack) Refresh(ctx context.Context, proj *workspace.Project, root string, m backend.UpdateMetadata,
	opts backend.UpdateOptions, scopes backend.CancellationScopeSource) (engine.ResourceChanges, error) {
	return backend.RefreshStack(ctx, s, proj, root, m, opts, scopes)
}

func (s *cloudStack) Destroy(ctx context.Context, proj *workspace.Project, root string, m backend.UpdateMetadata,
	opts backend.UpdateOptions, scopes backend.CancellationScopeSource) (engine.ResourceChanges, error) {
	return backend.DestroyStack(ctx, s, proj, root, m, opts, scopes)
}

func (s *cloudStack) GetLogs(ctx context.Context, query operations.LogQuery) ([]operations.LogEntry, error) {
	return backend.GetStackLogs(ctx, s, query)
}

func (s *cloudStack) ExportDeployment(ctx context.Context) (*apitype.UntypedDeployment, error) {
	return backend.ExportStackDeployment(ctx, s)
}

func (s *cloudStack) ImportDeployment(ctx context.Context, deployment *apitype.UntypedDeployment) error {
	return backend.ImportStackDeployment(ctx, s, deployment)
}

func (s *cloudStack) ConsoleURL() (string, error) {
	path, err := s.b.StackConsoleURL(s.Name())
	if err != nil {
		return "", nil
	}
	url := s.b.CloudConsoleURL(path)
	if url == "" {
		return "", errors.New("could not determine clould console URL")
	}
	return url, nil
}
