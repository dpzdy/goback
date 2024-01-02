package controllers

import "goback/services"

type FTempController struct {
	BaseController
}

// @Title DataNumLine
// @Description 获取数据折线图变化
// @router /dataNumLine [get]
func (ft *FTempController) DataNumLine() {

	res, err := services.FtempServiceMgr.DataNumLine(7)
	if err != nil {
		ft.RespMsg(FAIL, "get all students "+err.Error())
	}
	ft.RespData(SUCCESS, res)
}

// @Title GetNewsTotal
// @Description 获取当天新闻总数
// @router /todaynum [get]
func (ft *FTempController) GetNewsTotal() {
	total, err := services.FtempServiceMgr.GetNewsTotal()
	if err != nil {
		ft.RespMsg(FAIL, err.Error())
	}
	ft.RespData(SUCCESS, total)

}

// @Title GetNumsOnTopic
// @Description 获取不同主题新闻的数量
// @router /topicnum [get]
func (ft *FTempController) GetNumsOnTopic() {
	topics, nums, err := services.FtempServiceMgr.GetNumsOnTopic()
	res := make(map[string]interface{})
	res["topic"] = topics
	res["num"] = nums
	if err != nil {
		ft.RespMsg(FAIL, err.Error())
	}
	ft.RespData(SUCCESS, res)

}

// @Title GetNumsOnSource
// @Description 获取不同媒体的数量
// @router /sourcenum [get]
func (ft *FTempController) GetNumsOnSource() {
	sources, nums, err := services.FtempServiceMgr.GetNumsOnSource()
	res := make(map[string]interface{})
	res["source"] = sources
	res["num"] = nums
	if err != nil {
		ft.RespMsg(FAIL, err.Error())
	}
	ft.RespData(SUCCESS, res)

}

// @Title GetRealTimeNews
// @Description 获取实时新闻
// @router /realtimenews/:topic [get]
// http://localhost:8080/hotnews/realtimenews?topic=%E4%B8%AD%E7%BE%8E
func (ft *FTempController) GetRealTimeNews() {
	topic := ft.GetString("topic")
	infos, err := services.FtempServiceMgr.GetRealTimeNews(topic)
	if err != nil {
		ft.RespMsg(FAIL, err.Error())
	}
	ft.RespData(SUCCESS, infos)

}

// @Title GetDateTendencyLine
// @Description 获取今日时间段内数据变化
// @router /datetendency/:interval [get]
func (ft *FTempController) GetDateTendencyLine() {
	interval, err := ft.GetInt("interval", 6)
	if err != nil {
		ft.RespMsg(FAIL, err.Error())
	}
	weekdays, nums, err := services.FtempServiceMgr.GetDateTendencyLine(interval)
	res := make(map[string]interface{})
	res["weekday"] = weekdays
	res["num"] = nums
	if err != nil {
		ft.RespMsg(FAIL, err.Error())
	}
	ft.RespData(SUCCESS, res)

}
