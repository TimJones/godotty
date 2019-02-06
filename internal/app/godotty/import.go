package godotty

import (
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func (godotty *Godotty) Import(files []string) error {
	for _, file := range files {
		dottyfile, err := toDottyfile(file)
		if err != nil {
			return err
		}

		godotty.Config.Dottyfiles = append(godotty.Config.Dottyfiles, dottyfile)
	}
	return nil
}

func toDottyfile(pathfile string) (Dottyfile, error) {
	var src, dst string

	src, err := homedir.Expand(pathfile)
	if err != nil {
		return Dottyfile{}, err
	}

	homedirpath, err := homedir.Dir()
	if err != nil {
		return Dottyfile{}, err
	}

	if strings.HasPrefix(src, homedirpath) {
		src, err = filepath.Rel(homedirpath, src)
		if err != nil {
			return Dottyfile{}, err
		}
		dst = filepath.Join("~", src)
	} else {
		src, err = filepath.Abs(src)
		if err != nil {
			return Dottyfile{}, err
		}
		dst = src
		src, err = filepath.Rel(string(filepath.Separator), src)
		if err != nil {
			return Dottyfile{}, err
		}
	}
	src = strings.TrimPrefix(src, ".")
	return Dottyfile{Source: src, Destination: dst}, nil
}
