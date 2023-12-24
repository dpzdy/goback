package services

import (
	"fmt"
	"goback/utils"
	"time"
)

type FtempService interface {
	DataNumLine(timeGap int) (map[string]interface{}, error)
}
type FtempServiceMgr struct {
}
type LineResult struct {
	Day   string
	Count int
}

func (fs *FtempServiceMgr) DataNumLine(timeGap int) (map[string]interface{}, error) {
	// 获取一周前的时间
	db := utils.InitDB()
	beforeDays := 12
	// 获取一周内每天的数据量
	y, m, d := time.Now().AddDate(0, 0, -beforeDays).Date()
	before := fmt.Sprintf("%d-%02d-%02d", y, m, d)
	rows, err := db.Raw(`SELECT FORMAT(GetTime, 'yyyy-MM-dd') as day, COUNT(*) as count  FROM FTemp WHERE GetTime >= ? GROUP BY FORMAT(GetTime, 'yyyy-MM-dd')`,
		before).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 生成过去七天的日期
	var dates []string
	// 获取当前时间
	now := time.Now()
	// 生成过去七天的日期
	for i := 0; i < beforeDays; i++ {
		date := now.AddDate(0, 0, -i)
		// 将日期格式化为 "2006-01-02" 形式的字符串
		dateString := date.Format("2006-01-02")
		dates = append(dates, dateString)
	}
	for i := 0; i < len(dates)/2; i++ {
		dates[i], dates[len(dates)-1-i] = dates[len(dates)-1-i], dates[i]
	}
	resMap := make(map[string]int)
	for rows.Next() {
		var result LineResult
		err := rows.Scan(&result.Day, &result.Count)
		if err != nil {
			return nil, err
		}
		//results = append(results, result)
		resMap[result.Day] = result.Count
	}
	cnts := make([]int, 0)
	for _, date := range dates {
		cnts = append(cnts, resMap[date])
	}
	res := make(map[string]interface{}, 0)
	res["time"] = dates
	res["cnts"] = cnts
	return res, nil
}
