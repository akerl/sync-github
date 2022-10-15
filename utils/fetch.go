package utils

import (
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v45/github"
)

// HandleRepos handles all repos, updating or cloning as necessary
func HandleRepos(config Config, repos []*github.Repository) error {
	for _, repo := range repos {
		fmt.Println(*repo.FullName)
		repoPath := path.Join(config.TargetPath, *repo.FullName)
		gitPath := path.Join(repoPath, "objects")

		exists, err := repoExists(gitPath)
		if err != nil {
			return err
		}

		if exists {
			err = UpdateRepo(repoPath, config)
		} else {
			err = CloneRepo(repo, repoPath, config)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func repoExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// UpdateRepo updates an existing repo
func UpdateRepo(repoPath string, config Config) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	user, err := config.GetUser()
	if err != nil {
		return err
	}

	remote, err := r.Remote("origin")
	if err != nil {
		return err
	}

	err = remote.Fetch(&git.FetchOptions{
		RefSpecs: []gitconfig.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
		Force:    true,
		Auth:     &http.BasicAuth{Username: user, Password: config.AuthToken},
	})
	if err == git.NoErrAlreadyUpToDate {
		err = nil
	}
	return err
}

// CloneRepo clones an individual repo to a local path
func CloneRepo(repo *github.Repository, repoPath string, config Config) error {
	err := os.MkdirAll(repoPath, 0750)
	if err != nil {
		return nil
	}

	user, err := config.GetUser()
	if err != nil {
		return err
	}

	_, err = git.PlainClone(repoPath, true, &git.CloneOptions{
		URL:  *repo.CloneURL,
		Auth: &http.BasicAuth{Username: user, Password: config.AuthToken},
	})
	if err == transport.ErrEmptyRemoteRepository {
		err = nil
	}
	return err
}
