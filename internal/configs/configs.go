package configs

import (
	"github.com/spf13/viper"
	"os"
)

type Configs struct {
	Version      string   `mapstructure:"version" json:"version" yaml:"version"`
	ProtoVersion int      `mapstructure:"proto_version" json:"proto_version" yaml:"proto_version"`
	Mode         string   `mapstructure:"mode" json:"mode" yaml:"mode"`
	Modules      []string `mapstructure:"modules" json:"modules" yaml:"modules"`
	Port         string   `mapstructure:"port" json:"port" yaml:"port"`
}

func (cfg *Configs) Address() string {
	return ":" + cfg.Port
}

func (cfg *Configs) fillDefaults() {
	if cfg.Version == "" {
		cfg.Version = "6.0.3"
	}

	if cfg.ProtoVersion == 0 {
		cfg.ProtoVersion = 3
	}

	if cfg.Mode == "" {
		cfg.Mode = "standalone"
	}

	if cfg.Port == "" {
		cfg.Port = "6379"
	}
}

func LoadConfigsOrDefaults(path string) (*Configs, error) {
	if path == "" {
		return getDefaultConfigs()
	}
	if _, err := os.Stat(path); err != nil {
		return getDefaultConfigs()
	}

	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var configs Configs

	err = viper.Unmarshal(&configs)
	if err != nil {
		return nil, err
	}

	configs.fillDefaults()

	return &configs, nil
}

func getDefaultConfigs() (*Configs, error) {
	var configs Configs

	configs.fillDefaults()

	return &configs, nil
}
