package config

import (
	"regexp"
	"testing"

	"github.com/charlieegan3/tool-subpub/internal/pkg/destinations"
	"github.com/charlieegan3/tool-subpub/internal/pkg/jobs"
	"github.com/charlieegan3/tool-subpub/internal/pkg/sources"
	"github.com/charlieegan3/tool-subpub/internal/pkg/sub"
	"github.com/maxatome/go-testdeep/td"
)

func TestLoadConfig(t *testing.T) {
	exampleConfigFilePath := "example.yaml"

	config, err := Load(exampleConfigFilePath)
	if err != nil {
		t.Fatalf("failed to read example config file: %s", err)
	}

	expectedConfig := Config{
		Jobs: []jobs.Job{
			{
				Mappings: []jobs.Mapping{
					{
						Source: &sources.SourceHTTP{
							URL: "https://example.com/file.txt",
						},
						Destination: &destinations.DestinationBucketObject{
							Path: "gs://example-bucket-name/file.txt",
						},
					},
				},
				Substitutions: []sub.Substitution{
					&sub.SubstitutionString{
						Find:    "foo",
						Replace: "bar",
					},
					&sub.SubstitutionRegex{
						Find:    regexp.MustCompile("(a)(b)(c)"),
						Replace: "$3$2$1",
					},
				},
			},
		},
	}

	td.Cmp(t, config, expectedConfig)
}
