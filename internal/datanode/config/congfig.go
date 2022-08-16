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
	StoragePath  string         `yaml:"StoragePath"`
	NamenodeAddr string         `yaml:"NamenodeAddr"`
	Addr         *Address       `yaml:"Address"`
	Log          *log.LogConfig `yaml:"Log"`
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
	vp.SetConfigName("datanode")
	_ = vp.BindPFlag("port", pflag.Lookup("port"))
	err := vp.ReadInConfig()
	if err != nil {
		return err
	}

	if err := vp.UnmarshalKey("Address", &Cfg.Addr); err != nil {
		return err
	}
	Cfg.NamenodeAddr = vp.GetString("NamenodeAddr")
	Cfg.StoragePath = vp.GetString("StoragePath")

	if err := vp.UnmarshalKey("Log", &Cfg.Log); err != nil {
		return err
	}
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
