package entity

type JsonDecEntity struct {
	Code int				`json:"code"`
	Message string			`json:"message"`
	Data []HotSearchData	`json:"data"`
}

const REQ_NUMBER_LIMIT int8 = 2

type HotSearchData struct {
	//ID string				`json:"_id"`
	Topic string 			`json:"topic"`
	HotNumber int64			`json:"hotNumber"`
	WeiboURL string
	ReqNumber int8			//重新爬取热搜网页信息的次数
}
