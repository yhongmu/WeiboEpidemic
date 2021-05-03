package dao

import (
	"WeiboEpidemic/entity"
	"WeiboEpidemic/log"
	"bytes"
	"crypto/md5"
	"fmt"
	"os"
	"sync"
)

var pastSearchInstance *pastSearchDao
var once sync.Once

func GetPastDao() *pastSearchDao {
	once.Do(func() {
		pastSearchInstance = &pastSearchDao{}
	})
	return pastSearchInstance
}

type pastSearchDao struct {
}

func CreatePastSearchTable() {
	sql := `
	CREATE TABLE IF NOT EXISTS past_hot_search(
	    id MEDIUMINT UNSIGNED AUTO_INCREMENT,
	    search_id VARCHAR(40) NOT NULL COMMENT "每个热搜的专属id，由关键词和日期的十六进制字符串组成",
	    search_key VARCHAR(50) NOT NULL COMMENT "搜索关键词",
	    heat INT NOT NULL COMMENT "热度",
	    real_url VARCHAR(255) COMMENT "实际的微博url",
	    create_time VARCHAR(12) NOT NULL COMMENT "精确到天",
	    UNIQUE (search_id),
	    PRIMARY KEY (id)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT="历史新冠疫情热搜表";`

	_, err := DB.Exec(sql)
	if err != nil {
		log.GetLog().Error.Println("Create past_hot_search Table error, error info: " + err.Error())
		os.Exit(0)
	}
	log.GetLog().Info.Println("Create past_hot_search Table success")
}

func (p pastSearchDao) InsertAllPastHotSearch(hotSearchList []entity.HotSearchEntity) error {
	sqlInsert := `insert ignore into past_hot_search(search_id, search_key, heat, real_url, create_time) values `
	var buffer bytes.Buffer
	buffer.WriteString(sqlInsert)
	for index, hotSearch := range hotSearchList {
		//将热搜关键词和日期转十六进制，得到热搜的专属id
		byteData := []byte(hotSearch.SearchKey + hotSearch.CreateTime)
		md5Has := md5.Sum(byteData)
		searchID := fmt.Sprintf("%x", md5Has)
		s := fmt.Sprintf(`("%s", "%s", "%d", "%s", "%s")`,
			searchID, hotSearch.SearchKey, hotSearch.Heat, hotSearch.RealURL, hotSearch.CreateTime)
		if index == len(hotSearchList) - 1 {
			buffer.WriteString(s + ";")
		} else {
			buffer.WriteString(s + ",")
		}
	}
	str := buffer.String()
	_, err := DB.Exec(str)
	if err != nil {
		log.GetLog().Error.Println("大量的历史热搜数据插入失败！error info： " + err.Error())
		return err
	}
	log.GetLog().Info.Println("大量的历史热搜数据插入成功！")
	return nil
}

func (p pastSearchDao) SelectHotSearchByTimeLimit(startTime string, endTime string, heat int) ([]entity.HotSearchEntity, error) {
	sqlSelect := fmt.Sprintf("select search_key,heat,real_url,create_time from past_hot_search where create_time >= '%s' and create_time <= '%s' and heat >= %d;",
		startTime, endTime, heat)
	rows, err := DB.Query(sqlSelect)
	if err != nil {
		log.GetLog().Error.Println("查询历史热搜数据失败，error info：" + err.Error())
		return nil, err
	}
	var hotSearchList []entity.HotSearchEntity
	for rows.Next() {
		var search entity.HotSearchEntity
		err = rows.Scan(&search.SearchKey, &search.Heat, &search.RealURL, &search.CreateTime)
		if err != nil {
			log.GetLog().Error.Println("查询历史热搜数据失败，error info：" + err.Error())
			return nil, err
		}
		hotSearchList = append(hotSearchList, search)
	}
	log.GetLog().Info.Println("查询历史热搜数据成功！")
	return hotSearchList, nil
}
