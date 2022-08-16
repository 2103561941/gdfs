package log

import (
	"github.com/cyb0225/gdfs/pkg/log"
	"github.com/spf13/viper"
)

type LogConfig struct {
	Module     string `yaml:"Module"`
	LogPath    string `yaml:"LogPath"`
	MaxSize    int    `yaml:"MaxSize"`
	MaxBackups int    `yaml:"MaxBackups"`
	MaxAge     int    `yaml:"MaxAge"`
	Compress   bool   `yaml:"Compress"`
}

func Setup() error {
	vp := viper.New()
	vp.SetConfigName("public")
	vp.AddConfigPath(".")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return err
	}

	cfg := LogConfig{}
	if err := vp.UnmarshalKey("Log", &cfg); err != nil {
		return err
	}

	NewLogger(&cfg)

	return nil
}

func NewLogger(cfg *LogConfig) {
	level := levelTreansform(cfg.Module)

	tops := []log.TeeOption{
		{
			Filename: cfg.LogPath + "access.log",
			Ropt: log.RotateOptions{
				MaxSize:    cfg.MaxSize,
				MaxAge:     cfg.MaxAge,
				MaxBackups: cfg.MaxBackups,
				Compress:   cfg.Compress,
			},
			Lef: func(l log.Level) bool {
				return l <= log.InfoLevel && l >= level
			},
		},
		{
			Filename: cfg.LogPath + "error.log",
			Ropt: log.RotateOptions{
				MaxSize:    cfg.MaxSize,
				MaxAge:     cfg.MaxAge,
				MaxBackups: cfg.MaxBackups,
				Compress:   cfg.Compress,
			},
			Lef: func(l log.Level) bool {
				return l > log.InfoLevel && l >= level
			},
		},
	}

	log.NewLogger(tops, level)
}

func levelTreansform(level string) log.Level {
	if level == "debug" {
		return log.DebugLevel
	}
	if level == "release" || level == "info" {
		return log.InfoLevel
	}
	// default module: debug
	return log.DebugLevel
}
