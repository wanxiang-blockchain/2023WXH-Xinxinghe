package server

import (
	"fmt"
	"triple_star/config"
	"triple_star/controller"
	"triple_star/dao"
	"triple_star/logger"
	"triple_star/middleware"
	"triple_star/service"
	router "triple_star/util/util_router"
)

func httpServer() {
	r := router.NewRouter()
	r.Use(middleware.StopPanic)
	r.Use(middleware.HandleException)
	r.Use(middleware.HandleCorsRequest)
	r.Use(middleware.Auth)
	r.Add("/v1/triple_star/user/register", controller.User.Register)
	r.Add("/v1/triple_star/user/login", controller.User.Login)
	r.Add("/v1/triple_star/user/logout", controller.User.Logout)
	r.Add("/V1/triple_star/user/upload", controller.Data.Upload)
	r.Add("/v1/triple_star/data/add", controller.Data.Add)
	r.Add("/v1/triple_star/data/query", controller.Data.Query)
	r.Add("/v1/triple_star/data/buy", controller.Data.Buy)
	r.Add("/v1/triple_star/data/queryOneself", controller.Data.QueryOneSelf)

	r.Run(fmt.Sprintf(":%d", config.Config.Http.Port))
}

func Start(logLevel string) {
	logger.Init(logLevel)
	config.Init()
	dao.Init()
	service.Init()

	httpServer()
}
