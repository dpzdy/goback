package test

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"goback/services"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestPipeLine(t *testing.T) {
	services.InitChannel()
	services.RawToCh()
	services.Parse()
	for {

	}

}
func TestEmo(t *testing.T) {
	fmt.Println(services.GetEmotion("时光清浅新的一天总会如约而至白云轻轻的飘时光清浅新的一天总会如约而至白云轻轻的飘时光清浅新的一天总会如约而至白云轻轻的飘时光清浅新的一天总会如约而至白云轻轻的飘时光清浅新的一天总我去到底"))
}

// string appId = "20180628000181143";
// string password = "K4iQn27CSRm9EivZ6qhN";

const (
	APP_ID       = "20180628000181143"
	SECURITY_KEY = "K4iQn27CSRm9EivZ6qhN"
)

type TranslationResult struct {
	From        string `json:"from"`
	To          string `json:"to"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

func translate(q string, from, to string) (string, error) {
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
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result TranslationResult
	err = json.Unmarshal(body, &result)
	return result.TransResult[0].Dst, err
}

func md5V1(str ...string) string {
	hash := md5.New()
	for _, s := range str {
		hash.Write([]byte(s))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func TestBaidu(t *testing.T) {
	result, err := translate("China’s sustained efforts to optimize visa process enhance appeal for foreign investment", "auto", "zh")
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
