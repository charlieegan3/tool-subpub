package sub

import (
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestSubstitutionString(t *testing.T) {
	testCases := []struct {
		Description  string
		Substitution SubstitutionString
		Input        string
		Output       string
	}{
		{
			Description: "simple string substitution",
			Substitution: SubstitutionString{
				Find:    "foo",
				Replace: "bar",
			},
			Input:  "foo foo",
			Output: "bar bar",
		},
		{
			Description: "simple string substitution",
			Substitution: SubstitutionString{
				Find:    "bar",
				Replace: "foo",
			},
			Input:  "foobarfoo",
			Output: "foofoofoo",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			resultReader, err := testCase.Substitution.Run(io.NopCloser(strings.NewReader(testCase.Input)))
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			result, err := io.ReadAll(resultReader)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			td.Cmp(t, string(result), testCase.Output)
		})
	}
}

func TestSubstitutionRegex(t *testing.T) {
	testCases := []struct {
		Description  string
		Substitution SubstitutionRegex
		Input        string
		Output       string
	}{
		{
			Description: "group based replacement",
			Substitution: SubstitutionRegex{
				Find:    regexp.MustCompile("fo*"),
				Replace: "bar",
			},
			Input:  "foo foo",
			Output: "bar bar",
		},
		{
			Description: "group based replacement",
			Substitution: SubstitutionRegex{
				Find:    regexp.MustCompile("(a)(b)"),
				Replace: "$2$1",
			},
			Input:  "ab",
			Output: "ba",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			resultReader, err := testCase.Substitution.Run(io.NopCloser(strings.NewReader(testCase.Input)))
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			result, err := io.ReadAll(resultReader)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			td.Cmp(t, string(result), testCase.Output)
		})
	}
}
