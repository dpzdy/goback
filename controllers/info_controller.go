package controllers

import "goback/services"

type InfoController struct {
	BaseController
}

// @Title GetAll
// @router /all [get]
func (sc *InfoController) GetAll() {

	stus, err := services.AllInfo()
	if err != nil {
		sc.RespMsg(FAIL, "get all students "+err.Error())
	}
	sc.RespData(SUCCESS, stus, "get all students", len(stus))

}
