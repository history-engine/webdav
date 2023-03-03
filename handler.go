package main

import (
	"context"
	"fmt"
	"golang.org/x/net/webdav"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	prefix   = "/"
	datetime = "2006-01-02 15:04:05"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\t%s\n", r.Method, r.URL.Path)

	if !authHandler(w, r) {
		return
	}

	handler := webdav.Handler{
		Prefix:     prefix,
		LockSystem: webdav.NewMemLS(),
		FileSystem: webdav.Dir(data),
	}

	if r.Method == "GET" {
		getHandler(handler, w, r)
		return
	}

	fileDir := filepath.Dir(data + r.URL.Path)
	if !pathExists(fileDir) {
		_ = os.MkdirAll(fileDir, 0755)
	}

	go indexHookHandle(cloneHttpRequest(r))

	handler.ServeHTTP(w, r)
}

func authHandler(w http.ResponseWriter, r *http.Request) bool {
	_username, _password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	if _username != username || _password != "admin" {
		http.Error(w, "WebDAV: need authorized!", http.StatusUnauthorized)
		return false
	}

	return true
}

func getHandler(handler webdav.Handler, w http.ResponseWriter, r *http.Request) {
	f, err := handler.FileSystem.OpenFile(context.TODO(), r.URL.Path, os.O_RDONLY, 0664)
	if err != nil {
		log.Printf("%s\t%s not found\n", r.Method, r.URL.Path)
		return
	}
	defer f.Close()

	if n := fileHandler(f, w); n > 0 {
		return
	}

	listHandler(f, w, r)
}

func listHandler(file webdav.File, w http.ResponseWriter, r *http.Request) {
	files, err := file.Readdir(-1)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, "<html>\n")
	_, _ = fmt.Fprintf(w, "<head><title>Index of %s</title></head>\n", r.URL.Path)
	_, _ = fmt.Fprintf(w, "<body>\n")
	_, _ = fmt.Fprintf(w, "<h1>Index of %s</h1>\n", r.URL.Path)
	_, _ = fmt.Fprintf(w, "<hr>\n")
	_, _ = fmt.Fprintf(w, "<pre>\n")

	_, _ = fmt.Fprintf(w, "<table border=\"0\" width=\"%s\">\n", "90%")
	_, _ = fmt.Fprintf(w, "<tr>\n")
	_, _ = fmt.Fprintf(w, "<td><a href=\"../\">../</a></td>\n")
	_, _ = fmt.Fprintf(w, "<td></td>\n")
	_, _ = fmt.Fprintf(w, "<td></td>\n")
	_, _ = fmt.Fprintf(w, "</tr>\n")

	for _, d := range files {
		name := d.Name()
		size := byteCountIEC(d.Size())
		if d.IsDir() {
			name += "/"
			size = "-"
		}

		_, _ = fmt.Fprintf(w, "<tr>\n")
		_, _ = fmt.Fprintf(w, "<td><a href=\"%s\" title=\"%s\">%s</a></td>\n", name, name, name)
		_, _ = fmt.Fprintf(w, "<td>%s</td>\n", d.ModTime().Format(datetime))
		_, _ = fmt.Fprintf(w, "<td>%s</td>\n", size)
		_, _ = fmt.Fprintf(w, "</tr>\n")
	}

	_, _ = fmt.Fprintf(w, "</table>\n")
	_, _ = fmt.Fprintf(w, "</pre>\n")
	_, _ = fmt.Fprintf(w, "<hr>\n")
	_, _ = fmt.Fprintf(w, "</body>\n")
	_, _ = fmt.Fprintf(w, "</html>\n")
}

func fileHandler(file webdav.File, w http.ResponseWriter) int64 {
	stat, _ := file.Stat()
	if stat.IsDir() {
		return 0
	}

	w.Header().Set("Content-Length", strconv.FormatUint(uint64(stat.Size()), 10))
	w.Header().Set("Content-Type", fileMime(stat.Name()))
	n, err := io.Copy(w, file)
	if err != nil {
		log.Fatalf("file %s handler err %v\n", stat.Name(), err)
	}

	return n
}
