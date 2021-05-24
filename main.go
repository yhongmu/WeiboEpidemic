package main

import (
	"WeiboEpidemic/controller"
	"WeiboEpidemic/dao"
	"WeiboEpidemic/log"
	"WeiboEpidemic/network"
	"net/http"
	"time"
)

func main() {
	dao.InitDB()
	//dao.CreateTodaySearchTable()
	dao.CreatePastSearchTable()
	//utils.CreateExcelTable()
	server := http.Server{
		Addr: ":8090",
		Handler: network.GetRouterInstance(),
		ReadTimeout: 5 * time.Second,
	}
	//注册当日微博新冠疫情热搜API路由
	controller.GetTodaySearchInstance().Router(network.GetRouterInstance())
	//注册历史微博新冠疫情热搜API路由
	controller.GetPastSearchInstance().Router(network.GetRouterInstance())
	//注册报表下载的API路由
	controller.GetDownloadInstance().Router(network.GetRouterInstance())

	err := server.ListenAndServe()
	if err != nil {
		log.GetLog().Error.Println("start server error!")
	}
	log.GetLog().Info.Println("start server success!")
}
