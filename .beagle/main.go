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
					Token: "eyJhbGciOiJIUzUxMiIsImtpZCI6IjEyMyIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIzMzZmZTIxYmJlYTY5OTg0NWM1MjEyYjg1YjJhMWJhYWZjMzAyZDNlMDdkMjcwZjI1YmM5Njc5NTRlZDk5YmM1IiwiZXhwIjoxNjYzODM1MTE5LCJzdWIiOiJyb290In0.sAsMqmrfiiQYf3akTXBAFM2AbTiUnvewdqjoeYlsSIwOmlXH1I1SlSWrBvlP7tSC-Mk7x6H0e0AE0Ndcyw-W0A",
				},
			),
		},
	}
	ctx := context.Background()

	// Find
	user, _, err := client.Users.Find(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("UsersFind: ", user.Name, user.Login)

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

	// commit
	got, res, err := client.Git.FindCommit(context.Background(), "diaspora/diaspora", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d")
	if err != nil {
		fmt.Println(err)
	}
	//
	fmt.Println(got.Author, res)

}
