package controllers

import beego "github.com/beego/beego/v2/server/web"

type BaseController struct {
	beego.Controller
	controllerName string //当前控制名称
	actionName     string //当前action名称
}
type Result struct {
	Result string      `json:"result"`
	Data   interface{} `json:"data,omitempty" `
	Info   string      `json:"info"`
	Total  int         `json:"total,omitempty"`
}

const (
	SUCCESS    = "success"
	FAIL       = "fail"
	WARNING    = "warning"
	ERROR      = "error"
	PERMISSION = "permission"
	TIEMOUT    = "timeout" //登录超时专用
)

func (c *BaseController) JsonStop() {
	c.ServeJSON()
	c.StopRun()
}
func (c *BaseController) RespData(result string, data interface{}, msg string, total int) {
	c.Data["json"] = &Result{
		Result: result,
		Data:   data,
		Info:   msg,
		Total:  total,
	}
	c.JsonStop()
}
func (c *BaseController) RespMsg(result string, msg string) {
	c.Data["json"] = &Result{
		Result: result,
		Info:   msg,
	}
	c.JsonStop()
}
