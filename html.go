package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

var API = "http://history.fengqi.io:3456/content"

type Article struct {
	Title       string `json:"title"`
	Byline      string `json:"byline"`
	Dir         string `json:"dir"`
	Lang        string `json:"lang"`
	Content     string `json:"content"`
	TextContent string `json:"textContent"`
	Length      int    `json:"length"`
	SiteName    string `json:"siteName"`
}

type ParseRequest struct {
	Content string `json:"content"`
}

// 调用 mozilla/readability 解析 html
func parseHtml(html string) *Article {
	req := &ParseRequest{
		Content: html,
	}

	body, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Post(API, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	body, _ = io.ReadAll(res.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	article := &Article{}
	err = json.Unmarshal(body, article)
	if err != nil {
		panic(err)
	}

	return article
}
