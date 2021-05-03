package controller

import (
	"WeiboEpidemic/dao"
	"WeiboEpidemic/entity"
	"WeiboEpidemic/log"
	"WeiboEpidemic/network"
	"WeiboEpidemic/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	PAST_SEARCH_URL = "https://google-api.zhaoyizhe.com/google-api/index/mon/sec"
	WEIBO_INFO_URL = "https://s.weibo.com/weibo?q=%%23%s%%23"
)

var pastSearch *pastSearchController
var wg sync.WaitGroup
var mutex sync.Mutex
var flag bool	//微博内容的html解析失败后，再次递归请求的标志符，防止无限递归

func GetPastSearchInstance() *pastSearchController {
	if pastSearch == nil {
		pastSearch = &pastSearchController{}
	}
	return pastSearch
}

type pastSearchController struct {
}

func (t *pastSearchController) Router(router *network.RouterHandler) {
	router.Router("/weibo_epidemic/past_search", t.getPastSearchDB)
	router.Router("/weibo_epidemic/past_search_crawler", t.pastSearchCrawler)
	//router.Router("/weibo_epidemic/past_search_update", )
}

func (t *pastSearchController) getPastSearchDB(w http.ResponseWriter, r *http.Request)  {
	log.GetLog().Init("pastSearch")
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	heatStr := r.PostFormValue("heat")
	if heatStr == "" {
		heatStr = "0"
	}
	heat, err := strconv.Atoi(heatStr)
	if err != nil {
		network.ErrorShow(w, err)
		return
	}
	if startTime > endTime {
		network.ResultFail(w, -103, network.ErrorCode[-103])
		return
	}
	hotSearchList, err := dao.GetPastDao().SelectHotSearchByTimeLimit(startTime, endTime, heat)
	if err != nil {
		network.ErrorShow(w, err)
	} else {
		network.ResultOK(w, 0, "输入时间段的热搜时间查询成功！", hotSearchList)
	}
}

func (t *pastSearchController) pastSearchCrawler(w http.ResponseWriter, r *http.Request)  {
	log.GetLog().Init("pastSearchCrawler")
	password := r.PostFormValue("password")
	if password != "1017" {
		network.ResultFail(w, -102, network.ErrorCode[-102])
		return
	}
	t1 := time.Now()
	hotSearchList, err := t.crawlerRequest()
	if err != nil {
		network.ErrorShow(w, err)
		return
	}
	err = dao.GetPastDao().InsertAllPastHotSearch(hotSearchList)
	if err != nil {
		network.ErrorShow(w, err)
		return
	}
	elapsed := time.Since(t1)
	message := fmt.Sprintf("指定日期的新冠肺炎相关热搜成功存入数据库！耗时 %.2fs", float64(elapsed)/1e9)
	network.ResultOK(w, 0, message, hotSearchList)
}

func (t *pastSearchController) crawlerRequest() ([]entity.HotSearchEntity, error) {
	jsonMap, err := t.goRequestHotSearch()
	if err != nil {
		return nil, err
	}
	hotSearchList, err := t.goFilterSearch(jsonMap)
	if err != nil {
		return nil, err
	}
	return hotSearchList, nil
}

/*
 * 并发爬取多个日期的微博热搜数据
 *
 */
func (t *pastSearchController) goRequestHotSearch() (map[string]string, error) {
	// 获取新冠疫情开始日期到现在的前一天之间的日期列表
	dateList := utils.GetBetweenDates()
	jsonMap := make(map[string]string)
	header := make(map[string]string)
	header["Origin"] = "http://weibo.zhaoyizhe.com"
	debug.SetMaxThreads(1000)
	taskDate := make(chan string, 50)
	taskJSON := make(chan string, 50)
	wg.Add(len(dateList))
	for _, date := range dateList {
		date := date
		// 爬取某个日期的热搜数据
		go func() {
			defer wg.Done()
			urlValues := url.Values{}
			urlValues.Add("date", date)
			jsonStr, err := network.GetRequest(PAST_SEARCH_URL, urlValues, header)
			if err != nil {
				return
			}
			mutex.Lock()
			{
				taskDate <- date
				taskJSON <- jsonStr
			}
			mutex.Unlock()
		}()
	}
	for range dateList {
		jsonMap[<-taskDate] = <-taskJSON
	}
	close(taskDate)
	close(taskJSON)
	wg.Wait()
	if len(jsonMap) > 0 {
		return jsonMap, nil
	}
	return nil, network.Errors{Code: -402, Msg: network.ErrorCode[-402]}
}

/*
 * 先爬取微博的正文内容，通过内容筛选出新冠疫情相关热搜
 *
 */
func (t *pastSearchController) goFilterSearch(jsonMap map[string]string) ([]entity.HotSearchEntity, error){
	var hotSearchList []entity.HotSearchEntity
	for date, jsonStr := range jsonMap {
		var jsonDec entity.JsonDecEntity
		err := json.Unmarshal([]byte(jsonStr), &jsonDec)
		if err != nil || jsonDec.Code != 1 {
			println(fmt.Sprintf("date: %s 出现json数据错误！", date))
			continue
		}
		flag = true
		list := t.goFilterSearchMain(date, jsonDec.Data)
		hotSearchList = append(hotSearchList, list...)
	}
	if hotSearchList != nil {
		return hotSearchList, nil
	}
	return nil, network.Errors{Code: -202, Msg: network.ErrorCode[-202]}
}

func (t *pastSearchController) goFilterSearchMain(dateStr string, datas []entity.HotSearchData) []entity.HotSearchEntity {
	var hotSearchList []entity.HotSearchEntity
	var errorSearchList []entity.HotSearchData
	taskHTML := make(chan string, 50)
	taskData := make(chan entity.HotSearchData, 50)
	wg.Add(len(datas))
	for _, hotSearchData := range datas {
		hotSearchData := hotSearchData
		hotSearchData.WeiboURL = fmt.Sprintf(WEIBO_INFO_URL, url.QueryEscape(hotSearchData.Topic))
		// 获取热搜的正文内容
		go func() {
			defer wg.Done()
			htmlStr, err := network.GetRequest(hotSearchData.WeiboURL,nil, nil)
			if err != nil {
				log.GetLog().Error.Println(
					fmt.Sprintf(
						"date=%s，topic=%s，爬取微博正文内容报错，error=%s", dateStr, hotSearchData.Topic, err))
				return
			}
			mutex.Lock()
			{
				taskHTML <- htmlStr
				taskData <- hotSearchData
			}
			mutex.Unlock()
		}()
	}
	//记录html解析失败的原因
	var parseErrorTimes int
	// 解析获取到到热搜html，并筛选出新冠疫情相关到热搜
	for range datas {
		hotSearchData := <-taskData
		h := <-taskHTML
		text, err := network.PastHotSearchTextHTMLParse(h)
		if err != nil {
			if flag && strings.Contains(h, "正在检测访问环境") {
				parseErrorTimes++
			}
			if hotSearchData.ReqNumber < entity.REQ_NUMBER_LIMIT {
				hotSearchData.ReqNumber++
				errorSearchList = append(errorSearchList, hotSearchData)
			} else {
				// 检查html解析失败的原因是否是微博在进行安全检测
				if !flag {
					log.GetLog().Error.Println(
						fmt.Sprintf(
							"date=%s，topic=%s，连续%d次都无法解析出该热搜的正文内容，以下是它的html内容\n%s\n",
							dateStr, hotSearchData.Topic, entity.REQ_NUMBER_LIMIT+1, h))
				}
			}
			continue
		}
		if utils.ComplexFilter(text, hotSearchData) {
			hotSearch := entity.HotSearchEntity{}
			hotSearch.Heat = hotSearchData.HotNumber
			hotSearch.RealURL = hotSearchData.WeiboURL
			hotSearch.SearchKey = hotSearchData.Topic
			hotSearch.CreateTime = dateStr
			hotSearchList = append(hotSearchList, hotSearch)
		}
	}
	close(taskHTML)
	close(taskData)
	wg.Wait()
	// 如果html解析失败的原因是微博进行安全检测的话，那等待5分钟后再请求
	if parseErrorTimes > 5 || parseErrorTimes == len(datas){
		flag = false
		log.GetLog().Error.Println(fmt.Sprintf("date=%s，html解析失败的原因是微博进行安全检测，此时等待5分钟后再次请求", dateStr))
		time.Sleep(5 * time.Minute)
		return t.goFilterSearchMain(dateStr, datas)
	}
	sleepTime := time.Nanosecond
	switch {
	case len(datas) < 10:
		sleepTime = 1
	case len(datas) < 60:
		sleepTime = 2
	default:
		sleepTime = 3
	}
	log.GetLog().Error.Println(fmt.Sprintf("请求date=%s的热搜数据结束，请等待%d分钟", dateStr, sleepTime))
	time.Sleep(sleepTime * time.Minute)
	//递归，再次请求html解析失败度网页
	if errorSearchList != nil {
		log.GetLog().Error.Println(fmt.Sprintf("第%d次 有%d条数据需要重新请求", errorSearchList[0].ReqNumber, len(errorSearchList)))
		list := t.goFilterSearchMain(dateStr, errorSearchList)
		hotSearchList = append(hotSearchList, list...)
	}
	return hotSearchList
}

