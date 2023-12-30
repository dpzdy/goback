// @APIVersion 1.0.0
// @Title HotNewsAPI
// @Description HotNewsAPI
// @Contact zhang_yuu@outlook.com
package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"goback/controllers"
)

func init() {
	ns :=
		beego.NewNamespace("/hotnews",
			beego.NSInclude(
				&controllers.InfoController{},
			),
			beego.NSInclude(
				&controllers.FTempController{},
			),
		)
	beego.AddNamespace(ns)
}
