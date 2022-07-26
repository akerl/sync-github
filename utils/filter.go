package utils

import (
	"regexp"
)

// Filter simplifies running a repo name over a list of exclusion patterns
type Filter struct {
	Patterns []*regexp.Regexp
}

// NewFilter creates a new filter object with a list of regexp patterns
func NewFilter(config Config) (*Filter, error) {
	f := &Filter{}
	f.Patterns = make([]*regexp.Regexp, len(config.Excludes))
	for i, x := range config.Excludes {
		re, err := regexp.Compile(x)
		if err != nil {
			return f, err
		}
		f.Patterns[i] = re
	}
	return f, nil
}

// Match checks if the repo name matches any of the filter's patterns
func (f *Filter) Match(repo string) bool {
	for _, x := range f.Patterns {
		if x.MatchString(repo) {
			return true
		}
	}
	return false
}
