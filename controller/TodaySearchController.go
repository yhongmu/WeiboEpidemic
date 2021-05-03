package controller

import (
	"WeiboEpidemic/entity"
	"WeiboEpidemic/log"
	"WeiboEpidemic/network"
	"WeiboEpidemic/utils"
	"net/http"
)

const WEIBO_HOT_SEARCH_URL = "https://s.weibo.com/top/summary?cate=realtimehot"

var todaySearch *todaySearchController

func GetTodaySearchInstance() *todaySearchController {
	if todaySearch == nil {
		todaySearch = &todaySearchController{}
	}
	return todaySearch
}

type todaySearchController struct {
}

func (t *todaySearchController) Router(router *network.RouterHandler) {
	router.Router("/weibo_epidemic/today_search", t.getTodaySearch)
}

func (t *todaySearchController) getTodaySearch(w http.ResponseWriter, r *http.Request)  {
	log.GetLog().Init("TodaySearch")
	hotSearchList, err := t.todaySearchRequest()
	if err != nil {
		network.ErrorShow(w, err)
	} else {
		network.ResultOK(w, 0, "今日新冠肺炎相关热搜获取成功！", hotSearchList)
	}

}

func (t *todaySearchController) todaySearchRequest() ([]entity.HotSearchEntity, error) {
	html, err := network.GetRequest(WEIBO_HOT_SEARCH_URL, nil,nil)
	if err != nil {
		return nil, err
	}
	hotSearchList, err := network.TodayHotSearchHTMLParse(html)
	if err != nil {
		return nil, err
	}
	var searchListCopy []entity.HotSearchEntity
	//通过对搜索关键词进行简易筛选
	for _, hotSearch := range hotSearchList {
		if utils.SimpleFilter(hotSearch.SearchKey) {
			searchListCopy = append(searchListCopy, hotSearch)
		}
	}
	if len(searchListCopy) > 0 {
		return searchListCopy, err
	} else {
		return nil, network.Errors{Code: -501, Msg: network.ErrorCode[-501]}
	}

}
