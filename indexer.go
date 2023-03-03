package main

import (
	"io"
	"log"
	"net/http"
)

func indexHookHandle(r *http.Request) {
	switch r.Method {
	case "DELETE":
		//
		break
	case "PUT":
		indexPut(r)
	}
}

func indexPut(r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	// 解析html
	content := string(body)
	article := parseHtml(content)
	if article == nil {
		log.Printf("parse html err: %s\n", r.URL.Path)
		return
	}

	// 更新zin索引
	id := md5str(r.URL.Path)
	PutDoc(id, &ZincDocument{
		Filepath: r.URL.Path,
		Url:      extractSingleFileUrl(content),
		Title:    article.Title,
		Content:  article.TextContent,
		Size:     article.Length,
	})
}
