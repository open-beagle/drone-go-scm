// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package beagle

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/internal/null"
)

type repository struct {
	ID            int         `json:"id"`
	Path          string      `json:"path"`
	PathNamespace string      `json:"path_with_namespace"`
	DefaultBranch string      `json:"default_branch"`
	Visibility    string      `json:"visibility"`
	Archived      bool        `json:"archived"`
	WebURL        string      `json:"web_url"`
	SSHURL        string      `json:"ssh_url_to_repo"`
	HTTPURL       string      `json:"http_url_to_repo"`
	Namespace     namespace   `json:"namespace"`
	Permissions   permissions `json:"permissions"`
}

type namespace struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
}

type permissions struct {
	ProjectAccess access `json:"project_access"`
	GroupAccess   access `json:"group_access"`
}

type access struct {
	AccessLevel       int `json:"access_level"`
	NotificationLevel int `json:"notification_level"`
}

type hook struct {
	ID                    int       `json:"id"`
	URL                   string    `json:"url"`
	ProjectID             int       `json:"project_id"`
	PushEvents            bool      `json:"push_events"`
	IssuesEvents          bool      `json:"issues_events"`
	MergeRequestsEvents   bool      `json:"merge_requests_events"`
	TagPushEvents         bool      `json:"tag_push_events"`
	NoteEvents            bool      `json:"note_events"`
	JobEvents             bool      `json:"job_events"`
	PipelineEvents        bool      `json:"pipeline_events"`
	WikiPageEvents        bool      `json:"wiki_page_events"`
	EnableSslVerification bool      `json:"enable_ssl_verification"`
	CreatedAt             time.Time `json:"created_at"`
}

type repositoryService struct {
	client *wrapper
}

const analogdata string = `{
    "id": 178504,
    "description": "",
    "default_branch": "master",
    "tag_list": [],
    "ssh_url_to_repo": "git@gitlab.com:diaspora/diaspora.git",
    "http_url_to_repo": "https://gitlab.com/diaspora/diaspora.git",
    "web_url": "https://gitlab.com/diaspora/diaspora",
    "name": "Diaspora",
    "name_with_namespace": "diaspora / Diaspora",
    "path": "diaspora",
    "path_with_namespace": "diaspora/diaspora",
    "avatar_url": null,
    "star_count": 0,
    "forks_count": 0,
    "created_at": "2016-01-19T09:05:50.355Z",
    "last_activity_at": "2022-09-23T09:05:50.355Z",
    "_links": {
        "self": "http://gitlab.com/api/v4/projects/178504",
        "issues": "http://gitlab.com/api/v4/projects/178504/issues",
        "merge_requests": "http://gitlab.com/api/v4/projects/178504/merge_requests",
        "repo_branches": "http://gitlab.com/api/v4/projects/178504/repository/branches",
        "labels": "http://gitlab.com/api/v4/projects/178504/labels",
        "events": "http://gitlab.com/api/v4/projects/178504/events",
        "members": "http://gitlab.com/api/v4/projects/178504/members"
    },
    "archived": false,
    "visibility": "public",
    "resolve_outdated_diff_discussions": null,
    "container_registry_enabled": null,
    "issues_enabled": true,
    "merge_requests_enabled": true,
    "wiki_enabled": true,
    "jobs_enabled": true,
    "snippets_enabled": false,
    "shared_runners_enabled": true,
    "lfs_enabled": true,
    "creator_id": 57658,
    "namespace": {
        "id": 120836,
        "name": "diaspora",
        "path": "diaspora",
        "kind": "group",
        "full_path": "diaspora",
        "parent_id": null
    },
    "import_status": "finished",
    "open_issues_count": 0,
    "public_jobs": true,
    "ci_config_path": null,
    "shared_with_groups": [],
    "only_allow_merge_if_pipeline_succeeds": false,
    "request_access_enabled": true,
    "only_allow_merge_if_all_discussions_are_resolved": null,
    "printing_merge_request_link_enabled": true,
    "approvals_before_merge": 0,
    "permissions": {
        "project_access": null,
        "group_access": {
            "access_level": 40,
            "notification_level": 3
        }
    }
}`

// 模拟repo数据
func (s *repositoryService) Find(ctx context.Context, repo string) (*scm.Repository, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s", encode(repo))
	out := new(repository)
	err := json.Unmarshal([]byte(analogdata), out)

	// res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertRepository(out), nil, err
}

//
func (s *repositoryService) FindHook(ctx context.Context, repo string, id string) (*scm.Hook, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s/hooks/%s", encode(repo), id)
	// out := new(hook)
	// res, err := s.client.do(ctx, "GET", path, nil, out)
	// return convertHook(out), res, err
	return nil, nil, nil
}

func (s *repositoryService) FindPerms(ctx context.Context, repo string) (*scm.Perm, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s", encode(repo))
	out := new(repository)
	err := json.Unmarshal([]byte(analogdata), out)
	// res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertRepository(out).Perm, nil, err
}

func (s *repositoryService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects?%s", encodeMemberListOptions(opts))
	out := []*repository{}
	o := new(repository)
	err := json.Unmarshal([]byte(analogdata), o)
	out = append(out, o)
	res := new(scm.Response)
	// 根据opts值确定Next大小
	// if opts.Page == 1 {
	// 	res.Page.Next = 0
	// } else {
	// 	res.Page.Next = 1
	// }
	res.Page.Next = 0
	// res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertRepositoryList(out), res, err
}

// 模拟hook查询
func (s *repositoryService) ListHooks(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Hook, *scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s/hooks?%s", encode(repo), encodeListOptions(opts))
	// out := []*hook{}

	// res, err := s.client.do(ctx, "GET", path, nil, &out)
	return nil, nil, nil
}

func (s *repositoryService) ListStatus(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Status, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/projects/%s/repository/commits/%s/statuses?%s", encode(repo), ref, encodeListOptions(opts))
	out := []*status{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertStatusList(out), res, err
}

const hookdata string = `{
    "id": 1,
    "url": "http://example.com/hook",
    "project_id": 3,
    "push_events": true,
    "issues_events": false,
    "merge_requests_events": false,
    "tag_push_events": false,
    "note_events": false,
    "job_events": false,
    "pipeline_events": false,
    "wiki_page_events": false,
    "enable_ssl_verification": true,
    "created_at": "2022-09-23T09:05:50.355Z"
}`

// 模拟hook数据
func (s *repositoryService) CreateHook(ctx context.Context, repo string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	// params := url.Values{}
	// params.Set("url", input.Target)
	// if input.Secret != "" {
	// 	params.Set("token", input.Secret)
	// }
	// if input.SkipVerify {
	// 	params.Set("enable_ssl_verification", "false")
	// }
	// if input.Events.Branch {
	// 	// no-op
	// }
	// if input.Events.Issue {
	// 	params.Set("issues_events", "true")
	// }
	// if input.Events.IssueComment ||
	// 	input.Events.PullRequestComment {
	// 	params.Set("note_events", "true")
	// }
	// if input.Events.PullRequest {
	// 	params.Set("merge_requests_events", "true")
	// }
	// if input.Events.Push || input.Events.Branch {
	// 	params.Set("push_events", "true")
	// }
	// if input.Events.Tag {
	// 	params.Set("tag_push_events", "true")
	// }
	// // 模拟hook创建
	// fmt.Println(params.Encode())

	// // path := fmt.Sprintf("api/v4/projects/%s/hooks?%s", encode(repo), params.Encode())
	// out := new(hook)
	// err := json.Unmarshal([]byte(hookdata), out)
	// // res, err := s.client.do(ctx, "POST", path, nil, out)
	// return convertHook(out), nil, err
	return nil, nil, nil
}

const repostatus = `{
    "author": {
        "web_url": "https://gitlab.example.com/thedude",
        "name": "root",
        "avatar_url": "https://gitlab.example.com/uploads/user/avatar/28/The-Big-Lebowski-400-400.png",
        "username": "root",
        "state": "active",
        "id": 28
    },
    "name": "default",
    "sha": "18f3e63d05582537db6d183d9d557be09e1f90c8",
    "status": "running",
    "coverage": 100.0,
    "description": "the dude abides",
    "id": 93,
    "target_url": "https://gitlab.example.com/thedude/gitlab-ce/builds/91",
    "ref": null,
    "started_at": null,
    "created_at": "2016-01-19T09:05:50.355Z",
    "allow_failure": false,
    "finished_at": "2022-09-23T09:05:50.355Z"
}`

// 模拟status数据
func (s *repositoryService) CreateStatus(ctx context.Context, repo, ref string, input *scm.StatusInput) (*scm.Status, *scm.Response, error) {
	// params := url.Values{}
	// params.Set("state", convertFromState(input.State))
	// params.Set("name", input.Label)
	// params.Set("target_url", input.Target)
	// path := fmt.Sprintf("api/v4/projects/%s/statuses/%s?%s", encode(repo), ref, params.Encode())

	out := new(status)
	err := json.Unmarshal([]byte(repostatus), out)
	switch input.State {
	case scm.StatePending:
		out.Status = "pending"
	case scm.StateRunning:
		out.Status = "running"
	case scm.StateSuccess:
		out.Status = "success"
	default:
		out.Status = "error"
	}
	// res, err := s.client.do(ctx, "POST", path, nil, out)
	return convertStatus(out), nil, err
}

func (s *repositoryService) UpdateHook(ctx context.Context, repo string, id string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

// 模拟hook删除
func (s *repositoryService) DeleteHook(ctx context.Context, repo string, id string) (*scm.Response, error) {
	// path := fmt.Sprintf("api/v4/projects/%s/hooks/%s", encode(repo), id)
	// return s.client.do(ctx, "DELETE", path, nil, nil)
	return nil, nil
}

// helper function to convert from the gogs repository list to
// the common repository structure.
func convertRepositoryList(from []*repository) []*scm.Repository {
	to := []*scm.Repository{}
	for _, v := range from {
		to = append(to, convertRepository(v))
	}
	return to
}

// helper function to convert from the gogs repository structure
// to the common repository structure.
func convertRepository(from *repository) *scm.Repository {
	to := &scm.Repository{
		ID:         strconv.Itoa(from.ID),
		Namespace:  from.Namespace.Path,
		Name:       from.Path,
		Branch:     from.DefaultBranch,
		Archived:   from.Archived,
		Private:    convertPrivate(from.Visibility),
		Visibility: convertVisibility(from.Visibility),
		Clone:      from.HTTPURL,
		CloneSSH:   from.SSHURL,
		Link:       from.WebURL,
		Perm: &scm.Perm{
			Pull:  true,
			Push:  canPush(from),
			Admin: canAdmin(from),
		},
	}
	if path := from.Namespace.FullPath; path != "" {
		to.Namespace = path
	}
	if to.Namespace == "" {
		if parts := strings.SplitN(from.PathNamespace, "/", 2); len(parts) == 2 {
			to.Namespace = parts[1]
		}
	}
	return to
}

func convertHookList(from []*hook) []*scm.Hook {
	to := []*scm.Hook{}
	for _, v := range from {
		to = append(to, convertHook(v))
	}
	return to
}

func convertHook(from *hook) *scm.Hook {
	return &scm.Hook{
		ID:         strconv.Itoa(from.ID),
		Active:     true,
		Target:     from.URL,
		Events:     convertEvents(from),
		SkipVerify: !from.EnableSslVerification,
	}
}

type status struct {
	Name    string      `json:"name"`
	Desc    null.String `json:"description"`
	Status  string      `json:"status"`
	Sha     string      `json:"sha"`
	Ref     string      `json:"ref"`
	Target  null.String `json:"target_url"`
	Created time.Time   `json:"created_at"`
	Updated time.Time   `json:"updated_at"`
}

func convertStatusList(from []*status) []*scm.Status {
	to := []*scm.Status{}
	for _, v := range from {
		to = append(to, convertStatus(v))
	}
	return to
}

func convertStatus(from *status) *scm.Status {
	return &scm.Status{
		State:  convertState(from.Status),
		Label:  from.Name,
		Desc:   from.Desc.String,
		Target: from.Target.String,
	}
}

func convertEvents(from *hook) []string {
	var events []string
	if from.IssuesEvents {
		events = append(events, "issues")
	}
	if from.TagPushEvents {
		events = append(events, "tag")
	}
	if from.PushEvents {
		events = append(events, "push")
	}
	if from.NoteEvents {
		events = append(events, "comment")
	}
	if from.MergeRequestsEvents {
		events = append(events, "merge")
	}
	return events
}

func convertState(from string) scm.State {
	switch from {
	case "canceled":
		return scm.StateCanceled
	case "failed":
		return scm.StateFailure
	case "pending":
		return scm.StatePending
	case "running":
		return scm.StateRunning
	case "success":
		return scm.StateSuccess
	default:
		return scm.StateUnknown
	}
}

func convertFromState(from scm.State) string {
	switch from {
	case scm.StatePending:
		return "pending"
	case scm.StateRunning:
		return "running"
	case scm.StateSuccess:
		return "success"
	case scm.StateCanceled:
		return "canceled"
	default:
		return "failed"
	}
}

func convertPrivate(from string) bool {
	switch from {
	case "public", "":
		return false
	default:
		return true
	}
}

func convertVisibility(from string) scm.Visibility {
	switch from {
	case "public":
		return scm.VisibilityPublic
	case "private":
		return scm.VisibilityPrivate
	case "internal":
		return scm.VisibilityInternal
	default:
		return scm.VisibilityUndefined
	}
}

func canPush(proj *repository) bool {
	switch {
	case proj.Permissions.ProjectAccess.AccessLevel >= 30:
		return true
	case proj.Permissions.GroupAccess.AccessLevel >= 30:
		return true
	default:
		return false
	}
}

func canAdmin(proj *repository) bool {
	switch {
	case proj.Permissions.ProjectAccess.AccessLevel >= 40:
		return true
	case proj.Permissions.GroupAccess.AccessLevel >= 40:
		return true
	default:
		return false
	}
}
