package godotty

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
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
		importFile     string
		fileContent    []byte
		godotty        *Godotty
		expectedFile   string
		expectedConfig DottyConfig
	}{
		{
			importFile:  "~/.xinitrc",
			fileContent: []byte("xinitrc file content"),
			godotty: &Godotty{
				Fs: afero.NewMemMapFs(),
			},
			expectedFile: "xinitrc",
			expectedConfig: DottyConfig{
				Dottyfiles: []Dottyfile{
					{
						Source:      "xinitrc",
						Destination: "~/.xinitrc",
					},
				},
			},
		},
		{
			importFile:  "~/.bashrc",
			fileContent: []byte("bashrc file content"),
			godotty: &Godotty{
				Fs: afero.NewMemMapFs(),
				Config: DottyConfig{
					Dottyfiles: []Dottyfile{
						{
							Source:      "xinitrc",
							Destination: "~/.xinitrc",
						},
					},
				},
			},
			expectedFile: "bashrc",
			expectedConfig: DottyConfig{
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
		// Setup
		homedirpath, err := homedir.Dir()
		if err != nil {
			t.Fatal(err)
		}
		if err = test.godotty.Fs.MkdirAll(homedirpath, 0755); err != nil {
			t.Fatal(err)
		}
		if err = afero.WriteFile(test.godotty.Fs, test.importFile, test.fileContent, 0644); err != nil {
			t.Fatal(err)
		}

		// Test
		if err = test.godotty.Import(test.importFile); err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(test.expectedConfig, test.godotty.Config) {
			t.Errorf("importing %s failed:\n%s", test.importFile, cmp.Diff(test.expectedConfig, test.godotty.Config))
		}
		actualContent, err := afero.ReadFile(test.godotty.Fs, filepath.Join(test.godotty.Dir, test.expectedFile))
		if err = test.godotty.Import(test.importFile); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(test.fileContent, actualContent) {
			t.Errorf("importing %s failed, file content differs.\nExpected:\n%s\nGot:\n%s\n", test.importFile, test.fileContent, actualContent)
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
