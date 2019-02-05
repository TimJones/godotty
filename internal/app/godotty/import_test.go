package godotty

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestImport(t *testing.T) {
	testTable := []struct {
		importFiles []string
		godotty     *Godotty
		expected    DottyConfig
	}{
		{
			importFiles: []string{
				"~/.xinitrc",
			},
			godotty: &Godotty{},
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
			importFiles: []string{
				"~/.bashrc",
				"~/.zshrc",
			},
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
					{
						Source:      "zshrc",
						Destination: "~/.zshrc",
					},
				},
			},
		},
	}
	for _, test := range testTable {
		err := test.godotty.Import(test.importFiles)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(test.expected, test.godotty.Config) {
			t.Errorf("importing %s failed:\n%s", test.importFiles, cmp.Diff(test.expected, test.godotty.Config))
		}
	}
}

func TestSimplifyDotfilepath(t *testing.T) {
	testTable := []struct {
		dotfilepath string
		expected    string
	}{
		{
			dotfilepath: "~/.xinitrc",
			expected:    "xinitrc",
		},
		{
			dotfilepath: "~/.ssh/config",
			expected:    "ssh/config",
		},
	}
	for _, test := range testTable {
		actual, err := simplifyDotfilepath(test.dotfilepath)
		if err != nil {
			t.Error(err)
		}
		if actual != test.expected {
			t.Errorf("simplifyDotfilepath %s failed: expected %s, got %s", test.dotfilepath, test.expected, actual)
		}
	}
}
