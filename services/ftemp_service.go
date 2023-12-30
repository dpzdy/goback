package services

import (
	"fmt"
	"goback/models"
	"goback/utils"
	"time"
)

type FtempService interface {
	DataNumLine(timeGap int) (map[string]interface{}, error)
	GetNewsTotal() (int64, error)
	GetNumsOnTopic() ([]string, []int, error)
	GetNumsOnSource() ([]string, []int, error)
	GetRealTimeNews(topic string) ([]models.FTemp, error)
	GetDateTendencyLine(interval int) ([]string, []int, error)
}
type ftempServiceMgr struct {
}
type LineResult struct {
	Day   string
	Count int
}

func (fs *ftempServiceMgr) DataNumLine(timeGap int) (map[string]interface{}, error) {
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

func (fs *ftempServiceMgr) GetNewsTotal() (int64, error) {
	beginTime := time.Now().Format("2006-01-02")
	endTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	var total int64
	db := utils.InitDB()
	db.Model(models.FTemp{}).
		Where("GetTime > ? AND GetTime < ?", beginTime, endTime).
		Count(&total)
	return total, nil
}

type TopicGroup struct {
	Topic string
	Total int
}

func (fs *ftempServiceMgr) GetNumsOnTopic() ([]string, []int, error) {
	beginTime := time.Now().Format("2006-01-02")
	endTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	db := utils.InitDB()

	var res []TopicGroup
	db.Table("FTemp").
		Select("schoolName as Topic ,count(*) as Total").
		Where("GetTime > ? AND GetTime < ?", beginTime, endTime).
		Order("Total desc").
		Group("schoolName").Scan(&res)
	topics, nums := make([]string, 0), make([]int, 0)
	for _, item := range res {
		topics = append(topics, item.Topic)
		nums = append(nums, item.Total)
	}
	return topics, nums, nil

	//rows,err := db.Table("FTemp").
	//	Select("schoolName as Topic ,count(*) as Total").
	//	Where("GetTime > ? AND GetTime < ?", beginTime, endTime).
	//	Group("schoolName").
	//	Order("Total desc").
	//	Rows()
	//defer rows.Close()
	//if err != nil{
	//	fmt.Println(err)
	//}
	//for rows.Next(){
	//	var res TopicGroup
	//
	//	err := rows.Scan(&res.Topic,&res.Total)
	//	if err != nil{
	//		fmt.Println(err)
	//	}
	//	fmt.Println(res)
	//}
}

type SourceGroup struct {
	Source string
	Total  int
}

func (fs *ftempServiceMgr) GetNumsOnSource() ([]string, []int, error) {
	beginTime := time.Now().Format("2006-01-02")

	endTime := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	db := utils.InitDB()

	var res []SourceGroup
	db.Table("FTemp").
		Select("Source ,count(*) as Total").
		Where("GetTime > ? AND GetTime < ? AND Source <> ?", beginTime, endTime, "").
		Order("Total desc").
		Group("Source").
		Scan(&res)
	sources, nums := make([]string, 0), make([]int, 0)
	for _, item := range res {
		sources = append(sources, item.Source)
		nums = append(nums, item.Total)
	}
	return sources, nums, nil
}
func (fs *ftempServiceMgr) GetRealTimeNews(topic string) ([]models.FTemp, error) {
	var infos []models.FTemp
	db := utils.InitDB()
	db.Table("FTemp").Select([]string{"Title", "Summary", "schoolName", "GetTime"})
	if topic != "" {
		db.Where("schoolName = ?", topic)
	}
	db.Order("GetTime desc").Limit(10).Find(&infos)
	for _, item := range infos {
		fmt.Println(item.Title, item.Summary, item.SchoolName, item.GetTime)
	}
	return infos, nil
}

type TimeGroup struct {
	Date  time.Time
	Total int
}

func (fs *ftempServiceMgr) GetDateTendencyLine(interval int) ([]string, []int, error) {
	//生成当天时间区间
	year, month, day := time.Now().Date()
	location, _ := time.LoadLocation("Asia/Shanghai") // 这一步把错误忽略了，时区用Shanghai是因为没有Beijing
	today := time.Date(year, month, day, 0, 0, 0, 0, location)
	temp := today
	var weekDays []string
	for today.Day() == day {
		timeStr := today.Format(time.DateTime)
		weekDays = append(weekDays, timeStr)
		today = today.Add(time.Duration(interval) * time.Hour)
	}

	timeStr := temp.AddDate(0, 0, 1).Format(time.DateTime)
	weekDays = append(weekDays, timeStr)
	fmt.Println(weekDays)
	var infos []models.FTemp
	db := utils.InitDB()
	db.Model(models.FTemp{}).
		Where("GetTime > ? AND GetTime < ?", weekDays[0], weekDays[len(weekDays)-1]).
		Find(&infos)
	res := make([]int, len(weekDays)-1)
	for _, item := range infos {
		cur := item.GetTime.Format(time.DateTime)
		for i := 0; i < len(weekDays)-1; i++ {
			if cur > weekDays[i] && cur < weekDays[i+1] {
				res[i]++
			}
		}
	}
	return weekDays[:len(weekDays)-1], res, nil
}
