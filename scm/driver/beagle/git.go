// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package beagle

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
)

type gitService struct {
	client *wrapper
}

const commitData string = `{
    "id": "6104942438c14ec7bd21c6cd5bd995272b3faff6",
    "short_id": "6104942438c",
    "title": "Sanitize for network graph",
    "author_name": "root",
    "author_email": "root@wodcloud.com",
    "committer_name": "test",
    "committer_email": "test@wodcloud.com",
    "created_at": "2012-06-28T03:44:20-07:00",
    "message": "drone测试",
    "committed_date": "2012-06-28T03:44:20-07:00",
    "authored_date": "2012-06-28T03:44:20-07:00",
    "parent_ids": [
        "ae1d9fb46aa2b07ee9836d49862ec4e2c46fbbba"
    ],
    "last_pipeline": {
        "id": 8,
        "ref": "master",
        "sha": "2dc6aa325a317eda67812f05600bdf0fcdc70ab0",
        "status": "created"
    },
    "stats": {
        "additions": 15,
        "deletions": 10,
        "total": 25
    },
    "status": "running"
}`

func (s *gitService) CreateBranch(ctx context.Context, repo string, params *scm.ReferenceInput) (*scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/repository/branches", encode(repo))
	in := &createBranch{
		Branch: params.Name,
		Ref:    params.Sha,
	}
	return s.client.do(ctx, "POST", path, in, nil)
}

func (s *gitService) FindBranch(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/repository/branches/%s", encode(repo), name)
	out := new(branch)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertBranch(out), res, err
}

// 模拟commit数据
func (s *gitService) FindCommit(ctx context.Context, repo, ref string) (*scm.Commit, *scm.Response, error) {
	// if the reference is a branch, ensure forward slashes
	// in the branch name are escaped.
	out := new(commit)
	// if strings.Contains("ref", "/") {
	// 	ref = url.PathEscape(ref)
	// }
	// path := fmt.Sprintf("api/v4/projects/%s/repository/commits/%s", encode(repo), scm.TrimRef(ref))
	// out := new(commit)
	// res, err := s.client.do(ctx, "GET", path, nil, out)
	err := json.Unmarshal([]byte(commitData), out)
	return convertCommit(out), nil, err
}

func (s *gitService) FindTag(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/repository/tags/%s", encode(repo), name)
	out := new(branch)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertTag(out), res, err
}

func (s *gitService) ListBranches(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/repository/branches?%s", encode(repo), encodeListOptions(opts))
	out := []*branch{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertBranchList(out), res, err
}

func (s *gitService) ListCommits(ctx context.Context, repo string, opts scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/repository/commits?%s", encode(repo), encodeCommitListOptions(opts))
	out := []*commit{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertCommitList(out), res, err
}

func (s *gitService) ListTags(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/repository/tags?%s", encode(repo), encodeListOptions(opts))
	out := []*branch{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertTagList(out), res, err
}

func (s *gitService) ListChanges(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/repository/commits/%s/diff", encode(repo), ref)
	out := []*change{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertChangeList(out), res, err
}

func (s *gitService) CompareChanges(ctx context.Context, repo, source, target string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/repository/compare?from=%s&to=%s", encode(repo), source, target)
	out := new(compare)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertChangeList(out.Diffs), res, err
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
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Message        string    `json:"message"`
	AuthorName     string    `json:"author_name"`
	AuthorEmail    string    `json:"author_email"`
	AuthorDate     time.Time `json:"authored_date"`
	CommittedDate  time.Time `json:"committed_date"`
	CommitterName  string    `json:"committer_name"`
	CommitterEmail string    `json:"committer_email"`
	Created        time.Time `json:"created_at"`
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
		Sha:     from.ID,
		Author: scm.Signature{
			Login: from.AuthorName,
			Name:  from.AuthorName,
			Email: from.AuthorEmail,
			Date:  from.AuthorDate,
		},
		Committer: scm.Signature{
			Login: from.CommitterName,
			Name:  from.CommitterName,
			Email: from.CommitterEmail,
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
