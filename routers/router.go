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
		beego.NewNamespace("/hotNews",
			beego.NSInclude(
				&controllers.InfoController{},
			),
		)
	beego.AddNamespace(ns)
}
