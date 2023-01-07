package main

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "go-postgresql/docs"
	"go-postgresql/internal/api"
	"go-postgresql/internal/repository"
	"go-postgresql/internal/services"
)

//	@title			Coupon Api
//	@version		1.0
//	@description	Coupon Api
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@BasePath		/coupon/api
//	@schemes		http
func main() {
	e := echo.New()
	repository := repository.NewCouponRepository()
	service := services.NewCouponService(repository)
	api.NewHandler(e, service)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(
		e.Start(":1323"),
	)
}
