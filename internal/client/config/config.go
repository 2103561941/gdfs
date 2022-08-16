package config

import (
	"github.com/cyb0225/gdfs/internal/pkg/log"
	"github.com/spf13/viper"
)

var (
	Cfg *Config
)

type Config struct {
	NamenodeAddr string         `yaml:"NamenodeAddr"`
	Log          *log.LogConfig `yaml:"Log"`
	ChunkSize    int64          `yaml:"ChunkSize"`
}

// load configs
func NewConfig() error {
	Cfg = &Config{}

	vp := viper.New()
	vp.AddConfigPath(".")
	vp.AddConfigPath("config/")
	vp.SetConfigName("client")
	err := vp.ReadInConfig()
	if err != nil {
		return err
	}

	if err := vp.UnmarshalKey("Log", &Cfg.Log); err != nil {
		return err
	}

	Cfg.NamenodeAddr = vp.GetString("NamenodeAddr")
	Cfg.ChunkSize = vp.GetInt64("ChunkSize")

	return nil
}
