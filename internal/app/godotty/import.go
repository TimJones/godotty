package godotty

import (
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func (godotty *Godotty) Import(files []string) error {
	for _, file := range files {
		sourcefilepath, err := simplifyDotfilepath(file)
		if err != nil {
			return err
		}

		dottyfile := Dottyfile{
			Source:      sourcefilepath,
			Destination: file,
		}
		godotty.Config.Dottyfiles = append(godotty.Config.Dottyfiles, dottyfile)
	}
	return nil
}

func simplifyDotfilepath(dotfilepath string) (string, error) {
	dotfilepath, err := homedir.Expand(dotfilepath)
	if err != nil {
		return "", err
	}

	homedirpath, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(dotfilepath, homedirpath) {
		dotfilepath, err = filepath.Rel(homedirpath, dotfilepath)
		if err != nil {
			return "", err
		}
	}

	dotfilepath = strings.TrimPrefix(dotfilepath, ".")
	return dotfilepath, nil
}
