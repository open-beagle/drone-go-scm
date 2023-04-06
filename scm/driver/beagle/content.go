// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package beagle

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type contentService struct {
	client *wrapper
}

func (s *contentService) Find(ctx context.Context, repo, path, ref string) (*scm.Content, *scm.Response, error) {
	endpoint := fmt.Sprintf("awecloud/lzjciApi/devops/object/%s?ref=%s&path=%s", repo, ref, path)
	out := new(scm.Content)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	return out, res, err
}

func (s *contentService) Create(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	return nil, scm.ErrNotSupported

}

func (s *contentService) Update(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *contentService) Delete(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *contentService) List(ctx context.Context, repo, path, ref string, opts scm.ListOptions) ([]*scm.ContentInfo, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// type content struct {
// 	FilePath     string `json:"filePath"`
// 	Encoding     string `json:"encoding"`
// 	Content      string `json:"content"`
// 	Ref          string `json:"ref"`
// 	BlobID       string `json:"blobId"`
// 	CommitID     string `json:"commitId"`
// 	LastCommitID string `json:"lastCommitId"`
// }
