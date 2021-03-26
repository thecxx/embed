package structure

type EmbedConfig struct {
	Path     string      `mapstructure:"path" yaml:"path" json:"path"`
	Package  string      `mapstructure:"pkg" yaml:"pkg" json:"pkg"`
	Compress string      `mapstructure:"compress" yaml:"compress" json:"compress"`
	Items    []EmbedItem `mapstructure:"items" yaml:"items" json:"items"`
}

type EmbedItem struct {
	Name string `mapstructure:"name" yaml:"name" json:"name"`
	File string `mapstructure:"file" yaml:"file" json:"file"`
}
