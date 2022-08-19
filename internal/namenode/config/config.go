package config

import (
	"os"

	"github.com/cyb0225/gdfs/internal/pkg/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	Cfg *Config
)

type Config struct {
	BackupN     int            `yaml:"BackupN"`
	Timeout     int            `yaml:"Timeout"`
	ChunkSize   int64          `yaml:"ChunkSize"`
	Addr        *Address       `yaml:"Address"`
	Log         *log.LogConfig `yaml:"Log"`
	StoragePath string         `yaml:"StoragePath"`
}

type Address struct {
	IP   string `yaml:"IP"`
	Port string `yaml:"Port"`
}

// load configs
func NewConfig() error {
	Cfg = &Config{}

	vp := viper.New()
	vp.AddConfigPath(".")
	vp.AddConfigPath("config/")
	vp.SetConfigName("namenode")
	err := vp.ReadInConfig()
	if err != nil {
		return err
	}

	if err := vp.UnmarshalKey("Address", &Cfg.Addr); err != nil {
		return err
	}

	if err := vp.UnmarshalKey("Log", &Cfg.Log); err != nil {
		return err
	}

	Cfg.BackupN = vp.GetInt("BackupN")
	Cfg.Timeout = vp.GetInt("Timeout")
	Cfg.ChunkSize = vp.GetInt64("ChunkSize")
	Cfg.StoragePath = vp.GetString("StoragePath")

	flag()
	Cfg.Log.LogPath = Cfg.Log.LogPath + Cfg.Addr.Port + "/"
	Cfg.StoragePath = Cfg.StoragePath + Cfg.Addr.Port + "/"
	_ = os.Mkdir(Cfg.StoragePath, 0777)
	return nil
}

func flag() {
	port := pflag.StringP("port", "p", Cfg.Addr.Port, "datanode port, --port=50050")
	pflag.Parse()
	Cfg.Addr.Port = *port
}
