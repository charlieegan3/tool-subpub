package sub

import (
	"io"
	"regexp"
)

type Substitution interface {
	Run(io.ReadCloser) io.ReadCloser
}

type SubstitutionString struct {
	Find    string
	Replace string
}

func (s *SubstitutionString) Run(input io.ReadCloser) (output io.ReadCloser) {
	return nil
}

type SubstitutionRegex struct {
	Find    *regexp.Regexp
	Replace string
}

func (s *SubstitutionRegex) Run(input io.ReadCloser) (output io.ReadCloser) {
	return nil
}
