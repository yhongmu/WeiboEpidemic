package utils

import (
	"WeiboEpidemic/entity"
	"WeiboEpidemic/log"
	"fmt"
	"github.com/yanyiwu/gojieba"
	"strings"
)

const (
	//深度过滤的相关度阈值
	RELATED_THRESHOLD_VALUE = 0.078
	//单个词的相关度上限值
	KEYWORD_THRESHOLD_VALUE = 0.06
)

var keyList = []string{
	"疫情", "新冠", "核酸", "疫苗", "新增", "隔离", "阳性", "阴性", "境外输入", "接种", "口罩",
	"确诊", "肺炎", "病毒", "病区", "病例", "战疫", "感染", "传染", "钟南山", "在家办公",
}

func SimpleFilter(keyword string) bool {
	for _, key := range keyList {
		if strings.Contains(keyword, key) {
			return true
		}
	}
	return false
}

var keywordsRatio = map[string]int8 {
	"新冠": 10,
	"核酸": 8,
	"肺炎": 5,
	"钟南山": 5,
	"疫苗": 4,
	"疫情": 3,
	"战疫": 3,
	"抗疫": 3,
	"病毒": 3,
	"隔离": 3,
	"阳性": 3,
	"确诊": 3,
	"口罩": 2,
	"病例": 2,
	"新增": 1,
	"传染": 1,
	"感染": 1,
	"防护": 1,
	"无症状": 1,
}

func ComplexFilter(text string, hotSearchData entity.HotSearchData) bool {

	keywordTimes := map[string]int64 {
		"新冠": 0,
		"核酸": 0,
		"肺炎": 0,
		"钟南山": 0,
		"疫苗": 0,
		"疫情": 0,
		"战疫": 0,
		"抗疫": 0,
		"病毒": 0,
		"隔离": 0,
		"阳性": 0,
		"确诊": 0,
		"口罩": 0,
		"病例": 0,
		"新增": 0,
		"传染": 0,
		"感染": 0,
		"防护": 0,
		"无症状": 0,
	}
	x := gojieba.NewJieba()
	defer x.Free()
	words := x.Cut(text, true)
	i := 0
	for i < len(words) {
		if !IsChineseChar(words[i]) {
			words = append(words[:i], words[i+1:]...)
		} else {
			i++
		}
	}

	wordAll := x.CutAll(text)
	for _, word := range wordAll {
		if times, found := keywordTimes[word]; found {
			keywordTimes[word] = times + int64(keywordsRatio[word])
		}
	}
	var sum int64
	for _, times := range keywordTimes {
		if keywordRelated := int64(KEYWORD_THRESHOLD_VALUE * float64(len(words))); times > keywordRelated {
			times = keywordRelated
		}
		sum += times
	}
	kwdes := float64(sum) / float64(len(words))
	if kwdes >= RELATED_THRESHOLD_VALUE {
		return true
	}
	if SimpleFilter(hotSearchData.Topic) {
		log.GetLog().Warning.Println(fmt.Sprintf(
			"topic=%s，url=%s \n		该热搜关键词通过简易筛选，但正文内容的相关度只有 %.4f",
			hotSearchData.Topic, hotSearchData.WeiboURL, kwdes))
	}
	return false
}
