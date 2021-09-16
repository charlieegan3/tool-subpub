package jobs

import (
	"fmt"

	"github.com/charlieegan3/subpub/internal/pkg/destinations"
	"github.com/charlieegan3/subpub/internal/pkg/sources"
	"github.com/charlieegan3/subpub/internal/pkg/sub"
)

type Job struct {
	Mappings      []Mapping
	Substitutions []sub.Substitution
}

type Mapping struct {
	Source      sources.Source
	Destination destinations.Destination
}

func Run(j *Job) (err error) {
	for i, m := range j.Mappings {
		reader, err := m.Source.Fetch()
		if err != nil {
			return fmt.Errorf("mapping %d: failed to read from source: %s", i, err)
		}

		for _, s := range j.Substitutions {
			reader, err = s.Run(reader)
			if err != nil {
				return fmt.Errorf("mapping %d: make substitution: %s", i, err)
			}
		}

		err = m.Destination.Put(reader)
		if err != nil {
			return fmt.Errorf("mapping %d: failed put data to destination: %s", i, err)
		}
	}
	return nil
}
