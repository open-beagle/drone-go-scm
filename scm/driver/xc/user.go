// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xc

import (
	"context"
	"github.com/drone/go-scm/scm"
	"strings"
)

type userService struct {
	client *wrapper
}

func (s *userService) Find(ctx context.Context) (*scm.User, *scm.Response, error) {
	path := "xc/core/api/oauth/getUserInfo"
	out := new(user)

	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertUser(out), res, err
}

func (s *userService) FindLogin(ctx context.Context, login string) (*scm.User, *scm.Response, error) {
	path := "xc/core/api/oauth/getUserInfo"
	out := new(user)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	if err != nil {
		return nil, nil, err
	}
	if !strings.EqualFold(out.UserId, login) {
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
	path := "xc/core/api/oauth/getUserInfo"
	out := new(user)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	var o []*user
	o = append(o, out)
	return convertEmailList(o), res, err
}

type user struct {
	UserId   string `json:"user_id" xorm:"user_id pk"`
	UnitId   string `json:"unit_id"`
	UserType string `json:"user_type"`
	UserName string `json:"user_name"`
	Password string `json:"password,omitempty" xorm:"-"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

// helper function to convert from the gitlab user structure to
// the common user structure.
func convertUser(from *user) *scm.User {
	return &scm.User{
		Login: from.UserId,
		Name:  from.UserName,
		Email: from.Email,
	}
}

// helper function to convert from the gitlab email list to
// the common email structure.
func convertEmailList(from []*user) []*scm.Email {
	var to []*scm.Email
	for _, v := range from {
		to = append(to, convertEmail(v))
	}
	return to
}

// helper function to convert from the gitlab email structure to
// the common email structure.
func convertEmail(from *user) *scm.Email {
	return &scm.Email{
		Value: from.Email,
	}
}
