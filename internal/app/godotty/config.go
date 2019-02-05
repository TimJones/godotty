package godotty

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/afero"
)

func New() *Godotty {
	return &Godotty{
		Fs:   afero.NewOsFs(),
		Dir:  DefaultDirectory,
		File: DefaultConfigFilename,
	}
}

func (godotty *Godotty) LoadConfig() error {
	configFile := filepath.Join(godotty.Dir, godotty.File)
	if _, err := toml.DecodeFile(configFile, &godotty.Config); err != nil {
		return err
	}
	return nil
}

func (godotty *Godotty) SaveConfig() error {
	configFile := filepath.Join(godotty.Dir, godotty.File)
	file, err := godotty.Fs.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := toml.NewEncoder(file)
	return enc.Encode(godotty.Config)
}
