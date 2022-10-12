package tool

import (
	"database/sql"
	"embed"
	"fmt"
	"regexp"

	"github.com/Jeffail/gabs/v2"
	"github.com/charlieegan3/toolbelt/pkg/apis"
	"github.com/gorilla/mux"

	"github.com/charlieegan3/tool-subpub/internal/pkg/sub"
	"github.com/charlieegan3/tool-subpub/pkg/api"
	"github.com/charlieegan3/tool-subpub/pkg/tool/handlers"
)

// SubPub is a tool for substituting text in http accessible content and publishing it
type SubPub struct {
	targets map[string]api.Target
}

func (d *SubPub) Name() string {
	return "subpub"
}

func (d *SubPub) FeatureSet() apis.FeatureSet {
	return apis.FeatureSet{
		Config: true,
		HTTP:   true,
	}
}

func (d *SubPub) HTTPPath() string {
	return "subpub"
}

func (d *SubPub) SetConfig(config map[string]any) error {
	d.targets = make(map[string]api.Target)

	cfg := gabs.Wrap(config)

	configTargets, ok := cfg.Path("targets").Data().([]interface{})
	if !ok {
		return fmt.Errorf("missing required config path: targets (array)")
	}

	for _, t := range configTargets {
		t, ok := t.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid target config, not map[string]interface{}")
		}

		name, ok := t["name"].(string)
		if !ok {
			return fmt.Errorf("invalid target config, missing name (string)")
		}

		url, ok := t["url"].(string)
		if !ok {
			return fmt.Errorf("invalid target config, missing url (string)")
		}

		var loadedSubs []sub.Substitution
		substitutions, ok := t["substitutions"].([]interface{})
		if !ok {
			return fmt.Errorf("invalid target config, missing substitutions (array)")
		}

		for _, s := range substitutions {
			s, ok := s.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid substitution config, not map[string]interface{}")
			}

			subType, ok := s["type"].(string)
			if !ok {
				return fmt.Errorf("invalid substitution config, missing type (string)")
			}

			find, ok := s["find"].(string)
			if !ok {
				return fmt.Errorf("invalid substitution config, missing find (string)")
			}

			replace, ok := s["replace"].(string)
			if !ok {
				return fmt.Errorf("invalid substitution config, missing replace (string)")
			}

			if subType == "regex" {
				re, err := regexp.Compile(find)
				if err != nil {
					return fmt.Errorf("invalid substitution config, invalid regex %q: %w", find, err)
				}
				loadedSubs = append(loadedSubs, &sub.SubstitutionRegex{
					Find:    re,
					Replace: replace,
				})
			} else if subType == "string" {
				loadedSubs = append(loadedSubs, &sub.SubstitutionString{
					Find:    find,
					Replace: replace,
				})
			} else {
				return fmt.Errorf("invalid substitution config, invalid type %q", subType)
			}
		}

		d.targets[name] = api.Target{
			Name:          name,
			URL:           url,
			Substitutions: loadedSubs,
		}
	}

	return nil
}

func (d *SubPub) DatabaseMigrations() (*embed.FS, string, error) {
	return &embed.FS{}, "migrations", nil
}

func (d *SubPub) DatabaseSet(db *sql.DB) {}

func (d *SubPub) HTTPAttach(router *mux.Router) error {
	router.HandleFunc(
		"/target/{target}",
		handlers.BuildGetHandler(d.targets),
	).Methods("GET")

	return nil
}

func (d *SubPub) Jobs() ([]apis.Job, error) {
	return []apis.Job{}, nil
}
func (d *SubPub) ExternalJobsFuncSet(f func(job apis.ExternalJob) error) {}
