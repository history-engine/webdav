package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	ZincIndex     = "web_history"
	ZincDocPutApi = "http://zincsearch.fengqi.io/api/%s/_doc/%s"
)

type ZincDocument struct {
	Filepath string `json:"filepath"`
	Url      string `json:"url"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Size     int    `json:"size"`
}

// PutDoc 创建或更新文档
func PutDoc(id string, doc *ZincDocument) {
	body, _ := json.Marshal(doc)

	api := fmt.Sprintf(ZincDocPutApi, ZincIndex, id)
	req, _ := http.NewRequest("PUT", api, bytes.NewReader(body))

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "123456")

	res, err := http.DefaultClient.Do(req)
	if err != nil || res == nil {
		fmt.Printf("put zinc index err: %v\n", err)
		return
	}
	defer res.Body.Close()

	body, _ = io.ReadAll(res.Body)
	fmt.Printf("put zinc index: %s\n", res.Status)
}
