package godotty

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

const testDataDir = "fixtures"

func TestLoadConfig(t *testing.T) {
	testTable := []struct {
		godotty  *Godotty
		expected DottyConfig
	}{
		{
			godotty: &Godotty{
				Fs:   afero.NewOsFs(),
				Dir:  testDataDir,
				File: "single.toml",
			},
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
			godotty: &Godotty{
				Fs:   afero.NewOsFs(),
				Dir:  testDataDir,
				File: "multiple.toml",
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
		if err := test.godotty.LoadConfig(); err != nil {
			t.Error(err)
		}
		if !cmp.Equal(test.expected, test.godotty.Config) {
			t.Errorf("loading file %s failed:\n%s", test.godotty.File, cmp.Diff(test.expected, test.godotty.Config))
		}
	}
}
