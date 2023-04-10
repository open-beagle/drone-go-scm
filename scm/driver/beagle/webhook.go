package beagle

// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/drone/go-scm/scm"
)

type webhookService struct {
	client *wrapper
}

func (s *webhookService) Parse(req *http.Request, fn scm.SecretFunc) (scm.Webhook, error) {
	data, err := ioutil.ReadAll(
		io.LimitReader(req.Body, 10000000),
	)
	if err != nil {
		return nil, err
	}

	var hook scm.Webhook

	for {
		hook, err = convertPushHook(data)
		if err == nil {
			break
		}
		hook, err = converBranchHook(data)
		if err == nil {
			break
		}
		hook, err = convertCommentHook(data)
		if err == nil {
			break
		}
		hook, err = convertTagHook(data)
		if err == nil {
			break
		}
		hook, err = convertPullRequestHook(data)
		if err == nil {
			break
		}
		hook, err = convertDeploymentHook(data)
		if err == nil {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	token, err := fn(hook)
	if err != nil {
		return hook, err
	} else if token == "" {
		return hook, nil
	}
	return hook, nil
}

func convertPushHook(data []byte) (*scm.PushHook, error) {
	hook := new(scm.PushHook)
	err := json.Unmarshal(data, hook)
	return hook, err
}

func converBranchHook(data []byte) (*scm.BranchHook, error) {
	hook := new(scm.BranchHook)
	err := json.Unmarshal(data, hook)
	return hook, err
}

func convertCommentHook(data []byte) (*scm.IssueCommentHook, error) {
	hook := new(scm.IssueCommentHook)
	err := json.Unmarshal(data, hook)
	return hook, err
}

func convertTagHook(data []byte) (*scm.TagHook, error) {
	hook := new(scm.TagHook)
	err := json.Unmarshal(data, hook)
	return hook, err
}

func convertPullRequestHook(data []byte) (*scm.PullRequestHook, error) {
	hook := new(scm.PullRequestHook)
	err := json.Unmarshal(data, hook)
	return hook, err
}

func convertDeploymentHook(data []byte) (*scm.DeployHook, error) {
	hook := new(scm.DeployHook)
	err := json.Unmarshal(data, hook)
	return hook, err
}
