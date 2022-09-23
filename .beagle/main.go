package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/beagle"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

func main() {
	client := beagle.NewDefault()
	client.Client = &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(
				&scm.Token{
					Token: "eyJhbGciOiJIUzUxMiIsImtpZCI6IjEyMyIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIzMzZmZTIxYmJlYTY5OTg0NWM1MjEyYjg1YjJhMWJhYWZjMzAyZDNlMDdkMjcwZjI1YmM5Njc5NTRlZDk5YmM1IiwiZXhwIjoxNjYzOTMwMjIwLCJzdWIiOiJyb290In0.KpjdAfYtZAhE9n4QyC1gIvTLa_dLfYvbBH-N1iRHiGLFvqhDwx4JBgUzAvDBju5Esr-6cpA1U9dWZtJacljJfg",
				},
			),
		},
	}
	ctx := context.Background()

	// repo
	repo, _, _ := client.Repositories.Find(ctx, "test")
	fmt.Println(repo.Perm.Admin)

	// // Find
	// user, _, err := client.Users.Find(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("UsersFind: ", user.Name, user.Login)

	// // FindEmail
	// email, _, err := client.Users.FindEmail(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("FindEmail: ", email)
	// // FindLogin
	// user, _, err = client.Users.FindLogin(ctx, user.Login)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("FindLogin: ", user.Name, user.Login, user.Email)

	// // ListEmail
	// listemail, _, err := client.Users.ListEmail(ctx, scm.ListOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	// for _, l := range listemail {
	// 	fmt.Println("ListEmail: ", l.Value)
	// }

	// // commit
	// got, res, err := client.Git.FindCommit(context.Background(), "diaspora/diaspora", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// //
	// fmt.Println(got.Author, res)

	// // content
	// g, _, _ := client.Contents.Find(ctx, "test", "t", "tt")
	// fmt.Println(g.Data)
}
