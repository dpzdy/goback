// @APIVersion 1.0.0
// @Title HotNewsAPI
// @Description 国际热点十大新闻API
// @Contact zhnag_yuu@outlook.com
package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"goback/controllers"
)

func init() {
	ns :=
		beego.NewNamespace("/v1",
			beego.NSInclude(
				&controllers.InfoController{},
			),
		)
	beego.AddNamespace(ns)
}
