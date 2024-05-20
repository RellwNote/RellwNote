package config

import (
	"github.com/RellwNote/RellwNote/log"
	"gopkg.in/yaml.v3"
	"os"
)

const ConfigFilePath = "./config.yaml"
const TemplateDir = "./templates"

var Config Model

func init() {
	LoadConfig()
}

func LoadConfig() {
	configBytes, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		log.Error.Println("读取配置文件:", ConfigFilePath, "失败！", err)
		return
	}
	err = yaml.Unmarshal(configBytes, &Config)
	if err != nil {
		log.Error.Println("解析配置文件失败:", err)
		return
	}
}

type Model struct {
	PreviewServer string `yaml:"Preview server"`
	LibraryPath   string `yaml:"Library path"`
}
