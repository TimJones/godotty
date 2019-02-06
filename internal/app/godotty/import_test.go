package godotty

import (
	"bytes"
	"os"
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
		fileMode       os.FileMode
		godotty        *Godotty
		expectedFile   string
		expectedConfig DottyConfig
	}{
		{
			importFile:  "~/.xinitrc",
			fileContent: []byte("xinitrc file content"),
			fileMode:    0644,
			godotty: &Godotty{
				Fs:  afero.NewMemMapFs(),
				Dir: DefaultDirectory,
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
			fileMode:    0644,
			godotty: &Godotty{
				Fs:  afero.NewMemMapFs(),
				Dir: DefaultDirectory,
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
		{
			importFile:  "~/.ssh/config",
			fileContent: []byte("ssh_config file content"),
			fileMode:    0644,
			godotty: &Godotty{
				Fs:  afero.NewMemMapFs(),
				Dir: DefaultDirectory,
			},
			expectedFile: "ssh/config",
			expectedConfig: DottyConfig{
				Dottyfiles: []Dottyfile{
					{
						Source:      "ssh/config",
						Destination: "~/.ssh/config",
					},
				},
			},
		},
		{
			importFile:  "~/.local/bin/script",
			fileContent: []byte("script content"),
			fileMode:    0755,
			godotty: &Godotty{
				Fs:  afero.NewMemMapFs(),
				Dir: DefaultDirectory,
			},
			expectedFile: "local/bin/script",
			expectedConfig: DottyConfig{
				Dottyfiles: []Dottyfile{
					{
						Source:      "local/bin/script",
						Destination: "~/.local/bin/script",
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
		if err = afero.WriteFile(test.godotty.Fs, test.importFile, test.fileContent, test.fileMode); err != nil {
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
		fileInfo, err := test.godotty.Fs.Stat(filepath.Join(test.godotty.Dir, test.expectedFile))
		if err != nil {
			t.Fatal(err)
		}
		actualMode := fileInfo.Mode().Perm()
		if test.fileMode != actualMode {
			t.Errorf("importing %s failed, file mode differs.\nExpected:\n%s\nGot:\n%s\n", test.importFile, test.fileMode, actualMode)
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
