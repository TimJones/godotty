package godotty

type GodottyConfig struct {
	Dottyfiles []Dottyfile `toml:"dotfile"`
}

type Dottyfile struct {
	Source      string
	Destination string
}
