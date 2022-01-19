package main

import (
	_ "github.com/udistrital/alternancia_mid/routers"
	"github.com/udistrital/auditoria"
	"github.com/udistrital/utils_oas/customerrorv2"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	apistatus "github.com/udistrital/utils_oas/apiStatusLib"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	AllowedOrigins := []string{"*.udistrital.edu.co"}
	if beego.BConfig.RunMode != "prod" {
		AllowedOrigins = []string{"*"}
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: AllowedOrigins,
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Origin", "x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	beego.ErrorController(&customerrorv2.CustomErrorController{})
	apistatus.Init()
	auditoria.InitMiddleware()
	beego.Run()
}
