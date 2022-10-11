// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package beagle

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
	path := fmt.Sprintf("api/v4/projects/%s/repository/branches", encode(repo))
	in := &createBranch{
		Branch: params.Name,
		Ref:    params.Sha,
	}
	return s.client.do(ctx, "POST", path, in, nil)
}

func (s *gitService) FindBranch(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) FindCommit(ctx context.Context, repo, ref string) (*scm.Commit, *scm.Response, error) {
	path := fmt.Sprintf("awecloud/ciApi/devops/object/%s/%s/commits", repo, scm.TrimRef(ref))
	out := new(commit)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertCommit(out), res, err
}

func (s *gitService) FindTag(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s/repository/tags/%s", encode(repo), name)
	// out := new(branch)
	// res, err := s.client.do(ctx, "GET", path, nil, out)
	// return convertTag(out), res, err
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListBranches(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s/repository/branches?%s", encode(repo), encodeListOptions(opts))
	// out := []*branch{}
	// res, err := s.client.do(ctx, "GET", path, nil, &out)
	// return convertBranchList(out), res, err
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListCommits(ctx context.Context, repo string, opts scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s/repository/commits?%s", encode(repo), encodeCommitListOptions(opts))
	// out := []*commit{}
	// res, err := s.client.do(ctx, "GET", path, nil, &out)
	// return convertCommitList(out), res, err
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListTags(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s/repository/tags?%s", encode(repo), encodeListOptions(opts))
	// out := []*branch{}
	// res, err := s.client.do(ctx, "GET", path, nil, &out)
	// return convertTagList(out), res, err
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) ListChanges(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s/repository/commits/%s/diff", encode(repo), ref)
	// out := []*change{}
	// res, err := s.client.do(ctx, "GET", path, nil, &out)
	// return convertChangeList(out), res, err
	return nil, nil, scm.ErrNotSupported
}

func (s *gitService) CompareChanges(ctx context.Context, repo, source, target string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s/repository/compare?from=%s&to=%s", encode(repo), source, target)
	// out := new(compare)
	// res, err := s.client.do(ctx, "GET", path, nil, &out)
	// return convertChangeList(out.Diffs), res, err
	return nil, nil, scm.ErrNotSupported
}

type branch struct {
	Name   string `json:"name"`
	Commit struct {
		ID string `json:"id"`
	}
}

type createBranch struct {
	Branch string `json:"branch"`
	Ref    string `json:"ref"`
}

type commit struct {
	ID            int64     `json:"id"`
	Message       string    `json:"message"`
	AuthorName    string    `json:"author_name"`
	AuthorDate    time.Time `json:"authored_date"`
	CommittedDate time.Time `json:"committed_date"`
	CommitterName string    `json:"committer_name"`
	Created       time.Time `json:"created_at"`
}

type compare struct {
	Diffs []*change `json:"diffs"`
}

type change struct {
	OldPath string `json:"old_path"`
	NewPath string `json:"new_path"`
	Added   bool   `json:"new_file"`
	Renamed bool   `json:"renamed_file"`
	Deleted bool   `json:"deleted_file"`
}

func convertCommitList(from []*commit) []*scm.Commit {
	to := []*scm.Commit{}
	for _, v := range from {
		to = append(to, convertCommit(v))
	}
	return to
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

func convertBranchList(from []*branch) []*scm.Reference {
	to := []*scm.Reference{}
	for _, v := range from {
		to = append(to, convertBranch(v))
	}
	return to
}

func convertBranch(from *branch) *scm.Reference {
	return &scm.Reference{
		Name: scm.TrimRef(from.Name),
		Path: scm.ExpandRef(from.Name, "refs/heads/"),
		Sha:  from.Commit.ID,
	}
}

func convertTagList(from []*branch) []*scm.Reference {
	to := []*scm.Reference{}
	for _, v := range from {
		to = append(to, convertTag(v))
	}
	return to
}

func convertTag(from *branch) *scm.Reference {
	return &scm.Reference{
		Name: scm.TrimRef(from.Name),
		Path: scm.ExpandRef(from.Name, "refs/tags/"),
		Sha:  from.Commit.ID,
	}
}

func convertChangeList(from []*change) []*scm.Change {
	to := []*scm.Change{}
	for _, v := range from {
		to = append(to, convertChange(v))
	}
	return to
}

func convertChange(from *change) *scm.Change {
	to := &scm.Change{
		Path:    from.NewPath,
		Added:   from.Added,
		Deleted: from.Deleted,
		Renamed: from.Renamed,
	}
	if to.Path == "" {
		to.Path = from.OldPath
	}
	return to
}
