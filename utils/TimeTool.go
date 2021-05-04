package utils

import "time"

/*
 * 获取新冠疫情开始日期到现在的前一天之间的日期列表
 */
func GetBetweenDates(startTime, endTime string) []string {
	//START_TIME := "2019-10-01"
	var d []string
	timeFormatTpl := "2006-01-02"
	date, err := time.Parse(timeFormatTpl, startTime)
	if err !=  nil {
		// 时间解析，异常
		return d
	}

	//date2Str := time.Now().Format(timeFormatTpl)
		d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		if dateStr == endTime {
			break
		}
		d = append(d, dateStr)
	}
	return d
}
