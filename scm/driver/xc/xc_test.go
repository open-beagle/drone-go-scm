package xc

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

func TestXc(t *testing.T) {
	client := NewDefault()
	client.Client = &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(
				&scm.Token{
					Token: "eyJhdWQiOiJ4Yy1taWdyYXRpb24iLCJleHAiOjE3MDQzNDE3NDEsInN1YiI6InJvb3QifQ.2T-9p2D8BziBFnbNm198DnDBRvkiRXY-_IDLwfr313zBdTGr7m-U-1QdSnlsAv5JDTqzt3DqTxEkkGmjhj3gPA",
				},
			),
		},
	}
	ctx := context.Background()

	// repo
	repo, _, _ := client.Repositories.Find(ctx, "test")
	fmt.Println(repo.Perm.Admin)

	// Find
	user, _, err := client.Users.Find(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("UsersFind: ", user.Name, user.Login, user.Email)

	// FindEmail
	email, _, err := client.Users.FindEmail(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("FindEmail: ", email)
	// FindLogin
	user, _, err = client.Users.FindLogin(ctx, user.Login)
	if err != nil {
		panic(err)
	}
	fmt.Println("FindLogin: ", user.Name, user.Login, user.Email)

	// ListEmail
	listemail, _, err := client.Users.ListEmail(ctx, scm.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, l := range listemail {
		fmt.Println("ListEmail: ", l.Value)
	}

	// // commit
	// got, res, err := client.Git.FindCommit(context.Background(), "diaspora/diaspora", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// //
	// fmt.Println(got.Author, res)

	// content
	g, _, _ := client.Contents.Find(ctx, "test", "t", "tt")
	fmt.Println(g.Data)
}
