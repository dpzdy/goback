package controllers

import "goback/services"

type FTempController struct {
	BaseController
}

// @Title GetAll
// @router /all [get]
func (ft *FTempController) DataNumLine() {

	fsm := &services.FtempServiceMgr{}
	res, err := fsm.DataNumLine(7)
	if err != nil {
		ft.RespMsg(FAIL, "get all students "+err.Error())
	}
	ft.RespData(SUCCESS, res, "get all students", 0)
}
