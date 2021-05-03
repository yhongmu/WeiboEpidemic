package dao

import (
	"fmt"
	"os"
)

type TodaySearchDao struct {
}

func CreateTodaySearchTable() {
	sql := `
	CREATE TABLE IF NOT EXISTS today_hot_search(
	    id MEDIUMINT UNSIGNED AUTO_INCREMENT,
	    search_key VARCHAR(50) NOT NULL COMMENT "搜索关键词",
	    heat INT NOT NULL COMMENT "热度",
	    real_url VARCHAR(255) NOT NULL COMMENT "实际的微博url",
	    create_time DATETIME,
	    PRIMARY KEY (id)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT="当日新冠疫情热搜表";`

	_, err := DB.Exec(sql)
	if err != nil {
		fmt.Println("Create today_hot_search Table error, error info: " + err.Error())
		os.Exit(0)
	}
	fmt.Println("Create today_hot_search Table success")
}
