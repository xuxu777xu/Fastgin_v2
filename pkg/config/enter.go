package config

import "flag"

// Config 应用配置结构
type Config struct {
	DB     DBConfig  `yaml:"db"`
	Server Server    `yaml:"server"`
	Log    LogConfig `yaml:"log"`
}

// Server 服务器配置
type Server struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type FlagOptions struct {
	File string
}

var Options FlagOptions

func RunSettingFile() {

	flag.StringVar(&Options.File, "f", "settings.yaml", "配置文件路径")
	flag.Parse()
}
