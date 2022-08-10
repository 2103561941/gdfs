package client

import (
	"github.com/gin-gonic/gin"
)

// 注册路由
func initRouter(g *gin.Engine) {
	initMiddleware(g)
	initController(g)
}

// 注册全局中间件
func initMiddleware(g *gin.Engine) {

}

// 注册路由服务,分组中间件
func initController(g *gin.Engine) {

}
