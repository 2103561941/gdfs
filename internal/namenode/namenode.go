package namenode

import (
	"log"

	"github.com/cyb0225/gdfs/internal/namenode/server"
	"github.com/spf13/viper"
)


// load config file and start rpc server 
func Run() {
	port := viper.GetString("port")

	if err := server.RunServer(port); err != nil {
		log.Fatal("start rpc server failed: ", err.Error())
	}
}
