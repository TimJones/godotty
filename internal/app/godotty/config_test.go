package godotty

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

const testDataDir = "fixtures"

func helperLoadBytes(t *testing.T, name string) []byte {
	bytes, err := ioutil.ReadFile(filepath.Join(testDataDir, name))
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

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

func TestSaveConfig(t *testing.T) {
	testTable := []struct {
		godotty  *Godotty
		expected []byte
	}{
		{
			godotty: &Godotty{
				Fs:   afero.NewMemMapFs(),
				Dir:  ".",
				File: "single.toml",
				Config: DottyConfig{
					Dottyfiles: []Dottyfile{
						{
							Source:      "xinitrc",
							Destination: "~/.xinitrc",
						},
					},
				},
			},
			expected: helperLoadBytes(t, "single.toml"),
		},
		{
			godotty: &Godotty{
				Fs:   afero.NewMemMapFs(),
				Dir:  ".",
				File: "multiple.toml",
				Config: DottyConfig{
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
			expected: helperLoadBytes(t, "multiple.toml"),
		},
	}
	for _, test := range testTable {
		if err := test.godotty.SaveConfig(); err != nil {
			t.Error(err)
		}
		actual, err := afero.ReadFile(test.godotty.Fs, test.godotty.File)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(test.expected, actual) {
			t.Errorf("saving config file %s failed.\nExpected:\n%s\nGot:\n%s\n", test.godotty.File, test.expected, actual)
		}
	}
}
