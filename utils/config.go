package utils

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
)

// Config definition for syncgithub
type Config struct {
	AuthToken  string   `json:"auth_token"`
	TargetPath string   `json:"target_path"`
	Excludes   []string `json:"excludes"`
}

func getDefaultConfigPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(dir, ".syncgithub"), nil
}

// LoadConfig reads the config from a file
func LoadConfig(fileArg string) (Config, error) {
	var c Config
	var err error

	file := fileArg
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
