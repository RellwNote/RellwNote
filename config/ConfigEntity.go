package config

type PublicConfig struct {
	Template  Template  `yaml:"template"`
	Directory Directory `yaml:"directory"`
	Temp      Temp      `yaml:"temp"`
}

type Temp struct {
	Server Server `yaml:"server"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Template struct {
	FilePath string `yaml:"filepath"`
}

type Directory struct {
	FilePath string `yaml:"filepath"`
}
