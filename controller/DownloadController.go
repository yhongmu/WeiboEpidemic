package controller

import (
	"WeiboEpidemic/network"
	"io"
	"net/http"
	"os"
	"strconv"
)

var download *downloadController

func GetDownloadInstance() *downloadController {
	if download == nil {
		download = &downloadController{}
	}
	return download
}

type downloadController struct {
}

func (d *downloadController) Router(router *network.RouterHandler) {
	router.Router("/excel/AllHotSearchData", d.getAllData)
	router.Router("/excel/HotSearchKeyword", d.getKeywordData)
	router.Router("/excel/MonthHotSearchData", d.getMonthData)
}

func (d *downloadController) download(filename string, w http.ResponseWriter) {
	f, err := os.Open(filename)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	info, err := f.Stat()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	//_, contentType := getContentType(filename)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	//w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Type", "application/vnd.ms-excel")
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))

	f.Seek(0, 0)
	io.Copy(w, f)
}

func (d *downloadController) getAllData(w http.ResponseWriter, r *http.Request) {
	d.download("./excel/AllHotSearchData.xlsx", w)
}

func (d *downloadController) getMonthData(w http.ResponseWriter, r *http.Request) {
	d.download("./excel/MonthHotSearchDataAnalysis.xlsx", w)
}

func (d *downloadController) getKeywordData(w http.ResponseWriter, r *http.Request) {
	d.download("./excel/HotSearchKeywordFrequency.xlsx", w)
}
