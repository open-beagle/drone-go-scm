// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package beagle

import (
	"context"
	"strings"

	"github.com/drone/go-scm/scm"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	path := "awecloud/dex/oauth/getUserInfo"
	out := new(user)

	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	path := "awecloud/dex/oauth/getUserInfo"
	out := new(user)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	if err != nil {
		return nil, nil, err
	}
	if !strings.EqualFold(out.Metadata.Name, login) {
		return nil, nil, scm.ErrNotFound
	}
	return convertUser(out), res, err
}

func (s *userService) FindEmail(ctx context.Context) (string, *scm.Response, error) {
	user, res, err := s.Find(ctx)
	return user.Email, res, err
}

// 未使用
func (s *userService) ListEmail(ctx context.Context, opts scm.ListOptions) ([]*scm.Email, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/user/emails?%s", encodeListOptions(opts))
	path := "awecloud/dex/oauth/getUserInfo"
	out := new(user)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	o := []*spec{}
	o = append(o, &out.Spec)
	return convertEmailList(o), res, err
}

type user struct {
	Metadata metadata `json:"metadata"`
	Spec     spec     `json:"spec"`
}

type metadata struct {
	Name string `json:"name"`
}

type spec struct {
	Alias string `json:"alias"`
	Email string `json:"email"`
}

// helper function to convert from the gitlab user structure to
// the common user structure.
func convertUser(from *user) *scm.User {
	return &scm.User{
		Email: from.Spec.Email,
		Login: from.Metadata.Name,
		Name:  from.Spec.Alias,
	}
}

// helper function to convert from the gitlab email list to
// the common email structure.
func convertEmailList(from []*spec) []*scm.Email {
	to := []*scm.Email{}
	for _, v := range from {
		to = append(to, convertEmail(v))
	}
	return to
}

// helper function to convert from the gitlab email structure to
// the common email structure.
func convertEmail(from *spec) *scm.Email {
	return &scm.Email{
		Value: from.Email,
	}
}
