// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xc

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/drone/go-scm/scm"
)

type repositories struct {
	Data     []*repository `json:"data"`
	NextPage int           `json:"nextPage"`
}

type repository struct {
	Id            int64     `json:"id"`
	ProjectName   string    `json:"projectName"`
	Created       time.Time `json:"created"`
	Author        string    `json:"author"`
	Private       int       `json:"private"`
	GroupName     string    `json:"groupName"`
	DefaultBranch string    `json:"defaultBranch"`
	AccessLevel   access    `json:"access"`
}

type access struct {
	GroupAccessLevel   int64 `json:"groupAccessLevel"`   //组权限
	ProjectAccessLevel int64 `json:"projectAccessLevel"` //项目权限
}

type repositoryService struct {
	client *wrapper
}

func (s *repositoryService) Find(ctx context.Context, repo string) (*scm.Repository, *scm.Response, error) {
	path := fmt.Sprintf("awecloud/migrationApi/devops/project/drone/%s", repo)
	out := new(repository)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertRepository(out), res, err
}

func (s *repositoryService) FindHook(ctx context.Context, repo string, id string) (*scm.Hook, *scm.Response, error) {
	return nil, nil, nil
}

func (s *repositoryService) FindPerms(ctx context.Context, repo string) (*scm.Perm, *scm.Response, error) {
	path := fmt.Sprintf("awecloud/migrationApi/devops/project/drone/%s", repo)
	out := new(repository)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertRepository(out).Perm, res, err
}

func (s *repositoryService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	path := fmt.Sprintf("awecloud/migrationApi/devops/project/drone?%s", encodeMemberListOptions(opts))
	outs := new(repositories)
	out := []*repository{}
	res, err := s.client.do(ctx, "GET", path, nil, &outs)
	out = convertRepositories(outs)
	res.Page.Next = outs.NextPage
	return convertRepositoryList(out), res, err
}

func (s *repositoryService) ListHooks(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Hook, *scm.Response, error) {
	return nil, nil, nil
}

func (s *repositoryService) ListStatus(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) CreateHook(ctx context.Context, repo string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	return nil, nil, nil
}

func (s *repositoryService) CreateStatus(ctx context.Context, repo, ref string, input *scm.StatusInput) (*scm.Status, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) UpdateHook(ctx context.Context, repo string, id string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *repositoryService) DeleteHook(ctx context.Context, repo string, id string) (*scm.Response, error) {
	return nil, nil
}

// helper function to convert from the gogs repository list to
// the common repository structure.
func convertRepositoryList(from []*repository) []*scm.Repository {
	to := []*scm.Repository{}
	for _, v := range from {
		to = append(to, convertRepository(v))
	}
	return to
}

// helper function to convert from the gogs repository structure
// to the common repository structure.
func convertRepository(from *repository) *scm.Repository {
	to := &scm.Repository{
		ID:         strconv.Itoa(int(from.Id)),
		Namespace:  from.GroupName,
		Name:       from.ProjectName,
		Branch:     from.DefaultBranch,
		Archived:   false,
		Private:    convertPrivate(from.Private),
		Visibility: convertVisibility(from.Private),
		Perm: &scm.Perm{
			Pull:  true,
			Push:  canPush(from),
			Admin: canAdmin(from),
		},
	}
	return to
}

func convertPrivate(from int) bool {
	switch from {
	case 2:
		return false
	default:
		return true
	}
}

func convertVisibility(from int) scm.Visibility {
	switch from {
	case 2:
		return scm.VisibilityPublic
	case 1:
		return scm.VisibilityPrivate
	default:
		return scm.VisibilityUndefined
	}
}

func canPush(proj *repository) bool {
	switch {
	case proj.AccessLevel.ProjectAccessLevel >= 4:
		return true
	case proj.AccessLevel.GroupAccessLevel >= 4:
		return true
	default:
		return false
	}
}

func canAdmin(proj *repository) bool {
	switch {
	case proj.AccessLevel.ProjectAccessLevel == 5:
		return true
	case proj.AccessLevel.GroupAccessLevel == 5:
		return true
	default:
		return false
	}
}

func convertRepositories(from *repositories) []*repository {
	out := []*repository{}

	for _, o := range from.Data {
		out = append(out, o)
	}
	return out
}
