package models

import "time"

// FTemp represents the Beego model for the FTemp table.
type FTemp struct {
	//ID                int       `orm:"column(ID);auto;pk"`
	Title             string    `orm:"column(Title);size(500)"`
	TitleCh           string    `orm:"column(TitleCh);size(500)"`
	TitleScore        string    `orm:"column(Title_Score);size(50)"`
	ArticalKeyWords   string    `orm:"column(articalKeyWords);size(500)"`
	Txt               string    `orm:"column(Txt);type(text)"`
	TxtCh             string    `orm:"column(TxtCh);type(text)"`
	TxtScore          string    `orm:"column(Txt_Score);size(50)"`
	Url               string    `orm:"column(Url);size(500)"`
	Summary           string    `orm:"column(Summary);size(500)"`
	SummaryCh         string    `orm:"column(SummaryCh);size(500)"`
	TxtType           string    `orm:"column(TxtType);size(100)"`
	GetTime           time.Time `orm:"column(GetTime);type(datetime)"`
	PubTime           time.Time `orm:"column(PubTime);type(datetime)"`
	TopicType         string    `orm:"column(topicType);size(50)"`
	Source            string    `orm:"column(Source);size(50)"`
	NameList          string    `orm:"column(nameList);size(50)"`
	IsDown            int       `orm:"column(isdown)"`
	RelateNewsNumber  int       `orm:"column(relateNewsNumber)"`
	KeyWords          string    `orm:"column(KeyWords);size(50)"`
	SchoolName        string    `orm:"column(schoolName);size(50)"`
	Trace             int       `orm:"column(trace)"`
	Country           string    `orm:"column(country);size(500)"`
	Polar             string    `orm:"column(Polar);size(50)"`
	Author            string    `orm:"column(author);size(500)"`
	Language          string    `orm:"column(language);size(100)"`
	HtmlCode          string    `orm:"column(htmlCode);type(text)"`
	SemanticValue     float64   `orm:"column(semanticValue)"`
	MediaIntroduction string    `orm:"column(mediaIntroduction);size(500)"`
	ZuanfaNum         int       `orm:"column(zuanfaNum)"`
	PinlunNum         int       `orm:"column(pinlunNum)"`
	DianzanNum        int       `orm:"column(dianzanNum)"`
	ShouchangNum      int       `orm:"column(shouchangNum)"`
	TxtSeg            string    `orm:"column(txtSeg);type(text)"`
	IsTrans           int       `orm:"column(isTrans)"`
	PicUrl            string    `orm:"column(picUrl);type(text)"`
	Freq              string    `orm:"column(freq);type(text)"`
	PngUrl            string    `orm:"column(pngUrl);type(text)"`
	Domain            string    `orm:"column(domain);size(500)"`
	Continent         string    `orm:"column(continent);size(50)"`
}

// TableName sets the table name for the FTemp model.
func (m *FTemp) TableName() string {
	return "FTemp"
}
