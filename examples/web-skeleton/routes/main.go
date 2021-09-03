package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mix-go/web-skeleton/controllers"
	"github.com/mix-go/web-skeleton/middleware"
)

func Load(router *gin.Engine) {
	router.Use(gin.Recovery()) // error handle

	router.GET("hello",
		func(ctx *gin.Context) {
			hello := controllers.HelloController{}
			hello.Index(ctx)
		},
	)

	router.Any("users/add",
		middleware.SessionMiddleware(),
		func(ctx *gin.Context) {
			user := controllers.UserController{}
			user.Add(ctx)
		},
	)

	router.Any("login", func(ctx *gin.Context) {
		login := controllers.LoginController{}
		login.Index(ctx)
	})

	router.GET("websocket",
		func(ctx *gin.Context) {
			ws := controllers.WebSocketController{}
			ws.Index(ctx)
		},
	)
}
