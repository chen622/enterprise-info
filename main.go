package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func GetInfo(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("DNT", "1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://www.tianyancha.com/search?key=^%^E7^%^99^%^BE^%^E5^%^BA^%^A6^%^E7^%^BD^%^91^%^E8^%^AE^%^AF")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,zh-TW;q=0.6")
	req.Header.Set("Cookie", "")
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("请求异常: %s", err.Error()))
	}
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("请求状态异常: %d", resp.StatusCode))
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("解析数据异常: %s", err.Error()))
	}
	var result string
	doc.Find("div.sup-ie-company-header-child-1").Each(func(i int, s1 *goquery.Selection) {
		label := s1.ChildrenFiltered(".label").First().Text()
		if label == "电话：" {
			result += fmt.Sprintf("%s%s\n", label,
				s1.ChildrenFiltered(".link-hover-click").First().Text())
		} else if label == "网址：" {
			result += fmt.Sprintf("%s%s\n", label,
				s1.ChildrenFiltered(".company-link").First().Text())
		}
	})

	doc.Find("div.sup-ie-company-header-child-2").Each(func(i int, s1 *goquery.Selection) {
		label := s1.ChildrenFiltered(".label").First().Text()
		if label == "地址：" {
			str := strings.ReplaceAll(s1.ChildrenFiltered("#company_base_info_address").First().Text(), " ", "")
			result += fmt.Sprintf("%s%s\n", label,
				str)
		}
	})

	result += fmt.Sprintln("简介：" + strings.ReplaceAll(doc.Find("#company_base_info_detail").First().Text(), " ", ""))

	result += fmt.Sprintln("法人：" + doc.Find(".humancompany").Find(".link-click").Text())

	result += "股东：\n"
	//fmt.Println(doc.Find("[tyc-event-ch='CompangyDetail.gudong.ziranren']").First().Text())
	doc.Find("[tyc-event-ch='CompangyDetail.gudongxinxi']").First().Find("tbody").First().Children().Each(
		func(i int, tr *goquery.Selection) {
			var name string
			var tags []string
			var percent string
			tr.Children().Each(func(i2 int, td *goquery.Selection) {
				if i2 == 1 {
					name = tr.Find("[tyc-event-ch='CompangyDetail.gudong.ziranren']").First().Text()
					if name == "" {
						name = tr.Find("[tyc-event-ch='CompangyDetail.gudong.gongsi']").First().Text()
					}
					tr.Find(".tag-common").Each(func(i3 int, s3 *goquery.Selection) {
						tags = append(tags, s3.Text())
					})
				} else if i2 == 2 {
					percent = td.Find("span").First().Text()
				}
			})
			result += fmt.Sprintf("  - %s(%s): %s\n", name, percent, strings.Join(tags, " "))
		},
	)

	return result, nil
}

func main() {
	fmt.Printf("请输入公司名称：")
	var key string
	_, err := fmt.Scanf("%s", &key)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("搜索中...")
	searchResults, err := SearchTY(key)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("\n\n\n\n\n")
	for index, result := range searchResults {
		fmt.Printf("%d. 公司名称：%s\n", index+1, result.Name)
	}
	fmt.Printf("请选择一家公司: ")
	var index int
	_, err = fmt.Scanf("%d", &index)
	if err != nil {
		log.Fatalln(err)
	}
	if index <= 0 || index > len(searchResults) {
		log.Fatalln("选项不存在")
	}
	fmt.Println("搜索中...")
	info, err := GetInfo(searchResults[index-1].Url)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("\n\n\n\n\n\n名称：" + searchResults[index-1].Name)
	fmt.Println(info)
}
