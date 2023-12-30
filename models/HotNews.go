package models

import "time"

type HotNew struct {
	//Id       int
	Title     string    `orm:"column(Title);size(200)"`
	Txt       string    `orm:"column(Txt)"`
	Url       string    `orm:"column(Url);size(500)"`
	Summary   string    `orm:"column(Summary);size(500)"`
	GetTime   time.Time `orm:"column(GetTime);type(datetime)"`
	PubTime   time.Time `orm:"column(PubTime);type(datetime)"`
	Source    string    `orm:"column(Source);size(100)"`
	Freq      string    `orm:"column(Freq);size(500)"`
	Area      string    `orm:"column(Area);size(50)"`
	Labels    string    `orm:"column(Labels);size(100)"`
	Score     float64   `orm:"column(Score)"`
	KeyWords  string    `orm:"column(KeyWords);size(50)"`
	Topic     string    `orm:"column(Topic);size(50)"`
	PicUrl    string    `orm:"column(PicUrl)"`
	PncUrl    string    `orm:"column(PncUrl)"`
	TitleCh   string    `orm:"column(TitleCh);size(200)"`
	TxtCh     string    `orm:"column(TxtCh)"`
	SummaryCh string    `orm:"column(SummaryCh);size(500)"`
	Emotion   float64   `orm:"column(Emotion)"`
}

func (m *HotNew) TableName() string {
	return "HotNews"
}
