// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ciserver

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/drone/go-scm/scm"
)

type gitService struct {
	client *wrapper
}

func (s *gitService) CreateBranch(ctx context.Context, repo string, params *scm.ReferenceInput) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *gitService) FindBranch(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) FindCommit(ctx context.Context, repo, ref string) (*scm.Commit, *scm.Response, error) {
	path := fmt.Sprintf("awecloud/ciServer/devops/drone/commit/%s/%s", repo, scm.TrimRef(ref))
	out := new(commit)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertCommit(out), res, err
}

func (s *gitService) FindTag(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListBranches(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("awecloud/ciServer/devops/drone/branch/%s", encode(repo))
	out := []*branch{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertBranchList(out), res, err
}

func (s *gitService) ListCommits(ctx context.Context, repo string, opts scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListTags(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListChanges(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) CompareChanges(ctx context.Context, repo, source, target string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListGroup(ctx context.Context) ([]*scm.Group, *scm.Response, error) {
	path := fmt.Sprintf("awecloud/ciServer/devops/drone/group")
	out := []*group{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertGroupList(out), res, err
}

type commit struct {
	ID            int64     `json:"id"`
	Message       string    `json:"message"`
	AuthorName    string    `json:"authorName"`
	AuthorDate    time.Time `json:"authoredDate"`
	CommittedDate time.Time `json:"committedDate"`
	CommitterName string    `json:"committerName"`
}

type group struct {
	Id        int    `json:"id"`
	GroupName string `json:"name" `
}

func convertGroupList(from []*group) []*scm.Group {
	to := []*scm.Group{}
	for _, v := range from {
		to = append(to, convertGroup(v))
	}
	return to
}

func convertGroup(from *group) *scm.Group {
	return &scm.Group{
		Id:   from.Id,
		Name: from.GroupName,
	}
}

func convertCommit(from *commit) *scm.Commit {
	return &scm.Commit{
		Message: from.Message,
		Sha:     strconv.Itoa(int(from.ID)),
		Author: scm.Signature{
			Login: from.AuthorName,
			Name:  from.AuthorName,
			Date:  from.AuthorDate,
		},
		Committer: scm.Signature{
			Login: from.CommitterName,
			Name:  from.CommitterName,
			Date:  from.CommittedDate,
		},
	}
}

type branch struct {
	BranchName string `json:"branchName"`
}

func convertBranchList(from []*branch) []*scm.Reference {
	to := []*scm.Reference{}
	for _, v := range from {
		to = append(to, convertBranch(v))
	}
	return to
}

func convertBranch(from *branch) *scm.Reference {
	return &scm.Reference{
		Name: from.BranchName,
		Path: from.BranchName,
		Sha:  "",
	}
}
