package namenode

import "github.com/spf13/viper"

// load configs
func init() {
	// the config file path is based on the root directory.
	// and it can


	viper.Set("port", "50051")
}
