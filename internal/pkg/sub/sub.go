package sub

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Substitution interface {
	Run(io.ReadCloser) (io.ReadCloser, error)
}

type SubstitutionString struct {
	Find    string
	Replace string
}

func (s *SubstitutionString) Run(input io.ReadCloser) (output io.ReadCloser, err error) {
	bytes, err := io.ReadAll(input)
	if err != nil {
		return output, fmt.Errorf("failed to read input: %s", err)
	}
	err = input.Close()
	if err != nil {
		return output, fmt.Errorf("failed to close input: %s", err)
	}

	result := strings.ReplaceAll(string(bytes), s.Find, s.Replace)

	return io.NopCloser(strings.NewReader(result)), nil
}

type SubstitutionRegex struct {
	Find    *regexp.Regexp
	Replace string
}

func (s *SubstitutionRegex) Run(input io.ReadCloser) (output io.ReadCloser, err error) {
	bytes, err := io.ReadAll(input)
	if err != nil {
		return output, fmt.Errorf("failed to read input: %s", err)
	}
	err = input.Close()
	if err != nil {
		return output, fmt.Errorf("failed to close input: %s", err)
	}

	result := s.Find.ReplaceAllString(string(bytes), s.Replace)

	return io.NopCloser(strings.NewReader(result)), nil
}
