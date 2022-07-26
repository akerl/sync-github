package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"github.com/google/go-github/v45/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

type config struct {
	AuthToken string `json:"auth_token"`
}

func getDefaultConfigPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(dir, ".syncgithub"), nil
}

func loadConfig(file string) (config, error) {
	var c config
	var err error

	if file == "" {
		file, err = getDefaultConfigPath()
		if err != nil {
			return c, err
		}
	}

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(contents, &c)
	return c, err
}

func githubClient(string token) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tokenClient := oauth2.NewClient(context.Background(), tokenSource)
	return github.NewClient(tokenClient)
}

func updateRunner(cmd *cobra.Command, _ []string) error {
	c, err := cmd.Flags().GetBool("config")
	if err != nil {
		return err
	}

	config, err := loadConfig(c)
	if err != nil {
		return err
	}

	client := githubClient(config.AuthToken)

	options := &github.RepositoryListOptions{ListOptions: {PerPage: 100}}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, "", options)
		if err != nil {
			return err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, x := range allRepos {
		fmt.Println(x.CloneURL)
	}
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update local copy of repos",
	RunE:  updateRunner,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("config", "c", "", "Config file")
}
