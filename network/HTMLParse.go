package network

import (
	"WeiboEpidemic/entity"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"time"
)

func TodayHotSearchHTMLParse(html string) ([]entity.HotSearchEntity, error) {
	HOST := "https://s.weibo.com"
	var hotSearchList []entity.HotSearchEntity
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	doc.Find("#pl_top_realtimehot table tbody tr").Each(func(i int, s *goquery.Selection) {
		hotSearch := entity.HotSearchEntity{}
		hotSearch.SearchKey = s.Find(".td-02 a").Text()
		if urlStr, exists := s.Find(".td-02 a").Attr("href"); exists {
			hotSearch.RealURL = HOST + urlStr
		}
		heatStr := s.Find(".td-02 span").Text()
		if heatStr != "" {
			hotSearch.Heat, _ = strconv.ParseInt(heatStr, 10, 64)
			hotSearch.CreateTime = time.Now().Format("2006-01-02")
			//hotSearch.IsToday = true
			hotSearchList = append(hotSearchList, hotSearch)
		}
	})
	if hotSearchList == nil {
		return nil, Errors{Code: -201, Msg: ErrorCode[-201]}
	}
	return hotSearchList, nil
}

func PastHotSearchTextHTMLParse(htmlStr string) (string, error) {
	var text string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// 加上导语
	if daoyu := doc.Find("div.card-wrap").First().Text(); strings.Contains(daoyu, "导语") {
		daoyu = strings.Replace(daoyu, "\n", "", -1)
		daoyu = strings.TrimSpace(daoyu)
		text += daoyu
	}
	doc.Find("div.card-wrap[mid] .card .card-feed .content").Each(func(i int, s *goquery.Selection) {
		content := s.Find("p.txt[nick-name]").Last().Text()
		content = strings.Replace(content, "\n", "", -1)
		content = strings.TrimSpace(content)
		constWords := s.Find(".txt").Last().Find("a:has(i.wbicon)").Map(func(j int, s2 *goquery.Selection) string {
			return s2.Text()
		})
		for _, word := range constWords {
			content = strings.Replace(content, word, "", 1)
		}
		if content != "" {
			text += content
		}
	})
	if text == "" {
		return "", Errors{Code: -202, Msg: ErrorCode[-202]}
	}
	return text, nil
}
