package godotty

import (
	"os"

	"github.com/BurntSushi/toml"
)

func LoadFromFile(filepath string) (*GodottyConfig, error) {
	var config GodottyConfig
	if _, err := toml.DecodeFile(filepath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (config *GodottyConfig) SaveFile(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := toml.NewEncoder(file)
	return enc.Encode(config)
}
