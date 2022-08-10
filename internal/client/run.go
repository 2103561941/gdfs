package client

import "github.com/gin-gonic/gin"


func run() {
	g := gin.Default()
	initRouter(g)

	

}