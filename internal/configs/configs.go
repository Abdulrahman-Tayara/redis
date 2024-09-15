package configs

type Configs struct {
	Version      string   `mapstructure:"version"`
	ProtoVersion int      `mapstructure:"proto_version"`
	Mode         string   `mapstructure:"mode"`
	Modules      []string `mapstructure:"modules"`
}
