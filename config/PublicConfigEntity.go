package config

type PublicConfig struct {
	Template  Template  `yaml:"template"`
	Directory Directory `yaml:"directory"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Template struct {
	Server   Server `yaml:"server"`
	FilePath string `yaml:"filepath"`
}

type Directory struct {
	FilePath string `yaml:"filepath"`
}
