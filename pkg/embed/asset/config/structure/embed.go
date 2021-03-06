package structure

type EmbedConfig struct {
	Path     string `mapstructure:"path" yaml:"path" json:"path"`
	Package  string `mapstructure:"pkg" yaml:"pkg" json:"pkg"`
	Archive  string `mapstructure:"archive" yaml:"archive" json:"archive"`
	Compress string `mapstructure:"compress" yaml:"compress" json:"compress"`
	Items    []Item `mapstructure:"items" yaml:"items" json:"items"`
}

type Item struct {
	Name    string `mapstructure:"name" yaml:"name" json:"name"`
	File    string `mapstructure:"file" yaml:"file" json:"file"`
	Comment string `mapstructure:"comment" yaml:"comment" json:"comment"`
}
