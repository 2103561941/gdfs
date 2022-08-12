package namenode

import (
	"log"

	"github.com/cyb0225/gdfs/internal/namenode/server"
)

func Run() {
	if err := server.RunServer("50051"); err != nil {
		log.Fatal("start rpc server failed: ", err.Error())
	}
}
