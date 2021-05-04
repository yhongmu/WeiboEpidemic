package network

import "fmt"

type Errors struct {
	Code int16
	Msg string
}

func (e Errors) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

var ErrorCode = map[int16]string{
	0: "成功",
	-1: "表示系统错误",

	-100: "传参错误",
	-101: "搜索框未输入信息！",
	-102: "爬取所有历史热搜，输入的管理员密码错误！",
	-103: "输入的时间区间不正确！",

	-200: "html解析错误",
	-201: "未获取到当天热搜",
	-202: "未获取历史热搜的正文内容",

	-300: "数据库操作错误",
	-301: "数据插入错误",
	-302: "数据更新错误",
	-303: "未查询到章节URL",

	-400: "网络请求错误",
	-401: "同步网络请求错误",
	-402: "未获取到指定日期的热搜数据",

	-500: "数据筛选错误",
	-501: "今日暂无新冠疫情相关对热搜！",
}