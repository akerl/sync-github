package utils

import (
	"context"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

// GithubClient creates a Github Client from the provided OAuth token
func (c *Config) GithubClient() *github.Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.AuthToken})
	tokenClient := oauth2.NewClient(context.Background(), tokenSource)
	return github.NewClient(tokenClient)
}

// GetUser finds the current username
func (c *Config) GetUser() (string, error) {
	client := c.GithubClient()
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		return "", err
	}
	return *user.Login, nil
}

// FilteredRepos loads a list of repositories and removes ones that match a filter
func FilteredRepos(config Config) ([]*github.Repository, error) {
	filter, err := NewFilter(config)
	if err != nil {
		return []*github.Repository{}, err
	}

	client := config.GithubClient()

	var repos []*github.Repository
	options := &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 100}}
	for {
		repoPage, resp, err := client.Repositories.List(context.Background(), "", options)
		if err != nil {
			return []*github.Repository{}, err
		}

		for _, repo := range repoPage {
			if !filter.Match(*repo.FullName) {
				repos = append(repos, repo)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	return repos, nil
}
