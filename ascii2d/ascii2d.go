package ascii2d

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	xpath "github.com/antchfx/htmlquery"
)

type Result struct {
	PixivID    int64
	Title      string
	ImageURL   string
	MemberName string
	MemberID   int64
}

func Ascii2d(image string) (*Result, error) {
	var (
		api = "https://ascii2d.net/search/uri"
	)
	transport := http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: &transport,
	}

	// 包装请求参数
	data := url.Values{}
	data.Set("uri", image) // 图片链接
	fromData := strings.NewReader(data.Encode())

	// 网络请求
	req, _ := http.NewRequest("POST", api, fromData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// 色合检索改变到特征检索
	var bovwUrl = strings.ReplaceAll(resp.Request.URL.String(), "color", "bovw")
	bovwReq, _ := http.NewRequest("POST", bovwUrl, nil)
	bovwReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	bovwReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.36")
	bovwResp, err := client.Do(bovwReq)
	if err != nil {
		return nil, err
	}
	defer bovwResp.Body.Close()
	// 解析XPATH
	doc, err := xpath.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	// 取出每个返回的结果
	list := xpath.Find(doc, `//div[@class="row item-box"]`)
	var link string
	// 遍历取出第一个返回的PIXIV结果
	for _, n := range list {
		linkPath := xpath.Find(n, `//div[2]/div[3]/h6/a[1]`)
		picPath := xpath.Find(n, `//div[1]/img`)
		if len(linkPath) != 0 && len(picPath) != 0 {
			link = xpath.SelectAttr(linkPath[0], "href")
			if strings.Contains(link, "www.pixiv.net") {
				break
			}
		}
	}
	// 链接取出PIXIV id
	var index = strings.LastIndex(link, "/")
	if link == "" || index == -1 {
		return nil, fmt.Errorf("Ascii2d not found")
	}
	id, _ := strconv.ParseInt(link[index+1:], 10, 64)
	if id == 0 {
		return nil, fmt.Errorf("convert to pid error")
	}
	// 根据PID查询插图信息
	/*
		var illust = &pixiv.Illust{}
		if err := illust.IllustInfo(id); err != nil {
			return nil, err
		}
		if illust.AgeLimit != "all-age" {
			return nil, fmt.Errorf("Ascii2d not found")
		}
	*/
	// 待完善
	return nil, nil
}
