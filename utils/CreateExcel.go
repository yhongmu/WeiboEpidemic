package utils

import (
	"WeiboEpidemic/dao"
	"WeiboEpidemic/entity"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/yanyiwu/gojieba"
	"sort"
	"strconv"
)

func CreateExcelTable()  {
	hotSearchList, err := dao.GetPastDao().SelectHotSearchByTimeLimit("2020-01-01", "2021-05-01", 0)
	if err != nil {
		println(err)
		return
	}
	sort.Sort(entity.HotSearchList(hotSearchList)) //对切片数据排序
	createTable1(hotSearchList)
	createTable2(hotSearchList)
	createTable3(hotSearchList)
}

func createTable1(hotSearchList []entity.HotSearchEntity) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, _ = file.AddSheet("Sheet1")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "id"
	cell =row.AddCell()
	cell.Value = "热搜关键句"
	cell =row.AddCell()
	cell.Value = "热度"
	cell =row.AddCell()
	cell.Value = "微博内容url"
	cell =row.AddCell()
	cell.Value = "热搜日期"
	for i, hotSearch := range hotSearchList {
		sheet = file.Sheet["Sheet1"]
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = strconv.Itoa(i+1)
		cell = row.AddCell()
		cell.Value = hotSearch.SearchKey
		cell = row.AddCell()
		cell.Value = strconv.Itoa(int(hotSearch.Heat))
		cell = row.AddCell()
		cell.Value = hotSearch.RealURL
		cell = row.AddCell()
		cell.Value = hotSearch.CreateTime
	}
	err = file.Save("./excel/AllHotSearchData.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		fmt.Println("所有历史新冠热搜数据报表创建成功")
	}
}

func createTable2(hotSearchList []entity.HotSearchEntity) {
	type monthHotSearchInfo struct {
		monthStr string
		searchTime int
		keyword1 string
		word1Time int
		keyword2 string
		word2Time int
		keyword3 string
		word3Time int
		keyword4 string
		word4Time int
		keyword5 string
		word5Time int
		keyword6 string
		word6Time int
		keyword7 string
		word7Time int
		keyword8 string
		word8Time int
		keyword9 string
		word9Time int
	}
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, _ = file.AddSheet("Sheet2")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "当前月份"
	cell =row.AddCell()
	cell.Value = "新冠热搜条数"
	cell =row.AddCell()
	cell.Value = "热搜关键词1"
	cell =row.AddCell()
	cell.Value = "词1出现次数"
	cell =row.AddCell()
	cell.Value = "热搜关键词2"
	cell =row.AddCell()
	cell.Value = "词2出现次数"
	cell =row.AddCell()
	cell.Value = "热搜关键词3"
	cell =row.AddCell()
	cell.Value = "词3出现次数"
	cell =row.AddCell()
	cell.Value = "热搜关键词4"
	cell =row.AddCell()
	cell.Value = "词4出现次数"
	cell =row.AddCell()
	cell.Value = "热搜关键词5"
	cell =row.AddCell()
	cell.Value = "词5出现次数"
	cell =row.AddCell()
	cell.Value = "热搜关键词6"
	cell =row.AddCell()
	cell.Value = "词6出现次数"
	cell =row.AddCell()
	cell.Value = "热搜关键词7"
	cell =row.AddCell()
	cell.Value = "词7出现次数"
	cell =row.AddCell()
	cell.Value = "热搜关键词8"
	cell =row.AddCell()
	cell.Value = "词8出现次数"
	cell =row.AddCell()
	cell.Value = "热搜关键词9"
	cell =row.AddCell()
	cell.Value = "词9出现次数"
	month := "2020-01"
	j := 0
	var infoList []monthHotSearchInfo
	keywordMap := make(map[string]int)
	x := gojieba.NewJieba()
	defer x.Free()
	for i, hotSearch := range hotSearchList {
		j++
		words := x.Cut(hotSearch.SearchKey, true)
		for _, word := range words {
			if len(word) <= 3 {
				continue
			}
			if _, found := keywordMap[word]; found {
				keywordMap[word]++
			} else {
				keywordMap[word] = 1
			}
		}

		if month != hotSearch.CreateTime[0:7] || i == len(hotSearchList) - 1 {
			var times []int
			var keys []string
			for k := 0; k < 9; k++ {
				max := 0
				flag := ""
				for key, value := range keywordMap {
					if max < value {
						max = value
						flag = key
					}
				}
				times = append(times, max)
				keys = append(keys, flag)
				delete(keywordMap, flag)
			}
			monthInfo := monthHotSearchInfo{}
			monthInfo.monthStr = month
			monthInfo.searchTime = j
			monthInfo.keyword1 = keys[0]
			monthInfo.keyword2 = keys[1]
			monthInfo.keyword3 = keys[2]
			monthInfo.keyword4 = keys[3]
			monthInfo.keyword5 = keys[4]
			monthInfo.keyword6 = keys[5]
			monthInfo.keyword7 = keys[6]
			monthInfo.keyword8 = keys[7]
			monthInfo.keyword9 = keys[8]
			monthInfo.word1Time = times[0]
			monthInfo.word2Time = times[1]
			monthInfo.word3Time = times[2]
			monthInfo.word4Time = times[3]
			monthInfo.word5Time = times[4]
			monthInfo.word6Time = times[5]
			monthInfo.word7Time = times[6]
			monthInfo.word8Time = times[7]
			monthInfo.word9Time = times[8]
			infoList = append(infoList, monthInfo)
			keywordMap = make(map[string]int)
			month = hotSearch.CreateTime[0:7]
			j = 0
		}
	}

	for _, info := range infoList {
		sheet = file.Sheet["Sheet2"]
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = info.monthStr
		cell = row.AddCell()
		cell.Value = strconv.Itoa(info.searchTime)
		cell = row.AddCell()
		cell.Value = info.keyword1
		cell = row.AddCell()
		cell.Value = strconv.Itoa(info.word1Time)
		cell = row.AddCell()
		cell.Value = info.keyword2
		cell = row.AddCell()
		cell.Value = strconv.Itoa(info.word2Time)
		cell = row.AddCell()
		cell.Value = info.keyword3
		cell = row.AddCell()
		cell.Value = strconv.Itoa(info.word3Time)
		cell = row.AddCell()
		cell.Value = info.keyword4
		cell = row.AddCell()
		cell.Value = strconv.Itoa(info.word4Time)
		cell = row.AddCell()
		cell.Value = info.keyword5
		cell = row.AddCell()
		cell.Value = strconv.Itoa(info.word5Time)
		cell = row.AddCell()
		cell.Value = info.keyword6
		cell = row.AddCell()
		cell.Value = strconv.Itoa(info.word6Time)
		cell = row.AddCell()
		cell.Value = info.keyword7
		cell = row.AddCell()
		cell.Value = strconv.Itoa(info.word7Time)
		cell = row.AddCell()
		cell.Value = info.keyword8
		cell = row.AddCell()
		cell.Value = strconv.Itoa(info.word8Time)
		cell = row.AddCell()
		cell.Value = info.keyword9
		cell = row.AddCell()
		cell.Value = strconv.Itoa(info.word9Time)
	}

	err = file.Save("./excel/MonthHotSearchDataAnalysis.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		fmt.Println("每月的新冠热搜分析报表创建成功")
	}

}

func createTable3(hotSearchList []entity.HotSearchEntity) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, _ = file.AddSheet("Sheet2")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "热搜关键词"
	cell =row.AddCell()
	cell.Value = "关键词出现次数"
	keywordMap := make(map[string]int)
	x := gojieba.NewJieba()
	defer x.Free()
	for _, hotSearch := range hotSearchList {
		words := x.Cut(hotSearch.SearchKey, true)
		for _, word := range words {
			if len(word) <= 3 {
				continue
			}
			if _, found := keywordMap[word]; found {
				keywordMap[word]++
			} else {
				keywordMap[word] = 1
			}
		}
	}
	l := len(keywordMap)
	for k := 0; k < l; k++ {
		max := 0
		flag := ""
		for key, value := range keywordMap {
			if max < value {
				max = value
				flag = key
			}
		}
		if max > 10 {
			row = sheet.AddRow()
			cell = row.AddCell()
			cell.Value = flag
			cell =row.AddCell()
			cell.Value = strconv.Itoa(max)
		}
		delete(keywordMap, flag)
	}
	err = file.Save("./excel/HotSearchKeywordFrequency.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		fmt.Println("热搜关键词出现频率排序报表创建成功")
	}
}
