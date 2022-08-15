package namenode

import "github.com/spf13/viper"

// load configs
func init() {

	viper.Set("chunckSize", 1024) // Byte 
	viper.Set("port", "50051")
}
