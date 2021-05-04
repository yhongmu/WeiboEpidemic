package entity

type HotSearchEntity struct {
	SearchKey string		`json:"search_key"`					// 热搜关键词
	Heat int64				`json:"heat"`						// 热度
	//TemporaryURL string		`json:"temporary_url,omitempty"`	// 临时的url
	RealURL string			`json:"real_url"`					// 实际的微博热搜url
	CreateTime string		`json:"create_time,omitempty"`		// 热搜上榜的时间，精确到天
	//IsToday bool			`json:"is_today"`					// 是否是当天的热搜数据
}

type HotSearchList []HotSearchEntity

func (h HotSearchList) Len() int {
	return len(h)
}

func (h HotSearchList) Less(i, j int) bool {
	if h[i].CreateTime == h[j].CreateTime {
		return h[i].Heat > h[j].Heat
	}
	return h[i].CreateTime < h[j].CreateTime
}

func (h HotSearchList) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
