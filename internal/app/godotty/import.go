package godotty

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func Import(dottyDir string, files []string) error {
	godotty := New()
	godotty.Dir = dottyDir
	if err := godotty.LoadConfig(); err != nil {
		if err, ok := err.(*os.PathError); !ok {
			return err
		}
	}

	for _, file := range files {
		if err := godotty.Import(file); err != nil {
			return err
		}
	}

	return godotty.SaveConfig()
}

func (godotty *Godotty) Import(file string) error {
	dottyfile, err := toDottyfile(file)
	if err != nil {
		return err
	}

	// Source and Destination are swapped in this case as, relative to Godotty,
	// the Source is the Godotty repository and the Destination is the host
	sourceFile, err := godotty.Fs.Open(file)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	dest := filepath.Join(godotty.Dir, dottyfile.Source)
	if err := godotty.Fs.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	destFile, err := godotty.Fs.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}
	if err := godotty.Fs.Chmod(dest, sourceFileInfo.Mode()); err != nil {
		return err
	}

	replaced := false
	for i := range godotty.Config.Dottyfiles {
		if godotty.Config.Dottyfiles[i].Source == dottyfile.Source {
			godotty.Config.Dottyfiles[i] = dottyfile
			replaced = true
			break
		}
	}
	if !replaced {
		godotty.Config.Dottyfiles = append(godotty.Config.Dottyfiles, dottyfile)
	}

	return destFile.Sync()
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
