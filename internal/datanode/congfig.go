package datanode

import "github.com/spf13/viper"

// load configs
func init() {

	viper.Set("namenodeAddr", "127.0.0.1:50051")
	viper.Set("port", 50052)
	viper.Set("directory", "./storage")
	viper.Set("heartbeat", 5) // seconds
}
