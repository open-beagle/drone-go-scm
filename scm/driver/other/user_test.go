// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package other

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func testRequest(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.ID, "0d511a76-2ade-4c34-af0d-d17e84adb255"; got != want {
			t.Errorf("Want X-Request-Id: %q, got %q", want, got)
		}
	}
}

func testRate(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.Rate.Limit, 600; got != want {
			t.Errorf("Want RateLimit-Limit %d, got %d", want, got)
		}
		if got, want := res.Rate.Remaining, 599; got != want {
			t.Errorf("Want RateLimit-Remaining %d, got %d", want, got)
		}
		if got, want := res.Rate.Reset, int64(1512454441); got != want {
			t.Errorf("Want RateLimit-Reset %d, got %d", want, got)
		}
	}
}

func TestUserFind(t *testing.T) {
	defer gock.Off()

	gock.New("http://localhost:9096").
		Get("/test").
		Reply(200).
		Type("application/json").
		// SetHeaders(mockHeaders).
		File("testdata/user.json")

	client := NewDefault()
	got, res, err := client.Users.Find(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}

	want := new(scm.User)
	raw, _ := ioutil.ReadFile("testdata/user.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		json.NewEncoder(os.Stdout).Encode(got)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

// func TestUserLoginFind(t *testing.T) {
// 	defer gock.Off()

// 	gock.New("https://gitlab.com").
// 		Get("/api/v4/users").
// 		MatchParam("search", "john_smith").
// 		Reply(200).
// 		Type("application/json").
// 		SetHeaders(mockHeaders).
// 		File("testdata/user_search.json")

// 	client := NewDefault()
// 	got, res, err := client.Users.FindLogin(context.Background(), "john_smith")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	want := new(scm.User)
// 	raw, _ := ioutil.ReadFile("testdata/user_search.json.golden")
// 	json.Unmarshal(raw, &want)

// 	if diff := cmp.Diff(got, want); diff != "" {
// 		t.Errorf("Unexpected Results")
// 		t.Log(diff)
// 	}

// 	t.Run("Request", testRequest(res))
// 	t.Run("Rate", testRate(res))
// }

// func TestUserLoginFind_NotFound(t *testing.T) {
// 	defer gock.Off()

// 	gock.New("https://gitlab.com").
// 		Get("/api/v4/users").
// 		MatchParam("search", "jcitizen").
// 		Reply(200).
// 		Type("application/json").
// 		SetHeaders(mockHeaders).
// 		File("testdata/user_search.json")

// 	client := NewDefault()
// 	_, _, err := client.Users.FindLogin(context.Background(), "jcitizen")
// 	if err != scm.ErrNotFound {
// 		t.Errorf("Want Not Found Error, got %s", err)
// 	}
// }

// func TestUserLoginFind_NotAuthorized(t *testing.T) {
// 	defer gock.Off()

// 	gock.New("https://gitlab.com").
// 		Get("/api/v4/users").
// 		MatchParam("search", "jcitizen").
// 		Reply(401).
// 		Type("application/json").
// 		SetHeaders(mockHeaders).
// 		BodyString(`{"message":"401 Unauthorized"}`)

// 	client := NewDefault()
// 	_, _, err := client.Users.FindLogin(context.Background(), "jcitizen")
// 	if err == nil {
// 		t.Errorf("Want 401 Unauthorized")
// 		return
// 	}
// 	if got, want := err.Error(), "401 Unauthorized"; got != want {
// 		t.Errorf("Want %s, got %s", want, got)
// 	}
// }

// func TestUserEmailFind(t *testing.T) {
// 	defer gock.Off()

// 	gock.New("https://gitlab.com").
// 		Get("/api/v4/user").
// 		Reply(200).
// 		Type("application/json").
// 		SetHeaders(mockHeaders).
// 		File("testdata/user.json")

// 	client := NewDefault()
// 	got, res, err := client.Users.FindEmail(context.Background())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if got, want := got, "john@example.com"; got != want {
// 		t.Errorf("Want user Email %q, got %q", want, got)
// 	}

// 	t.Run("Request", testRequest(res))
// 	t.Run("Rate", testRate(res))
// }

// func TestUserEmailList(t *testing.T) {
// 	defer gock.Off()

// 	gock.New("https://gitlab.com").
// 		Get("/api/v4/user/emails").
// 		MatchParam("page", "1").
// 		MatchParam("per_page", "30").
// 		Reply(200).
// 		Type("application/json").
// 		SetHeaders(mockHeaders).
// 		SetHeaders(mockPageHeaders).
// 		File("testdata/emails.json")

// 	client := NewDefault()
// 	got, res, err := client.Users.ListEmail(context.Background(), scm.ListOptions{Size: 30, Page: 1})
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	want := []*scm.Email{}
// 	raw, _ := ioutil.ReadFile("testdata/emails.json.golden")
// 	_ = json.Unmarshal(raw, &want)

// 	if diff := cmp.Diff(got, want); diff != "" {
// 		t.Errorf("Unexpected Results")
// 		t.Log(diff)
// 	}

// 	t.Run("Request", testRequest(res))
// 	t.Run("Rate", testRate(res))
// 	t.Run("Page", testPage(res))
// }
