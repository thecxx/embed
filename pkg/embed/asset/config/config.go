package config

import (
	"encoding/json"
	"errors"

	"github.com/thecxx/embed/pkg/embed/asset/config/structure"
	"gopkg.in/yaml.v2"
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
	switch format {
	// yaml
	case "yaml":
		return yaml.Unmarshal(buffer, rawVal)
	// json
	case "json":
		return json.Unmarshal(buffer, rawVal)
	}
	return errors.New("format not supported")
}
