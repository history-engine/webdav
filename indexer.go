package main

import (
	"fmt"
	"io"
	"net/http"
)

func indexHookHandle(r *http.Request) {
	// SingleFile只有PUT MKCOL PROPFIND 3种调用
	switch r.Method {
	case "PUT":
		indexPut(r)
		break

	default:
		// todo
		break
	}
}

func indexPut(r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read html content err: %v\n", err)
		return
	}
	defer r.Body.Close()

	// 解析html
	content := string(body)
	article := parseHtml(content)
	if article == nil {
		fmt.Printf("parse html content err: %s\n", r.URL.Path)
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
