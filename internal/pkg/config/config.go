package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/charlieegan3/subpub/internal/pkg/destinations"
	"github.com/charlieegan3/subpub/internal/pkg/sources"
	"github.com/charlieegan3/subpub/internal/pkg/sub"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Jobs []Job
}

type Job struct {
	Mappings      []Mapping
	Substitutions []sub.Substitution
}

type Mapping struct {
	Source      sources.Source
	Destination destinations.Destination
}

func Load(configFilePath string) (cfg Config, err error) {
	bytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %v", err)
	}

	rawData := make(map[string]interface{})
	err = yaml.Unmarshal(bytes, &rawData)
	if err != nil {
		return cfg, fmt.Errorf("cannot unmarshal data to raw map[string]interface{}: %v", err)
	}

	err = loadJobs(rawData, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to unmarshal jobs: %s", err)
	}

	return cfg, err
}

func loadJobs(rawData map[string]interface{}, cfg *Config) (err error) {
	jobs, ok := rawData["jobs"].([]interface{})
	if !ok {
		return fmt.Errorf("cannot unmarshal jobs to []interface{}: %#v", rawData)
	}

	for _, j := range jobs {
		rawJob, ok := j.(map[string]interface{})
		if !ok {
			return fmt.Errorf("cannot unmarshal job to map[string]interface{}: %#v", rawData)
		}

		var job Job

		err = loadSubstitutions(rawJob, &job)
		if err != nil {
			return fmt.Errorf("failed to unmarshal substitutions: %s", err)
		}
		err = loadMappings(rawJob, &job)
		if err != nil {
			return fmt.Errorf("failed to unmarshal mappings: %s", err)
		}

		cfg.Jobs = append(cfg.Jobs, job)
	}

	return nil
}

func loadMappings(rawData map[string]interface{}, job *Job) (err error) {
	rawMappings, ok := rawData["mappings"].([]interface{})
	if !ok {
		return fmt.Errorf("cannot unmarshal mappings to []interface{}: %#v", rawData)
	}

	for _, m := range rawMappings {
		rawMapping, ok := m.(map[string]interface{})
		if !ok {
			return fmt.Errorf("cannot unmarshal mapping to map[string]interface{}: %#v", m)
		}

		source, ok := rawMapping["source"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid or missing value for mapping.source: %#v", rawMapping)
		}

		destination, ok := rawMapping["destination"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid or missing value for mapping.destination: %#v", rawMapping)
		}

		var mapping Mapping
		err = loadSource(source, &mapping)
		if err != nil {
			return fmt.Errorf("failed to load mapping source: %s", err)
		}
		err = loadDestination(destination, &mapping)
		if err != nil {
			return fmt.Errorf("failed to load mapping destination: %s", err)
		}

		job.Mappings = append(job.Mappings, mapping)
	}

	return nil
}

func loadSource(rawData map[string]interface{}, mapping *Mapping) (err error) {
	sourceType, ok := rawData["type"].(string)
	if !ok {
		return fmt.Errorf("non string value for source.type: %#v", rawData)
	}

	if sourceType != "http" {
		return fmt.Errorf("source.type %s is not supported", sourceType)
	}

	sourceURL, ok := rawData["url"].(string)
	if !ok {
		return fmt.Errorf("non string value for source.url: %#v", rawData)
	}

	mapping.Source = &sources.SourceHTTP{URL: sourceURL}

	return nil
}

func loadDestination(rawData map[string]interface{}, mapping *Mapping) (err error) {
	destinationType, ok := rawData["type"].(string)
	if !ok {
		return fmt.Errorf("non string value for destinationType.type: %#v", rawData)
	}

	if destinationType != "gcsBucketObject" {
		return fmt.Errorf("source.type %s is not supported", destinationType)
	}

	bucketName, ok := rawData["bucket_name"].(string)
	if !ok {
		return fmt.Errorf("non string value for destination.bucket_name: %#v", rawData)
	}

	objectName, ok := rawData["object_name"].(string)
	if !ok {
		return fmt.Errorf("non string value for destination.object_name: %#v", rawData)
	}

	mapping.Destination = &destinations.DestinationGCSBucketObject{
		BucketName: bucketName,
		ObjectName: objectName,
	}

	return nil
}

func loadSubstitutions(rawData map[string]interface{}, job *Job) (err error) {
	subs, ok := rawData["substitutions"].([]interface{})
	if !ok {
		return fmt.Errorf("cannot unmarshal substitutions to []interface{}: %#v", rawData)
	}

	for _, s := range subs {
		rawSub, ok := s.(map[string]interface{})
		if !ok {
			return fmt.Errorf("cannot unmarshal substitution to map[string]interface{}: %#v", s)
		}

		subType, ok := rawSub["type"].(string)
		if !ok {
			return fmt.Errorf("non string value for substitution.type: %#v", s)
		}

		subFind, ok := rawSub["find"].(string)
		if !ok {
			return fmt.Errorf("non string value for substitution.replace: %#v", s)
		}

		subReplace, ok := rawSub["replace"].(string)
		if !ok {
			return fmt.Errorf("non string value for substitution.replace: %#v", s)
		}

		switch subType {
		case "string":
			job.Substitutions = append(job.Substitutions, &sub.SubstitutionString{
				Find:    subFind,
				Replace: subReplace,
			})
		case "regex":
			re, err := regexp.Compile(subFind)
			if err != nil {
				return fmt.Errorf("failed to parse find value as regexp: %s", subFind)
			}
			job.Substitutions = append(job.Substitutions, &sub.SubstitutionRegex{
				Find:    re,
				Replace: subReplace,
			})
		default:
			return fmt.Errorf("value %s for subType is not supported", subType)
		}
	}

	return nil
}
