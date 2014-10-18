package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func upload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := http.Post(
		"https://dubhacks-worker-1.ngrok.com/create_compile_ipa",
		"application/octet-stream", file)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		log.Println(err)
		return
	}

	sum := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":"+port, nil)
}
