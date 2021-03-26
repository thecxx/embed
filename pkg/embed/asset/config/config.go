package config

import (
	"bytes"

	"github.com/spf13/viper"
	"github.com/thecxx/embed/pkg/embed/asset/config/structure"
)

var (
	Embed = (*structure.EmbedConfig)(nil)
)

// InitEmbedConfig initializes embed configuration.
func InitEmbedConfig(buffer []byte, format string) error {
	return loadConfig(buffer, format, &Embed)
}

// loadConfig loads configuration from the `buffer`, and then save it to `rawVal`.
func loadConfig(buffer []byte, format string, rawVal interface{}) error {
	v := viper.New()
	{
		v.SetConfigType(format)
		v.ReadConfig(bytes.NewReader(buffer))
	}
	return v.Unmarshal(&rawVal)
}
