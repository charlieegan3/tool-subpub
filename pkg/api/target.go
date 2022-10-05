package api

import "github.com/charlieegan3/tool-subpub/internal/pkg/sub"

type Target struct {
	Name          string
	URL           string
	Substitutions []sub.Substitution
}
