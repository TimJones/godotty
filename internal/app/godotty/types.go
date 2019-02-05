package godotty

import "github.com/spf13/afero"

const DefaultDirectory = "."
const DefaultConfigFilename = "godotty.toml"

type Godotty struct {
	Fs     afero.Fs
	Dir    string
	File   string
	Config DottyConfig
}

type DottyConfig struct {
	Dottyfiles []Dottyfile `toml:"dotfile"`
}

type Dottyfile struct {
	Source      string `toml:"source"`
	Destination string `toml:"destination"`
}
