package godotty

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const testDataDir = "fixtures"

func TestLoadConfig(t *testing.T) {
	testTable := []struct {
		dataFile string
		expected *GodottyConfig
	}{
		{
			dataFile: filepath.Join(testDataDir, "single.toml"),
			expected: &GodottyConfig{
				Dottyfiles: []Dottyfile{
					{
						Source:      "xinitrc",
						Destination: "~/.xinitrc",
					},
				},
			},
		},
		{
			dataFile: filepath.Join(testDataDir, "multiple.toml"),
			expected: &GodottyConfig{
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
		actual, err := LoadFromFile(test.dataFile)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(test.expected, actual) {
			t.Errorf("loading file %s failed:\n%s", test.dataFile, cmp.Diff(test.expected, actual))
		}
	}
}
