package services

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"goback/models"
	"goback/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	DataIndex      int
	RAW_TO_TRANS   chan models.HotNew
	TRANS_TO_LABEL chan models.HotNew
	LABEL_TO_AREA  chan models.HotNew
	AREA_TO_EMO    chan models.HotNew
	EMO_TO_SCORE   chan models.HotNew
	SCORE_TO_STORE chan models.HotNew
)

func InitChannel() {
	RAW_TO_TRANS = make(chan models.HotNew, 200)
	TRANS_TO_LABEL = make(chan models.HotNew, 1)
	LABEL_TO_AREA = make(chan models.HotNew, 1)
	AREA_TO_EMO = make(chan models.HotNew, 1)
	EMO_TO_SCORE = make(chan models.HotNew, 1)
	SCORE_TO_STORE = make(chan models.HotNew, 1)
}

// FTemp -> models.HotNews
func FToH(temp models.FTemp) models.HotNew {
	h := models.HotNew{}
	h.Title = temp.Title
	h.Txt = temp.Txt
	h.Url = temp.Url
	h.Summary = temp.Summary
	h.GetTime = temp.GetTime
	h.PubTime = temp.PubTime
	h.Source = temp.Source
	h.Freq = temp.Freq
	h.KeyWords = temp.KeyWords
	h.Topic = temp.SchoolName
	h.PicUrl = temp.PicUrl
	h.PncUrl = temp.PngUrl
	return h
}
func RawToCh() {
	//确定从那个ID开始进行读取数据，获得数据库最大的id
	//批量读取数据500
	var temps []models.FTemp
	db := utils.InitDB()
	db.Model(models.FTemp{}).Order("GetTime desc").Limit(5).Find(&temps)
	for _, item := range temps {
		item.Freq = strings.Replace(item.Freq, "/n", "", -1)
		RAW_TO_TRANS <- FToH(item)
	}
}

// txt summary 翻译1
const (
	APP_ID       = "20180628000181143"
	SECURITY_KEY = "K4iQn27CSRm9EivZ6qhN"
	TRANS_AUTO   = "auto"
	TRANS_ZH     = "zh"
)

type TranslationResult struct {
	From        string `json:"from"`
	To          string `json:"to"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

func translate(q string, from, to string) (*TranslationResult, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sign := md5V1(APP_ID, q, timestamp, SECURITY_KEY)

	v := url.Values{}
	v.Set("q", q)
	v.Set("from", from)
	v.Set("to", to)
	v.Set("appid", APP_ID)
	v.Set("salt", timestamp)
	v.Set("sign", sign)

	urlStr := "http://api.fanyi.baidu.com/api/trans/vip/translate?" + v.Encode()

	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result TranslationResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	fmt.Println(len(result.TransResult))
	return &result, nil
}

func md5V1(str ...string) string {
	hash := md5.New()
	for _, s := range str {
		hash.Write([]byte(s))
	}
	return hex.EncodeToString(hash.Sum(nil))
}
func TransPipeLine() {
	for true {
		select {
		case msg := <-RAW_TO_TRANS:
			fmt.Println("获得到数据")
			if msg.Title != "" {
				if res, err := translate(msg.Title, TRANS_AUTO, TRANS_ZH); err != nil {
					msg.TitleCh = ""
				} else {
					msg.TitleCh = res.TransResult[0].Dst
				}
			}
			if msg.Txt != "" {
				if res, err := translate(msg.Txt, TRANS_AUTO, TRANS_ZH); err != nil {
					msg.TxtCh = ""
				} else {
					msg.TxtCh = res.TransResult[0].Dst
				}
			}
			if msg.Summary != "" {
				if res, err := translate(msg.Summary, TRANS_AUTO, TRANS_ZH); err != nil {
					msg.SummaryCh = ""
				} else {
					if len(res.TransResult) > 0 {
						msg.SummaryCh = res.TransResult[0].Dst
					}
					msg.SummaryCh = ""

				}
			}
			TRANS_TO_LABEL <- msg
		case <-QuitChan:
			logs.Info("翻译流程退出！！")
			return
		}
	}
}

// 话题标签生成2
func GetLabels(txt string) string {
	return "  "
}
func LabelsPipeLine() {
	for true {
		select {
		case msg := <-TRANS_TO_LABEL:
			msg.Labels = GetLabels(msg.Txt)
			LABEL_TO_AREA <- msg
		case <-QuitChan:

			logs.Info("话题标签生成流程退出！！")
			return
		}
	}

}

// 区域标签生成3
func GetArea(txt string) string {
	return "   "
}
func AreaPipeLine() {
	for true {
		select {
		case msg := <-LABEL_TO_AREA:
			msg.Area = GetArea(msg.Txt)
			AREA_TO_EMO <- msg
		case <-QuitChan:
			logs.Info("区域标签生成流程退出！！")
			return
		}
	}
}

// 情感标签生成4
func GetEmotion(sour string) (float64, error) {
	jsonStr := []byte(fmt.Sprintf(`{ "textType": "string", "token": "test","sour":"%s"}`, sour))
	url := "https://eae266ec46b040f9afb1ae22bef2676e.apig.cn-north-4.huaweicloudapis.com/v1/infers/240dd325-dfaf-4950-81a1-992f3aae0164/api/Mod"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Apig-AppCode", "2fbd1dee3ec64bf3a35c860027f00d84faa45118659841f3a28153759f78e2cc")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	statuscode := resp.StatusCode
	body, _ := ioutil.ReadAll(resp.Body)
	res := make(map[string]interface{}, 0)
	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println(err)
		return 0, err
	}
	modValue, err := strconv.ParseFloat(res["modValue"].(string), 2)
	fmt.Println(statuscode)
	return modValue, nil
}
func EmotionPipeLine() {
	for true {
		select {
		case msg := <-AREA_TO_EMO:
			if score, err := GetEmotion(msg.TxtCh); err != nil {
				msg.Score = 0
			} else {
				msg.Score = score
			}
			EMO_TO_SCORE <- msg
		case <-QuitChan:
			logs.Info("情感标签生成流程退出！！")
			return
		}
	}
}

// 评分5
func GetScore(msg models.HotNew) float64 {
	return 0.1
}
func ScorePipeLine() {
	for true {
		select {
		case msg := <-EMO_TO_SCORE:
			msg.Score = GetScore(msg)
			SCORE_TO_STORE <- msg
		case <-QuitChan:
			logs.Info("热度值生成流程退出！！")
			return
		}
	}
}

// 存储到models.HotNews
func StorePipeLine() {
	for true {
		select {
		case msg := <-SCORE_TO_STORE:
			//存储到models.HotNews create
			fmt.Println(msg.TitleCh, msg.TitleCh, msg.TxtCh)
			fmt.Println(msg.Freq)
		case <-QuitChan:
			logs.Info("存储流程退出！！")
			return
		}
	}
}

// https://zhuanlan.zhihu.com/p/589067307?utm_id=0  rpc
func Parse() {

	//FTemp -> models.HotNews
	//txt summary 翻译1
	for i := 0; i < 5; i++ {
		go TransPipeLine()
	}
	//话题标签生成2
	for i := 0; i < 5; i++ {
		go LabelsPipeLine()
	}
	//区域标签生成3
	for i := 0; i < 5; i++ {
		go AreaPipeLine()
	}
	//情感标签生成4
	for i := 0; i < 5; i++ {
		go EmotionPipeLine()
	}
	//评分5
	for i := 0; i < 5; i++ {
		go ScorePipeLine()
	}
	//存储到models.HotNews
	for i := 0; i < 5; i++ {
		go StorePipeLine()
	}
}

var QuitChan = make(chan int) //退出管道  监听退出信号后 退出客户端 和 goroutine
// 程序退出监听  程序退出前关闭所有Clinet 和 线程
func CloseAll() {
	//创建结束监听，以便及时关闭所有goroutine
	logs.Info("==================退出监听启动=========================")
	c := make(chan os.Signal, 1)
	//监听linux的 ctrl+c 信号  以及程序结束信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	for sig := range c {
		switch sig {
		// 获取退出信号时，关闭globalQuit, 让所有监听者退出
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL:
			close(QuitChan)
			time.Sleep(1 * time.Second)
			return
		}
	}
}
