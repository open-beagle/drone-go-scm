// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package beagle

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"

	"github.com/drone/go-scm/scm"
)

type contentService struct {
	client *wrapper
}

// 增加模拟数据
func (s *contentService) Find(ctx context.Context, repo, path, ref string) (*scm.Content, *scm.Response, error) {
	endpoint := fmt.Sprintf("awecloud/lzjciApi/devops/object/%s?ref=%s&path=%s", repo, ref, path)
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
	endpoint := fmt.Sprintf("api/v4/projects/%s/repository/files/%s", encode(repo), encodePath(path))
	in := &createUpdateContent{
		Branch:        params.Branch,
		Content:       params.Data,
		CommitMessage: params.Message,
		Encoding:      "base64",
		AuthorName:    params.Signature.Name,
		AuthorEmail:   params.Signature.Email,
	}
	res, err := s.client.do(ctx, "POST", endpoint, in, nil)
	return res, err

}

func (s *contentService) Update(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	endpoint := fmt.Sprintf("api/v4/projects/%s/repository/files/%s", encode(repo), encodePath(path))
	in := &createUpdateContent{
		Branch:        params.Branch,
		Content:       params.Data,
		CommitMessage: params.Message,
		Encoding:      "base64",
		AuthorName:    params.Signature.Name,
		AuthorEmail:   params.Signature.Email,
		LastCommitID:  params.Sha,
	}
	res, err := s.client.do(ctx, "PUT", endpoint, in, nil)
	return res, err
}

func (s *contentService) Delete(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	endpoint := fmt.Sprintf("api/v4/projects/%s/repository/files/%s", encode(repo), encodePath(path))
	in := &createUpdateContent{
		Branch:        params.Branch,
		CommitMessage: params.Message,
		Encoding:      "base64",
		AuthorName:    params.Signature.Name,
		AuthorEmail:   params.Signature.Email,
		LastCommitID:  params.Sha,
	}
	res, err := s.client.do(ctx, "DELETE", endpoint, in, nil)
	return res, err
}

func (s *contentService) List(ctx context.Context, repo, path, ref string, opts scm.ListOptions) ([]*scm.ContentInfo, *scm.Response, error) {
	endpoint := fmt.Sprintf("api/v4/projects/%s/repository/tree?path=%s&ref=%s&%s", encode(repo), url.QueryEscape(path), ref, encodeListOptions(opts))
	out := []*object{}
	res, err := s.client.do(ctx, "GET", endpoint, nil, &out)
	return convertContentInfoList(out), res, err
}

type content struct {
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
	Size         int    `json:"size"`
	Encoding     string `json:"encoding"`
	Content      string `json:"content"`
	Ref          string `json:"ref"`
	BlobID       string `json:"blob_id"`
	CommitID     string `json:"commit_id"`
	LastCommitID string `json:"last_commit_id"`
}

type createUpdateContent struct {
	Branch        string `json:"branch"`
	Content       []byte `json:"content"`
	CommitMessage string `json:"commit_message"`
	Encoding      string `json:"encoding"`
	AuthorEmail   string `json:"author_email"`
	AuthorName    string `json:"author_name"`
	LastCommitID  string `json:"last_commit_id"`
}

type object struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
}

func convertContentInfoList(from []*object) []*scm.ContentInfo {
	to := []*scm.ContentInfo{}
	for _, v := range from {
		to = append(to, convertContentInfo(v))
	}
	return to
}

func convertContentInfo(from *object) *scm.ContentInfo {
	to := &scm.ContentInfo{Path: from.Path}
	// See the following link for supported file modes:
	// https://godoc.org/gopkg.in/src-d/go-git.v4/plumbing/filemode
	switch mode, _ := strconv.ParseInt(from.Mode, 8, 32); mode {
	case 0100644, 0100664, 0100755:
		to.Kind = scm.ContentKindFile
	case 0040000:
		to.Kind = scm.ContentKindDirectory
	case 0120000:
		to.Kind = scm.ContentKindSymlink
	case 0160000:
		to.Kind = scm.ContentKindGitlink
	default:
		to.Kind = scm.ContentKindUnsupported
	}
	return to
}
