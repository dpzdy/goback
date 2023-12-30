package main

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "goback/routers"
	"goback/services"
)

func main() {
	services.ServiceInit()
	//fmt.Println(beego.BConfig.Listen.EnableAdmin)
	fmt.Println(beego.BConfig.RunMode)

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.BConfig.WebConfig.EnableDocs = true

	// 跨域解决方案
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		// 允许访问所有源
		AllowAllOrigins: true,
		// 可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// 指的是允许的Header的种类
		AllowHeaders: []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		// 公开的HTTP标头列表
		ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		// 如果设置，则允许共享身份验证凭据，例如cookie
		AllowCredentials: true,
	}))
	//beego.InsertFilter("*", beego.AfterExec, ResFilter)

	beego.Run()
}

//0821 https://www.books;tack.cn/read/beego-2.0-zh/mvc-controller-config.md
//8022 https://www.bookstack.cn/read/beego-2.0-zh/mvc-model-models.md#emgt3z

var ResFilter = func(ctx *context.Context) {
	res := ctx.ResponseWriter
	res.Header().Set("X-Frame-Options", "SAMEORIGIN")
	res.Header().Set("X-Content-Type-Options", "nosniff")
	res.Header().Set("X-XSS-Protection", "1; mode=block")
	res.Header().Set("Content-Security-Policy", "object-src 'self'")
	res.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubdomains; preload")
	res.Header().Set("Referrer-Policy", "origin")
	res.Header().Set("X-Download-Options", "noopen")
	res.Header().Set("X-Permitted-Cross-Domain-Policies", "master-only")
	res.Header().Set("SameSite", "none")
	res.Header().Set("Secure", "true")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	cookies := ctx.Request.Cookies()
	if cookies != nil {
		str := ""
		for _, cookie := range cookies {
			str += cookie.Name + "=" + cookie.Value + ";" + "Secure=true;" + "HttpOnly;SameSite=none" //Cookie设置Secure标识  //Cookie设置HttpOnly
		}
		ctx.SetCookie("Set-Cookie", str)
	}
}
