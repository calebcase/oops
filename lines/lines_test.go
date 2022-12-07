package lines_test

import (
	"testing"

	"github.com/calebcase/oops/lines"
	"github.com/stretchr/testify/require"
)

func TestPrintf(t *testing.T) {
	ls := lines.Sprintf("a\nb\nc\n")
	require.Equal(t, []string{"a", "b", "c", ""}, ls)
}

func TestIndent(t *testing.T) {
	type TC struct {
		Name   string
		Lines  []string
		Prefix string
		Skip   int
		Expect []string
	}

	tcs := []TC{
		{
			Name: "indent",
			Lines: []string{
				"a",
				"b",
				"c",
			},
			Prefix: "..",
			Skip:   0,
			Expect: []string{
				"..a",
				"..b",
				"..c",
			},
		},
		{
			Name: "skip 1",
			Lines: []string{
				"a",
				"b",
				"c",
			},
			Prefix: "..",
			Skip:   1,
			Expect: []string{
				"a",
				"..b",
				"..c",
			},
		},
		{
			Name: "skip 2",
			Lines: []string{
				"a",
				"b",
				"c",
			},
			Prefix: "..",
			Skip:   2,
			Expect: []string{
				"a",
				"b",
				"..c",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			require.Equal(t, tc.Expect, lines.Indent(tc.Lines, tc.Prefix, tc.Skip))
		})
	}
}

func TestTrimTrailing(t *testing.T) {
	type TC struct {
		Name   string
		Lines  []string
		Cutset string
		Expect []string
	}

	tcs := []TC{
		{
			Name:   "len 0",
			Lines:  []string{},
			Cutset: "",
			Expect: []string{},
		},
		{
			Name: "nothing",
			Lines: []string{
				"a",
				"b",
				"c",
			},
			Cutset: "",
			Expect: []string{
				"a",
				"b",
				"c",
			},
		},
		{
			Name: "empty 1",
			Lines: []string{
				"a",
				"b",
				"c",
				"",
			},
			Cutset: "",
			Expect: []string{
				"a",
				"b",
				"c",
			},
		},
		{
			Name: "empty 2",
			Lines: []string{
				"a",
				"b",
				"c",
				"",
				"",
			},
			Cutset: "",
			Expect: []string{
				"a",
				"b",
				"c",
			},
		},
		{
			Name: "whitespace",
			Lines: []string{
				"a",
				"b",
				"c",
				" \t",
			},
			Cutset: " \t",
			Expect: []string{
				"a",
				"b",
				"c",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			require.Equal(t, tc.Expect, lines.TrimTrailing(tc.Lines, tc.Cutset))
		})
	}
}
