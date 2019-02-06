package godotty

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func helperRelToAbs(t *testing.T, path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

func helperRelToRoot(t *testing.T, path string) string {
	rel, err := filepath.Rel(string(filepath.Separator), helperRelToAbs(t, path))
	if err != nil {
		t.Fatal(err)
	}
	return rel
}

func TestImport(t *testing.T) {
	testTable := []struct {
		importFile string
		godotty    *Godotty
		expected   DottyConfig
	}{
		{
			importFile: "~/.xinitrc",
			godotty:    &Godotty{},
			expected: DottyConfig{
				Dottyfiles: []Dottyfile{
					{
						Source:      "xinitrc",
						Destination: "~/.xinitrc",
					},
				},
			},
		},
		{
			importFile: "~/.bashrc",
			godotty: &Godotty{
				Config: DottyConfig{
					Dottyfiles: []Dottyfile{
						{
							Source:      "xinitrc",
							Destination: "~/.xinitrc",
						},
					},
				},
			},
			expected: DottyConfig{
				Dottyfiles: []Dottyfile{
					{
						Source:      "xinitrc",
						Destination: "~/.xinitrc",
					},
					{
						Source:      "bashrc",
						Destination: "~/.bashrc",
					},
				},
			},
		},
	}
	for _, test := range testTable {
		err := test.godotty.Import(test.importFile)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(test.expected, test.godotty.Config) {
			t.Errorf("importing %s failed:\n%s", test.importFile, cmp.Diff(test.expected, test.godotty.Config))
		}
	}
}

func TestToDottyfile(t *testing.T) {
	testTable := []struct {
		filepath string
		expected Dottyfile
	}{
		{
			filepath: "~/.xinitrc",
			expected: Dottyfile{
				Source:      "xinitrc",
				Destination: "~/.xinitrc",
			},
		},
		{
			filepath: "~/.ssh/config",
			expected: Dottyfile{
				Source:      "ssh/config",
				Destination: "~/.ssh/config",
			},
		},
		{
			filepath: "/absolute/file",
			expected: Dottyfile{
				Source:      "absolute/file",
				Destination: "/absolute/file",
			},
		},
		{
			filepath: "../relative/file",
			expected: Dottyfile{
				Source:      helperRelToRoot(t, "../relative/file"),
				Destination: helperRelToAbs(t, "../relative/file"),
			},
		},
	}
	for _, test := range testTable {
		actual, err := toDottyfile(test.filepath)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(actual, test.expected) {
			t.Errorf("toDottyfile %s failed:\n%s", test.filepath, cmp.Diff(test.expected, actual))
		}
	}
}
