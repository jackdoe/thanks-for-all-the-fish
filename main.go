package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"io/ioutil"
	"log"
	"os"
)

func clone(repo *github.Repository, githubUser, srhtUser string) {
	path, err := ioutil.TempDir("", *repo.Name)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(path)
	log.Printf("%s -> %s", *repo.FullName, path)
	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:      *repo.CloneURL,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Fatal(err)
	}
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal(err)
	}

	err = r.DeleteRemote("origin")
	remote := fmt.Sprintf("git@git.sr.ht:~%s/~%s", srhtUser, *repo.Name)
	log.Printf("replacing remote to %s", remote)
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remote},
	})

	if err != nil {
		log.Fatal(err)
	}

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	})

	if err != nil {
		log.Fatal(err)
	}

	os.RemoveAll(path)

}

func main() {
	var pgithubUser = flag.String("github", "", "github username")
	var psrhtUser = flag.String("srht", "", "sr.ht username")
	flag.Parse()

	if *pgithubUser == "" {
		log.Fatal("need github user, use -h for help")
	}
	if *psrhtUser == "" {
		psrhtUser = pgithubUser
	}

	ctx := context.Background()
	client := github.NewClient(nil)
	repos, _, err := client.Repositories.List(ctx, *pgithubUser, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, repo := range repos {
		clone(repo, *pgithubUser, *psrhtUser)
		os.Exit(1)
	}
}
