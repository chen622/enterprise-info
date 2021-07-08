package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

type SearchResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func SearchTY(key string) ([]*SearchResult, error) {

	client := &http.Client{}
	u, err := url.Parse("https://www.tianyancha.com/search")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("key", key)
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("请求异常: %s", err.Error()))
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("DNT", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.tianyancha.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,zh-TW;q=0.6")
	req.Header.Set("Cookie", "")
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("请求异常: %s", err.Error()))
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("请求状态异常: %d", resp.StatusCode))
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("解析数据异常: %s", err.Error()))
	}
	var results []*SearchResult
	doc.Find("a.select-none").Each(func(i int, selection *goquery.Selection) {
		href, exist := selection.Attr("href")
		if exist {
			results = append(results, &SearchResult{Name: selection.Text(), Url: href})
		}
	})
	return results, nil
}
