package utils

import "time"

/*
 * 获取新冠疫情开始日期到现在的前一天之间的日期列表
 */
func GetBetweenDates() []string {
	//START_TIME := "2019-10-01"
	START_TIME := "2020-01-01"
	var d []string
	timeFormatTpl := "2006-01-02"
	date, err := time.Parse(timeFormatTpl, START_TIME)
	if err !=  nil {
		// 时间解析，异常
		return d
	}

	//date2Str := time.Now().Format(timeFormatTpl)
	date2Str := "2020-01-11"
		d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		if dateStr == date2Str {
			break
		}
		d = append(d, dateStr)
	}
	return d
}
