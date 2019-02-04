package godotty

import (
	"github.com/BurntSushi/toml"
)

func LoadFromFile(file string) (*GodottyConfig, error) {
	var config GodottyConfig
	if _, err := toml.DecodeFile(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
