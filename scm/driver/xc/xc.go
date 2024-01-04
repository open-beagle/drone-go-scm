// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gitlab implements a GitLab client.
package xc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/drone/go-scm/scm"
)

type Token struct {
	Access  string
	Refresh string
	Expires time.Time
}

type key int

const (
	tokenKey key = iota
	errorKey
)

// New returns a new GitLab API client.
func New(uri string) (*scm.Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}
	client := &wrapper{new(scm.Client)}
	client.BaseURL = base
	// initialize services
	client.Driver = scm.DriverXc
	// client.Linker = &linker{base.String()}
	client.Contents = &contentService{client}
	client.Git = &gitService{client}
	// client.Issues = &issueService{client}
	// client.Organizations = &organizationService{client}
	// client.Milestones = &milestoneService{client}
	// client.PullRequests = &pullService{client}
	client.Repositories = &repositoryService{client}
	// client.Releases = &releaseService{client}
	// client.Reviews = &reviewService{client}
	client.Users = &userService{client}
	// client.Webhooks = &webhookService{client}
	return client.Client, nil
}

// NewDefault returns a new GitLab API client using the
// default gitlab.com address.

func NewDefault() *scm.Client {
	client, _ := New("https://xc-cloud.dev.wodcloud.com")
	return client
}

// wraper wraps the Client to provide high level helper functions
// for making http requests and unmarshaling the response.
type wrapper struct {
	*scm.Client
}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response.
func (c *wrapper) do(ctx context.Context, method, path string, in, out interface{}) (*scm.Response, error) {
	req := &scm.Request{
		Method: method,
		Path:   path,
	}

	token, _ := ctx.Value(tokenKey).(*Token)

	// if we are posting or putting data, we need to
	// write it to the body of the request.
	if in != nil {
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(in)
		//添加获取用户信息header
		req.Header = map[string][]string{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer " + token.Access},
		}
		req.Body = buf
	}

	// execute the http request
	res, err := c.Client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.Status > 300 {
		err := new(Error)
		json.NewDecoder(res.Body).Decode(err)
		return res, err
	}

	if out == nil {
		return res, nil
	}

	// if a json response is expected, parse and return
	// the json response.
	return res, json.NewDecoder(res.Body).Decode(out)
}

// Error represents a GitLab error.
type Error struct {
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}
