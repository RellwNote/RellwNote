package config

import (
	"github.com/RellwNote/RellwNote/log"
	"gopkg.in/yaml.v3"
	"os"
)

const publicConfigFilePath = "./config.yml"

var GetPublicConfig PublicConfig

func init() {
	configBytes, err := os.ReadFile(publicConfigFilePath)
	if err != nil {
		log.Error.Println("读取配置文件:", publicConfigFilePath, "失败！", err)
		return
	}
	err = yaml.Unmarshal(configBytes, &GetPublicConfig)
	if err != nil {
		log.Error.Println("解析配置文件失败:", err)
		return
	}
}
