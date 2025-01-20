package core

import (
	"FastGin/pkg/config"
	"gopkg.in/yaml.v3"
	"os"
)

func ReadConfig(FilePath string) (cfg *config.Config) {
	cfg = new(config.Config)
	//fmt.Println(FilePath) 打印配置文件名
	byteData, err := os.ReadFile(FilePath)
	if err != nil {
		//logg.Log.Errorf("读取配置文件失败: %w", err)
		return
	}

	err = yaml.Unmarshal(byteData, cfg)
	if err != nil {
		//logg.Log.Errorf("解析配置文件失败: %w", err)
		return
	}

	return cfg
}
