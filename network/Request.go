package network

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

var options = cookiejar.Options{
	PublicSuffixList: publicsuffix.List,
}

var jar, _ = cookiejar.New(&options)

var Client = &http.Client{
	Timeout: 15 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
	Jar: jar,
}

func GetRequest(requestURL string, urlValues url.Values, header map[string]string ) (string, error) {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if urlValues != nil {
		reqBody := urlValues.Encode()
		req.URL.RawQuery = reqBody
	}
	for key, value := range header {
		req.Header.Set(key, value)
	}
	return Response(req)
}

func PostRequest(requestURL string, urlValues url.Values) (string, error) {
	reqBody := urlValues.Encode()
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(reqBody))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//req.URL.RawQuery = reqBody
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return Response(req)
}

func Response(req *http.Request) (string, error) {
	resp, err := Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
