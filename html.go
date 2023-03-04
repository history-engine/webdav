package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	body, _ := json.Marshal(req)
	res, err := http.DefaultClient.Post(API, "application/json", bytes.NewReader(body))
	if err != nil || res == nil {
		fmt.Printf("request readability err: %v\n", err)
		return nil
	}
	defer res.Body.Close()

	body, _ = io.ReadAll(res.Body)
	article := &Article{}
	_ = json.Unmarshal(body, article)

	return article
}
