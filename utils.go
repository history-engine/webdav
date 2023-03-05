package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"regexp"
)

func cloneHttpRequest(r *http.Request) *http.Request {
	body, _ := io.ReadAll(r.Body)
	req := r.Clone(r.Context())
	r.Body = io.NopCloser(bytes.NewReader(body))
	req.Body = io.NopCloser(bytes.NewReader(body))
	return req
}

func md5str(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func extractSingleFileUrl(content string) string {
	regex, _ := regexp.Compile(`(?s)<!--.*?(htt.+://\S+).*?saved\sdate.*?-->`)
	matches := regex.FindStringSubmatch(content)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}
