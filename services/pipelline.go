package services

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"goback/models"
	"os"
	"os/signal"
	"syscall"
	"testing"
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
	for _, item := range temps {
		RAW_TO_TRANS <- FToH(item)
	}

}

// txt summary 翻译1
func trans(eng string) string {
	//调用翻译接口
	ch := eng + "/*"
	return ch
}
func TransPipeLine() {
	for true {
		select {
		case msg := <-RAW_TO_TRANS:
			if msg.Title != "" {
				msg.TitleCh = trans(msg.Title)
			}
			if msg.Txt != "" {
				msg.Txt = trans(msg.TxtCh)
			}
			if msg.Summary != "" {
				msg.Txt = trans(msg.SummaryCh)
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
func GetEmotion(txt string) float64 {
	//调用情感分析接口
	return 0.1
}
func EmotionPipeLine() {
	for true {
		select {
		case msg := <-AREA_TO_EMO:
			msg.Emotion = GetEmotion(msg.Txt)
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
			fmt.Println(msg)
		case <-QuitChan:
			logs.Info("存储流程退出！！")
			return
		}
	}
}

// https://zhuanlan.zhihu.com/p/589067307?utm_id=0  rpc
func parse() {
	InitChannel()
	//FTemp -> models.HotNews
	//GetPipeLine()
	//txt summary 翻译1
	//TransPipeLine()
	//话题标签生成2
	//LabelsPipeLine()
	//区域标签生成3
	//AreaPipeLine()
	//情感标签生成4
	//EmotionPipeLine()
	//评分5
	//ScorePipeLine()
	//存储到models.HotNews
	//StorePipeLine()
}

func TestPipeLine(t *testing.T) {
	parse()
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
