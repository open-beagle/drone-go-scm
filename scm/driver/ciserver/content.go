// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ciserver

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/drone/go-scm/scm"
)

type contentService struct {
	client *wrapper
}

func (s *contentService) Find(ctx context.Context, repo, path, ref string) (*scm.Content, *scm.Response, error) {
	var endpoint string
	if strings.Contains(repo, "/") {
		endpoint = fmt.Sprintf("awecloud/ciServer/devops/drone/content/repo/%s?ref=%s&path=%s", repo, ref, path)
	} else {
		endpoint = fmt.Sprintf("awecloud/ciServer/devops/drone/content/project/%s?ref=%s&path=%s", repo, ref, path)
	}
	out := new(content)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	raw, berr := base64.StdEncoding.DecodeString(out.Content)
	if berr != nil {
		return nil, nil, err
	}
	return &scm.Content{
		Path:   out.FilePath,
		Data:   raw,
		Sha:    out.LastCommitID,
		BlobID: out.BlobID,
	}, res, err
}

func (s *contentService) Create(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	endpoint := fmt.Sprintf("awecloud/ciServer/devops/drone/content/project/%s", encode(repo))
	in := &createUpdateContent{
		FilePath:      path,
		Branch:        params.Branch,
		Content:       params.Data,
		CommitMessage: params.Message,
	}
	res, err := s.client.do(ctx, "POST", endpoint, in, nil)
	return res, err
}

func (s *contentService) Update(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	endpoint := fmt.Sprintf("awecloud/ciServer/devops/drone/content/project/%s", encode(repo))
	in := &createUpdateContent{
		FilePath:      path,
		Branch:        params.Branch,
		Content:       params.Data,
		CommitMessage: params.Message,
	}
	res, err := s.client.do(ctx, "POST", endpoint, in, nil)
	return res, err
}

func (s *contentService) Delete(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *contentService) List(ctx context.Context, repo, path, ref string, opts scm.ListOptions) ([]*scm.ContentInfo, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

type content struct {
	FilePath     string `json:"filePath"`
	Encoding     string `json:"encoding"`
	Content      string `json:"content"`
	Ref          string `json:"ref"`
	BlobID       string `json:"blobId"`
	CommitID     string `json:"commitId"`
	LastCommitID string `json:"lastCommitId"`
}

type createUpdateContent struct {
	FilePath      string `json:"filePath"`
	Branch        string `json:"branch"`
	Content       []byte `json:"content"`
	CommitMessage string `json:"commit_message"`
}
